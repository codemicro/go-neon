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

		if bytes.HasPrefix(token.cont, []byte("{{")) {
			opWordB, operandB := chopToken(token.cont)
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
					codeNode, err := parseCodeToken(fs, token)
					if err != nil {
						return nil, err
					}
					returnValue.ChildNodes = append(returnValue.ChildNodes, codeNode)
				default:
					return nil, fmt.Errorf("%s: unsupported opword %q inside of function", fs.ResolvePosition(token.pos), opWord)
				}

			}

		} else {
			returnValue.ChildNodes = append(returnValue.ChildNodes, &ast.PlaintextNode{
				Plaintext: string(token.cont),
			})
		}

	}

	return returnValue, nil
}

func parseCodeToken(fs *FileSet, token *rawToken) (*ast.RawCodeNode, error) {
	returnValue := new(ast.RawCodeNode)
	opWord, operand := chopToken(token.cont)
	if !bytes.Equal(opWord, []byte("code")) {
		panic("impossible state: cannot parse a code block from something that isn't a code block")
	}
	returnValue.GoCode = strings.Trim(string(operand), "\n")
	return returnValue, nil
}
