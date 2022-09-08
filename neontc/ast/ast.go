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

type SubstitutionNode struct {
	Pos
	Expression string
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
