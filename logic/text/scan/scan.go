package scan

import (
	"fmt"
	positions "go/token"

	"github.com/arneph/mercury/logic/text/tokens"
)

type Scanner struct {
	file *positions.File
	src  []byte

	offset int
}

func NewScanner(file *positions.File, src []byte) *Scanner {
	return &Scanner{
		file: file,
		src:  src,
	}
}

func (s *Scanner) File() *positions.File {
	return s.file
}

type NewLineMode int

const (
	EMIT_NEW_LINES NewLineMode = iota
	SKIP_NEW_LINES
)

func (s *Scanner) Peek(newLineMode NewLineMode) (pos positions.Pos, tok tokens.Token, lit string) {
	previousOffset := s.offset
	pos, tok, lit = s.Scan(newLineMode)
	s.offset = previousOffset
	return
}

func (s *Scanner) Scan(newLineMode NewLineMode) (pos positions.Pos, tok tokens.Token, lit string) {
	switch newLineMode {
	case EMIT_NEW_LINES:
		s.skipWhitespace()
	case SKIP_NEW_LINES:
		s.skipWhitespaceAndNewLines()
	default:
		panic(fmt.Errorf("unknown NewLineMode: %v", newLineMode))
	}
	pos = s.file.Pos(s.offset)
	if s.offset >= len(s.src) {
		return pos, tokens.EOF, ""
	}
	switch ch := s.src[s.offset]; ch {
	case '\n':
		tok = tokens.NEWLINE
		lit = "\n"
	case '+':
		tok = tokens.ADD
		lit = "+"
	case '-':
		tok = tokens.SUB
		lit = "-"
	case '*':
		tok = tokens.MUL
		lit = "*"
	case '/':
		tok = tokens.QUO
		lit = "/"
	case '%':
		tok = tokens.REM
		lit = "%"
	case '(':
		tok = tokens.LPAREN
		lit = "("
	case '[':
		tok = tokens.LBRACK
		lit = "["
	case '{':
		tok = tokens.LBRACE
		lit = "{"
	case ')':
		tok = tokens.RPAREN
		lit = ")"
	case ']':
		tok = tokens.RBRACK
		lit = "]"
	case '}':
		tok = tokens.RBRACE
		lit = "}"
	case ',':
		tok = tokens.COMMA
		lit = ","
	case ':':
		tok = tokens.COLON
		lit = ":"
	default:
		lit = string(ch)
		if isIdentifierStart(ch) {
			for i := s.offset + 1; i < len(s.src) && isIdentifierCharacter(s.src[i]); i++ {
				lit += string(s.src[i])
			}
			switch lit {
			case "component":
				tok = tokens.COMPONENT
			case "test":
				tok = tokens.TEST
			case "define":
				tok = tokens.DEFINE
			case "set":
				tok = tokens.SET
			case "assert":
				tok = tokens.ASSERT
			case "expect":
				tok = tokens.EXPECT
			case "is":
				tok = tokens.IS
			case "for":
				tok = tokens.FOR
			case "from":
				tok = tokens.FROM
			case "to":
				tok = tokens.TO
			default:
				tok = tokens.IDENTIFIER
			}
		} else if isNumberCharacter(ch) {
			i := s.offset + 1
			ok := isNumberCharacter
			hex := ch == '0' && i < len(s.src) && s.src[i] == 'x'
			if hex {
				lit += "x"
				i++
				ok = isHexCharacter
			}
			for ; i < len(s.src) && ok(s.src[i]); i++ {
				lit += string(s.src[i])
			}
			if hex && len(lit) < 3 {
				tok = tokens.ERROR
			} else {
				tok = tokens.NUMBER
			}
		} else {
			tok = tokens.ERROR
		}
	}
	s.offset += len(lit)
	return
}

func isIdentifierStart(b byte) bool {
	return b == '\'' || b == '_' || ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z')
}

func isIdentifierCharacter(b byte) bool {
	return isIdentifierStart(b) || isNumberCharacter(b)
}

func isNumberCharacter(b byte) bool {
	return ('0' <= b && b <= '9')
}
func isHexCharacter(b byte) bool {
	return isNumberCharacter(b) || ('a' <= b && b <= 'f')
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t'
}

func isWhitespaceOrNewLine(b byte) bool {
	return isWhitespace(b) || b == '\n'
}

func (s *Scanner) skipWhitespace() {
	for s.offset < len(s.src) && isWhitespace(s.src[s.offset]) {
		s.offset += 1
	}
}

func (s *Scanner) skipWhitespaceAndNewLines() {
	for s.offset < len(s.src) && isWhitespaceOrNewLine(s.src[s.offset]) {
		s.offset += 1
	}
}
