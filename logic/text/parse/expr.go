package parse

import (
	"fmt"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseExpr() ast.Expr {
	return p.parseExprWithPrecedence(precedence(0))
}

type precedence int

func operatorPrecedence(tok tokens.Token) precedence {
	switch tok {
	case tokens.ADD, tokens.SUB:
		return 1
	case tokens.MUL, tokens.QUO, tokens.REM:
		return 2
	default:
		panic(fmt.Errorf("unexpected operator token: %v", tok))
	}
}

func (p *parser) parseExprWithPrecedence(pre precedence) (expr ast.Expr) {
	switch pos, tok, lit := p.scanner.Peek(scan.EMIT_NEW_LINES); tok {
	case tokens.IDENTIFIER, tokens.NUMBER:
		expr = p.parseLiteral()
	case tokens.ADD, tokens.SUB:
		expr = p.parseUnaryExpr()
	default:
		p.scanner.Scan(scan.EMIT_NEW_LINES)
		p.errs.Add(p.file().Position(pos), fmt.Sprintf("expected expression, got: %s", lit))
		return nil
	}
	if expr == nil {
		return
	}
parseLoop:
	for {
		switch _, tok, _ := p.scanner.Peek(scan.EMIT_NEW_LINES); tok {
		case tokens.ADD, tokens.SUB, tokens.MUL, tokens.QUO, tokens.REM:
			if operatorPrecedence(tok) <= pre {
				break parseLoop
			}
			opStart, op, _ := p.scanner.Scan(scan.EMIT_NEW_LINES)
			rhsOperand := p.parseExprWithPrecedence(pre + 1)
			expr = &ast.BinaryExpr{
				Operator:      op,
				OperatorStart: opStart,
				LhsOperand:    expr,
				RhsOperand:    rhsOperand,
			}
		default:
			break parseLoop
		}
	}
	return
}

func (p *parser) parseUnaryExpr() *ast.UnaryExpr {
	opStart, op, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	switch op {
	case tokens.ADD, tokens.SUB:
		break
	default:
		p.errs.Add(p.file().Position(opStart), fmt.Sprintf("expected '+' or '-', got: %s", lit))
		return nil
	}
	operand := p.parseLiteral()
	if operand == nil {
		return nil
	}
	return &ast.UnaryExpr{
		Operator:      op,
		OperatorStart: opStart,
		Operand:       operand,
	}
}
