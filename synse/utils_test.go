package synse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeURI(t *testing.T) {
	tests := []struct {
		in       []string
		expected string
	}{
		{
			in:       []string{""},
			expected: "",
		},
		{
			in:       []string{"foo"},
			expected: "foo",
		},
		{
			in:       []string{"foo", "bar"},
			expected: "foo/bar",
		},
		{
			in:       []string{"foo", "bar", "baz/"},
			expected: "foo/bar/baz/",
		},
		{
			in:       []string{"foo/", "/bar"},
			expected: "foo///bar",
		},
	}

	for _, tt := range tests {
		out := makeURI(tt.in...)
		assert.Equal(t, tt.expected, out)
	}
}

func TestStructToMapString(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected map[string]string
	}{
		{
			struct{}{},
			map[string]string{},
		},
		{
			struct {
				foo string
			}{
				foo: "bar",
			},
			map[string]string{
				"foo": "bar",
			},
		},
		{
			struct {
				Foo string
			}{
				Foo: "bar",
			},
			map[string]string{
				"foo": "bar",
			},
		},
		{
			struct {
				foo string
			}{
				foo: "Bar",
			},
			map[string]string{
				"foo": "Bar",
			},
		},
		{
			struct {
				foo int
			}{
				foo: int(1),
			},
			map[string]string{
				"foo": "1",
			},
		},
		{
			struct {
				foo bool
			}{
				foo: true,
			},
			map[string]string{
				"foo": "true",
			},
		},
		{
			struct {
				foo []string
			}{
				foo: []string{"foo", "bar"},
			},
			map[string]string{
				"foo": "foo,bar",
			},
		},
		{
			struct {
				foo []string
				bar string
			}{
				foo: []string{"foo", "bar", "foobar"},
				bar: "bar",
			},
			map[string]string{
				"foo": "foo,bar,foobar",
				"bar": "bar",
			},
		},
	}

	for _, tt := range tests {
		out := structToMapString(tt.in)
		assert.Equal(t, tt.expected, out)
	}
}
