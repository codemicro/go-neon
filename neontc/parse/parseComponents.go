package parse

import (
	"bytes"
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"strings"
)

func parseFuncTokens(fs *FileSet, tokens *tokenSet) (*ast.FuncDeclNode, error) {
	returnValue := new(ast.FuncDeclNode)

	// TODO: It would be helpful to actually parse the function signature here.
	token := tokens.Next()
	opWord, operand := chopToken(token.cont)
	if !bytes.Equal(opWord, []byte("func")) {
		panic("impossible state: cannot parse a function from something that isn't a function declaration")
	}
	returnValue.Signature = string(operand)

tokenLoop:
	for token := tokens.Next(); token != nil; token = tokens.Next() {

		if bytes.HasPrefix(token.cont, []byte("{[")) {
			// This is a substitution (I hope)
			returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.SubstitutionNode{
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
			default:
				return nil, fmt.Errorf("%s: unsupported opword %q inside of function", fs.ResolvePosition(token.pos), opWord)
			}

		} else {
			returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.PlaintextNode{
				Plaintext: string(token.cont),
			})
		}

	}

	return returnValue, nil
}

func parseCodeTokens(fs *FileSet, tokens *tokenSet) (*ast.RawCodeNode, error) {
	returnValue := new(ast.RawCodeNode)

	token := tokens.Next()
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
