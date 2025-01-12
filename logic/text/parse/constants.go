package parse

import (
	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseConstants() *ast.Constants {
	value := p.parseNumber()
	if value == nil {
		return nil
	}
	values := []*ast.Number{value}
	for {
		_, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES)
		if tok != tokens.COMMA {
			break
		}
		p.scanner.Scan(scan.EMIT_NEW_LINES)
		value := p.parseNumber()
		if value == nil {
			return nil
		}
		values = append(values, value)
	}
	return &ast.Constants{
		Values: values,
	}
}
