package parse

import (
	"fmt"
	positions "go/token"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseComponent() *ast.Component {
	component, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.COMPONENT {
		p.errs.Add(p.file().Position(component), fmt.Sprintf("expected 'component', got: %s", lit))
		return nil
	}
	name := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if name == nil {
		return nil
	}
	inputsInfo := p.parseComponentInputsOrOutputs()
	if inputsInfo.lParen == positions.NoPos {
		return nil
	}
	outputsInfo := p.parseComponentInputsOrOutputs()
	if outputsInfo.lParen == positions.NoPos {
		return nil
	}
	bodyInfo := p.parseComponentBody()
	if bodyInfo.lBrace == positions.NoPos {
		return nil
	}
	return &ast.Component{
		Component:    component,
		Name:         name,
		InputLParen:  inputsInfo.lParen,
		Inputs:       inputsInfo.definitions,
		InputRParen:  inputsInfo.rParen,
		OutputLParen: outputsInfo.lParen,
		Outputs:      outputsInfo.definitions,
		OutputRParen: outputsInfo.rParen,
		LBrace:       bodyInfo.lBrace,
		Entries:      bodyInfo.entries,
		RBRace:       bodyInfo.rBrace,
	}
}

type componentInputsOrOutputs struct {
	lParen      positions.Pos
	definitions *ast.BusDefinitionList
	rParen      positions.Pos
}

func (p *parser) parseComponentInputsOrOutputs() componentInputsOrOutputs {
	lParen, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.LPAREN {
		p.errs.Add(p.file().Position(lParen), fmt.Sprintf("expected '(', got: %s", lit))
		return componentInputsOrOutputs{
			lParen:      positions.NoPos,
			definitions: nil,
			rParen:      positions.NoPos,
		}
	}
	var defs *ast.BusDefinitionList
	if _, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES); tok != tokens.RPAREN {
		defs = p.parseBusDefinitionList()
		if defs == nil {
			return componentInputsOrOutputs{
				lParen:      positions.NoPos,
				definitions: nil,
				rParen:      positions.NoPos,
			}
		}
	}
	rParen, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.RPAREN {
		p.errs.Add(p.file().Position(rParen), fmt.Sprintf("expected ')', got: %s", lit))
		return componentInputsOrOutputs{
			lParen:      positions.NoPos,
			definitions: nil,
			rParen:      positions.NoPos,
		}
	}
	return componentInputsOrOutputs{
		lParen:      lParen,
		definitions: defs,
		rParen:      rParen,
	}
}

type componentBody struct {
	lBrace  positions.Pos
	entries []ast.ComponentEntry
	rBrace  positions.Pos
}

func (p *parser) parseComponentBody() componentBody {
	lBrace, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.LBRACE {
		p.errs.Add(p.file().Position(lBrace), fmt.Sprintf("expected '{', got: %s", lit))
		return componentBody{
			lBrace:  positions.NoPos,
			entries: nil,
			rBrace:  positions.NoPos,
		}
	}
	var entries []ast.ComponentEntry
	var rBrace positions.Pos
parseLoop:
	for {
		pos, tok, lit := p.scanner.Peek(scan.SKIP_NEW_LINES)
		switch tok {
		case tokens.RBRACE:
			p.scanner.Scan(scan.EMIT_NEW_LINES)
			rBrace = pos
			break parseLoop
		case tokens.DEFINE:
			defEntry := p.parseBusDefinitionEntry()
			if defEntry == nil {
				break
			}
			entries = append(entries, defEntry)
			if ok := p.parseNewLine(); ok {
				continue parseLoop
			}
		case tokens.FOR:
			forLoop := p.parseForLoop()
			if forLoop == nil {
				break
			}
			entries = append(entries, forLoop)
		case tokens.IDENTIFIER:
			instance := p.parseInstance()
			if instance == nil {
				break
			}
			entries = append(entries, instance)
			if ok := p.parseNewLine(); ok {
				continue parseLoop
			}
		default:
			p.scanner.Scan(scan.EMIT_NEW_LINES)
			p.errs.Add(p.file().Position(pos), fmt.Sprintf("unexpected token: %s", lit))
		}
		if ok := p.recoverToNewLine(); !ok {
			return componentBody{
				lBrace:  positions.NoPos,
				entries: nil,
				rBrace:  positions.NoPos,
			}
		}
	}
	return componentBody{
		lBrace:  lBrace,
		entries: entries,
		rBrace:  rBrace,
	}
}

func (p *parser) parseBusDefinitionEntry() *ast.BusDefinitionEntry {
	define, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.DEFINE {
		p.errs.Add(p.file().Position(define), fmt.Sprintf("expected 'define', got: %s", lit))
		return nil
	}
	defs := p.parseBusDefinitionList()
	if defs == nil {
		return nil
	}
	return &ast.BusDefinitionEntry{
		Define:      define,
		Definitions: defs,
	}
}

