package parse

import (
	"bytes"
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"strings"
)

func parseFuncTokens(tokens *tokenSet) (*ast.FuncDeclNode, error) {
	returnValue := new(ast.FuncDeclNode)

	// TODO: It would be helpful to actually parse the function signature here.
	token := tokens.Next()
	opWord, operand := chopToken(token)
	if !bytes.Equal(opWord, []byte("func")) {
		panic("impossible state: cannot parse a function from something that isn't a function declaration")
	}
	returnValue.Signature = string(operand)

tokenLoop:
	for token := tokens.Next(); token != nil; token = tokens.Next() {

		if bytes.HasPrefix(token, []byte("{{")) {
			opWordB, operandB := chopToken(token)
			opWord, operand := string(opWordB), string(operandB)

			if operand == "" {
				if opWord == "endfunc" {
					break tokenLoop
				}

				// This is a substitution (I hope)
				returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.SubstitutionNode{
					Expression: opWord,
				})
			} else {

				switch opWord {
				case "code":
					codeNode, err := parseCodeToken(token)
					if err != nil {
						return nil, err
					}
					returnValue.ChildNodes = append(returnValue.ChildNodes, codeNode)
				default:
					// TODO: error position??
					return nil, fmt.Errorf("unsupported opword %q inside of function", opWord)
				}

			}

		} else {
			returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.PlaintextNode{
				Plaintext: string(token),
			})
		}

	}

	return returnValue, nil
}

func parseCodeToken(token []byte) (*ast.RawCodeNode, error) {
	returnValue := new(ast.RawCodeNode)
	opWord, operand := chopToken(token)
	if !bytes.Equal(opWord, []byte("code")) {
		panic("impossible state: cannot parse a code block from something that isn't a code block")
	}
	returnValue.GoCode = strings.Trim(string(operand), "\n")
	return returnValue, nil
}
