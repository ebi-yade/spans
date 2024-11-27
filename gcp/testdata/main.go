package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	exporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/cockroachdb/errors"
	"github.com/ebi-yade/spans"
	"github.com/ebi-yade/spans/gcp"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/option"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	if err := main_(); err != nil {
		slog.Error(fmt.Sprintf("error: %v", err))
	}
}

func main_() error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	exp, err := exporter.New(
		exporter.WithTraceClientOptions([]option.ClientOption{option.WithTelemetryDisabled()}), // avoid recursive spans
	)
	if err != nil {
		return errors.Wrap(err, "error NewExporter")
	}

	batchProcessor := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(gcp.NewProcessor(batchProcessor)), // <= IMPORTANT!
		sdktrace.WithSampler(sdktrace.AlwaysSample()),                // not recommended for production
	)
	defer tp.ForceFlush(ctx)
	otel.SetTracerProvider(tp)

	if err := doSomething(context.Background()); err != nil {
		return errors.Wrap(err, "error doSomething")
	}
	return nil
}

var tracer = otel.Tracer("testdata")

type Something struct {
	Bool   bool
	Int    int
	Float  float64
	String string

	SecretString string `otel:"-"`

	IntSlice    []int
	StringSlice []string
}

func newSomething() Something {
	return Something{
		Bool:   true,
		Int:    123,
		Float:  3.14,
		String: "hello",

		SecretString: "secret",

		IntSlice:    []int{1, 2, 3},
		StringSlice: []string{"hello", "world"},
	}
}

func doSomething(ctx context.Context) error {
	sth := newSomething()
	ctx, span := tracer.Start(ctx, "doSomething", spans.WithAttrs(
		spans.ObjectAttr("something", sth),
	))
	defer span.End()

	time.Sleep(10 * time.Millisecond)

	if err := child(ctx); err != nil {
		return errors.Wrap(err, "error child")
	}

	time.Sleep(20 * time.Millisecond)

	return nil
}

func child(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "child")
	defer span.End()

	time.Sleep(50 * time.Millisecond)

	return nil
}
