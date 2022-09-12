package parse

import (
	"bytes"
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"regexp"
	"strings"
)

var funcSignatureRegexp = regexp.MustCompile(`(\(\w+ \*?\[?\]?\w+?\) )?(\w+)(\[(?:\w+ \*?\[?\]?\w+,? ?)+\])?(\((?:\w+ \*?\[?\]?\w+,? ?)*\)) ?(\[?\]?\w+|\((?:\*?\[?\]?\w+,? ?)+\))?`)

func parseFuncTokens(fs *FileSet, tokens *tokenSet) (*ast.FuncDeclNode, error) {
	returnValue := new(ast.FuncDeclNode)

	token := tokens.Next()
	returnValue.Pos = ast.Pos(token.pos)

	opWord, operand := chopToken(token.cont)
	if !bytes.Equal(opWord, []byte("func")) {
		panic("impossible state: cannot parse a function from something that isn't a function declaration")
	}
	returnValue.Signature = string(operand)

	if !funcSignatureRegexp.MatchString(returnValue.Signature) {
		return nil, fmt.Errorf("%s: invalid function signature", fs.ResolvePosition(token.pos))
	}

	submatches := funcSignatureRegexp.FindStringSubmatch(returnValue.Signature)
	if len(submatches[5]) != 0 {
		return nil, fmt.Errorf("%s: functions may not have return values", fs.ResolvePosition(token.pos))
	}

	returnValue.Identifier = submatches[2]

	type (
		loopStart struct {
			pos  int64
			expr string
		}
		loopEnd struct {
			pos int64
		}

		ifStart struct {
			pos  int64
			expr string
		}
		elifNode struct {
			pos  int64
			expr string
		}
		elseNode struct {
			pos int64
		}
		ifEnd struct {
			pos int64
		}
	)

tokenLoop:
	for token := tokens.Next(); token != nil; token = tokens.Next() {

		if bytes.HasPrefix(token.cont, []byte("{[")) {
			// This is a substitution (I hope)
			subToken, err := parseSubsitution(fs, token)
			if err != nil {
				return nil, err
			}

			returnValue.ChildNodes = append(returnValue.ChildNodes, subToken)
		} else if bytes.HasPrefix(token.cont, []byte("{{")) {
			opWordB, operandB := chopToken(token.cont)
			opWord, operand := string(opWordB), string(operandB)

			switch opWord {
			case "endfunc":
				if operand != "" {
					return nil, fmt.Errorf("%s: endfunc takes no arguments", fs.ResolvePosition(token.pos))
				}
				break tokenLoop
			case "code":
				tokens.Rewind()
				codeNode, err := parseCodeTokens(fs, tokens)
				if err != nil {
					return nil, err
				}
				returnValue.ChildNodes = append(returnValue.ChildNodes, codeNode)
			case "for":
				returnValue.ChildNodes = append(returnValue.ChildNodes, &loopStart{pos: token.pos, expr: operand})
			case "endfor":
				returnValue.ChildNodes = append(returnValue.ChildNodes, &loopEnd{pos: token.pos})
			case "if":
				returnValue.ChildNodes = append(returnValue.ChildNodes, &ifStart{pos: token.pos, expr: operand})
			case "elif":
				returnValue.ChildNodes = append(returnValue.ChildNodes, &elifNode{pos: token.pos, expr: operand})
			case "else":
				returnValue.ChildNodes = append(returnValue.ChildNodes, &elseNode{pos: token.pos})
			case "endif":
				returnValue.ChildNodes = append(returnValue.ChildNodes, &ifEnd{pos: token.pos})
			default:
				return nil, fmt.Errorf("%s: unsupported opword %q inside of function", fs.ResolvePosition(token.pos), opWord)
			}

		} else {
			returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.PlaintextNode{
				Pos:       ast.Pos(token.pos),
				Plaintext: string(token.cont),
			})
		}

	}

	// balance loop starts and ends
	{
		// FIFO stack
		var (
			loopNodes   []*ast.LoopNode
			ifNodes     []*ast.ConditionalNode
			n           int
			addToOutput = func(node ast.Node) {
				var idx int

				if len(loopNodes) != 0 && len(ifNodes) != 0 {
					lastLoop := loopNodes[len(loopNodes)-1]
					lastIf := ifNodes[len(ifNodes)-1]

					if lastIf.Pos == lastLoop.Pos {
						panic("impossible state: if and loop in the same position")
					} else if lastIf.Pos > lastLoop.Pos {
						goto addToIf
					}
					goto addToLoop
				} else if len(loopNodes) != 0 {
					goto addToLoop
				} else if len(ifNodes) != 0 {
					goto addToIf
				} else {
					goto addToStandard
				}

			addToIf:
				idx = len(ifNodes) - 1
				ifNodes[idx].ChildNodes = append(ifNodes[idx].ChildNodes, node)
				goto resume
			addToLoop:
				idx = len(loopNodes) - 1
				loopNodes[idx].ChildNodes = append(loopNodes[idx].ChildNodes, node)
				goto resume
			addToStandard:
				returnValue.ChildNodes[n] = node
				n += 1
			resume:
			}
		)

		for _, node := range returnValue.ChildNodes {
			switch typedNode := node.(type) {
			case *loopStart:
				newNode := &ast.LoopNode{
					Pos:            ast.Pos(typedNode.pos),
					LoopExpression: typedNode.expr,
					ChildNodes:     nil,
				}

				addToOutput(newNode)
				loopNodes = append(loopNodes, newNode)
			case *loopEnd:
				if len(loopNodes) == 0 {
					return nil, fmt.Errorf("%s: unexpected endfor", fs.ResolvePosition(typedNode.pos))
				}
				loopNodes = loopNodes[:len(loopNodes)-1]
			case *ifStart:
				newNode := &ast.ConditionalNode{
					Pos:        ast.Pos(typedNode.pos),
					Expression: typedNode.expr,
				}

				addToOutput(newNode)
				ifNodes = append(ifNodes, newNode)
			case *elifNode:
				if len(ifNodes) == 0 {
					return nil, fmt.Errorf("%s: elif without parent if statement", fs.ResolvePosition(typedNode.pos))
				}
				newNode := &ast.ConditionalNode{
					Pos:        ast.Pos(typedNode.pos),
					Expression: typedNode.expr,
				}
				idx := len(ifNodes) - 1
				ifNodes[idx].Else = newNode
				ifNodes[idx] = newNode
			case *elseNode:
				if len(ifNodes) == 0 {
					return nil, fmt.Errorf("%s: else without parent if statement", fs.ResolvePosition(typedNode.pos))
				}
				newNode := &ast.ConditionalNode{
					Pos: ast.Pos(typedNode.pos),
				}
				idx := len(ifNodes) - 1
				ifNodes[idx].Else = newNode
				ifNodes[idx] = newNode
			case *ifEnd:
				if len(ifNodes) == 0 {
					return nil, fmt.Errorf("%s: unexpected ifend", fs.ResolvePosition(typedNode.pos))
				}
				ifNodes = ifNodes[:len(ifNodes)-1]
			default:
				addToOutput(node)
			}
		}
		if len(loopNodes) != 0 {
			return nil, fmt.Errorf(
				"%s: unclosed for statement",
				fs.ResolvePosition(
					int64(loopNodes[len(loopNodes)-1].Pos),
				),
			)
		}

		returnValue.ChildNodes = returnValue.ChildNodes[:n]
	}

	return returnValue, nil
}

