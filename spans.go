package spans

import (
	"fmt"
	"log/slog"

	pkgotel "github.com/ebi-yade/spans/pkg/otel"
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
			results, err := pkgotel.MarshalOtelAttributes(v)
			if err != nil { // 解釈に失敗した場合はログを吐いて無視するようにしている
				slog.Error(fmt.Sprintf("error MarshalOtelAttributes: key=>%s, type=>%T, error=>%v", attr.key, v, err))
			}
			for _, res := range results {
				attributes = append(attributes, attribute.KeyValue{
					Key:   attribute.Key(fmt.Sprintf("%s.%s", attr.key, res.Key)),
					Value: res.Value,
				})
			}
		}
	}

	return attributes
}

// ============================================================================
// Extended APIs to support flexible data types while keeping code simple
// ============================================================================

// WithAttrs can be used in place of trace.WithAttributes to set multiple attributes on a span at the time of creation.
func WithAttrs(attrs ...KeyValue) trace.SpanStartEventOption {
	attributes := getStandardAttributes(attrs)
	return trace.WithAttributes(attributes...)
}

// SetAttrs can be used in place of span.SetAttributes to set multiple attributes on a span after it has been created.
func SetAttrs(span trace.Span, attrs ...KeyValue) {
	attributes := getStandardAttributes(attrs)
	span.SetAttributes(attributes...)
}

// ObjectAttr は構造体や map などの複合型を属性として設定するための KeyValue を生成します。
// 綺麗な型制約や命名を提供できなかったのはご愛嬌として、以下の点にご注意ください。
//   - プリミティブ型以外の配列、またはそれを子として持つ構造体や map は構造上受け入れられない（JSONとして出力される）
//   - map のキーの基底型が string でない場合は <T value> が属性のキーとして採用される
func ObjectAttr(k string, v interface{}) KeyValue {
	return newKeyValue(k, v)
}

// ============================================================================
// Compatible APIs to initialize KeyValue
// ============================================================================

// BoolAttr is comparable to attribute.Bool.
func BoolAttr(k string, v bool) KeyValue {
	return newKeyValue(k, attribute.BoolValue(v))
}

// BoolSliceAttr is comparable to attribute.BoolSlice.
func BoolSliceAttr(k string, v []bool) KeyValue {
	return newKeyValue(k, attribute.BoolSliceValue(v))
}

// IntAttr is comparable to attribute.Int.
func IntAttr(k string, v int) KeyValue {
	return newKeyValue(k, attribute.IntValue(v))
}

// IntSliceAttr is comparable to attribute.IntSlice.
func IntSliceAttr(k string, v []int) KeyValue {
	return newKeyValue(k, attribute.IntSliceValue(v))
}

// Int64Attr is comparable to attribute.Int64.
func Int64Attr(k string, v int64) KeyValue {
	return newKeyValue(k, attribute.Int64Value(v))
}

// Int64SliceAttr is comparable to attribute.Int64Slice.
func Int64SliceAttr(k string, v []int64) KeyValue {
	return newKeyValue(k, attribute.Int64SliceValue(v))
}

// Float64Attr is comparable to attribute.Float64.
func Float64Attr(k string, v float64) KeyValue {
	return newKeyValue(k, attribute.Float64Value(v))
}

// Float64SliceAttr is comparable to attribute.Float64Slice.
func Float64SliceAttr(k string, v []float64) KeyValue {
	return newKeyValue(k, attribute.Float64SliceValue(v))
}

// StringAttr is comparable to attribute.String.
func StringAttr(k, v string) KeyValue {
	return newKeyValue(k, attribute.StringValue(v))
}

// StringSliceAttr is comparable to attribute.StringSlice.
func StringSliceAttr(k string, v []string) KeyValue {
	return newKeyValue(k, attribute.StringSliceValue(v))
}

// StringerAttr is comparable to attribute.Stringer.
func StringerAttr(k string, v fmt.Stringer) KeyValue {
	return newKeyValue(k, attribute.StringValue(v.String()))
}
