package tc

import (
	"fmt"
	"go/parser"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/codegen"
	"github.com/codemicro/go-neon/neontc/util"
	"golang.org/x/tools/go/loader"
)

func DetermineSubstitutionTypes(
	modulePath string,
	packageName string,
	baseDirectory string,
	files []*ast.TemplateFile,
	deleteTypecheckingFiles bool,
) (map[*ast.SubstitutionNode]types.Type, error) {
	tempPackageName := util.GenerateRandomIdentifier()
	tempModulePath := modulePath + "/ntc-tc-" + tempPackageName

	ids := make(map[string]*ast.SubstitutionNode)

	newDirectory := filepath.Join(baseDirectory, "ntc-tc-"+tempPackageName)
	if err := os.MkdirAll(newDirectory, os.ModeDir); err != nil {
		return nil, err
	}

	// copy .go files into the temporary directory for typechecking purposes
	dirEntries, err := os.ReadDir(baseDirectory)
	if err != nil {
		return nil, err
	}
	for _, de := range dirEntries {
		if de.IsDir() || !strings.EqualFold(filepath.Ext(de.Name()), ".go") {
			continue
		}

		f1, err := os.Open(filepath.Join(baseDirectory, de.Name()))
		if err != nil {
			return nil, err
		}

		f2, err := os.OpenFile(filepath.Join(newDirectory, de.Name()), os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			_ = f1.Close()
			return nil, err
		}

		_, err = io.Copy(f2, f1)

		_ = f1.Close()
		_ = f2.Close()

		if err != nil {
			return nil, err
		}
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

		renderedBytes, err := generator.Render(packageName)
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

	if deleteTypecheckingFiles {
		if err := os.RemoveAll(newDirectory); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: could not remove temporary directory %s\n", newDirectory)
		}
	}

	return expressionTypes, nil
}
