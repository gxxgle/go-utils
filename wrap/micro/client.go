package micro

import (
	"context"

	"github.com/gxxgle/go-utils/errors"
	"github.com/gxxgle/go-utils/tracing"
	"github.com/micro/go-micro/client"
)

type ClientWrap struct {
	client.Client
}

func (w *ClientWrap) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx, span, err := tracing.SpanFromContext(ctx, req.Method())
	if err != nil {
		return err
	}

	if span != nil {
		defer span.Finish()
	}

	err = w.Client.Call(ctx, req, rsp, opts...)
	if err == nil {
		return nil
	}

	e := errors.Parse(err)
	er, ok := errors.Errors[e.Code]
	if ok {
		return er
	}

	return err
}

func NewClientWrap(c client.Client) client.Client {
	return &ClientWrap{c}
}
