package conver

import (
	"fmt"
	"strconv"
	"strings"
)

func String(val interface{}) (string, error) {
	switch result := val.(type) {
	case string:
		return result, nil
	case []byte:
		return string(result), nil
	case fmt.Stringer:
		return result.String(), nil
	case error:
		return result.Error(), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", result), nil
	default:
		return fmt.Sprintf("%+v", val), nil
	}
}

func StringMust(val interface{}, def ...string) string {
	result, err := String(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Bool(val interface{}) (bool, error) {
	if val == nil {
		return false, nil
	}

	switch result := val.(type) {
	case bool:
		return result, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return result != 0, nil
	case []byte:
		return stringToBool(string(result))
	case string:
		return stringToBool(result)
	default:
		return false, converError(val, "bool")
	}
}

func BoolMust(val interface{}, def ...bool) bool {
	result, err := Bool(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Bytes(val interface{}) ([]byte, error) {
	switch result := val.(type) {
	case []byte:
		return result, nil
	default:
		str, err := String(val)
		return []byte(str), err
	}
}

func BytesMust(val interface{}, def ...[]byte) []byte {
	result, err := Bytes(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Float64(val interface{}) (float64, error) {
	switch result := val.(type) {
	case float64:
		return result, nil
	case int:
		return float64(result), nil
	case int8:
		return float64(result), nil
	case int16:
		return float64(result), nil
	case int32:
		return float64(result), nil
	case int64:
		return float64(result), nil
	case uint:
		return float64(result), nil
	case uint8:
		return float64(result), nil
	case uint16:
		return float64(result), nil
	case uint32:
		return float64(result), nil
	case uint64:
		return float64(result), nil
	case float32:
		return float64(result), nil
	case bool:
		if result {
			return 1, nil
		}
		return 0, nil
	default:
		str := strings.ReplaceAll(strings.TrimSpace(StringMust(val)), " ", "")
		return strconv.ParseFloat(str, 64)
	}
}

func Float64Must(val interface{}, def ...float64) float64 {
	result, err := Float64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Float32(val interface{}) (float32, error) {
	result, err := Float64(val)
	return float32(result), err
}

func Float32Must(val interface{}, def ...float32) float32 {
	result, err := Float32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Int64(val interface{}) (int64, error) {
	switch result := val.(type) {
	case int64:
		return result, nil
	case int:
		return int64(result), nil
	case int8:
		return int64(result), nil
	case int16:
		return int64(result), nil
	case int32:
		return int64(result), nil
	case uint:
		return int64(result), nil
	case uint8:
		return int64(result), nil
	case uint16:
		return int64(result), nil
	case uint32:
		return int64(result), nil
	case uint64:
		return int64(result), nil
	case bool:
		if result {
			return 1, nil
		}
		return 0, nil
	default:
		fval, err := Float64(val)
		if err != nil {
			return 0, err
		}

		return int64(fval), nil
	}
}

func Int64Must(val interface{}, def ...int64) int64 {
	result, err := Int64(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Int(val interface{}) (int, error) {
	i, err := Int64(val)
	return int(i), err
}

func IntMust(val interface{}, def ...int) int {
	result, err := Int(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func Int32(val interface{}) (int32, error) {
	i, err := Int64(val)
	return int32(i), err
}

func Int32Must(val interface{}, def ...int32) int32 {
	result, err := Int32(val)
	if err != nil && len(def) > 0 {
		return def[0]
	}
	return result
}

func converError(val interface{}, tpy string) error {
	return fmt.Errorf("the %T{%v} can not conver to a %v", val, val, tpy)
}

func stringToBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "TRUE", "True", "ok", "OK", "yes", "YES":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "":
		return false, nil
	}
	return false, converError(val, "bool")
}
