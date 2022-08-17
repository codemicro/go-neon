package tc

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/parse"
	"github.com/codemicro/go-neon/neontc/util"
	"os"
	"path/filepath"
	"strings"
)

func RunOnDirectory(directory string) error {

	// Find and parse module name out of go.mod

	goModDir, err := util.FindGoModDir(directory)
	if err != nil {
		return err
	}

	modulePath, err := util.ExtractModuleNameFromGoMod(filepath.Join(goModDir, "go.mod"))
	if err != nil {
		return err
	}

	// requiredPathTranslation is the directories that need to be added to go
	// from the module path to the target module
	requiredPathTranslation, err := filepath.Rel(goModDir, directory)
	if err != nil {
		return err
	}

	// list `ntc` files
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	// parse those files
	var files []*ast.TemplateFile
	for _, de := range dirEntries {
		fmt.Println(strings.EqualFold(filepath.Ext(de.Name()), ".ntc"), de.Name()+" -> "+filepath.Ext(de.Name()))
		if de.IsDir() || !strings.EqualFold(filepath.Ext(de.Name()), ".ntc") {
			continue
		}

		fullPath := filepath.Join(directory, de.Name())

		cont, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		tf, err := parse.File(fullPath, cont)
		if err != nil {
			return err
		}

		files = append(files, tf)
	}

	// generate typechecking package
	_, err = DetermineSubstitutionTypes(modulePath+"/"+filepath.ToSlash(requiredPathTranslation), directory, files)
	if err != nil {
		return err
	}

	// generate output code
	if err := OutputGeneratorCode("FIXME", directory, files); err != nil {
		return err
	}

	return nil
}