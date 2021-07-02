package util

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"net/http"
)

type Trace struct {
	OpenTracer opentracing.Tracer
	tracerCloser io.Closer
	tracerSkipper bool
}

var Tracer Trace

func TracerInit(service string) {
	Tracer = Trace{}
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println(err.Error())
		Tracer.tracerSkipper = true
		return
	}

	cfg.ServiceName = service
	cfg.Sampler.Type = jaeger.SamplerTypeConst
	cfg.Sampler.Param = 1
	cfg.Reporter.LogSpans = true
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	Tracer.OpenTracer, Tracer.tracerCloser, err = cfg.NewTracer(
		config.Logger(jLogger),
		config.Metrics(jMetricsFactory),
	)
	if err != nil {
		Tracer.tracerSkipper = true
		fmt.Printf("ERROR: cannot init Jaeger: %v\n", err)
		return
	}
	opentracing.SetGlobalTracer(Tracer.OpenTracer)
}

// Inject injects the outbound HTTP request with the given span's context to ensure
// correct propagation of span context throughout the trace.
func (tracer *Trace) Inject(span opentracing.Span, request *http.Request) error {
	if tracer.tracerSkipper {
		return nil
	}
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

// Extract extracts the inbound HTTP request to obtain the parent span's context to ensure
// correct propagation of span context throughout the trace.
func (tracer *Trace) Extract(r *http.Request) (opentracing.SpanContext, error) {
	if tracer.tracerSkipper {
		return nil, nil
	}
	return tracer.OpenTracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}

// StartSpanFromRequest extracts the parent span context from the inbound HTTP request
// and starts a new child span if there is a parent span.
func (tracer *Trace) StartSpanFromRequest(spanName string, r *http.Request) opentracing.Span {
	if tracer.tracerSkipper {
		return nil
	}
	spanCtx, _ := tracer.Extract(r)
	return tracer.OpenTracer.StartSpan(spanName, ext.RPCServerOption(spanCtx))
}

func (tracer *Trace) StartSpanFromContext(ctx context.Context, spanName string) opentracing.Span {
	if tracer.tracerSkipper {
		return nil
	}
	span, _ := opentracing.StartSpanFromContext(ctx, spanName)
	return span
}

func (tracer *Trace) ContextWithSpan(ctx context.Context, span opentracing.Span) context.Context {
	if tracer.tracerSkipper {
		return nil
	}
	return opentracing.ContextWithSpan(ctx, span)
}

func (tracer *Trace) LogFields(span opentracing.Span, key string, value string) {
	if tracer.tracerSkipper {
		return
	}
	span.LogFields(logString(key, value))
}

func logString(key string, value string) log.Field {
	return log.String(key, value)
}

func (tracer *Trace) LogError(span opentracing.Span, err error, fields ...log.Field) {
	if tracer.tracerSkipper {
		return
	}
	ext.LogError(span, err, fields...)
}

func (tracer *Trace) CloseTracer() error {
	if tracer.tracerSkipper{
		return nil
	}
	return tracer.tracerCloser.Close()
}

func (tracer *Trace) FinishSpan(span opentracing.Span) {
	if tracer.tracerSkipper{
		return
	}
	span.Finish()
}