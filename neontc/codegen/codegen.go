package codegen

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/util"
	"go/types"
	"io"
	"math/rand"
)

func GenerateHeader(writer io.Writer, packageName string) error {
	x := fmt.Sprintf("package %s\n", packageName)
	_, err := writer.Write([]byte(x))
	return err
}

func GenerateImports(writer io.Writer) error {
	_, err := writer.Write([]byte(`import (
	"bytes"
	"strconv"
)
`))
	return err
}

func GenerateRawCode(writer io.Writer, rawCodeNode *ast.RawCodeNode) error {
	_, err := writer.Write([]byte(rawCodeNode.GoCode + "\n"))
	return err
}

func GenerateTypecheckingFunction(writer io.Writer, funcDecl *ast.FuncDeclNode) (map[string]*ast.SubstitutionNode, error) {

	if _, err := writer.Write([]byte("func ")); err != nil {
		return nil, err
	}

	if _, err := writer.Write([]byte(funcDecl.Signature)); err != nil {
		return nil, err
	}

	if _, err := writer.Write([]byte(" {\n")); err != nil {
		return nil, err
	}

	ids := make(map[string]*ast.SubstitutionNode)
	for _, childNode := range funcDecl.ChildNodes {
		switch node := childNode.(type) {
		case *ast.RawCodeNode:
			if err := GenerateRawCode(writer, node); err != nil {
				return nil, err
			}
		case *ast.SubstitutionNode:
			name := "ntc_" + util.GenerateRandomIdentifier()
			ids[name] = node
			code := "var " + name + " = " + node.Expression + "; _ = " + name + "\n"
			if _, err := writer.Write([]byte(code)); err != nil {
				return nil, err
			}
		}
	}

	if _, err := writer.Write([]byte("\n}\n")); err != nil {
		return nil, err
	}

	return ids, nil
}

func GenerateFunction(writer io.Writer, funcDecl *ast.FuncDeclNode, nodeTypes map[*ast.SubstitutionNode]types.Type) error {
	var i int64
	for _, char := range funcDecl.Signature {
		i += int64(char)
	}
	writerID := "ntc" + util.GenerateRandomIdentifierWithRand(rand.New(rand.NewSource(64)))
	if _, err := writer.Write([]byte("func " + funcDecl.Signature + " string {\n")); err != nil {
		return err
	}

	if _, err := writer.Write([]byte(writerID + " := new(bytes.Buffer)\n")); err != nil {
		return err
	}

	for _, childNode := range funcDecl.ChildNodes {
		switch node := childNode.(type) {
		case *ast.PlaintextNode:
			q := fmt.Sprintf("_, _ = %s.Write([]byte(%q))\n", writerID, node.Plaintext)
			if _, err := writer.Write([]byte(q)); err != nil {
				return err
			}
		case *ast.RawCodeNode:
			if err := GenerateRawCode(writer, node); err != nil {
				return err
			}
		case *ast.SubstitutionNode:
			// TODO: actual types
			nodeType, found := nodeTypes[node]
			if !found {
				panic("impossible state: substitution expression was not type checked")
			}

			underlyingType := nodeType.Underlying()
			var basicType *types.Basic

			for _, tp := range types.Typ {
				if types.Identical(underlyingType, tp) {
					basicType = tp
					break
				}
			}

			if basicType == nil {
				return fmt.Errorf("unsupported type %s", underlyingType.String())
			}

			starter := "_, _ = %s.Write([]byte("

			switch basicType.Kind() {
			case types.Int:
				fallthrough
			case types.Int8:
				fallthrough
			case types.Int16:
				fallthrough
			case types.Int32:
				fallthrough
			case types.Int64:
				starter += "strconv.FormatInt(int64(%s), 2)"

			case types.Uint:
				fallthrough
			case types.Uint8:
				fallthrough
			case types.Uint16:
				fallthrough
			case types.Uint32:
				fallthrough
			case types.Uint64:
				starter += "strconv.FormatUint(uint64(%s), 2)"

			case types.Float32:
				starter += "strconv.FormatFloat(float64(%s), 'f', -1, 32)"
			case types.Float64:
				starter += "strconv.FormatFloat(float64(%s), 'f', -1, 64)"

			case types.Bool:
				starter += "strconv.FormatBool(%s)"

			case types.String:
				starter += "%s"

			default:
				return fmt.Errorf("unsupported type %s", basicType.Name())
			}

			fmt.Println(node.Expression, nodeType.String())

			// We need to convert the type into a string of some sort, and then into bytes.
			// Strings will need HTTP escaping applied to them.
			q := fmt.Sprintf(starter+"))\n", writerID, node.Expression)
			if _, err := writer.Write([]byte(q)); err != nil {
				return err
			}
		}
	}

	if _, err := writer.Write([]byte("return " + writerID + ".String()\n}\n")); err != nil {
		return err
	}

	return nil
}
