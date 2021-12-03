package parser

import (
	"fmt"

	"github.com/chai2010/ugo/ast"
	"github.com/chai2010/ugo/lexer"
	"github.com/chai2010/ugo/token"
)

func ParseFile(filename, src string) (*ast.File, error) {
	p := NewParser(filename, src)
	return p.ParseFile()
}

type Parser struct {
	filename string
	src      string

	*TokenStream
	file *ast.File
	err  error
}

func NewParser(filename, src string) *Parser {
	return &Parser{filename: filename, src: src}
}

func (p *Parser) ParseFile() (file *ast.File, err error) {
	defer func() {
		file, err = p.file, p.err
	}()

	tokens, comments := lexer.Lex(p.filename, p.src)
	for _, tok := range tokens {
		if tok.Type == token.ERROR {
			p.errorf(tok.Pos, "invalid token: %s", tok.Literal)
		}
	}

	p.TokenStream = NewTokenStream(p.filename, p.src, tokens, comments)
	p.parseFile()

	return
}

func (p *Parser) errorf(pos int, format string, args ...interface{}) {
	p.err = fmt.Errorf("%s: %s",
		lexer.PosString(p.filename, p.src, pos),
		fmt.Sprintf(format, args...),
	)
	panic(p.err)
}
