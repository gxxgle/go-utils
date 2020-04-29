package iriser

import (
	"reflect"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func Register(r router.Party, h interface{}) {
	hValue := reflect.ValueOf(h)
	for i := 0; i < hValue.NumMethod(); i++ {
		mValue := hValue.Method(i)
		if mValue.Type().NumOut() != 0 {
			continue
		}

		if mValue.Type().NumIn() != 1 {
			continue
		}

		if !strings.Contains(mValue.Type().In(0).String(), ".Context") {
			continue
		}

		mType := reflect.TypeOf(h).Method(i)
		r.Post("/"+mType.Name, func(ctx iris.Context) {
			mValue.Call([]reflect.Value{reflect.ValueOf(ctx)})
		})
	}
}
