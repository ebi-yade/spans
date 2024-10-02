package otelattr

import (
	"reflect"
	"strings"
	"unicode"
)

type structFiled struct {
	attributeName   string
	filedName       string
	filedIndex      int
	omitEmpty       bool
	attributePrefix string
}

var structFieldsCache = newCache[[]structFiled]()

func getStructFields(t reflect.Type) []structFiled {
	if v, ok := structFieldsCache.get(t); ok {
		return v
	}

	var fields []structFiled
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		tag := f.Tag.Get("otelattr")
		tagParts := strings.Split(tag, ",")
		if tagParts[0] == "-" {
			continue
		}
		attributeName := tagParts[0]
		attributePrefix := tagParts[0] + "."
		if attributeName == "" {
			attributeName = camelToSnake(f.Name)
			attributePrefix = ""
		}
		var omitEmpty bool
		for _, part := range tagParts[1:] {
			if part == "omitempty" {
				omitEmpty = true
			}
		}

		fields = append(fields, structFiled{
			attributeName:   attributeName,
			filedName:       f.Name,
			filedIndex:      i,
			omitEmpty:       omitEmpty,
			attributePrefix: attributePrefix,
		})
	}

	structFieldsCache.set(t, fields)
	return fields
}

func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				if !unicode.IsUpper(rune(s[i-1])) || (i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
