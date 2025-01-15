package scan

import (
	positions "go/token"
	"testing"

	"github.com/arneph/mercury/logic/text/tokens"
	"github.com/stretchr/testify/assert"
)

func TestEmptyFile(t *testing.T) {
	fileSet := positions.NewFileSet()
	file := fileSet.AddFile("fake.mercury", fileSet.Base(), 0)
	s := NewScanner(file, nil)
	if f := s.File(); f != file {
		t.Errorf("Scanner.File() = %p; want %p", f, file)
	}
	pos, tok, lit := s.Peek(EMIT_NEW_LINES)
	if pos != file.Pos(0) {
		t.Errorf("pos = %v; want %v", file.Pos(0), pos)
	}
	if tok != tokens.EOF {
		t.Errorf("tok = %v; want EOF", tok)
	}
	if lit != "" {
		t.Errorf("lit = %q; want empty string", lit)
	}

	pos, tok, lit = s.Scan(EMIT_NEW_LINES)
	if pos != file.Pos(0) {
		t.Errorf("pos = %v; want %v", file.Pos(0), pos)
	}
	if tok != tokens.EOF {
		t.Errorf("tok = %v; want EOF", tok)
	}
	if lit != "" {
		t.Errorf("lit = %q; want empty string", lit)
	}

	pos, tok, lit = s.Scan(EMIT_NEW_LINES)
	if pos != file.Pos(0) {
		t.Errorf("pos = %v; want %v", file.Pos(0), pos)
	}
	if tok != tokens.EOF {
		t.Errorf("tok = %v; want EOF", tok)
	}
	if lit != "" {
		t.Errorf("lit = %q; want empty string", lit)
	}

	pos, tok, lit = s.Scan(SKIP_NEW_LINES)
	if pos != file.Pos(0) {
		t.Errorf("pos = %v; want %v", file.Pos(0), pos)
	}
	if tok != tokens.EOF {
		t.Errorf("tok = %v; want EOF", tok)
	}
	if lit != "" {
		t.Errorf("lit = %q; want empty string", lit)
	}
}

func TestIndividualTokens(t *testing.T) {
	testcases := []struct {
		src         []byte
		expectedTok tokens.Token
	}{
		{
			src:         []byte("\n"),
			expectedTok: tokens.NEWLINE,
		},
		{
			src:         []byte("+"),
			expectedTok: tokens.ADD,
		},
		{
			src:         []byte("-"),
			expectedTok: tokens.SUB,
		},
		{
			src:         []byte("*"),
			expectedTok: tokens.MUL,
		},
		{
			src:         []byte("/"),
			expectedTok: tokens.QUO,
		},
		{
			src:         []byte("%"),
			expectedTok: tokens.REM,
		},
		{
			src:         []byte("("),
			expectedTok: tokens.LPAREN,
		},
		{
			src:         []byte(")"),
			expectedTok: tokens.RPAREN,
		},
		{
			src:         []byte("["),
			expectedTok: tokens.LBRACK,
		},
		{
			src:         []byte("]"),
			expectedTok: tokens.RBRACK,
		},
		{
			src:         []byte("{"),
			expectedTok: tokens.LBRACE,
		},
		{
			src:         []byte("}"),
			expectedTok: tokens.RBRACE,
		},
		{
			src:         []byte(","),
			expectedTok: tokens.COMMA,
		},
		{
			src:         []byte(":"),
			expectedTok: tokens.COLON,
		},
		{
			src:         []byte("$"),
			expectedTok: tokens.ERROR,
		},
		{
			src:         []byte("component"),
			expectedTok: tokens.COMPONENT,
		},
		{
			src:         []byte("test"),
			expectedTok: tokens.TEST,
		},
		{
			src:         []byte("define"),
			expectedTok: tokens.DEFINE,
		},
		{
			src:         []byte("set"),
			expectedTok: tokens.SET,
		},
		{
			src:         []byte("assert"),
			expectedTok: tokens.ASSERT,
		},
		{
			src:         []byte("expect"),
			expectedTok: tokens.EXPECT,
		},
		{
			src:         []byte("is"),
			expectedTok: tokens.IS,
		},
		{
			src:         []byte("for"),
			expectedTok: tokens.FOR,
		},
		{
			src:         []byte("from"),
			expectedTok: tokens.FROM,
		},
		{
			src:         []byte("to"),
			expectedTok: tokens.TO,
		},
		{
			src:         []byte("0"),
			expectedTok: tokens.NUMBER,
		},
		{
			src:         []byte("7"),
			expectedTok: tokens.NUMBER,
		},
		{
			src:         []byte("583563"),
			expectedTok: tokens.NUMBER,
		},
		{
			src:         []byte("0x0"),
			expectedTok: tokens.NUMBER,
		},
		{
			src:         []byte("0xfedcba9876543210"),
			expectedTok: tokens.NUMBER,
		},
		{
			src:         []byte("0x"),
			expectedTok: tokens.ERROR,
		},
	}
	for _, testcase := range testcases {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(testcase.src))
		s := NewScanner(file, testcase.src)
		pos, tok, lit := s.Scan(EMIT_NEW_LINES)
		if pos != file.Pos(0) {
			t.Errorf("pos = %v; want %v", file.Pos(0), pos)
		}
		if tok != testcase.expectedTok {
			t.Errorf("tok = %v; want %v", tok, testcase.expectedTok)
		}
		if lit != string(testcase.src) {
			t.Errorf("lit = %q; want %q", lit, testcase.src)
		}
	}
}

