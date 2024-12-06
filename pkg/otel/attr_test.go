package otel

import (
	"cmp"
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

type structNoTags struct {
	BoolValue   bool
	BoolSlice   []bool
	IntValue    int
	IntSlice    []int
	Int64Value  int64
	Int64Slice  []int64
	FloatValue  float64
	FloatSlice  []float64
	StringValue string
	StringSlice []string
}

func TestMarshalOtelAttributes__PrimitiveTypes(t *testing.T) {
	args := structNoTags{
		BoolValue:   true,
		IntValue:    1,
		Int64Value:  2,
		FloatValue:  3.14,
		StringValue: "hello",
		StringSlice: []string{"hello", "world"},
	}
	want := []attribute.KeyValue{
		attribute.Bool("bool_value", true),
		attribute.BoolSlice("bool_slice", []bool{}),
		attribute.Int64("int_value", 1),
		attribute.Int64Slice("int_slice", []int64{}),
		attribute.Int64("int64_value", 2),
		attribute.Int64Slice("int64_slice", []int64{}),
		attribute.Float64("float_value", 3.14),
		attribute.Float64Slice("float_slice", []float64{}),
		attribute.String("string_value", "hello"),
		attribute.StringSlice("string_slice", []string{"hello", "world"}),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__Time(t *testing.T) {
	const layout = "2006-01-02T15:04:05.999999999"
	const testTimeStr = "2021-01-01T00:00:00.123456789"
	testTime, err := time.ParseInLocation(layout, testTimeStr, time.FixedZone("Asia/Tokyo", 9*60*60))
	require.NoError(t, err)

	type Wrapper time.Time
	args := struct {
		StdTime     time.Time
		WrapperTime Wrapper
		Map         map[string]interface{}
	}{
		StdTime:     testTime,
		WrapperTime: Wrapper(testTime),
		Map: map[string]interface{}{
			"std_time":     testTime,
			"wrapper_time": Wrapper(testTime),
		},
	}

	const expectedStr = "2021-01-01T00:00:00.123456789+09:00"
	want := []attribute.KeyValue{
		attribute.String("std_time", expectedStr),
		attribute.String("wrapper_time", expectedStr),
		attribute.String("map.std_time", expectedStr),
		attribute.String("map.wrapper_time", expectedStr),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithTypedNil(t *testing.T) {
	var args *structNoTags
	want := make([]attribute.KeyValue, 0)
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithNil(t *testing.T) {
	var args interface{}
	want := make([]attribute.KeyValue, 0)
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

type structWithNameTags struct {
	BoolValue   bool      `otel:"b"`
	BoolSlice   []bool    `otel:"bs"`
	IntValue    int       `otel:"i"`
	IntSlice    []int     `otel:"is"`
	Int64Value  int64     `otel:"i64"`
	Int64Slice  []int64   `otel:"is64"`
	FloatValue  float64   `otel:"f"`
	FloatSlice  []float64 `otel:"fs"`
	StringValue string    `otel:"s"`
	StringSlice []string  `otel:"ss"`
}

func TestMarshalOtelAttributes__WithTags(t *testing.T) {
	args := structWithNameTags{
		BoolValue:   true,
		IntValue:    1,
		Int64Value:  2,
		FloatValue:  3.14,
		StringValue: "hello",
	}
	want := []attribute.KeyValue{
		attribute.Bool("b", true),
		attribute.BoolSlice("bs", []bool{}),
		attribute.Int64("i", 1),
		attribute.Int64Slice("is", []int64{}),
		attribute.Int64("i64", 2),
		attribute.Int64Slice("is64", []int64{}),
		attribute.Float64("f", 3.14),
		attribute.Float64Slice("fs", []float64{}),
		attribute.String("s", "hello"),
		attribute.StringSlice("ss", []string{}),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

type structWithOmitemptyTags struct {
	BoolValue   bool      `otel:",omitempty"`
	BoolSlice   []bool    `otel:",omitempty"`
	IntValue    int       `otel:",omitempty"`
	IntSlice    []int     `otel:",omitempty"`
	Int64Value  int64     `otel:",omitempty"`
	Int64Slice  []int64   `otel:",omitempty"`
	FloatValue  float64   `otel:",omitempty"`
	FloatSlice  []float64 `otel:",omitempty"`
	StringValue string    `otel:",omitempty"`
	StringSlice []string  `otel:",omitempty"`
}

func TestMarshalOtelAttributes__WithOmitempty(t *testing.T) {
	args := structWithOmitemptyTags{}
	want := []attribute.KeyValue{}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

type structWithMarshaller struct {
	Value int
}

func (t structWithMarshaller) MarshalOtelAttributes() ([]attribute.KeyValue, error) {
	return []attribute.KeyValue{
		attribute.Int("http.staus_code", t.Value),
	}, nil
}

func TestMarshalOtelAttributes__WithStructMarshaller(t *testing.T) {
	args := structWithMarshaller{Value: 200}
	want := []attribute.KeyValue{
		attribute.Int("http.staus_code", 200),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithStructPointerMarshaller(t *testing.T) {
	args := &structWithMarshaller{Value: 200}
	want := []attribute.KeyValue{
		attribute.Int("http.staus_code", 200),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__NestedStructWithMarshallerMember(t *testing.T) {
	args := structWithMarshaller{Value: 200}
	want := []attribute.KeyValue{
		attribute.Int("http.staus_code", 200), // `struct.` prefix is omitted
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

type structWithNameAndOmitemptyTags struct {
	BoolValue   bool      `otel:"b,omitempty"`
	BoolSlice   []bool    `otel:"bs,omitempty"`
	IntValue    int       `otel:"i,omitempty"`
	IntSlice    []int     `otel:"is,omitempty"`
	Int64Value  int64     `otel:"i64,omitempty"`
	Int64Slice  []int64   `otel:"is64,omitempty"`
	FloatValue  float64   `otel:"f,omitempty"`
	FloatSlice  []float64 `otel:"fs,omitempty"`
	StringValue string    `otel:"s,omitempty"`
	StringSlice []string  `otel:"ss,omitempty"`
}

func TestMarshalOtelAttributes__WithStructInStruct(t *testing.T) {
	args := struct {
		Struct structWithNameAndOmitemptyTags
	}{
		Struct: structWithNameAndOmitemptyTags{
			BoolSlice:   []bool{true},
			IntSlice:    []int{1},
			Int64Slice:  []int64{2},
			FloatSlice:  []float64{3.14, 2.71},
			StringSlice: []string{"hello", "world"},
		},
	}
	want := []attribute.KeyValue{
		attribute.BoolSlice("struct.bs", []bool{true}),
		attribute.IntSlice("struct.is", []int{1}),
		attribute.Int64Slice("struct.is64", []int64{2}),
		attribute.Float64Slice("struct.fs", []float64{3.14, 2.71}),
		attribute.StringSlice("struct.ss", []string{"hello", "world"}),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithStructInStructPointer(t *testing.T) {
	args := &structWithNameAndOmitemptyTags{
		BoolSlice:   []bool{true},
		IntSlice:    []int{1},
		Int64Slice:  []int64{2},
		FloatSlice:  []float64{3.14, 2.71},
		StringSlice: []string{"hello", "world"},
	}
	want := []attribute.KeyValue{
		attribute.BoolSlice("bs", []bool{true}),
		attribute.IntSlice("is", []int{1}),
		attribute.Int64Slice("is64", []int64{2}),
		attribute.Float64Slice("fs", []float64{3.14, 2.71}),
		attribute.StringSlice("ss", []string{"hello", "world"}),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithStructInStructWithPrefixTag(t *testing.T) {
	args := struct {
		Struct structWithNameAndOmitemptyTags `otel:"test"`
	}{
		structWithNameAndOmitemptyTags{
			BoolSlice:   []bool{true},
			IntSlice:    []int{1},
			Int64Slice:  []int64{2},
			FloatSlice:  []float64{3.14, 2.71},
			StringSlice: []string{"hello", "world"},
		},
	}
	want := []attribute.KeyValue{
		attribute.BoolSlice("test.bs", []bool{true}),
		attribute.IntSlice("test.is", []int{1}),
		attribute.Int64Slice("test.is64", []int64{2}),
		attribute.Float64Slice("test.fs", []float64{3.14, 2.71}),
		attribute.StringSlice("test.ss", []string{"hello", "world"}),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithIgnoreField(t *testing.T) {
	args := struct {
		IgnoreField int `otel:"-"`
	}{
		IgnoreField: 100,
	}
	want := []attribute.KeyValue{}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func TestMarshalOtelAttributes__WithMap(t *testing.T) {
	args := map[string]interface{}{
		"bool_value":   true,
		"int_value":    1,
		"int64_value":  2,
		"float_value":  3.14,
		"string_value": "hello",
	}
	want := []attribute.KeyValue{
		attribute.Bool("bool_value", true),
		attribute.Int64("int_value", 1),
		attribute.Int64("int64_value", 2),
		attribute.Float64("float_value", 3.14),
		attribute.String("string_value", "hello"),
	}
	got, err := MarshalOtelAttributes(args)
	assert.NoError(t, err)
	assertAttributes(t, want, got)
}

func assertAttributes(tb testing.TB, want, got []attribute.KeyValue, msgAndArgs ...interface{}) bool {
	tb.Helper()

	sortKV := func(i, j attribute.KeyValue) int {
		return cmp.Compare(i.Key, j.Key)
	}
	slices.SortFunc(want, sortKV)
	slices.SortFunc(got, sortKV)

	if !assert.ObjectsAreEqualValues(want, got) {
		return assert.Fail(tb,
			fmt.Sprintf(
				"not equal:\n\twant: %s\n\t got: %s",
				attributesToString(want),
				attributesToString(got),
			),
			msgAndArgs...,
		)
	}
	return true
}

func attributesToString(kvs []attribute.KeyValue) string {
	var s string
	for i, kv := range kvs {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%s=%s", kv.Key, kv.Value.Emit())
	}
	return "[" + s + "]"
}
