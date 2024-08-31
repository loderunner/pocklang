package pock

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	type testCase struct {
		input    string
		expected Expr
	}
	cases := []testCase{
		{
			input: "1 == 2",
			expected: BinaryExpr{
				Op:    Eq,
				Left:  LiteralExpr{Token: Token{Type: Integer, Lexeme: "1", IntegerValue: 1}},
				Right: LiteralExpr{Token: Token{Type: Integer, Lexeme: "2", IntegerValue: 2}},
			},
		},
	}

	t.Parallel()
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c.input))
			require.NoError(t, err)
			expr, err := Parse(tokens)
			require.NoError(t, err)
			require.Equal(t, c.expected, expr)
		})
	}
}

func TestParserSnapshots(t *testing.T) {
	cases := []string{
		"hello.world > 3",
		`"hello" != "world"`,
		"((3+2) - 14) == -19",
		`123.45 * "d" < asdrg`,
		"true && false || null == (42 / 2)",
	}
	t.Parallel()
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c))
			require.NoError(t, err)
			expr, err := Parse(tokens)
			require.NoError(t, err)
			snaps.MatchSnapshot(t, expr)
		})
	}
}

func TestParserErrors(t *testing.T) {
	cases := []string{
		"",
		"hello.",
		".hello",
		"12 <",
		"12.hello",
		"4 << 54",
		"(41 + d",
		`(""+)`,
		"--3",
		"3*",
		"true && || false",
		"true || && false",
	}
	t.Parallel()
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c))
			require.NoError(t, err)
			_, err = Parse(tokens)
			require.Error(t, err)
			snaps.MatchSnapshot(t, err)
		})
	}
}

var benchmarkExpr Expr

func BenchmarkParser(b *testing.B) {
	for i := range 6 {
		count := 1 << (i * 2)
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			input := strings.Join(slices.Repeat(
				[]string{`(!((hello.world + 3.0) == (true && false) || "hello"))`},
				count,
			),
				" && ",
			)

			b.ResetTimer()
			for range b.N {
				benchmarkTokens, _ = Scan(strings.NewReader(input))
				benchmarkExpr, _ = Parse(benchmarkTokens)
			}
		})
	}
}
