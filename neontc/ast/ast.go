package ast

type Node any

type Pos int64

type TemplateFile struct {
	Filepath string
	Nodes    []Node
}

type FuncDeclNode struct {
	Pos
	Signature  string
	Identifier string
	ChildNodes []Node
}

type SubstitutionModifier uint

// If SubModNone and the other submods are in one const block, `iota` will pick up on the zero and start at 1

const SubModNone SubstitutionModifier = 0

const (
	SubModUnsafe SubstitutionModifier = 1 << iota
)

func (s SubstitutionModifier) Includes(i SubstitutionModifier) bool {
	return s&i != 0
}

var SubModMapping = map[string]SubstitutionModifier{
	"unsafe": SubModUnsafe,
}

type SubstitutionNode struct {
	Pos
	Expression string
	Modifier   SubstitutionModifier
}

type PlaintextNode struct {
	Pos
	Plaintext string
}

type RawCodeNode struct {
	Pos
	GoCode string
}

type ImportNode struct {
	Pos
	Alias      string
	ImportPath string
}

type LoopNode struct {
	Pos
	LoopExpression string
	ChildNodes     []Node
}

type ConditionalNode struct {
	Pos
	Expression string
	ChildNodes []Node
	Else       *ConditionalNode
}
