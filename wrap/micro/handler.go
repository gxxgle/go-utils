package micro

import (
	"context"
	"time"

	"github.com/gxxgle/go-utils/errors"
	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/log"
	"github.com/gxxgle/go-utils/tracing"
	"github.com/micro/go-micro/server"
)

type HandlerWrapConfig struct {
	IsLogRequest  bool
	IsLogResponse bool
}

func HandlerWrap(config map[string]HandlerWrapConfig) server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			method := req.Method()
			cf, ok := config[method]
			if !ok {
				log.Errorw("handler wrap no method found", "method", method)
				return errors.Internal("method config not found")
			}

			ctx, span, err := tracing.SpanFromContext(ctx, method)
			if err != nil {
				return err
			}

			if span != nil {
				defer span.Finish()
			}

			start := time.Now()
			err = fn(ctx, req, rsp)
			kv := []interface{}{
				"method", method,
				"cost", time.Since(start),
			}

			if cf.IsLogRequest {
				kv = append(kv, "req", json.MustMarshalToString(req.Body()))
			}

			if cf.IsLogResponse {
				kv = append(kv, "resp", json.MustMarshalToString(rsp))
			}

			kv = append(kv)

			msg := "api access"
			if err != nil {
				msg = "api error"
				kv = append(kv, "err", err)
			}

			log.Infow(msg, kv...)
			return err
		}
	}
}
