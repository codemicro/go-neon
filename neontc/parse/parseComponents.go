package parse

import (
	"bytes"
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"regexp"
	"strings"
)

var funcSignatureRegexp = regexp.MustCompile(`(\(\w+ \*?\w+?\) )?(\w+)(\[(?:\w+ \w+,? ?)+\])?(\((?:\w+ \*?\w+,? ?)*\)) ?(\w+|\((?:\*?\w+,? ?)+\))?`)

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

	type (
		loopStart struct {
			pos  int64
			expr string
		}
		loopEnd struct {
			pos int64
		}
	)

tokenLoop:
	for token := tokens.Next(); token != nil; token = tokens.Next() {

		if bytes.HasPrefix(token.cont, []byte("{[")) {
			// This is a substitution (I hope)
			returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.SubstitutionNode{
				Pos:        ast.Pos(token.pos),
				Expression: strings.Trim(string(token.cont), " {}[]"),
			})
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
		// FIFO stack holding the indexes of each loop node once created
		var (
			loopNodePositions []int
			n                 int
		)
		for _, node := range returnValue.ChildNodes {
			switch typedNode := node.(type) {
			case *loopStart:
				returnValue.ChildNodes[n] = &ast.LoopNode{
					Pos:            ast.Pos(typedNode.pos),
					LoopExpression: typedNode.expr,
					ChildNodes:     nil,
				}
				loopNodePositions = append(loopNodePositions, n)
				n += 1
			case *loopEnd:
				if len(loopNodePositions) == 0 {
					return nil, fmt.Errorf("%s: unexpected endfor", fs.ResolvePosition(typedNode.pos))
				}
				loopNodePositions = loopNodePositions[:len(loopNodePositions)-1]
			default:
				if len(loopNodePositions) != 0 {
					idx := loopNodePositions[len(loopNodePositions)-1]
					returnValue.ChildNodes[idx].(*ast.LoopNode).ChildNodes = append(returnValue.ChildNodes[idx].(*ast.LoopNode).ChildNodes, node)
				} else {
					returnValue.ChildNodes[n] = node
					n += 1
				}
			}
		}
		if len(loopNodePositions) != 0 {
			return nil, fmt.Errorf(
				"%s: unclosed for statement",
				fs.ResolvePosition(
					int64(returnValue.ChildNodes[loopNodePositions[len(loopNodePositions)-1]].(*ast.LoopNode).Pos),
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
