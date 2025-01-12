package parse

import (
	"fmt"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseLiteral() ast.Expr {
	pos, tok, lit := p.scanner.Peek(scan.EMIT_NEW_LINES)
	switch tok {
	case tokens.IDENTIFIER:
		return p.parseIdentifier(scan.EMIT_NEW_LINES)
	case tokens.NUMBER:
		return p.parseNumber()
	default:
		p.errs.Add(p.file().Position(pos), fmt.Sprintf("expected identifier or number, got: %s", lit))
		return nil
	}
}

func (p *parser) parseIdentifier(newLineMode scan.NewLineMode) *ast.Identifier {
	pos, tok, lit := p.scanner.Scan(newLineMode)
	if tok != tokens.IDENTIFIER {
		p.errs.Add(p.file().Position(pos), fmt.Sprintf("expected identifier, got: %s", lit))
		return nil
	}
	return &ast.Identifier{
		Name:  lit,
		Start: pos,
	}
}

func (p *parser) parseNumber() *ast.Number {
	pos, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.NUMBER {
		p.errs.Add(p.file().Position(pos), fmt.Sprintf("expected number, got: %s", lit))
		return nil
	}
	return &ast.Number{
		Value: lit,
		Start: pos,
	}
}
