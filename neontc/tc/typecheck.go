package tc

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/codegen"
	"github.com/codemicro/go-neon/neontc/util"
	"go/parser"
	"go/types"
	"golang.org/x/tools/go/loader"
	"os"
	"path/filepath"
)

func DetermineSubstitutionTypes(
	modulePath string,
	baseDirectory string,
	files []*ast.TemplateFile,
) (map[*ast.SubstitutionNode]types.Type, error) {
	tempPackageName := util.GenerateRandomIdentifier()
	tempModulePath := modulePath + "/ntc-tc-" + tempPackageName

	ids := make(map[string]*ast.SubstitutionNode)

	newDirectory := filepath.Join(baseDirectory, "ntc-tc-"+tempPackageName)
	if err := os.MkdirAll(newDirectory, os.ModeDir); err != nil {
		return nil, err
	}

	for _, templateFile := range files {
		newFilename := filepath.Join(newDirectory, filepath.Base(templateFile.Filepath)) + ".go"

		generator := codegen.NewGenerator(4)

		for _, childNode := range templateFile.Nodes {
			switch node := childNode.(type) {
			case *ast.FuncDeclNode:
				additionalIDs, err := generator.GenerateTypecheckingFunction(node)
				if err != nil {
					return nil, err
				}

				for k, v := range additionalIDs {
					ids[k] = v
				}
			case *ast.RawCodeNode:
				if err := generator.GenerateRawCode(node); err != nil {
					return nil, err
				}
			case *ast.ImportNode:
				generator.AddPackageImportWithAlias(node.ImportPath, node.Alias)
			default:
				panic(fmt.Errorf("unexpected node type %T", node))
			}
		}

		f, err := os.OpenFile(newFilename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		renderedBytes, err := generator.Render(tempPackageName)
		if err != nil {
			return nil, err
		}

		if _, err := f.Write(renderedBytes); err != nil {
			return nil, err
		}

		_ = f.Close()
	}

	// run the things!!!

	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(tempModulePath)
	lprog, err := conf.Load()
	if err != nil {
		return nil, err // load error
	}

	pkg := lprog.Package(tempModulePath)

	expressionTypes := make(map[*ast.SubstitutionNode]types.Type)

	for id, obj := range pkg.Info.Defs {
		sdn, found := ids[id.Name]
		if found {
			expressionTypes[sdn] = obj.Type()
		}
	}

	_ = os.RemoveAll(newDirectory)

	return expressionTypes, nil
}
