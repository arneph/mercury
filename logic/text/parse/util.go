package parse

import (
	"fmt"

	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) parseNewLine() bool {
	pos, tok, lit := p.scanner.Scan(scan.EMIT_NEW_LINES)
	if tok != tokens.NEWLINE {
		p.errs.Add(p.file().Position(pos), fmt.Sprintf("expected new line, got: %s", lit))
		return false
	} else {
		return true
	}
}
