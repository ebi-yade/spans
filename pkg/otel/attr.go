package otel

import (
	"encoding/json"
	"fmt"

	"reflect"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

type Marshaler interface {
	MarshalOtelAttributes() ([]attribute.KeyValue, error)
}

func MarshalOtelAttributes(v interface{}) ([]attribute.KeyValue, error) {
	if m, ok := v.(Marshaler); ok {
		return m.MarshalOtelAttributes()
	}
	rv := reflect.ValueOf(v)
	return marshalOtelAttributes(rv)
}

func marshalOtelAttributes(rv reflect.Value) ([]attribute.KeyValue, error) {
	switch rv.Kind() {
	case reflect.Struct:
		return marshalStruct(rv)
	case reflect.Ptr:
		return marshalOtelAttributes(rv.Elem())
	case reflect.Interface:
		return marshalOtelAttributes(rv.Elem())
	case reflect.Map:
		return marshalMap(rv)
	default:
		return nil, fmt.Errorf("unsupported type %s", rv.Type())
	}
}

func marshalMap(rv reflect.Value) ([]attribute.KeyValue, error) {
	keys := rv.MapKeys()
	if len(keys) == 0 {
		return []attribute.KeyValue{}, nil
	}
	if keys[0].Kind() != reflect.String {
		return nil, fmt.Errorf("unsupported map key type %s", keys[0].Type())
	}
	attrs := make([]attribute.KeyValue, 0, len(keys))
	for index, key := range keys {
		mv := rv.MapIndex(key)
		keyString := key.String()
		switch value := mv.Interface().(type) {
		case bool:
			attrs = append(attrs, attribute.Bool(keyString, value))
		case int:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case int8:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case int16:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case int32:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case int64:
			attrs = append(attrs, attribute.Int64(keyString, value))
		case uint:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case uint8:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case uint16:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case uint32:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case uint64:
			attrs = append(attrs, attribute.Int64(keyString, (int64)(value)))
		case float32:
			attrs = append(attrs, attribute.Float64(keyString, (float64)(value)))
		case float64:
			attrs = append(attrs, attribute.Float64(keyString, value))
		case string:
			attrs = append(attrs, attribute.String(keyString, value))
		case []bool:
			attrs = append(attrs, attribute.BoolSlice(keyString, value))
		case []int:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []int8:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []int16:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []int32:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []int64:
			attrs = append(attrs, attribute.Int64Slice(keyString, value))
		case []uint:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []uint8:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []uint16:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []uint32:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []uint64:
			s := make([]int64, len(value))
			for i, v := range value {
				s[i] = int64(v)
			}
			attrs = append(attrs, attribute.Int64Slice(keyString, s))
		case []float32:
			s := make([]float64, len(value))
			for i, v := range value {
				s[i] = float64(v)
			}
			attrs = append(attrs, attribute.Float64Slice(keyString, s))
		case []float64:
			attrs = append(attrs, attribute.Float64Slice(keyString, value))
		case []string:
			attrs = append(attrs, attribute.StringSlice(keyString, value))
		default:
			if mv.Type().ConvertibleTo(timeType) {
				t := mv.Convert(timeType).Interface().(time.Time)
				attrs = append(attrs, attribute.String(keyString, t.Format(time.RFC3339Nano)))
				continue
			}

			kvs, err := marshalField(structFiled{
				attributeName:   keyString,
				filedIndex:      index,
				attributePrefix: keyString + ".",
			}, reflect.ValueOf(value))
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, kvs...)
		}
	}
	return attrs, nil
}

func marshalStruct(rv reflect.Value) ([]attribute.KeyValue, error) {
	t := rv.Type()
	fields := getStructFields(t)
	kvs := make([]attribute.KeyValue, 0, len(fields))
	for _, f := range fields {
		fv := rv.Field(f.filedIndex)
		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		filedValue, err := marshalField(f, fv)
		if err != nil {
			return nil, err
		}
		kvs = append(kvs, filedValue...)
	}
	return kvs, nil
}

