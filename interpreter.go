package pock

import (
	"fmt"
)

type (
	Value        interface{}
	IntValue     int64
	DecimalValue float64
	StringValue  string
	BoolValue    bool
	NullValue    struct{}
)

var nullV = NullValue{}

type Interpreter struct {
	variables map[string]any
}

func NewInterpreter() *Interpreter {
	return &Interpreter{variables: map[string]any{}}
}

func NewInterpreterWithState(state map[string]any) (*Interpreter, error) {
	i := NewInterpreter()
	err := i.LoadState(state)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Interpreter) LoadState(state map[string]any) error {
	return loadState(s.variables, state)
}

func loadState(base, state map[string]any) error {
	for k, v := range state {
		switch v := v.(type) {
		case int64, float64, string, bool, nil:
			base[k] = v
		case int:
			base[k] = int64(v)
		case int8:
			base[k] = int64(v)
		case int16:
			base[k] = int64(v)
		case int32:
			base[k] = int64(v)
		case uint:
			base[k] = int64(v)
		case uint8:
			base[k] = int64(v)
		case uint16:
			base[k] = int64(v)
		case uint32:
			base[k] = int64(v)
		case uint64:
			base[k] = int64(v)
		case float32:
			base[k] = float64(v)
		case map[string]any:
			base[k] = map[string]any{}
			loadState(base[k].(map[string]any), v)
		default:
			return fmt.Errorf("invalid type: %T", v)
		}
	}
	return nil
}

func (s *Interpreter) LoadInt(name string, value int64) {
	s.variables[name] = value
}

func (s *Interpreter) LoadDecimal(name string, value float64) {
	s.variables[name] = value
}

func (s *Interpreter) LoadBool(name string, value bool) {
	s.variables[name] = value
}

func (s *Interpreter) LoadNull(name string) {
	s.variables[name] = nil
}

func (s *Interpreter) LoadMap(name string, value map[string]any) error {
	s.variables[name] = map[string]any{}
	return loadState(s.variables[name].(map[string]any), value)
}

func (s Interpreter) Evaluate(expr Expr) (Value, error) {
	switch expr := expr.(type) {
	case BinaryExpr:
		return s.evaluateBinary(expr)
	case UnaryExpr:
		return s.evaluateUnary(expr)
	case GroupExpr:
		return s.evaluateGroup(expr)
	case GetExpr:
		return s.evaluateGet(expr)
	case LiteralExpr:
		return s.evaluateLiteral(expr)
	}
	panic("invalid expression")
}

func (s Interpreter) evaluateBinary(expr BinaryExpr) (Value, error) {
	left, err := s.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := s.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case Or:
		left, right, ok := checkBinary[bool, bool](left, right)
		if !ok {
			return nil, fmt.Errorf("`||` operands must be boolean")
		}
		return left || right, nil
	case And:
		left, right, ok := checkBinary[bool, bool](left, right)
		if !ok {
			return nil, fmt.Errorf("`&&` operands must be boolean")
		}
		return left && right, nil
	case Lt:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left < right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) < right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left < float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left < right, nil
		}
		return nil, fmt.Errorf("`<` operands must be integer or decimal")
	case Lte:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left <= right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) <= right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left <= float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left <= right, nil
		}
		return nil, fmt.Errorf("`<=` operands must be integer or decimal")
	case Gt:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left > right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) > right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left > float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left > right, nil
		}
		return nil, fmt.Errorf("`>` operands must be integer or decimal")
	case Gte:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left >= right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) >= right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left >= float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left >= right, nil
		}
		return nil, fmt.Errorf("`>=` operands must be integer or decimal")
	case Eq:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left == right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) == right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left == float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left == right, nil
		}
		if left, right, ok := checkBinary[bool, bool](left, right); ok {
			return left == right, nil
		}
		if left, right, ok := checkBinary[string, string](left, right); ok {
			return left == right, nil
		}
		if left, right, ok := checkBinary[NullValue, NullValue](left, right); ok {
			return left == right, nil
		}
		return nil, fmt.Errorf(
			"`==` operands mismatch: %s and %s",
			typeName(left),
			typeName(right),
		)
	case Neq:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left != right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) != right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left != float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left != right, nil
		}
		if left, right, ok := checkBinary[bool, bool](left, right); ok {
			return left != right, nil
		}
		if left, right, ok := checkBinary[string, string](left, right); ok {
			return left != right, nil
		}
		if left, right, ok := checkBinary[NullValue, NullValue](left, right); ok {
			return left != right, nil
		}
		return nil, fmt.Errorf(
			"`!=` operands mismatch: %s and %s",
			typeName(left),
			typeName(right),
		)
	case Plus:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left + right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) + right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left + float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left + right, nil
		}
		return nil, fmt.Errorf("`+` operands must be integer or decimal")
	case Minus:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left - right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) - right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left - float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left - right, nil
		}
		return nil, fmt.Errorf("`-` operands must be integer or decimal")
	case Star:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left * right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) * right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left * float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left * right, nil
		}
		return nil, fmt.Errorf("`*` operands must be integer or decimal")
	case Slash:
		if left, right, ok := checkBinary[int64, int64](left, right); ok {
			return left / right, nil
		}
		if left, right, ok := checkBinary[int64, float64](left, right); ok {
			return float64(left) / right, nil
		}
		if left, right, ok := checkBinary[float64, int64](left, right); ok {
			return left / float64(right), nil
		}
		if left, right, ok := checkBinary[float64, float64](left, right); ok {
			return left / right, nil
		}
		return nil, fmt.Errorf("`/` operands must be integer or decimal")
	}
	panic(fmt.Sprintf("invalid binary operator: %s", expr.Op))
}