func parseCodeTokens(fs *FileSet, tokens *tokenSet) (*ast.RawCodeNode, error) {
	returnValue := new(ast.RawCodeNode)

	token := tokens.Next()
	returnValue.Pos = ast.Pos(token.pos)

	opWord, operand := chopToken(token.cont)
	if !bytes.Equal(opWord, []byte("code")) {
		panic("impossible state: cannot extract raw code from something that isn't a code block")
	}

	if len(operand) != 0 {
		return nil, fmt.Errorf("%s: code block declarations cannot contain inline code", fs.ResolvePosition(token.pos))
	}

	var rawCode []byte
	for {
		token := tokens.Next()
		opword, operand := chopToken(token.cont)
		if bytes.Equal(opword, []byte("endcode")) {
			if len(operand) != 0 {
				return nil, fmt.Errorf("%s: code block declarations cannot contain inline code", fs.ResolvePosition(token.pos))
			}
			break
		}
		rawCode = append(rawCode, token.cont...)
	}

	returnValue.GoCode = strings.Trim(string(rawCode), "\n")
	return returnValue, nil
}

var importStatementRegexp = regexp.MustCompile(`import (.|[A-Za-z]\w+)? ?"([\w./]+)"`)

func parseImportToken(fs *FileSet, token *rawToken) (*ast.ImportNode, error) {
	returnValue := new(ast.ImportNode)
	returnValue.Pos = ast.Pos(token.pos)

	val := unwrapToken(token.cont)

	if !importStatementRegexp.Match(val) {
		return nil, fmt.Errorf("%s: could not parse import statement", fs.ResolvePosition(token.pos))
	}

	submatches := importStatementRegexp.FindStringSubmatch(string(val))

	returnValue.Alias = submatches[1]
	returnValue.ImportPath = submatches[2]

	return returnValue, nil
}

var substitutionModifierRegexp = regexp.MustCompile(` #([a-zA-Z,]+)$`) // includes support for comma seperated vals

func parseSubsitution(_ *FileSet, token *rawToken) (*ast.SubstitutionNode, error) {
	returnValue := &ast.SubstitutionNode{
		Pos: ast.Pos(token.pos),
	}

	expr := string(token.cont)
	expr = strings.TrimPrefix(expr, "{[")
	expr = strings.TrimSuffix(expr, "]}")
	expr = strings.TrimSpace(expr)

	submatches := substitutionModifierRegexp.FindStringSubmatch(expr)
	if len(submatches) != 0 {
		expr = substitutionModifierRegexp.ReplaceAllString(expr, "")

		rawModifier := strings.ToLower(submatches[1])
		returnValue.Modifier = ast.SubModMapping[rawModifier]
	}

	returnValue.Expression = expr

	return returnValue, nil
}
