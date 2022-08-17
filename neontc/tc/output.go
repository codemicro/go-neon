package tc

import (
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/codegen"
	"os"
	"path/filepath"
)

func OutputGeneratorCode(
	packageName string,
	directory string,
	files []*ast.TemplateFile,
) error {

	for _, templateFile := range files {
		newFilename := filepath.Join(directory, filepath.Base(templateFile.Filepath)) + ".go"

		f, err := os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		if err := codegen.GenerateHeader(f, packageName); err != nil {
			return err
		}

		if err := codegen.GenerateImports(f); err != nil {
			return err
		}

		for _, childNode := range templateFile.Nodes {
			switch node := childNode.(type) {
			case *ast.FuncDeclNode:
				if err := codegen.GenerateFunction(f, node); err != nil {
					return err
				}
			case *ast.RawCodeNode:
				if err := codegen.GenerateRawCode(f, node); err != nil {
					return err
				}
			}
		}

		_ = f.Close()
	}

	return nil
}
