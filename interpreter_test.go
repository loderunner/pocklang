package pock

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
)

func TestInterpreter(t *testing.T) {
	type testCase struct {
		state    map[string]any
		input    string
		expected any
	}
	cases := []testCase{
		{input: "true", expected: true},
		{input: "false", expected: false},
		{input: "null", expected: null},
		{input: "3", expected: 3},
		{input: "3.14", expected: 3.14},
		{input: `"Hello World!"`, expected: "Hello World!"},
		{input: "1138 <= 1138", expected: true},
		{input: "1138.0 <= 1138", expected: true},
		{input: "1138 <= 1138.0", expected: true},
		{input: "1138.0 <= 1138.0", expected: true},
		{input: "1138 <= 9999999", expected: true},
		{input: "1138.0 <= 9999999", expected: true},
		{input: "1138 <= 9999999.0", expected: true},
		{input: "1138.0 <= 9999999.0", expected: true},
		{input: "1138 <= 0", expected: false},
		{input: "1138.0 <= 0", expected: false},
		{input: "1138 <= 0.0", expected: false},
		{input: "1138.0 <= 0.0", expected: false},
		{input: "1138 < 1138", expected: false},
		{input: "1138.0 < 1138", expected: false},
		{input: "1138 < 1138.0", expected: false},
		{input: "1138.0 < 1138.0", expected: false},
		{input: "1138 < 9999999", expected: true},
		{input: "1138.0 < 9999999", expected: true},
		{input: "1138 < 9999999.0", expected: true},
		{input: "1138.0 < 9999999.0", expected: true},
		{input: "1138 < 0", expected: false},
		{input: "1138.0 < 0", expected: false},
		{input: "1138 < 0.0", expected: false},
		{input: "1138.0 < 0.0", expected: false},
		{input: "1138 >= 1138", expected: true},
		{input: "1138.0 >= 1138", expected: true},
		{input: "1138 >= 1138.0", expected: true},
		{input: "1138.0 >= 1138.0", expected: true},
		{input: "1138 >= 9999999", expected: false},
		{input: "1138.0 >= 9999999", expected: false},
		{input: "1138 >= 9999999.0", expected: false},
		{input: "1138.0 >= 9999999.0", expected: false},
		{input: "1138 >= 0", expected: true},
		{input: "1138.0 >= 0", expected: true},
		{input: "1138 >= 0.0", expected: true},
		{input: "1138.0 >= 0.0", expected: true},
		{input: "1138 > 1138", expected: false},
		{input: "1138.0 > 1138", expected: false},
		{input: "1138 > 1138.0", expected: false},
		{input: "1138.0 > 1138.0", expected: false},
		{input: "1138 > 9999999", expected: false},
		{input: "1138.0 > 9999999", expected: false},
		{input: "1138 > 9999999.0", expected: false},
		{input: "1138.0 > 9999999.0", expected: false},
		{input: "1138 > 0", expected: true},
		{input: "1138.0 > 0", expected: true},
		{input: "1138 > 0.0", expected: true},
		{input: "1138.0 > 0.0", expected: true},
		{input: "1138 == 1138", expected: true},
		{input: "1138.0 == 1138", expected: true},
		{input: "1138 == 1138.0", expected: true},
		{input: "1138.0 == 1138.0", expected: true},
		{input: "1138 == 0", expected: false},
		{input: "1138.0 == 0", expected: false},
		{input: "1138 == 0.0", expected: false},
		{input: "1138.0 == 0.0", expected: false},
		{input: `"hello" == "hello"`, expected: true},
		{input: `"hello" == "world"`, expected: false},
		{input: "true == true", expected: true},
		{input: "true == false", expected: false},
		{input: "false == true", expected: false},
		{input: "false == false", expected: true},
		{input: "null == null", expected: true},
		{input: "1138 != 1138", expected: false},
		{input: "1138.0 != 1138", expected: false},
		{input: "1138 != 1138.0", expected: false},
		{input: "1138.0 != 1138.0", expected: false},
		{input: "1138 != 0", expected: true},
		{input: "1138.0 != 0", expected: true},
		{input: "1138 != 0.0", expected: true},
		{input: "1138.0 != 0.0", expected: true},
		{input: `"hello" != "hello"`, expected: false},
		{input: `"hello" != "world"`, expected: true},
		{input: "true != true", expected: false},
		{input: "true != false", expected: true},
		{input: "false != true", expected: true},
		{input: "false != false", expected: false},
		{input: "null != null", expected: false},
		{input: "-1138", expected: -1138},
		{input: "-1138.0", expected: -1138.0},
		{input: "!true", expected: false},
		{input: "!false", expected: true},
		{input: "1138 + 10", expected: 1148},
		{input: "1138.0 + 10", expected: 1148.0},
		{input: "1138 + 10.0", expected: 1148.0},
		{input: "1138.0 + 10.0", expected: 1148.0},
		{input: "1138 - 10", expected: 1128},
		{input: "1138.0 - 10", expected: 1128.0},
		{input: "1138 - 10.0", expected: 1128.0},
		{input: "1138.0 - 10.0", expected: 1128.0},
		{input: "1138 * 10", expected: 11380},
		{input: "1138.0 * 10", expected: 11380.0},
		{input: "1138 * 10.0", expected: 11380.0},
		{input: "1138.0 * 10.0", expected: 11380.0},
		{input: "1138 / 10", expected: 113},
		{input: "1138.0 / 10", expected: 113.8},
		{input: "1138 / 10.0", expected: 113.8},
		{input: "1138.0 / 10.0", expected: 113.8},
		{input: "false && false", expected: false},
		{input: "false && true", expected: false},
		{input: "true && false", expected: false},
		{input: "true && true", expected: true},
		{input: "false || false", expected: false},
		{input: "false || true", expected: true},
		{input: "true || false", expected: true},
		{input: "true || true", expected: true},
		{input: "(1 == 2) == false", expected: true},
		{
			state:    map[string]any{"hello": "world"},
			input:    "hello",
			expected: "world",
		},
		{
			state:    map[string]any{"hello": map[string]any{"world": 1138}},
			input:    "hello.world",
			expected: 1138,
		},
		{
			state: map[string]any{
				"hello": map[string]any{
					"world": "Hello World!",
				},
				"THX": 1138,
			},
			input:    `(THX - 1 == 2) || (hello.world == "Hello World!")`,
			expected: true,
		},
	}

	t.Parallel()
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c.input))
			require.NoError(t, err)
			expr, err := Parse(tokens)
			require.NoError(t, err)
			var i *Interpreter
			if c.state == nil {
				i = NewInterpreter()
			} else {
				i, err = NewInterpreterWithState(c.state)
				require.NoError(t, err)
			}
			val, err := i.Evaluate(expr)
			require.NoError(t, err)
			require.EqualValues(t, c.expected, val)
		})
	}
}

