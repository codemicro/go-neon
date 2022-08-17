package util

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var ErrNoGoMod = errors.New("no go.mod file found")

// FindGoModDir will find the directory containing the go.mod file for a
// project, starting with `startDir`.
//
// If the file cannot be found in the current directory, it will traverse up
// each directory until it gets to the root directory, looking in each
// subsequent dir.
func FindGoModDir(startDir string) (string, error) {
	if found, err := doesFileExist("go.mod"); err != nil {
		return "", err
	} else if found {
		return startDir, nil
	}

	absStartDir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	pathComponents := strings.Split(filepath.ToSlash(absStartDir), "/")
	fmt.Println(pathComponents)
	for i := len(pathComponents) - 1; i > 0; i -= 1 {
		dir := strings.Join(pathComponents[0:i], string(os.PathSeparator))

		// when probing the root of the path (`C:\` or `/`), we need the
		// trailing slash, but strings.Join won't add it since it's only a list
		// of items of length one.
		if i == 1 {
			dir += string(os.PathSeparator)
		}

		probe := filepath.Join(dir, "go.mod")
		if found, err := doesFileExist(probe); err != nil {
			return "", err
		} else if found {
			absPath, err := filepath.Abs(dir)
			if err != nil {
				return "", err
			}
			return absPath, nil
		}
	}

	return "", ErrNoGoMod
}

func doesFileExist(fname string) (bool, error) {
	_, err := os.Stat(fname)
	if err == nil {
		return true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	return false, nil
}

var ErrCannotParseModuleName = errors.New("cannot parse module name from go.mod file")

func ExtractModuleNameFromGoMod(goModPath string) (string, error) {
	fcont, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	if x := modfile.ModulePath(fcont); x == "" {
		return "", ErrCannotParseModuleName
	} else {
		return x, nil
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomIdentifier() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
