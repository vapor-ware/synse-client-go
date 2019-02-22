package synse

// utils.go provides function utilities for the client.

import (
	"fmt"
	"reflect"
	"strings"
)

// makeURI joins the given components into a string, delimited with '/' which
// can then be used as the URI for API requests.
func makeURI(components ...string) string {
	return strings.Join(components, "/")
}

// structToMapString decodes a struct value into a map[string]string that can
// be used as query parameters. It assumes that the struct fields follow one of
// these types: bool, string, int, float, slice.
func structToMapString(s interface{}) map[string]string {
	out := map[string]string{}
	v := ""

	fields := reflect.TypeOf(s)
	values := reflect.ValueOf(s)

	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)
		value := values.Field(i)

		switch value.Kind() {
		case reflect.Slice:
			s := []string{}
			for i := 0; i < value.Len(); i++ {
				s = append(s, fmt.Sprint(value.Index(i)))
			}

			v = strings.Join(s, ",")
		default:
			v = fmt.Sprint(value)
		}

		out[strings.ToLower(field.Name)] = v
	}

	return out
}
