package iriser

import (
	"reflect"
	"strings"
	"time"

	"github.com/gxxgle/go-utils/log"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/logger"
)

func Register(r router.Party, h interface{}) {
	prefix := r.GetRelPath()
	hValue := reflect.ValueOf(h)
	for i := 0; i < hValue.NumMethod(); i++ {
		mValue := hValue.Method(i)
		if mValue.Type().NumOut() != 0 {
			continue
		}

		if mValue.Type().NumIn() != 1 {
			continue
		}

		if mValue.Type().In(0).String() != "context.Context" {
			continue
		}

		mType := reflect.TypeOf(h).Method(i)
		r.Post("/"+mType.Name, func(ctx iris.Context) {
			mValue.Call([]reflect.Value{reflect.ValueOf(ctx)})
		})

		log.L.WithField("path", prefix+"/"+mType.Name).Debug("register handler")
	}
}

func NewLogger(ignore func(method, path string) bool) iris.Handler {
	cfg := logger.DefaultConfig()
	cfg.LogFuncCtx = func(ctx iris.Context, latency time.Duration) {
		if ctx.Method() == iris.MethodOptions {
			return
		}
		if ignore != nil && ignore(ctx.Method(), ctx.Path()) {
			return
		}
		fields := log.F{
			"method":     ctx.Method(),
			"path":       ctx.Path(),
			"latency_ms": latency.Milliseconds(),
		}
		ctx.Values().Visit(func(key string, value interface{}) {
			if strings.HasPrefix(key, "log.") {
				fields[strings.TrimPrefix(key, "log.")] = value
			}
		})
		log.L.WithFields(fields).Info("api request")
	}
	return logger.New(cfg)
}