func marshalField(f structFiled, fv reflect.Value) ([]attribute.KeyValue, error) {
	switch fv.Kind() {
	case reflect.Bool:
		return []attribute.KeyValue{attribute.Bool(f.attributeName, fv.Bool())}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []attribute.KeyValue{attribute.Int64(f.attributeName, fv.Int())}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []attribute.KeyValue{attribute.Int64(f.attributeName, int64(fv.Uint()))}, nil
	case reflect.Float32, reflect.Float64:
		return []attribute.KeyValue{attribute.Float64(f.attributeName, fv.Float())}, nil
	case reflect.String:
		return []attribute.KeyValue{attribute.String(f.attributeName, fv.String())}, nil
	case reflect.Slice, reflect.Array:
		return marshalSlice(f, fv)
	case reflect.Struct:
		// convert time.Time to string
		if fv.Type().ConvertibleTo(timeType) {
			t := fv.Convert(timeType).Interface().(time.Time)
			return []attribute.KeyValue{attribute.String(f.attributeName, t.Format(time.RFC3339Nano))}, nil
		}
		return f.marshalNested(fv)

	case reflect.Ptr:
		if fv.IsNil() {
			return nil, nil
		}
		return marshalField(f, fv.Elem())

	case reflect.Map:
		return f.marshalNested(fv)

	default:
		bs, err := json.Marshal(fv.Interface())
		if err != nil {
			return []attribute.KeyValue{}, err
		}
		return []attribute.KeyValue{attribute.String(f.attributeName, string(bs))}, nil
	}
}

func (s structFiled) marshalNested(fv reflect.Value) ([]attribute.KeyValue, error) {
	attrs, err := MarshalOtelAttributes(fv.Interface())
	if err != nil {
		return []attribute.KeyValue{}, err
	}
	for i := range attrs {
		attrs[i].Key = attribute.Key(s.attributePrefix) + attrs[i].Key
	}
	return attrs, nil
}

func marshalSlice(f structFiled, fv reflect.Value) ([]attribute.KeyValue, error) {
	switch fv.Type().Elem().Kind() {
	case reflect.Bool:
		return []attribute.KeyValue{attribute.BoolSlice(f.attributeName, reflectValueToSlice[bool](fv))}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []attribute.KeyValue{attribute.Int64Slice(f.attributeName, reflectValueToSlice[int64](fv))}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []attribute.KeyValue{attribute.Int64Slice(f.attributeName, reflectValueToSlice[int64](fv))}, nil
	case reflect.Float32, reflect.Float64:
		return []attribute.KeyValue{attribute.Float64Slice(f.attributeName, reflectValueToSlice[float64](fv))}, nil
	case reflect.String:
		return []attribute.KeyValue{attribute.StringSlice(f.attributeName, reflectValueToSlice[string](fv))}, nil
	case reflect.Struct:
		if fv.Type().Elem().ConvertibleTo(timeType) {
			strs := make([]string, fv.Len())
			for i := 0; i < fv.Len(); i++ {
				t := fv.Index(i).Convert(timeType).Interface().(time.Time)
				strs[i] = t.Format(time.RFC3339Nano)
			}
			return []attribute.KeyValue{attribute.StringSlice(f.attributeName, strs)}, nil
		}
		fallthrough // OpenTelemetry プロトコルのレベルで、そもそも合成型の配列は非対応のため、いずれにせよ文字列化以外の選択肢は与えない
	default:
		bs, err := json.Marshal(fv.Interface())
		if err != nil {
			return []attribute.KeyValue{}, err
		}
		return []attribute.KeyValue{attribute.String(f.attributeName, string(bs))}, nil
	}
}

func reflectValueToSlice[T any](v reflect.Value) []T {
	slice := make([]T, v.Len())
	var zero T
	rt := reflect.TypeOf(zero)
	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Convert(rt).Interface().(T)
	}
	return slice
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice:
		return v.Len() == 0
	default:
		return false
	}
}
