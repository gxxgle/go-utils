package iriser

import (
	"reflect"
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

func NewLogger() iris.Handler {
	cfg := logger.DefaultConfig()
	cfg.LogFuncCtx = func(ctx iris.Context, latency time.Duration) {
		log.L.WithFields(log.F{
			"method":     ctx.Method(),
			"path":       ctx.Path(),
			"latency_ms": latency.Milliseconds(),
		}).Info("api request")
	}
	return logger.New(cfg)
}