func TestInterpreterError(t *testing.T) {
	type testCase struct {
		state map[string]any
		input string
	}
	cases := []testCase{
		{input: "1 || 0"},
		{input: `1 && "hello"`},
		{input: `"hello" <3`},
		{input: `"hello" <=3`},
		{input: `"hello" > 3`},
		{input: `"hello" >= 3`},
		{input: `1 == "hello"`},
		{input: `1 != "hello"`},
		{input: `"hello" + "world"`},
		{input: `"hello" - "world"`},
		{input: `"hello" * "world"`},
		{input: `"hello" / "world"`},
		{input: "!1"},
		{input: "-true"},
		{input: "-true"},
		{input: "hello"},
		{state: map[string]any{"hello": true}, input: "world"},
		{state: map[string]any{"hello": true}, input: "hello.world"},
		{state: map[string]any{"hello": map[string]any{}}, input: "hello"},
		{state: map[string]any{"hello": map[string]any{}}, input: "hello.world"},
	}

	t.Parallel()
	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			tokens, err := Scan(strings.NewReader(c.input))
			require.NoError(t, err)
			expr, err := Parse(tokens)
			require.NoError(t, err)
			var i *Interpreter
			if c.state == nil {
				i = NewInterpreter()
			} else {
				i, err = NewInterpreterWithState(c.state)
				require.NoError(t, err)
			}
			_, err = i.Evaluate(expr)
			require.Error(t, err)
			snaps.MatchSnapshot(t, err)
		})
	}
}

var benchmarkValue Value

func BenchmarkInterpreter(b *testing.B) {
	for i := range 6 {
		count := 1 << (i * 2)
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			input := strings.Join(
				slices.Repeat(
					[]string{"(hello.world + 3 == 0) || (1.0 + 1 == 2.0)"},
					count,
				),
				" && ",
			)

			b.ResetTimer()
			for range b.N {
				benchmarkTokens, _ = Scan(strings.NewReader(input))
				benchmarkExpr, _ = Parse(benchmarkTokens)
				s, _ := NewInterpreterWithState(map[string]any{"hello": map[string]any{"world": 1138}})
				benchmarkValue, _ = s.Evaluate(benchmarkExpr)
			}
		})
	}
}