func (s Interpreter) evaluateUnary(expr UnaryExpr) (Value, error) {
	val, err := s.Evaluate(expr.Expr)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case Not:
		if val, ok := val.(bool); ok {
			return !val, nil
		}
		return nil, fmt.Errorf("`!` operand must be boolean")
	case Minus:
		switch val := val.(type) {
		case int64:
			return -val, nil
		case float64:
			return -val, nil
		}
		return nil, fmt.Errorf("`-` operand must be integer or decimal")
	}
	panic(fmt.Sprintf("invalid unary operator: %s", expr.Op))
}

func (s Interpreter) evaluateGroup(expr GroupExpr) (Value, error) {
	return s.Evaluate(expr.Expr)
}

func (s Interpreter) evaluateGet(expr GetExpr) (Value, error) {
	if len(expr.Names) < 0 {
		panic("empty get expression")
	}

	name := expr.Names[0]
	val, ok := s.variables[name]
	if !ok {
		return nil, fmt.Errorf("unknown variable '%s'", name)
	}
	for i := 1; i < len(expr.Names); i++ {
		var obj map[string]any
		name = expr.Names[i]
		if obj, ok = val.(map[string]any); !ok {
			return nil, fmt.Errorf("%s is not a map", name)
		}
		val, ok = obj[name]
		if !ok {
			return nil, fmt.Errorf("unknown key '%s'", name)
		}
	}

	if _, ok := val.(map[string]any); ok {
		return nil, fmt.Errorf("%s is not a primitive value", name)
	}

	return val, nil
}

func (s Interpreter) evaluateLiteral(expr LiteralExpr) (Value, error) {
	switch expr.Token.Type {
	case True:
		return true, nil
	case False:
		return false, nil
	case Null:
		return nullV, nil
	case Integer:
		return expr.Token.IntegerValue, nil
	case Decimal:
		return expr.Token.DecimalValue, nil
	case String:
		return expr.Token.StringValue, nil
	}
	panic(
		fmt.Sprintf(
			"invalid literal expression: Token{%s, %s}",
			expr.Token.Type,
			expr.Token.Lexeme,
		),
	)
}

func checkBinary[L, R interface{}](left, right any) (L, R, bool) {
	var zeroL L
	var zeroR R
	leftL, ok := left.(L)
	if !ok {
		return zeroL, zeroR, false
	}
	rightR, ok := right.(R)
	if !ok {
		return zeroL, zeroR, false
	}
	return leftL, rightR, true
}

func typeName(v any) string {
	switch v.(type) {
	case bool:
		return "boolean"
	case int64:
		return "integer"
	case float64:
		return "decimal"
	case string:
		return "string"
	case NullValue:
		return "null"
	case map[string]any:
		return "map"
	}
	panic(fmt.Sprintf("invalid type: %T", v))
}
