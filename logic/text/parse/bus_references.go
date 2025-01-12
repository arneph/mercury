package parse

import (
	"fmt"
	positions "go/token"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseBusReferenceList() *ast.BusReferenceList {
	busRef := p.parseBusReference()
	if busRef == nil {
		return nil
	}
	busRefs := []*ast.BusReference{busRef}
	for {
		_, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES)
		if tok != tokens.COMMA {
			break
		}
		p.scanner.Scan(scan.EMIT_NEW_LINES)
		busRef := p.parseBusReference()
		if busRef == nil {
			return nil
		}
		busRefs = append(busRefs, busRef)
	}
	return &ast.BusReferenceList{
		References: busRefs,
	}
}

func (p *parser) parseBusReference() *ast.BusReference {
	name := p.parseIdentifier(scan.SKIP_NEW_LINES)
	if name == nil {
		return nil
	}
	lBrack, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES)
	if tok != tokens.LBRACK {
		return &ast.BusReference{
			Name:      name,
			LBrack:    positions.NoPos,
			WireIndex: nil,
			RBrack:    positions.NoPos,
		}
	}
	p.scanner.Scan(scan.EMIT_NEW_LINES)
	wireIndex := p.parseExpr()
	if wireIndex == nil {
		return nil
	}
	rBrack, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.RBRACK {
		p.errs.Add(p.file().Position(rBrack), fmt.Sprintf("expected ']', got: %s", lit))
		return nil
	}
	return &ast.BusReference{
		Name:      name,
		LBrack:    lBrack,
		WireIndex: wireIndex,
		RBrack:    rBrack,
	}
}
