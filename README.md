# spans

Provides an easy way to create OpenTelemetry spans with structured attributes.

## Installation

```sh
go get github.com/ebi-yade/spans
```

## Usage

```go
package main

import (
	"github.com/ebi-yade/spans"
	"go.opentelemetry.io/otel"
)

type HTTPContext struct {
	Status     int    `otel:"status_code"` // you can explicitly specify the key suffix
	Method     string `otel:"method"`
	Path       string // => key suffix is "path"
	RemoteAddr string // => key suffix is "remote_addr"

	AuthHeader string `otel:"-"` // you also can ignore the field
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

## Acknowledgements

The prototype of this library was created by @mashiike.