func TestSkipsWhitespace(t *testing.T) {
	testcases := [][]byte{
		[]byte(" +"),
		[]byte("       +"),
		[]byte("\t+"),
		[]byte("\t\t\t\t+"),
		[]byte(" \t    \t+"),
		[]byte("\t\t\t \t    +"),
	}
	for _, testcase := range testcases {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(testcase))
		s := NewScanner(file, testcase)
		pos, tok, lit := s.Scan(EMIT_NEW_LINES)
		if pos != file.Pos(len(testcase)-1) {
			t.Errorf("pos = %v; want %v", file.Pos(len(testcase)-1), pos)
		}
		if tok != tokens.ADD {
			t.Errorf("tok = %v; want tokens.ADD", tok)
		}
		if lit != "+" {
			t.Errorf("lit = %q; want '+'", lit)
		}
	}
}

func TestSkipsWhitespaceAndNewLines(t *testing.T) {
	testcases := [][]byte{
		[]byte(" +"),
		[]byte("       +"),
		[]byte("\t+"),
		[]byte("\t\t\t\t+"),
		[]byte(" \t    \t+"),
		[]byte("\t\t\t \t    +"),
		[]byte("\n+"),
		[]byte("\n\n\n\n+"),
		[]byte("\n +"),
		[]byte("   \n    +"),
		[]byte("\n\t+"),
		[]byte("\t\t\n\t\t+"),
		[]byte(" \t \n   \t+"),
		[]byte("\t\n\t\t \t    \n+"),
	}
	for _, testcase := range testcases {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(testcase))
		s := NewScanner(file, testcase)
		pos, tok, lit := s.Scan(SKIP_NEW_LINES)
		if pos != file.Pos(len(testcase)-1) {
			t.Errorf("pos = %v; want %v", file.Pos(len(testcase)-1), pos)
		}
		if tok != tokens.ADD {
			t.Errorf("tok = %v; want tokens.ADD", tok)
		}
		if lit != "+" {
			t.Errorf("lit = %q; want '+'", lit)
		}
	}
}

func TestPanicsForUnknownNewLineMode(t *testing.T) {
	assert.Panics(t, func() {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), 0)
		s := NewScanner(file, nil)
		s.Scan(NewLineMode(-1))
	})
}

