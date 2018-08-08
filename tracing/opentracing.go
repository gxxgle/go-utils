package tracing

import (
	"context"
	"time"

	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	tracer opentracing.Tracer
)

// Init creates a new instance of Jaeger tracer
func Init(service, url string) error {
	if service == "" || url == "" {
		return nil
	}

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	sender, err := jaeger.NewUDPTransport(url, 0)
	if err != nil {
		return err
	}

	tracer, _, err = cfg.New(
		service,
		config.Reporter(jaeger.NewRemoteReporter(
			sender,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
		)),
	)
	return err
}

// SpanFromContext get tracing span from context
func SpanFromContext(ctx context.Context, method string) (context.Context, opentracing.Span, error) {
	if tracer == nil {
		return ctx, nil, nil
	}

	var (
		span             opentracing.Span
		md, ok           = metadata.FromContext(ctx)
		wireContext, err = tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	)

	if !ok {
		md = make(map[string]string)
	}

	if err != nil {
		span = tracer.StartSpan(method)
	} else {
		span = tracer.StartSpan(method, opentracing.ChildOf(wireContext))
	}

	err = span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	if err != nil {
		return nil, nil, err
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)
	return ctx, span, nil
}
