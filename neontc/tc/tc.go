package tc

import (
	"errors"
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

	if len(files) == 0 {
		return errors.New("no matching input files")
	}

	// generate typechecking package
	subsitutionTypes, err := DetermineSubstitutionTypes(modulePath+"/"+filepath.ToSlash(requiredPathTranslation), directory, files)
	if err != nil {
		return err
	}

	// generate output code
	// TODO: make this not based on $GOPACKAGE
	gopkg := os.Getenv("GOPACKAGE")
	if gopkg == "" {
		return errors.New("$GOPACKAGE empty (did you run this with go generate?)")
	}
	if err := OutputGeneratorCode(gopkg, directory, files, subsitutionTypes); err != nil {
		return err
	}

	return nil
}
