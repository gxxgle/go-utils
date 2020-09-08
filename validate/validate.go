package validate

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gxxgle/go-utils/conver"
)

var (
	V = validator.New()
)

func init() {
	_ = V.RegisterValidation("default", tagDefault)
}

func tagDefault(fl validator.FieldLevel) bool {
	var (
		param = fl.Param()
		value = fl.Field()
	)

	switch value.Kind() {
	case reflect.String:
		if len(value.Interface().(string)) == 0 {
			value.SetString(param)
		}
	case reflect.Int64, reflect.Int32, reflect.Int:
		if value.Int() == 0 {
			value.SetInt(conver.Int64Must(param))
		}
	case reflect.Uint64, reflect.Uint32, reflect.Uint:
		if value.Uint() == 0 {
			value.SetUint(uint64(conver.Int64Must(param)))
		}
	case reflect.Float64, reflect.Float32:
		if value.Float() == 0 {
			value.SetFloat(conver.Float64Must(param))
		}
	}

	return true
}
