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
	unwrappedToken := unwrapToken(token)

	firstSpace := bytes.IndexRune(unwrappedToken, ' ')
	firstNewline := bytes.IndexRune(unwrappedToken, '\n')

	if firstSpace == -1 && firstNewline == -1 {
		opWord = unwrappedToken
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

		opWord = unwrappedToken[:minimum]
		operand = unwrappedToken[minimum+1:]
	}

	return
}

func unwrapToken(token []byte) []byte {
	trimmedToken := bytes.TrimPrefix(
		bytes.TrimSuffix(token, []byte("}}")),
		[]byte("{{"),
	)
	trimmedToken = bytes.TrimSpace(trimmedToken)
	return trimmedToken
}

func File(fs *FileSet, fpath string, input []byte) (*ast.TemplateFile, error) {
	tf := new(ast.TemplateFile)
	tf.Filepath = fpath

	baseFileSetPosition := fs.AddFile(fpath, input)

	tokens, err := tokens(fs, baseFileSetPosition, input)
	if err != nil {
		return nil, err
	}

	for token := tokens.Next(); token != nil; token = tokens.Next() {

		if bytes.HasPrefix(token.cont, []byte("{{")) {
			// curly block

			opWordB, _ := chopToken(token.cont)
			opWord := string(opWordB)

			switch opWord {
			case "func":
				tokens.Rewind()
				funcDecl, err := parseFuncTokens(fs, tokens)
				if err != nil {
					return nil, err
				}
				tf.Nodes = append(tf.Nodes, funcDecl)
			case "code":
				tokens.Rewind()
				codeNode, err := parseCodeTokens(fs, tokens)
				if err != nil {
					return nil, err
				}
				tf.Nodes = append(tf.Nodes, codeNode)
			case "import":
				importNode, err := parseImportToken(fs, token)
				if err != nil {
					return nil, err
				}
				tf.Nodes = append(tf.Nodes, importNode)
			default:
				return nil, fmt.Errorf("%s: unsupported opword %q", fs.ResolvePosition(token.pos), opWord)
			}

		}

	}

	return tf, nil
}
