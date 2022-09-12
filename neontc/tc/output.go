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
	"regexp"
)

var funcCallIdentifierRegexp = regexp.MustCompile(`^(\w+)\(.*\)$`)

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

	markTrustedFunctions(files)

	for _, templateFile := range files {
		newFilename := filepath.Join(directory, filepath.Base(templateFile.Filepath)) + ".go"

		generator := codegen.NewGenerator(
			util.Int64FromString(filepath.Base(templateFile.Filepath)),
		)

		for _, childNode := range templateFile.Nodes {
			switch node := childNode.(type) {
			case *ast.FuncDeclNode:
				if err := generator.GenerateFunction(fs, node, nodeTypes); err != nil {
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

// markTrustedFunctions marks other generated functions as trusted/enables unsafe mode on them
func markTrustedFunctions(files []*ast.TemplateFile) {
	var nodesToVisit []ast.Node

	trustedIdentifiers := make(map[string]struct{})
	for _, templateFile := range files {
		for _, childNode := range templateFile.Nodes {
			nodesToVisit = append(nodesToVisit, childNode)
			if childNode, ok := childNode.(*ast.FuncDeclNode); ok {
				trustedIdentifiers[childNode.Identifier] = struct{}{}
			}
		}
	}

	for i := 0; i < len(nodesToVisit); i += 1 {
		node := nodesToVisit[i]
		switch node := node.(type) {
		case *ast.SubstitutionNode:
			submatches := funcCallIdentifierRegexp.FindStringSubmatch(node.Expression)
			if len(submatches) != 0 {
				if _, found := trustedIdentifiers[submatches[1]]; found {
					(*node).Modifier = node.Modifier | ast.SubModUnsafe
				}
			}
		case *ast.FuncDeclNode:
			nodesToVisit = append(nodesToVisit, node.ChildNodes...)
		case *ast.LoopNode:
			nodesToVisit = append(nodesToVisit, node.ChildNodes...)
		case *ast.ConditionalNode:
			nodesToVisit = append(nodesToVisit, node.ChildNodes...)
			if node.Else != nil {
				nodesToVisit = append(nodesToVisit, node.Else)
			}
		}
	}
}
