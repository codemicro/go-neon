package tc

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/config"
	"github.com/codemicro/go-neon/neontc/parse"
	"github.com/codemicro/go-neon/neontc/util"
)

func RunOnDirectory(conf *config.Config, directory string) error {

	if conf.Package == "" {
		return errors.New("unspecified package - try running with go generate or using the --package flag")
	}

	var err error
	directory, err = filepath.Abs(directory)
	if err != nil {
		return err
	}

	// Find and parse module name out of go.mod

	goModDir, err := util.FindGoModDir(directory)
	if err != nil {
		return err
	}

	modulePath, err := util.ExtractModuleNameFromGoMod(filepath.Join(goModDir, "go.mod"))
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(os.Stderr, "Found module %s\n", modulePath)

	// requiredPathTranslation is the directories that need to be added to go
	// from the module path to the target module
	requiredPathTranslation, err := filepath.Rel(goModDir, directory)
	if err != nil {
		return err
	}
	requiredPathTranslation = filepath.ToSlash(requiredPathTranslation)

	// list input files
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	fset := parse.NewFileSet()

	// parse those files
	var files []*ast.TemplateFile
	for _, de := range dirEntries {
		if de.IsDir() || !strings.HasSuffix(de.Name(), "."+conf.FileExtension) {
			continue
		}

		fullPath := filepath.Join(directory, de.Name())

		cont, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(os.Stderr, "Parsing %s\n", fullPath)

		tf, err := parse.File(fset, fullPath, cont)
		if err != nil {
			return err
		}

		files = append(files, tf)
	}

	if len(files) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "No input files matching criteria")
		os.Exit(2)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Typechecking %s\n", modulePath+"/"+requiredPathTranslation)

	// generate typechecking package
	subsitutionTypes, err := DetermineSubstitutionTypes(modulePath+"/"+requiredPathTranslation, conf.Package, directory, files, !conf.KeepTempFiles)
	if err != nil {
		return err
	}

	// generate output code

	_, _ = fmt.Fprintf(os.Stderr, "Generating output in %s\n", conf.OutputDirectory)

	if err := OutputGeneratorCode(fset, conf.Package, conf.OutputDirectory, files, subsitutionTypes); err != nil {
		return err
	}

	return nil
}
