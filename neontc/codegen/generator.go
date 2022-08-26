package codegen

import (
	"bytes"
	"fmt"
	"github.com/codemicro/go-neon/neontc/util"
	"math/rand"
	"strings"
)

type Generator struct {
	rand        *rand.Rand
	builder     *strings.Builder
	importNames map[string]string
}

func NewGenerator(seed int64) *Generator {
	return &Generator{
		rand:        rand.New(rand.NewSource(seed)),
		builder:     new(strings.Builder),
		importNames: make(map[string]string),
	}
}

func (g *Generator) AddPackageImportWithAlias(packageName string, importAlias string) {
	g.importNames[packageName] = importAlias
}

func (g *Generator) AddPackageImport(packageName string) {
	g.AddPackageImportWithAlias(packageName, "")
}

func (g *Generator) Render(packageName string) ([]byte, error) {
	writer := new(bytes.Buffer)

	if _, err := writer.Write([]byte(fmt.Sprintf("package %s\nimport (\n", packageName))); err != nil {
		return nil, err
	}

	for pkg, name := range g.importNames {
		if _, err := writer.Write([]byte(fmt.Sprintf("%s %q\n", name, pkg))); err != nil {
			return nil, err
		}
	}

	if _, err := writer.Write([]byte(")\n" + g.builder.String())); err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func (g *Generator) getImportName(pkg string) string {
	if v, found := g.importNames[pkg]; found {
		return v
	}
	id := "ntc" + util.GenerateRandomIdentifierWithRand(g.rand)
	g.importNames[pkg] = id
	return id
}