func TestSeveralTokens(t *testing.T) {
	type scanResult struct {
		posOffset int
		tok       tokens.Token
		lit       string
	}
	testcases := []struct {
		src      []byte
		expected []scanResult
	}{
		{
			src: []byte("component Multiplex(input[16], s[4]) (output) {}"),
			expected: []scanResult{
				{0, tokens.COMPONENT, "component"},
				{10, tokens.IDENTIFIER, "Multiplex"},
				{19, tokens.LPAREN, "("},
				{20, tokens.IDENTIFIER, "input"},
				{25, tokens.LBRACK, "["},
				{26, tokens.NUMBER, "16"},
				{28, tokens.RBRACK, "]"},
				{29, tokens.COMMA, ","},
				{31, tokens.IDENTIFIER, "s"},
				{32, tokens.LBRACK, "["},
				{33, tokens.NUMBER, "4"},
				{34, tokens.RBRACK, "]"},
				{35, tokens.RPAREN, ")"},
				{37, tokens.LPAREN, "("},
				{38, tokens.IDENTIFIER, "output"},
				{44, tokens.RPAREN, ")"},
				{46, tokens.LBRACE, "{"},
				{47, tokens.RBRACE, "}"},
			},
		},
		{
			src: []byte("test MyTest {\n component: MyComp\n set a, b: 5, 7\n expect x, y: 8, 2\n}"),
			expected: []scanResult{
				{0, tokens.TEST, "test"},
				{5, tokens.IDENTIFIER, "MyTest"},
				{12, tokens.LBRACE, "{"},
				{13, tokens.NEWLINE, "\n"},
				{15, tokens.COMPONENT, "component"},
				{24, tokens.COLON, ":"},
				{26, tokens.IDENTIFIER, "MyComp"},
				{32, tokens.NEWLINE, "\n"},
				{34, tokens.SET, "set"},
				{38, tokens.IDENTIFIER, "a"},
				{39, tokens.COMMA, ","},
				{41, tokens.IDENTIFIER, "b"},
				{42, tokens.COLON, ":"},
				{44, tokens.NUMBER, "5"},
				{45, tokens.COMMA, ","},
				{47, tokens.NUMBER, "7"},
				{48, tokens.NEWLINE, "\n"},
				{50, tokens.EXPECT, "expect"},
				{57, tokens.IDENTIFIER, "x"},
				{58, tokens.COMMA, ","},
				{60, tokens.IDENTIFIER, "y"},
				{61, tokens.COLON, ":"},
				{63, tokens.NUMBER, "8"},
				{64, tokens.COMMA, ","},
				{66, tokens.NUMBER, "2"},
				{67, tokens.NEWLINE, "\n"},
				{68, tokens.RBRACE, "}"},
			},
		},
	}
	for _, testcase := range testcases {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(testcase.src))
		s := NewScanner(file, testcase.src)
		for _, expected := range testcase.expected {
			pos, tok, lit := s.Scan(EMIT_NEW_LINES)
			if pos != file.Pos(expected.posOffset) {
				t.Errorf("pos = %v; want %v", pos, file.Pos(expected.posOffset))
			}
			if tok != expected.tok {
				t.Errorf("tok = %v; want %v", tok, expected.tok)
			}
			if lit != string(expected.lit) {
				t.Errorf("lit = %q; want %q", lit, expected.lit)
			}
		}
		pos, tok, lit := s.Scan(SKIP_NEW_LINES)
		if pos != file.Pos(len(testcase.src)) {
			t.Errorf("pos = %v; want %v", pos, file.Pos(len(testcase.src)))
		}
		if tok != tokens.EOF {
			t.Errorf("tok = %v; want EOF", tok)
		}
		if lit != "" {
			t.Errorf("lit = %q; want empty string", lit)
		}
	}
}

func FuzzScanner(f *testing.F) {
	f.Fuzz(func(t *testing.T, in []byte) {
		fileSet := positions.NewFileSet()
		file := fileSet.AddFile("fake.mercury", fileSet.Base(), len(in))
		s := NewScanner(file, in)
		for {
			pos, tok, _ := s.Scan(EMIT_NEW_LINES)
			if pos < file.Pos(0) || pos > file.Pos(len(in)) {
				t.Fatalf("pos = %v; want between %v and %v", pos, file.Pos(0), file.Pos(len(in)))
			}
			if tok < tokens.ERROR || tok > tokens.TO {
				t.Fatalf("tok = %v; want defined token value", tok)
			}
			if tok == tokens.EOF {
				break
			}
		}
	})
}
