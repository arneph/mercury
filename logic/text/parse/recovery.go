package parse

import (
	"github.com/arneph/mercury/logic/text/scan"
	"github.com/arneph/mercury/logic/text/tokens"
)

func (p *parser) recoverToNewLine() bool {
recoveryLoop:
	for {
		switch _, tok, _ := p.scanner.Scan(scan.EMIT_NEW_LINES); tok {
		case tokens.NEWLINE:
			break recoveryLoop
		case tokens.EOF:
			return false
		}
	}
	return true
}
