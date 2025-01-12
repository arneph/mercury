package parse

import (
	errors "go/scanner"
	positions "go/token"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
)

func ParseFile(posFile *positions.File, src []byte) (*ast.File, errors.ErrorList) {
	p := &parser{
		scanner: scan.NewScanner(posFile, src),
	}
	astFile := p.parseFile()
	return astFile, p.errs
}

type parser struct {
	scanner *scan.Scanner
	errs    errors.ErrorList
}

func (p *parser) file() *positions.File {
	return p.scanner.File()
}
