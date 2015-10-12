package ql

import (
	"fmt"
	"math"
	"strconv"
)

const (
	MAX_INT = 9007199254740991
	MIN_INT = -9007199254740991
)

func coerceInt(v interface{}) interface{} {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return uint64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 0, 64)
		return i
	case bool:
		if v {
			return int64(1)
		}
		return int64(0)
	default:
		return int64(0)
	}
}

func coerceFloat(v interface{}) interface{} {
	switch v := v.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case bool:
		if v {
			return float64(1)
		}
		return float64(0)
	default:
		return float64(0)
	}
}

func coerceString(v interface{}) interface{} {
	return fmt.Sprint(v)
}

func coerceBool(v interface{}) interface{} {
	switch v := v.(type) {
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0 && !math.IsNaN(float64(v))
	case float64:
		return v != 0 && !math.IsNaN(v)
	case string:
		b, _ := strconv.ParseBool(v)
		return b
	case bool:
		return v
	default:
		return false
	}
}

var Int = Scalar{
	Name:       "Int",
	Serialize:  coerceInt,
	ParseValue: coerceInt,
	ParseLiteral: func(kind, value string) interface{} {
		if kind != INT {
			return nil
		}
		num, _ := strconv.ParseInt(value, 10, 64)
		return num
	},
}

var Float = Scalar{
	Name:       "Float",
	Serialize:  coerceFloat,
	ParseValue: coerceFloat,
	ParseLiteral: func(kind, value string) interface{} {
		if kind != FLOAT && kind != INT {
			return nil
		}
		num, _ := strconv.ParseFloat(value, 64)
		return num
	},
}

var String = Scalar{
	Name:       "String",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(kind, value string) interface{} {
		if kind != STRING {
			return nil
		}
		return value
	},
}

var Boolean = Scalar{
	Name:       "Boolean",
	Serialize:  coerceBool,
	ParseValue: coerceBool,
	ParseLiteral: func(kind, value string) interface{} {
		if kind != BOOLEAN {
			return nil
		}
		num, _ := strconv.ParseBool(value)
		return num
	},
}

var ID = Scalar{
	Name:       "ID",
	Serialize:  coerceString,
	ParseValue: coerceString,
	ParseLiteral: func(kind, value string) interface{} {
		if kind != STRING || kind != INT {
			return nil
		}
		return value
	},
}
