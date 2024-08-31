package pock

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScannerEmpty(t *testing.T) {
	tokens, err := Scan(strings.NewReader(""))
	require.NoError(t, err)
	require.Len(t, tokens, 0)
}

func TestScannerTokenType(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected TokenType
	}
	cases := []testCase{
		{name: "Or", input: "||", expected: Or},
		{name: "And", input: "&&", expected: And},
		{name: "Lt", input: "<", expected: Lt},
		{name: "Lte", input: "<=", expected: Lte},
		{name: "Gt", input: ">", expected: Gt},
		{name: "Gte", input: ">=", expected: Gte},
		{name: "Eq", input: "==", expected: Eq},
		{name: "Neq", input: "!=", expected: Neq},
		{name: "Plus", input: "+", expected: Plus},
		{name: "Minus", input: "-", expected: Minus},
		{name: "Star", input: "*", expected: Star},
		{name: "Slash", input: "/", expected: Slash},
		{name: "Not", input: "!", expected: Not},
		{name: "LeftParen", input: "(", expected: LeftParen},
		{name: "RightParen", input: ")", expected: RightParen},
		{name: "Dot", input: ".", expected: Dot},
		{name: "True", input: "true", expected: True},
		{name: "False", input: "false", expected: False},
		{name: "Null", input: "null", expected: Null},
		{name: "Integer", input: "123", expected: Integer},
		{name: "Decimal", input: "123.45", expected: Decimal},
		{name: "String", input: `"Hello World!"`, expected: String},
		{name: "Identifier", input: "hello_world", expected: Identifier},
	}

	t.Parallel()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c.input))
			require.NoError(t, err)
			require.Len(t, tokens, 1)
			require.Equal(t, c.expected, tokens[0].Type)
		})
	}
}

func TestScannerIntegerValue(t *testing.T) {
	tokens, err := Scan(strings.NewReader("123"))
	require.NoError(t, err)
	require.EqualValues(t, 123, tokens[0].IntegerValue)
}

func TestScannerDecimalValue(t *testing.T) {
	tokens, err := Scan(strings.NewReader("123.45"))
	require.NoError(t, err)
	require.EqualValues(t, 123.45, tokens[0].DecimalValue)
}

func TestScannerStringValue(t *testing.T) {
	tokens, err := Scan(strings.NewReader(`"Hello World!"`))
	require.NoError(t, err)
	require.EqualValues(t, "Hello World!", tokens[0].StringValue)
}

func TestScannerIdentifierValue(t *testing.T) {
	tokens, err := Scan(strings.NewReader("hello_world"))
	require.NoError(t, err)
	require.EqualValues(t, "hello_world", tokens[0].IdentifierValue)
}

func TestScannerSequence(t *testing.T) {
	type testCase struct {
		input    string
		expected []Token
	}
	cases := []testCase{
		{
			input: "1 == 2",
			expected: []Token{
				{Type: Integer, IntegerValue: 1},
				{Type: Eq},
				{Type: Integer, IntegerValue: 2},
			},
		},
		{
			input: `hello.world <= "Hello World!"`,
			expected: []Token{
				{Type: Identifier, IdentifierValue: "hello"},
				{Type: Dot},
				{Type: Identifier, IdentifierValue: "world"},
				{Type: Lte},
				{Type: String, StringValue: "Hello World!"},
			},
		},
		{
			input: `!((1+2) != 3.0) is true`,
			expected: []Token{
				{Type: Not},
				{Type: LeftParen},
				{Type: LeftParen},
				{Type: Integer, IntegerValue: 1},
				{Type: Plus},
				{Type: Integer, IntegerValue: 2},
				{Type: RightParen},
				{Type: Neq},
				{Type: Decimal, DecimalValue: 3.0},
				{Type: RightParen},
				{Type: Identifier, IdentifierValue: "is"},
				{Type: True},
			},
		},
	}

	t.Parallel()
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c.input))
			require.NoError(t, err)
			for i, expected := range c.expected {
				compareTokens(t, expected, tokens[i])
			}
		})
	}
}

func TestScannerErrors(t *testing.T) {
	cases := []string{
		"hello | world",
		"hello & world",
		"a = 1",
		`"hello world`,
		"123.4.5.6",
	}
	t.Parallel()
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			_, err := Scan(strings.NewReader(c))
			require.Error(t, err)
		})
	}
}

func compareTokens(t *testing.T, expected, actual Token) {
	t.Helper()

	if !assert.Equal(t, expected.Type, actual.Type, "Token types not equal") {
		return
	}

	switch expected.Type {
	case Integer:
		if !assert.Equal(t, expected.IntegerValue, actual.IntegerValue, "Token values not equal") {
			return
		}
	case Decimal:
		if !assert.Equal(t, expected.DecimalValue, actual.DecimalValue, "Token values not equal") {
			return
		}
	case String:
		if !assert.Equal(t, expected.StringValue, actual.StringValue, "Token values not equal") {
			return
		}
	case Identifier:
		if !assert.Equal(t, expected.IdentifierValue, actual.IdentifierValue, "Token values not equal") {
			return
		}
	}
}

var benchmarkTokens []Token

func BenchmarkScanner(b *testing.B) {
	for i := range 5 {
		count := 16 << (i * 2)
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			input := strings.Repeat(
				`hello.world 12 == 4.0 "Hello World!" != <= true false !true () null `,
				count,
			)

			b.ResetTimer()
			for range b.N {
				benchmarkTokens, _ = Scan(strings.NewReader(input))
			}
		})
	}
}
