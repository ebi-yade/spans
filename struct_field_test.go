package otelattr

import (
	"testing"
)

func Test_camelToSnake(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"fooID", "foo_id"},
		{"fooBar", "foo_bar"},
		{"Baz", "baz"},
		{"FizzBuzz", "fizz_buzz"},
		{"FooIDBar", "foo_id_bar"},
		{"fooIDBar", "foo_id_bar"},
		{"foo", "foo"},
		{"FOO", "foo"},
		{"", ""},
		{"a", "a"},
		{"foo123Bar", "foo123_bar"},
		{"IDFooBar", "id_foo_bar"},
		{"MyHTTPRequest", "my_http_request"},
	}

	for _, c := range cases {
		got := camelToSnake(c.input)
		if got != c.expected {
			t.Errorf("toSnakeCase(%q) == %q, want %q", c.input, got, c.expected)
		}
	}
}
