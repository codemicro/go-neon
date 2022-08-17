package ast

type Node any

type TemplateFile struct {
	Filepath string
	Nodes    []Node
}

type FuncDeclNode struct {
	Signature  string
	ChildNodes []Node
}

type SubstitutionNode struct {
	Expression string
}

type PlaintextNode struct {
	Plaintext string
}

type RawCodeNode struct {
	GoCode string
}
