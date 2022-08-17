package codegen

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/util"
	"io"
)

func GenerateHeader(writer io.Writer, packageName string) error {
	x := fmt.Sprintf("package %s\n", packageName)
	_, err := writer.Write([]byte(x))
	return err
}

func GenerateImports(writer io.Writer) error {
	_, err := writer.Write([]byte("import \"bytes\"\n"))
	return err
}

func GenerateRawCode(writer io.Writer, rawCodeNode *ast.RawCodeNode) error {
	_, err := writer.Write([]byte(rawCodeNode.GoCode + "\n"))
	return err
}

func GenerateTypecheckingFunction(writer io.Writer, funcDecl *ast.FuncDeclNode) (map[string]*ast.SubstitutionNode, error) {

	if _, err := writer.Write([]byte("func ")); err != nil {
		return nil, err
	}

	if _, err := writer.Write([]byte(funcDecl.Signature)); err != nil {
		return nil, err
	}

	if _, err := writer.Write([]byte(" {\n")); err != nil {
		return nil, err
	}

	ids := make(map[string]*ast.SubstitutionNode)
	for _, childNode := range funcDecl.ChildNodes {
		switch node := childNode.(type) {
		case *ast.RawCodeNode:
			if err := GenerateRawCode(writer, node); err != nil {
				return nil, err
			}
		case *ast.SubstitutionNode:
			name := "ntc_" + util.GenerateRandomIdentifier()
			ids[name] = node
			code := "var " + name + " = " + node.Expression + "; _ = " + name + "\n"
			if _, err := writer.Write([]byte(code)); err != nil {
				return nil, err
			}
		}
	}

	if _, err := writer.Write([]byte("\n}\n")); err != nil {
		return nil, err
	}

	return ids, nil
}

func GenerateFunction(writer io.Writer, funcDecl *ast.FuncDeclNode) error {
	// TODO: derive this identifier from the function somehow
	writerID := "ntc" + util.GenerateRandomIdentifier()
	if _, err := writer.Write([]byte("func " + funcDecl.Signature + " string {\n")); err != nil {
		return err
	}

	if _, err := writer.Write([]byte(writerID + " := new(bytes.Buffer)\n")); err != nil {
		return err
	}

	for _, childNode := range funcDecl.ChildNodes {

		switch node := childNode.(type) {
		case *ast.PlaintextNode:
			q := fmt.Sprintf("_, _ = %s.Write([]byte(%q))\n", writerID, node.Plaintext)
			if _, err := writer.Write([]byte(q)); err != nil {
				return err
			}
		case *ast.RawCodeNode:
			if err := GenerateRawCode(writer, node); err != nil {
				return err
			}
		case *ast.SubstitutionNode:
			// TODO: actual types
			q := fmt.Sprintf("_, _ = %s.Write([]byte(%s))\n", writerID, node.Expression)
			if _, err := writer.Write([]byte(q)); err != nil {
				return err
			}
		}
	}

	if _, err := writer.Write([]byte("return " + writerID + ".String()\n}\n")); err != nil {
		return err
	}

	return nil
}
