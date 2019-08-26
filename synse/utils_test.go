package synse

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		in       map[string]interface{}
		expected string
	}{
		{
			in: map[string]interface{}{
				"scheme": "http",
				"host":   "localhost:5000",
				"path":   []string{""},
			},
			expected: "http://localhost:5000",
		},
		{
			in: map[string]interface{}{
				"scheme": "http",
				"host":   "localhost:5000",
				"path":   []string{"test"},
			},
			expected: "http://localhost:5000/test",
		},
		{
			in: map[string]interface{}{
				"scheme": "http",
				"host":   "localhost:5000",
				"path":   []string{"v3", "test"},
			},
			expected: "http://localhost:5000/v3/test",
		},
		{
			in: map[string]interface{}{
				"scheme": "https",
				"host":   "localhost:5000",
				"path":   []string{"v3", "test"},
			},
			expected: "https://localhost:5000/v3/test",
		},
	}

	for _, tt := range tests {
		out := buildURL(tt.in["scheme"].(string), tt.in["host"].(string), tt.in["path"].([]string)...)
		assert.Equal(t, tt.expected, out)
	}
}

func TestMakePath(t *testing.T) {
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
			in:       []string{"/foo"},
			expected: "/foo",
		},
		{
			in:       []string{"/foo", "bar"},
			expected: "/foo/bar",
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
		out := makePath(tt.in...)
		assert.Equal(t, tt.expected, out)
	}
}

func TestStructToURLValues(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected url.Values
	}{
		{
			struct{}{},
			url.Values{},
		},
		{
			struct {
				foo string
			}{
				foo: "bar",
			},
			url.Values{
				"foo": []string{"bar"},
			},
		},
		{
			struct {
				Foo string
			}{
				Foo: "bar",
			},
			url.Values{
				"foo": []string{"bar"},
			},
		},
		{
			struct {
				foo string
			}{
				foo: "Bar",
			},
			url.Values{
				"foo": []string{"Bar"},
			},
		},
		{
			struct {
				foo int
			}{
				foo: int(1),
			},
			url.Values{
				"foo": []string{"1"},
			},
		},
		{
			struct {
				foo bool
			}{
				foo: true,
			},
			url.Values{
				"foo": []string{"true"},
			},
		},
		{
			struct {
				foo []string
			}{
				foo: []string{"foo", "bar"},
			},
			url.Values{
				"foo": []string{"foo", "bar"},
			},
		},
		{
			struct {
				foo []string
			}{
				foo: []string{"foo,bar", "bar"},
			},
			url.Values{
				"foo": []string{"foo,bar", "bar"},
			},
		},
	}

	for _, tt := range tests {
		out := structToURLValues(tt.in)
		assert.Equal(t, tt.expected, out)
	}
}
