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
