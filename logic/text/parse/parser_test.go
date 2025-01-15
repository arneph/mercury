package parse

import (
	positions "go/token"
	"testing"
)

func FuzzParser(f *testing.F) {
	f.Add([]byte("component Nand(a, b) (r) {\nr: rand(a, b)\n}"))
	f.Add([]byte("test NandTest {\ncomponent: Nand\nset a, b: 0, 1\nexpect r is 42\n}"))
	f.Fuzz(func(t *testing.T, in []byte) {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(in))
		ast, errs := ParseFile(file, in)
		if ast == nil && errs.Len() == 0 {
			t.Errorf("Expected ParseFile to return an AST or errors, got neither.")
		}
	})
}
