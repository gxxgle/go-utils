package fn

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func Name(v interface{}) (string, error) {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Func {
		return "", fmt.Errorf("%T is not func type", v)
	}

	name := runtime.FuncForPC(value.Pointer()).Name()
	if i := strings.LastIndex(name, "."); i != -1 {
		return strings.TrimSuffix(name[i+1:], "-fm"), nil
	}

	return "", fmt.Errorf("not found func name")
}
