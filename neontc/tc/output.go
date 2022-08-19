package tc

import (
	"bytes"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/codegen"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
)

func OutputGeneratorCode(
	packageName string,
	directory string,
	files []*ast.TemplateFile,
	nodeTypes map[*ast.SubstitutionNode]types.Type,
) error {

	for _, templateFile := range files {
		newFilename := filepath.Join(directory, filepath.Base(templateFile.Filepath)) + ".go"

		sb := new(bytes.Buffer)

		if err := codegen.GenerateHeader(sb, packageName); err != nil {
			return err
		}

		if err := codegen.GenerateImports(sb); err != nil {
			return err
		}

		for _, childNode := range templateFile.Nodes {
			switch node := childNode.(type) {
			case *ast.FuncDeclNode:
				if err := codegen.GenerateFunction(sb, node, nodeTypes); err != nil {
					return err
				}
			case *ast.RawCodeNode:
				if err := codegen.GenerateRawCode(sb, node); err != nil {
					return err
				}
			}
		}

		formatted, err := format.Source(sb.Bytes())
		if err != nil {
			return err
		}

		f, err := os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		if _, err := f.Write(formatted); err != nil {
			return err
		}
		_ = f.Close()
	}

	return nil
}
