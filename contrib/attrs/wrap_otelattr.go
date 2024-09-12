// Package attrs extends "go.opentelemetry.io/otel/attribute" to provide a user-friendly API for setting structured attributes on spans.
package attrs

import (
	"fmt"
	"log/slog"

	"github.com/mashiike/otelattr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type KeyValue struct {
	key   attribute.Key
	value any
}

func newKeyValue(k string, v any) KeyValue {
	return KeyValue{key: attribute.Key(k), value: v}
}

func getStandardAttributes(attrs []KeyValue) []attribute.KeyValue {
	attributes := make([]attribute.KeyValue, 0, len(attrs))
	for _, attr := range attrs {
		switch v := attr.value.(type) {
		case attribute.Value:
			attributes = append(attributes, attribute.KeyValue{Key: attr.key, Value: v})
		default:
			res, err := otelattr.MarshalOtelAttributes(v)
			if err != nil {
				slog.Error(fmt.Sprintf("error MarshalOtelAttributes: key=>%s, type=>%T, error=>%v", attr.key, v, err))
			}
			attributes = append(attributes, res...)
		}
	}

	return attributes
}

// OnStart can be used in place of trace.WithAttributes to set multiple attributes on a span at the time of creation.
func OnStart(attrs ...KeyValue) trace.SpanStartEventOption {
	attributes := getStandardAttributes(attrs)
	return trace.WithAttributes(attributes...)
}

// Set can be used in place of span.SetAttributes to set multiple attributes on a span after it has been created.
func Set(span trace.Span, attrs ...KeyValue) {
	attributes := getStandardAttributes(attrs)
	span.SetAttributes(attributes...)
}

func Any(k string, v any) KeyValue {
	if s, ok := v.(fmt.Stringer); ok {
		return Stringer(k, s)
	}

	switch tv := v.(type) {
	case bool:
		return Bool(k, tv)
	case []bool:
		return BoolSlice(k, tv)
	case int:
		return Int(k, tv)
	case []int:
		return IntSlice(k, tv)
	case int64:
		return Int64(k, tv)
	case []int64:
		return Int64Slice(k, tv)
	case float64:
		return Float64(k, tv)
	case []float64:
		return Float64Slice(k, tv)
	case string:
		return String(k, tv)
	case []string:
		return StringSlice(k, tv)
	default:
		return newKeyValue(k, tv)
	}
}

// Bool is comparable to attribute.Bool.
func Bool(k string, v bool) KeyValue {
	return newKeyValue(k, attribute.BoolValue(v))
}

// BoolSlice is comparable to attribute.BoolSlice.
func BoolSlice(k string, v []bool) KeyValue {
	return newKeyValue(k, attribute.BoolSliceValue(v))
}

// Int is comparable to attribute.Int.
func Int(k string, v int) KeyValue {
	return newKeyValue(k, attribute.IntValue(v))
}

// IntSlice is comparable to attribute.IntSlice.
func IntSlice(k string, v []int) KeyValue {
	return newKeyValue(k, attribute.IntSliceValue(v))
}

// Int64 is comparable to attribute.Int64.
func Int64(k string, v int64) KeyValue {
	return newKeyValue(k, attribute.Int64Value(v))
}

// Int64Slice is comparable to attribute.Int64Slice.
func Int64Slice(k string, v []int64) KeyValue {
	return newKeyValue(k, attribute.Int64SliceValue(v))
}

// Float64 is comparable to attribute.Float64.
func Float64(k string, v float64) KeyValue {
	return newKeyValue(k, attribute.Float64Value(v))
}

// Float64Slice is comparable to attribute.Float64Slice.
func Float64Slice(k string, v []float64) KeyValue {
	return newKeyValue(k, attribute.Float64SliceValue(v))
}

// String is comparable to attribute.String.
func String(k, v string) KeyValue {
	return newKeyValue(k, attribute.StringValue(v))
}

// StringSlice is comparable to attribute.StringSlice.
func StringSlice(k string, v []string) KeyValue {
	return newKeyValue(k, attribute.StringSliceValue(v))
}

// Stringer is comparable to attribute.Stringer.
func Stringer(k string, v fmt.Stringer) KeyValue {
	return newKeyValue(k, attribute.StringValue(v.String()))
}
