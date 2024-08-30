package parser

import (
	"errors"
	"fmt"
	"io"
)

// Syntactical grammar:
// Expr    -> Or ;
// Or      -> And ("||" And)* ;
// And     -> Comp ("&&" Comp)* ;
// Comp    -> Term (("<" | ">" | ">=" | "<=" | "==" | "!=") Term) ;
// Term    -> Factor (("+" | "-") Factor)* ;
// Factor  -> Unary (("*" | "/") Unary)* ;
// Unary   -> ("!" | "-") Primary ;
// Primary -> "true" | "false" | "null" | INTEGER | DECIMAL | STRING | "(" Expression ")" | IDENTIFIER ("." IDENTIFIER)* ;

var ParseError = errors.New("parse error")

func Parse(tokens []Token) (Expr, error) {
	var err error
	p := parser{tokens: tokens}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if !p.eof() {
		return nil,
			fmt.Errorf(
				"%w at '%s': expected end of expression",
				ParseError,
				p.peek().Lexeme,
			)
	}
	return expr, nil
}

type parser struct {
	current int
	tokens  []Token
}

func (p parser) eof() bool {
	return p.current >= len(p.tokens)
}

func (p parser) peek() Token {
	if p.eof() {
		return Token{}
	}
	return p.tokens[p.current]
}

func (p *parser) advance() (Token, error) {
	p.current++
	if p.eof() {
		return Token{}, io.EOF
	}
	return p.peek(), nil
}

func (p *parser) parseExpr() (Expr, error) {
	expr, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *parser) parseOr() (Expr, error) {
	expr, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.peek().Type == Or {
		op, _ := p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{Op: op.Type, Left: expr, Right: right}
	}

	return expr, nil
}

func (p *parser) parseAnd() (Expr, error) {
	expr, err := p.parseComp()
	if err != nil {
		return nil, err
	}

	for p.peek().Type == And {
		op, _ := p.advance()
		right, err := p.parseComp()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{Op: op.Type, Left: expr, Right: right}
	}

	return expr, nil
}

func (p *parser) parseComp() (Expr, error) {
	expr, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	peekType := p.peek().Type
	if peekType == Lt ||
		peekType == Lte ||
		peekType == Gt ||
		peekType == Gte ||
		peekType == Eq ||
		peekType == Neq {
		p.advance()
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		return BinaryExpr{Op: peekType, Left: expr, Right: right}, nil
	}

	return expr, nil
}

func (p *parser) parseTerm() (Expr, error) {
	expr, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	peekType := p.peek().Type
	if peekType == Plus || peekType == Minus {
		p.advance()
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		return BinaryExpr{Op: peekType, Left: expr, Right: right}, nil
	}

	return expr, nil
}

func (p *parser) parseFactor() (Expr, error) {
	expr, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	peekType := p.peek().Type
	if peekType == Star || peekType == Slash {
		p.advance()
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return BinaryExpr{Op: peekType, Left: expr, Right: right}, nil
	}

	return expr, nil
}

func (p *parser) parseUnary() (Expr, error) {
	peekType := p.peek().Type
	if peekType == Not || peekType == Minus {
		p.advance()
		expr, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return UnaryExpr{Op: peekType, Expr: expr}, nil
	}

	return p.parsePrimary()
}

func (p *parser) parsePrimary() (Expr, error) {
	if p.eof() {
		return nil, fmt.Errorf("%w: unexpected end of expression", ParseError)
	}

	tok := p.peek()
	switch tok.Type {
	case True, False, Null, Integer, Decimal, String:
		p.advance()
		return LiteralExpr{Token: tok}, nil
	case LeftParen:
		return p.parseGroup()
	case Identifier:
		return p.parseGet()
	}

	return nil, fmt.Errorf("%w at '%s': unexpected token", ParseError, tok.Lexeme)
}

func (p *parser) parseGroup() (Expr, error) {
	p.advance()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.peek().Type != RightParen {
		return nil, fmt.Errorf("%w: missing closing parenthesis", ParseError)
	}
	p.advance()
	return GroupExpr{Expr: expr}, nil
}

func (p *parser) parseGet() (Expr, error) {
	names := []string{p.peek().Lexeme}
	for p.advance(); p.peek().Type == Dot; p.advance() {
		p.advance()
		tok := p.peek()
		if tok.Type != Identifier {
			return nil, fmt.Errorf("%w at '%s': expected identifier after '.'", ParseError, tok.Lexeme)
		}
		names = append(names, tok.Lexeme)
	}
	return GetExpr{Names: names}, nil
}
