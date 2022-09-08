package tc

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/codegen"
	"github.com/codemicro/go-neon/neontc/parse"
	"github.com/codemicro/go-neon/neontc/util"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
)

func OutputGeneratorCode(
	fs *parse.FileSet,
	packageName string,
	directory string,
	files []*ast.TemplateFile,
	nodeTypes map[*ast.SubstitutionNode]types.Type,
) error {

	if err := os.MkdirAll(directory, os.ModeDir); err != nil {
		return err
	}

	otherFuncDecls := make(map[string]struct{})
	for _, templateFile := range files {
		for _, childNode := range templateFile.Nodes {
			if childNode, ok := childNode.(*ast.FuncDeclNode); ok {
				otherFuncDecls[childNode.Identifier] = struct{}{}
			}
		}
		
	}
	
	for _, templateFile := range files {
		newFilename := filepath.Join(directory, filepath.Base(templateFile.Filepath)) + ".go"

		generator := codegen.NewGenerator(
			util.Int64FromString(filepath.Base(templateFile.Filepath)),
		)

		for _, childNode := range templateFile.Nodes {
			switch node := childNode.(type) {
			case *ast.FuncDeclNode:
				fmt.Println(node.Identifier)
				if err := generator.GenerateFunction(fs, node, nodeTypes, otherFuncDecls); err != nil {
					return err
				}
			case *ast.RawCodeNode:
				if err := generator.GenerateRawCode(node); err != nil {
					return err
				}
			case *ast.ImportNode:
				generator.AddPackageImportWithAlias(node.ImportPath, node.Alias)
			default:
				panic(fmt.Errorf("unexpected node type %T", node))
			}
		}

		renderedBytes, err := generator.Render(packageName)
		if err != nil {
			return err
		}

		formatted, err := format.Source(renderedBytes)
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
