package gin

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gxxgle/go-utils/json"
	"github.com/gxxgle/go-utils/log"
	"github.com/gxxgle/go-utils/tracing"
	"github.com/opentracing/opentracing-go"
)

func NewHandlerWrap() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err    error
			req    = c.Request
			start  = time.Now()
			ctx    = req.Context()
			msg    = "api access"
			span   opentracing.Span
			method = strings.TrimPrefix(req.URL.Path, "/")
		)

		kv := []interface{}{
			"method", method,
		}

		defer func() {
			if err != nil {
				msg = "api error"
				kv = append(kv, "err", err)
				c.String(http.StatusOK, err.Error())
				c.Abort()
			}

			log.Infow(msg, kv...)
		}()

		ctx, span, err = tracing.SpanFromContext(ctx, method)
		if err != nil {
			return
		}

		if span != nil {
			defer span.Finish()
		}

		if ctx != nil {
			*req = *req.WithContext(ctx)
		}

		c.Next()

		kv = append(kv,
			"cost", time.Since(start),
			"req_query", json.MustMarshalToString(req.URL.Query()),
			"req_params", json.MustMarshalToString(c.Params),
		)
	}
}