func (p *parser) parseForLoop() *ast.ForLoop {
	for_, tok, lit := p.scanner.Scan(scan.SKIP_NEW_LINES)
	if tok != tokens.FOR {
		p.errs.Add(p.file().Position(for_), fmt.Sprintf("expected 'for', got: %s", lit))
		return nil
	}
	variable := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if variable == nil {
		return nil
	}
	from, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.FROM {
		p.errs.Add(p.file().Position(from), fmt.Sprintf("expected 'from', got: %s", lit))
		return nil
	}
	first := p.parseExpr()
	if first == nil {
		return nil
	}
	to, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.TO {
		p.errs.Add(p.file().Position(to), fmt.Sprintf("expected 'to', got: %s", lit))
		return nil
	}
	last := p.parseExpr()
	if last == nil {
		return nil
	}
	bodyInfo := p.parseComponentBody()
	if bodyInfo.lBrace == positions.NoPos {
		return nil
	}
	return &ast.ForLoop{
		For:      for_,
		Variable: variable,
		From:     from,
		First:    first,
		To:       to,
		Last:     last,
		LBrace:   bodyInfo.lBrace,
		Entries:  bodyInfo.entries,
		RBrace:   bodyInfo.rBrace,
	}
}

func (p *parser) parseInstance() ast.Instance {
	outputs := p.parseBusReferenceList()
	if outputs == nil {
		return nil
	}
	colon, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.COLON {
		p.errs.Add(p.file().Position(colon), fmt.Sprintf("expected ':', got: %s", lit))
		return nil
	}
	pos, tok, lit := p.scanner.Peek(scan.EMIT_NEW_LINES)
	switch tok {
	case tokens.IDENTIFIER:
		return p.parseComponentInstance(outputs, colon)
	case tokens.NUMBER:
		return p.parseConstantsInstance(outputs, colon)
	default:
		p.scanner.Scan(scan.EMIT_NEW_LINES)
		p.errs.Add(p.file().Position(pos), fmt.Sprintf("expected identifier or number, got: %s", lit))
		return nil
	}
}

func (p *parser) parseConstantsInstance(outputs *ast.BusReferenceList, colon positions.Pos) *ast.ConstantsInstace {
	constants := p.parseConstants()
	if constants == nil {
		return nil
	}
	return &ast.ConstantsInstace{
		Outputs:   outputs,
		Colon:     colon,
		Constants: constants,
	}
}

func (p *parser) parseComponentInstance(outputs *ast.BusReferenceList, colon positions.Pos) *ast.ComponentInstance {
	defName := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if defName == nil {
		return nil
	}
	lParen, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.LPAREN {
		p.errs.Add(p.file().Position(lParen), fmt.Sprintf("expected '(', got: %s", lit))
		return nil
	}
	var inputs *ast.BusReferenceList
	if _, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES); tok != tokens.RPAREN {
		inputs = p.parseBusReferenceList()
		if inputs == nil {
			return nil
		}
	}
	rParen, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.RPAREN {
		p.errs.Add(p.file().Position(rParen), fmt.Sprintf("expected ')', got: %s", lit))
		return nil
	}
	return &ast.ComponentInstance{
		Outputs:        outputs,
		Colon:          colon,
		DefinitionName: defName,
		LParen:         lParen,
		Inputs:         inputs,
		RParen:         rParen,
	}
}

func (p *parser) parseBusDefinitionList() *ast.BusDefinitionList {
	defs := make([]*ast.BusDefinition, 0, 1)
	for {
		def := p.parseBusDefinition()
		if def == nil {
			return nil
		}
		defs = append(defs, def)
		_, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES)
		if tok == tokens.COMMA {
			p.scanner.Scan(scan.EMIT_NEW_LINES)
		} else {
			break
		}
	}
	return &ast.BusDefinitionList{
		Defintions: defs,
	}
}

func (p *parser) parseBusDefinition() *ast.BusDefinition {
	name := p.parseIdentifier(scan.EMIT_NEW_LINES)
	if name == nil {
		return nil
	}
	lBrack, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES)
	if tok != tokens.LBRACK {
		return &ast.BusDefinition{
			Name:      name,
			LBrack:    positions.NoPos,
			WireCount: nil,
			RBrack:    positions.NoPos,
		}
	}
	p.scanner.Scan(scan.EMIT_NEW_LINES)
	wireCount := p.parseNumber()
	if wireCount == nil {
		return nil
	}
	rBrack, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.RBRACK {
		p.errs.Add(p.file().Position(rBrack), fmt.Sprintf("expected ']', got: %s", lit))
		return nil
	}
	return &ast.BusDefinition{
		Name:      name,
		LBrack:    lBrack,
		WireCount: wireCount,
		RBrack:    rBrack,
	}
}
