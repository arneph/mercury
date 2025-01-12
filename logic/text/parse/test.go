package parse

import (
	"fmt"
	positions "go/token"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseTest() *ast.Test {
	test, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.TEST {
		p.errs.Add(p.file().Position(test), fmt.Sprintf("expected 'test', got: %s", lit))
	}
	name := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if name == nil {
		return nil
	}
	lBrace, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.LBRACE {
		p.errs.Add(p.file().Position(lBrace), fmt.Sprintf("expected '{', got: %s", lit))
		return nil
	}
	var entries []ast.TestEntry
	var rBrace positions.Pos
parseLoop:
	for {
		pos, tok, lit := p.scanner.Peek(scan.SKIP_NEW_LINES)
		switch tok {
		case tokens.RBRACE:
			p.scanner.Scan(scan.EMIT_NEW_LINES)
			rBrace = pos
			break parseLoop
		case tokens.COMPONENT:
			componentDecl := p.parseComponentDecl()
			if componentDecl == nil {
				break
			}
			entries = append(entries, componentDecl)
			if ok := p.parseNewLine(); ok {
				continue parseLoop
			}
		case tokens.SET:
			setInstr := p.parseSetInstr()
			if setInstr == nil {
				break
			}
			entries = append(entries, setInstr)
			if ok := p.parseNewLine(); ok {
				continue parseLoop
			}
		case tokens.ASSERT:
			assertion := p.parseAssertion()
			if assertion == nil {
				break
			}
			entries = append(entries, assertion)
			if ok := p.parseNewLine(); ok {
				continue parseLoop
			}
		case tokens.EXPECT:
			expectation := p.parseExpectation()
			if expectation == nil {
				break
			}
			entries = append(entries, expectation)
			if ok := p.parseNewLine(); ok {
				continue parseLoop
			}
		default:
			p.scanner.Scan(scan.EMIT_NEW_LINES)
			p.errs.Add(p.file().Position(pos), fmt.Sprintf("unexpected token: %s", lit))
		}
		if ok := p.recoverToNewLine(); !ok {
			return nil
		}
	}
	return &ast.Test{
		Test:    test,
		Name:    name,
		LBrace:  lBrace,
		Entries: entries,
		RBrace:  rBrace,
	}
}

func (p *parser) parseComponentDecl() *ast.ComponentDecl {
	component, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.COMPONENT {
		p.errs.Add(p.file().Position(component), fmt.Sprintf("expected 'component', got: %s", lit))
		return nil
	}
	colon, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.COLON {
		p.errs.Add(p.file().Position(colon), fmt.Sprintf("expected ':', got: %s", lit))
		return nil
	}
	componentName := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if componentName == nil {
		return nil
	}
	return &ast.ComponentDecl{
		Component:     component,
		Colon:         colon,
		ComponentName: componentName,
	}
}

func (p *parser) parseSetInstr() *ast.SetInstr {
	set, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.SET {
		p.errs.Add(p.file().Position(set), fmt.Sprintf("expected 'set', got: %s", lit))
		return nil
	}
	inputs := p.parseBusReferenceList()
	if inputs == nil {
		return nil
	}
	colon, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.COLON {
		p.errs.Add(p.file().Position(colon), fmt.Sprintf("expected ':', got: %s", lit))
		return nil
	}
	constants := p.parseConstants()
	if constants == nil {
		return nil
	}
	return &ast.SetInstr{
		Set:       set,
		Inputs:    inputs,
		Colon:     colon,
		Constants: constants,
	}
}

func (p *parser) parseAssertion() *ast.Assertion {
	assert, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.ASSERT {
		p.errs.Add(p.file().Position(assert), fmt.Sprintf("expected 'assert', got: %s", lit))
		return nil
	}
	outputs := p.parseBusReferenceList()
	if outputs == nil {
		return nil
	}
	is, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.IS {
		p.errs.Add(p.file().Position(is), fmt.Sprintf("expected 'is', got: %s", lit))
		return nil
	}
	constants := p.parseConstants()
	if constants == nil {
		return nil
	}
	return &ast.Assertion{
		Assert:    assert,
		Outputs:   outputs,
		Is:        is,
		Constants: constants,
	}
}

func (p *parser) parseExpectation() *ast.Expectation {
	expect, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.EXPECT {
		p.errs.Add(p.file().Position(expect), fmt.Sprintf("expected 'expect', got: %s", lit))
		return nil
	}
	outputs := p.parseBusReferenceList()
	if outputs == nil {
		return nil
	}
	is, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.IS {
		p.errs.Add(p.file().Position(is), fmt.Sprintf("expected 'is', got: %s", lit))
		return nil
	}
	constants := p.parseConstants()
	if constants == nil {
		return nil
	}
	return &ast.Expectation{
		Expect:    expect,
		Outputs:   outputs,
		Is:        is,
		Constants: constants,
	}
}
