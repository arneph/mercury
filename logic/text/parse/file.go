package parse

import (
	"fmt"

	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseFile() *ast.File {
	var nodes []ast.FileNode
parseLoop:
	for {
		pos, tok, lit := p.scanner.Peek(scan.SKIP_NEW_LINES)
		switch tok {
		case tokens.EOF:
			p.scanner.Scan(scan.SKIP_NEW_LINES)
			break parseLoop
		case tokens.COMPONENT:
			component := p.parseComponent()
			if component != nil {
				nodes = append(nodes, component)
			}
		case tokens.TEST:
			test := p.parseTest()
			if test != nil {
				nodes = append(nodes, test)
			}
		default:
			p.scanner.Scan(scan.SKIP_NEW_LINES)
			p.errs.Add(p.file().Position(pos), fmt.Sprintf("unexpected token: %s", lit))
			if ok := p.recoverToNewLine(); !ok {
				return nil
			}
		}
	}
	return &ast.File{
		FileStart: p.file().Pos(0),
		FileEnd:   p.file().Pos(p.file().Size()),
		Nodes:     nodes,
	}
}
