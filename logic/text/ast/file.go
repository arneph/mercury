package ast

import (
	positions "go/token"
)

type File struct {
	FileStart, FileEnd positions.Pos
	Nodes              []FileNode
}

func (f *File) Pos() positions.Pos {
	return f.FileStart
}

func (f *File) End() positions.Pos {
	return f.FileEnd
}

type FileNode interface {
	Node
	fileNode()
}
