package parse

import (
	positions "go/token"
	"testing"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
)

func TestParseLiteralFails(t *testing.T) {
	src := []byte("\t\tcomponent")
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
	p := parser{scanner: scan.NewScanner(file, src)}
	expr := p.parseLiteral()
	if p.errs.Len() != 1 {
		t.Errorf("Expected parse errors; got %v", p.errs)
	}
	if expr != nil {
		t.Errorf("Expected parse failure; got %v", expr)
	}
}

func TestParsesIdentifierLiteral(t *testing.T) {
	src := []byte("     hello")
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
	p := parser{scanner: scan.NewScanner(file, src)}
	expr := p.parseLiteral()
	if p.errs.Len() != 0 {
		t.Errorf("Expected no parse errors; got %v", p.errs)
	}
	if expr == nil {
		t.Fatalf("Expected parsed ast.Identifier; got nil")
	}
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected parsed ast.Expr to be an ast.Identifier; got %t", expr)
	}
	if ident.Name != "hello" {
		t.Errorf("Expected parsed ast.Identifier.Name to be 'hello'; got %v", ident.Name)
	}
	if ident.Pos() != file.Pos(5) {
		t.Errorf("Expected parsed ast.Identifier at Pos(5); got %v", ident.Pos())
	}
}

func TestParsesNumberLiteral(t *testing.T) {
	src := []byte("     0xcafe")
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
	p := parser{scanner: scan.NewScanner(file, src)}
	expr := p.parseLiteral()
	if p.errs.Len() != 0 {
		t.Errorf("Expected no parse errors; got %v", p.errs)
	}
	if expr == nil {
		t.Fatalf("Expected parsed ast.Identifier; got nil")
	}
	num, ok := expr.(*ast.Number)
	if !ok {
		t.Fatalf("Expected parsed ast.Expr to be an ast.Number; got %t", expr)
	}
	if num.Value != "0xcafe" {
		t.Errorf("Expected parsed ast.Number.Value to be '42'; got %v", num.Value)
	}
	if num.Pos() != file.Pos(5) {
		t.Errorf("Expected parsed ast.Number at Pos(5); got %v", num.Pos())
	}
}

func TestParseIdentifierFails(t *testing.T) {
	src := []byte("   42")
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
	p := parser{scanner: scan.NewScanner(file, src)}
	ident := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if p.errs.Len() != 1 {
		t.Errorf("Expected parse errors; got %v", p.errs)
	}
	if ident != nil {
		t.Errorf("Expected parse failure; got %v", ident)
	}
}

func TestParsesIdentifier(t *testing.T) {
	for _, nlm := range []scan.NewLineMode{scan.EMIT_NEW_LINES, scan.SKIP_NEW_LINES} {
		src := []byte("   hello")
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
		p := parser{scanner: scan.NewScanner(file, src)}
		ident := p.parseIdentifier(nlm)
		if p.errs.Len() != 0 {
			t.Errorf("Expected no parse errors; got %v", p.errs)
		}
		if ident == nil {
			t.Fatalf("Expected parsed ast.Identifier; got nil")
		}
		if ident.Name != "hello" {
			t.Errorf("Expected parsed ast.Identifier.Name to be 'hello'; got %v", ident.Name)
		}
		if ident.Pos() != file.Pos(3) {
			t.Errorf("Expected parsed ast.Identifier at Pos(3); got %v", ident.Pos())
		}
	}
}

func TestParseNumberFails(t *testing.T) {
	src := []byte("   hello")
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
	p := parser{scanner: scan.NewScanner(file, src)}
	num := p.parseNumber()
	if p.errs.Len() != 1 {
		t.Errorf("Expected parse errors; got %v", p.errs)
	}
	if num != nil {
		t.Errorf("Expected parse failure; got %v", num)
	}
}

func TestParsesNumber(t *testing.T) {
	src := []byte("   42")
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(src))
	p := parser{scanner: scan.NewScanner(file, src)}
	num := p.parseNumber()
	if p.errs.Len() != 0 {
		t.Errorf("Expected no parse errors; got %v", p.errs)
	}
	if num == nil {
		t.Fatalf("Expected parsed ast.Number; got nil")
	}
	if num.Value != "42" {
		t.Errorf("Expected parsed ast.Number.Value to be '42'; got %v", num.Value)
	}
	if num.Pos() != file.Pos(3) {
		t.Errorf("Expected parsed ast.Number at Pos(3); got %v", num.Pos())
	}
}
