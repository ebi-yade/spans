# spans

Provides an easy way to create OpenTelemetry spans with structured attributes.

## Usage

```go
package main

import (
	"github.com/ebi-yade/spans"
	"go.opentelemetry.io/otel"
)

type HTTPContext struct {
	Status int    `otel:"status_code"`
	Method string `otel:"method"`
	Path   string `otel:"path"`
}

func main() {
	httpCtx := HTTPContext{
		Status: 200,
		Method: "GET",
		Path:   "/",
	}
	ctx, spans := otel.Tracer("handler").Start("foo", spans.WithAttrs(
		spans.ObjectAttr("http", httpCtx),
	))
}

```

LICENSE: MIT
