package synse

// utils.go provides function utilities for the client.

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/creasty/defaults"
	"github.com/pkg/errors"
)

// setDefaults setups default options.
func setDefaults(opts *Options) error {
	if opts == nil {
		return errors.New("options can not be nil")
	}

	if opts.Address == "" {
		return errors.New("no address is specified")
	}

	if err := defaults.Set(opts); err != nil {
		return errors.New("failed to set default configs")
	}

	return nil
}

// setTLS registers the certificates with configured optionss.
func setTLS(opts *Options) (tls.Certificate, error) {
	if opts.TLS.CertFile == "" && opts.TLS.KeyFile == "" {
		return tls.Certificate{}, errors.New("no certificates are specified")
	}

	cert, err := tls.LoadX509KeyPair(opts.TLS.CertFile, opts.TLS.KeyFile)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "failed to set client certificates")
	}

	return cert, nil
}

// buildURL builds up a complete URL from given scheme, host and path.
func buildURL(scheme string, host string, path ...string) string {
	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   strings.Join(path, "/"),
	}

	return u.String()
}

// makePath joins the given components into a string, delimited with '/' which
// can then be used as the URI for API requests.
func makePath(components ...string) string {
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
