package codegen

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/ast"
	"github.com/codemicro/go-neon/neontc/util"
	"go/types"
	"math/rand"
)

func (g *Generator) GenerateRawCode(rawCodeNode *ast.RawCodeNode) error {
	_, err := g.builder.WriteString(rawCodeNode.GoCode + "\n")
	return err
}

func (g *Generator) GenerateTypecheckingFunction(funcDecl *ast.FuncDeclNode) (map[string]*ast.SubstitutionNode, error) {

	if _, err := g.builder.WriteString("func "); err != nil {
		return nil, err
	}

	if _, err := g.builder.WriteString(funcDecl.Signature); err != nil {
		return nil, err
	}

	if _, err := g.builder.WriteString(" {\n"); err != nil {
		return nil, err
	}

	ids := make(map[string]*ast.SubstitutionNode)
	for _, childNode := range funcDecl.ChildNodes {
		if err := g.generateTypecheckingNode(childNode, &ids); err != nil {
			return nil, err
		}
	}

	if _, err := g.builder.WriteString("\n}\n"); err != nil {
		return nil, err
	}

	return ids, nil
}

func (g *Generator) generateTypecheckingNode(node ast.Node, ids *map[string]*ast.SubstitutionNode) error {
	switch node := node.(type) {
	case *ast.RawCodeNode:
		if err := g.GenerateRawCode(node); err != nil {
			return err
		}
	case *ast.SubstitutionNode:
		name := "ntc_" + util.GenerateRandomIdentifier()
		(*ids)[name] = node
		code := "var " + name + " = " + node.Expression + "; _ = " + name + "\n"
		if _, err := g.builder.WriteString(code); err != nil {
			return err
		}
	case *ast.PlaintextNode:
		// intentional no-op
	case *ast.LoopNode:
		_, err := g.builder.WriteString("for " + node.LoopExpression + " {\n")
		if err != nil {
			return err
		}
		for _, childNode := range node.ChildNodes {
			if err := g.generateTypecheckingNode(childNode, ids); err != nil {
				return err
			}
		}
		_, err = g.builder.WriteString("}\n")
		if err != nil {
			return err
		}
	case *ast.ConditionalNode:
		currentNode := node

		for currentNode != nil {
			{
				var x string
				if currentNode == node {
					// This is the first if node
					x = "if " + currentNode.Expression + " {\n"
				} else {
					x = "else"
					if currentNode.Expression != "" {
						x += " if " + currentNode.Expression
					}
					x += " {\n"
				}

				_, err := g.builder.WriteString(x)
				if err != nil {
					return err
				}

			}

			for _, childNode := range currentNode.ChildNodes {
				if err := g.generateTypecheckingNode(childNode, ids); err != nil {
					return err
				}
			}

			_, err := g.builder.WriteString("} ")
			if err != nil {
				return err
			}

			currentNode = currentNode.Else
		}

		_, err := g.builder.WriteRune('\n')
		if err != nil {
			return err
		}

	default:
		panic(fmt.Errorf("unexpected node type %T", node))
	}
	return nil
}

func (g *Generator) GenerateFunction(funcDecl *ast.FuncDeclNode, nodeTypes map[*ast.SubstitutionNode]types.Type) error {
	var i int64
	for _, char := range funcDecl.Signature {
		i += int64(char)
	}
	writerID := "ntc" + util.GenerateRandomIdentifierWithRand(rand.New(rand.NewSource(64)))
	if _, err := g.builder.WriteString("func " + funcDecl.Signature + " string {\n"); err != nil {
		return err
	}

	{
		x := fmt.Sprintf("%s := new(%s.Builder)\n", writerID, g.getImportName("strings"))
		if _, err := g.builder.WriteString(x); err != nil {
			return err
		}
	}

	for _, childNode := range funcDecl.ChildNodes {
		if err := g.generateNode(childNode, writerID, nodeTypes); err != nil {
			return err
		}
	}

	if _, err := g.builder.WriteString("return " + writerID + ".String()\n}\n"); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateNode(node ast.Node, writerID string, nodeTypes map[*ast.SubstitutionNode]types.Type) error {
	switch node := node.(type) {
	case *ast.PlaintextNode:
		q := fmt.Sprintf("_, _ = %s.WriteString(%q)\n", writerID, node.Plaintext)
		if _, err := g.builder.WriteString(q); err != nil {
			return err
		}
	case *ast.RawCodeNode:
		if err := g.GenerateRawCode(node); err != nil {
			return err
		}
	case *ast.SubstitutionNode:
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

		starter := "_, _ = %s.WriteString("

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
			starter += g.getImportName("strconv") + ".FormatInt(int64(%s), 10)"

		case types.Uint:
			fallthrough
		case types.Uint8:
			fallthrough
		case types.Uint16:
			fallthrough
		case types.Uint32:
			fallthrough
		case types.Uint64:
			starter += g.getImportName("strconv") + ".FormatUint(uint64(%s), 10)"

		case types.Float32:
			starter += g.getImportName("strconv") + ".FormatFloat(float64(%s), 'f', -1, 32)"
		case types.Float64:
			starter += g.getImportName("strconv") + ".FormatFloat(float64(%s), 'f', -1, 64)"

		case types.Bool:
			starter += g.getImportName("strconv") + ".FormatBool(%s)"

		case types.String:
			starter += g.getImportName("html") + ".EscapeString(%s)"

		default:
			return fmt.Errorf("unsupported type %s", basicType.Name())
		}

		// Strings will need HTTP escaping applied to them.
		q := fmt.Sprintf(starter+")\n", writerID, node.Expression)
		if _, err := g.builder.WriteString(q); err != nil {
			return err
		}
	case *ast.LoopNode:
		_, err := g.builder.WriteString("for " + node.LoopExpression + " {\n")
		if err != nil {
			return err
		}
		for _, childNode := range node.ChildNodes {
			if err := g.generateNode(childNode, writerID, nodeTypes); err != nil {
				return err
			}
		}
		_, err = g.builder.WriteString("}\n")
		if err != nil {
			return err
		}
	case *ast.ConditionalNode:
		currentNode := node

		for currentNode != nil {
			{
				var x string
				if currentNode == node {
					// This is the first if node
					x = "if " + currentNode.Expression + " {\n"
				} else {
					x = "else"
					if currentNode.Expression != "" {
						x += " if " + currentNode.Expression
					}
					x += " {\n"
				}

				_, err := g.builder.WriteString(x)
				if err != nil {
					return err
				}

			}

			for _, childNode := range currentNode.ChildNodes {
				if err := g.generateNode(childNode, writerID, nodeTypes); err != nil {
					return err
				}
			}

			_, err := g.builder.WriteString("} ")
			if err != nil {
				return err
			}

			currentNode = currentNode.Else
		}

		_, err := g.builder.WriteRune('\n')
		if err != nil {
			return err
		}
	default:
		panic(fmt.Errorf("unexpected node type %T", node))
	}
	return nil
}
