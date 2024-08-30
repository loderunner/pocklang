package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"unicode"
)

type TokenType int

const (
	// Logical operators
	Or TokenType = iota
	And

	// Comparison operators
	Lt
	Lte
	Gt
	Gte
	Eq
	Neq

	// Mathematical operators
	Plus
	Minus
	Star
	Slash
	Not

	// Misc
	LeftParen
	RightParen
	Dot

	// Keywords
	True
	False
	Null

	// Literal types
	Integer
	Decimal
	String
	Identifier

	// EOF
	EOF
)

var reservedRunes = []rune{'|', '&', '<', '>', '=', '+', '-', '*', '/', '!', '"', '.', '(', ')'}

var (
	ScanError       = errors.New("scan error")
	whitespaceError = errors.New("whitespace")
)

type Token struct {
	Type   TokenType
	Lexeme string

	IntegerValue    int64
	DecimalValue    float64
	StringValue     string
	IdentifierValue string
}

type scanner struct {
	io.RuneScanner

	current rune
	buf     *bytes.Buffer
}

func (s scanner) advance() (rune, error) {
	r, sz, err := s.ReadRune()
	if err != nil {
		return 0, err
	}
	if r == 0xfffd && sz == 1 {
		return 0, fmt.Errorf("invalid UTF-8 sequence")
	}
	s.buf.WriteRune(r)
	return r, nil
}

func (s scanner) backtrack() error {
	l := s.buf.Len()
	err := s.UnreadRune()
	if err != nil {
		return err
	}
	s.buf.Truncate(l - 1)
	return nil
}

func (s scanner) match(expected rune) (bool, error) {
	r, err := s.advance()
	if err != nil {
		return false, err
	}
	if r != expected {
		err = s.backtrack()
		if err != nil {
			return false, err
		}
		return false, nil
	}

	return true, nil
}

func Scan(rs io.RuneScanner) ([]Token, error) {
	var tok Token
	var err error
	s := scanner{RuneScanner: rs, buf: new(bytes.Buffer)}
	tokens := make([]Token, 0)

	for err == nil || err == whitespaceError {
		tok, err = scanToken(s)
		if err == nil {
			tokens = append(tokens, tok)
		}
	}

	if !errors.Is(err, io.EOF) {
		return nil, err
	}

	tokens = append(tokens, Token{Type: EOF})
	return tokens, nil
}

func scanToken(s scanner) (Token, error) {
	defer s.buf.Reset()

	r, err := s.advance()
	if err != nil {
		return Token{}, fmt.Errorf("%w: %w", ScanError, err)
	}

	switch r {
	case '+':
		return Token{Type: Plus, Lexeme: s.buf.String()}, nil
	case '-':
		return Token{Type: Minus, Lexeme: s.buf.String()}, nil
	case '*':
		return Token{Type: Star, Lexeme: s.buf.String()}, nil
	case '/':
		return Token{Type: Slash, Lexeme: s.buf.String()}, nil
	case '(':
		return Token{Type: LeftParen, Lexeme: s.buf.String()}, nil
	case ')':
		return Token{Type: RightParen, Lexeme: s.buf.String()}, nil
	case '.':
		return Token{Type: Dot, Lexeme: s.buf.String()}, nil
	case '|':
		ok, err := s.match('|')
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if ok {
			return Token{Type: Or, Lexeme: s.buf.String()}, nil
		}
		return Token{}, fmt.Errorf("%w: expected '|' after '|'", ScanError)
	case '&':
		ok, err := s.match('&')
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if ok {
			return Token{Type: And, Lexeme: s.buf.String()}, nil
		}
		return Token{}, fmt.Errorf("%w: expected '&&' after '&&'", ScanError)
	case '=':
		ok, err := s.match('=')
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if ok {
			return Token{Type: Eq, Lexeme: s.buf.String()}, nil
		}
		return Token{}, fmt.Errorf("%w: expected '=' after '='", ScanError)
	case '!':
		ok, err := s.match('=')
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if ok {
			return Token{Type: Neq, Lexeme: s.buf.String()}, nil
		}
		return Token{Type: Not, Lexeme: s.buf.String()}, nil
	case '<':
		ok, err := s.match('=')
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if ok {
			return Token{Type: Lte, Lexeme: s.buf.String()}, nil
		}
		return Token{Type: Lt, Lexeme: s.buf.String()}, nil
	case '>':
		ok, err := s.match('=')
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if ok {
			return Token{Type: Gte, Lexeme: s.buf.String()}, nil
		}
		return Token{Type: Gt, Lexeme: s.buf.String()}, nil
	case '"':
		for r, err = s.advance(); r != '"' && err == nil; r, err = s.advance() {
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{}, fmt.Errorf("%w: unterminated string", ScanError)
			}
			return Token{}, err
		}
		lex := s.buf.String()
		return Token{
			Type:        String,
			Lexeme:      lex,
			StringValue: lex[1 : len(lex)-1],
		}, nil
	default:
		if isDigit(r) {
			isDecimal := false
			for r, err = s.advance(); (isDigit(r) || r == '.') && err == nil; r, err = s.advance() {
				if r == '.' {
					isDecimal = true
				}
			}
			if err != nil && !errors.Is(err, io.EOF) {
				return Token{}, err
			}
			if !errors.Is(err, io.EOF) {
				s.backtrack()
			}
			lex := s.buf.String()
			if isDecimal {
				val, err := strconv.ParseFloat(lex, 64)
				if err != nil {
					return Token{}, fmt.Errorf("%w: invalid number: %w", ScanError, err)
				}
				return Token{
					Type:         Decimal,
					Lexeme:       lex,
					DecimalValue: val,
				}, nil
			} else {
				val, err := strconv.ParseInt(lex, 10, 64)
				if err != nil {
					return Token{}, fmt.Errorf("%w: invalid number: %w", ScanError, err)
				}
				return Token{
					Type:         Integer,
					Lexeme:       lex,
					IntegerValue: val,
				}, nil
			}
		}
		if unicode.IsSpace(r) {
			for r, err = s.advance(); unicode.IsSpace(r) && err == nil; r, err = s.advance() {
			}
			if err != nil && !errors.Is(err, io.EOF) {
				return Token{}, err
			}
			if !errors.Is(err, io.EOF) {
				s.backtrack()
			}
			return Token{}, whitespaceError
		}

		for r, err = s.advance(); isIdent(r) && err == nil; r, err = s.advance() {
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return Token{}, err
		}
		if !errors.Is(err, io.EOF) {
			s.backtrack()
		}
		lex := s.buf.String()
		switch lex {
		case "true":
			return Token{Type: True, Lexeme: lex}, nil
		case "false":
			return Token{Type: False, Lexeme: lex}, nil
		case "null":
			return Token{Type: Null, Lexeme: lex}, nil
		default:
			return Token{Type: Identifier, Lexeme: lex, IdentifierValue: lex}, nil
		}
	}
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdent(r rune) bool {
	return !slices.Contains(reservedRunes, r) && !unicode.IsSpace(r)
}
