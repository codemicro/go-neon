package parse

import (
	"bytes"
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
)

func getAtIndex[T any](x []T, i int) T {
	var out T
	if 0 <= i && i < len(x) {
		out = x[i]
	}
	return out
}

func chopToken(token []byte) (opWord, operand []byte) {
	trimmedToken := bytes.TrimPrefix(
		bytes.TrimSuffix(token, []byte("}}")),
		[]byte("{{"),
	)
	trimmedToken = bytes.TrimSpace(trimmedToken)

	firstSpace := bytes.IndexRune(trimmedToken, ' ')
	firstNewline := bytes.IndexRune(trimmedToken, '\n')

	if firstSpace == -1 && firstNewline == -1 {
		opWord = trimmedToken
	} else {
		var minimum int
		if firstSpace == -1 && firstNewline != -1 {
			minimum = firstNewline
		} else if firstNewline == -1 && firstSpace != -1 {
			minimum = firstSpace
		} else if firstSpace > firstNewline {
			minimum = firstNewline
		} else {
			minimum = firstSpace
		}

		opWord = trimmedToken[:minimum]
		operand = trimmedToken[minimum+1:]
	}

	return
}

// func File(fs *FileSet, fpath string, input []byte) (*ast.TemplateFile, error) {
func File(fpath string, input []byte) (*ast.TemplateFile, error) {
	tf := new(ast.TemplateFile)
	tf.Filepath = fpath

	// baseFileSetPosition := fs.AddFile(fpath, input)

	tokens, err := tokens(input)
	if err != nil {
		return nil, err
	}

	var i int
	for token := tokens.Next(); token != nil; token = tokens.Next() {

		if bytes.HasPrefix(token, []byte("{{")) {
			// curly block

			opWordB, _ := chopToken(token)
			opWord := string(opWordB)

			switch opWord {
			case "func":
				tokens.Rewind()
				funcDecl, err := parseFuncTokens(tokens)
				if err != nil {
					return nil, err
				}
				tf.Nodes = append(tf.Nodes, funcDecl)
			case "code":
				codeNode, err := parseCodeToken(token)
				if err != nil {
					return nil, err
				}
				tf.Nodes = append(tf.Nodes, codeNode)
			default:
				// TODO: error position??
				return nil, fmt.Errorf("unsupported opword %q inside of function", opWord)
			}

		}

		i += len(token)
	}

	return tf, nil
}
