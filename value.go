package pock

type Value interface {
	GetInteger() (int64, bool)
	GetDecimal() (float64, bool)
	GetString() (string, bool)
	GetBool() (bool, bool)
	GetNull() (interface{}, bool)
}

type IntValue int64

func (v IntValue) GetInteger() (int64, bool) {
	return int64(v), true
}

func (v IntValue) GetDecimal() (float64, bool) {
	return 0.0, false
}

func (v IntValue) GetString() (string, bool) {
	return "", false
}

func (v IntValue) GetBool() (bool, bool) {
	return false, false
}

func (v IntValue) GetNull() (interface{}, bool) {
	return nil, false
}

type DecimalValue float64

func (v DecimalValue) GetInteger() (int64, bool) {
	return 0, false
}

func (v DecimalValue) GetDecimal() (float64, bool) {
	return float64(v), true
}

func (v DecimalValue) GetString() (string, bool) {
	return "", false
}

func (v DecimalValue) GetBool() (bool, bool) {
	return false, false
}

func (v DecimalValue) GetNull() (interface{}, bool) {
	return nil, false
}

type StringValue string

func (v StringValue) GetInteger() (int64, bool) {
	return 0, false
}

func (v StringValue) GetDecimal() (float64, bool) {
	return 0.0, false
}

func (v StringValue) GetString() (string, bool) {
	return string(v), true
}

func (v StringValue) GetBool() (bool, bool) {
	return false, false
}

func (v StringValue) GetNull() (interface{}, bool) {
	return nil, false
}

type BoolValue bool

func (v BoolValue) GetInteger() (int64, bool) {
	return 0, false
}

func (v BoolValue) GetDecimal() (float64, bool) {
	return 0.0, false
}

func (v BoolValue) GetString() (string, bool) {
	return "", false
}

func (v BoolValue) GetBool() (bool, bool) {
	return bool(v), true
}

func (v BoolValue) GetNull() (interface{}, bool) {
	return nil, false
}

type NullValue struct{}

func (v NullValue) GetInteger() (int64, bool) {
	return 0, false
}

func (v NullValue) GetDecimal() (float64, bool) {
	return 0.0, false
}

func (v NullValue) GetString() (string, bool) {
	return "", false
}

func (v NullValue) GetBool() (bool, bool) {
	return false, false
}

func (v NullValue) GetNull() (interface{}, bool) {
	return nil, true
}

var null = NullValue{}
