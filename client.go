package client

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	DefaultTimeout = 30 * time.Second
	DocsURL        = "https://observatory.manrs.org/api/v2/docs"
	EndpointURL    = "https://observatory.manrs.org/api/v2"
	Version        = "1.0.0"
)

var (
	jsonMarshalerType = reflect.TypeFor[json.Marshaler]()
	textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()
)

// Client represents the MANRS Observatory API client.
type Client struct {
	apiKey string
	client *retryablehttp.Client
	debug  bool
}

// New returns a client for the MANRS Observatory API.
func New(apiKey string, opts ...Option) *Client {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil
	httpClient.HTTPClient.Timeout = DefaultTimeout
	httpClient.CheckRetry = checkRetry

	c := &Client{
		apiKey: apiKey,
		client: httpClient,
		debug:  false,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetAPIKey replaces the existing API key in use.
func (c *Client) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// SetTimeout updates the client timeout.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.client.HTTPClient.Timeout = timeout
}

// Option configures a Client.
type Option func(*Client)

// WithClient allows providing a custom *http.Client.
func WithClient(client *http.Client) Option {
	return func(c *Client) {
		if client == nil {
			client = &http.Client{Timeout: DefaultTimeout}
		}

		c.client.HTTPClient = client
	}
}

// WithDebug enables or disables debug mode.
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithTimeout sets a custom timeout on the underlying http.Client.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.client == nil {
			c.client = retryablehttp.NewClient()
		}

		c.client.HTTPClient.Timeout = timeout
	}
}

func (c *Client) marshal(v any) ([]byte, error) {
	b, err := json.Marshal(stripReadOnly(reflect.ValueOf(v)))
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return b, nil
}

// stripReadOnly rebuilds v with every readOnly:"true" field removed. Those fields are assigned
// by the server, so sending them back is at best ignored and at worst rejected: the API validates
// Favorite.ID as a UUID, which a zero-valued string is not.
//
// Values that serialise themselves (time.Time and anything else implementing json.Marshaler or
// encoding.TextMarshaler) are returned untouched rather than walked, since their fields are
// unexported and carry no meaning once decomposed.
func stripReadOnly(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}

	if selfMarshaling(v) {
		return v.Interface()
	}

	switch v.Kind() { //nolint:exhaustive // Scalar kinds fall through to the default, which returns them as-is.
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return nil
		}

		return stripReadOnly(v.Elem())

	case reflect.Map:
		if v.IsNil() {
			return nil
		}

		out := make(map[string]any, v.Len())
		for _, key := range v.MapKeys() {
			out[fmt.Sprint(key.Interface())] = stripReadOnly(v.MapIndex(key))
		}

		return out

	case reflect.Slice, reflect.Array:
		// Byte slices encode as base64 strings, not as arrays of numbers.
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return v.Interface()
		}

		if v.Kind() == reflect.Slice && v.IsNil() {
			return nil
		}

		out := make([]any, v.Len())
		for i := range v.Len() {
			out[i] = stripReadOnly(v.Index(i))
		}

		return out

	case reflect.Struct:
		out := make(map[string]any)
		addStructFields(v, out)

		return out

	default:
		return v.Interface()
	}
}

// addStructFields writes v's serialisable fields into out, skipping read-only ones and promoting
// embedded structs into the parent object the way encoding/json does.
func addStructFields(v reflect.Value, out map[string]any) {
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() || field.Tag.Get("readOnly") == "true" {
			continue
		}

		name, opts, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" && opts == "" {
			continue
		}

		value := v.Field(i)

		if field.Anonymous && name == "" {
			embedded := value
			for embedded.Kind() == reflect.Pointer {
				if embedded.IsNil() {
					break
				}

				embedded = embedded.Elem()
			}

			if embedded.Kind() == reflect.Struct && !selfMarshaling(embedded) {
				addStructFields(embedded, out)

				continue
			}
		}

		if omitsField(opts, value) {
			continue
		}

		if name == "" {
			name = field.Name
		}

		out[name] = stripReadOnly(value)
	}
}

// omitsField reports whether a json tag's omitempty or omitzero option drops this value.
func omitsField(opts string, value reflect.Value) bool {
	if opts == "" {
		return false
	}

	for opt := range strings.SplitSeq(opts, ",") {
		switch opt {
		case "omitzero":
			if value.IsZero() {
				return true
			}
		case "omitempty":
			if isEmptyValue(value) {
				return true
			}
		}
	}

	return false
}

// isEmptyValue mirrors encoding/json's definition of empty for the omitempty option, which covers
// false, 0, nil pointers and interfaces, and any empty array, map, slice, or string.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() { //nolint:exhaustive // Kinds with no encoding/json notion of empty fall through to the default.
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	default:
		return false
	}
}

// selfMarshaling reports whether v encodes itself, in which case its fields must not be walked.
func selfMarshaling(v reflect.Value) bool {
	t := v.Type()

	if t.Implements(jsonMarshalerType) || t.Implements(textMarshalerType) {
		return true
	}

	// A pointer receiver still counts when the value is addressable, matching encoding/json.
	if v.CanAddr() {
		pt := reflect.PointerTo(t)
		if pt.Implements(jsonMarshalerType) || pt.Implements(textMarshalerType) {
			return true
		}
	}

	return false
}

type methodContextKey struct{}

// WithMethod attaches the HTTP method to ctx. On a transport error retryablehttp's CheckRetry
// receives no request and a nil response, so ctx is the only way to recover the method.
func withMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodContextKey{}, method)
}

// CheckRetry retries GET only. A write may have been applied before a timeout hid the response,
// so retrying POST, PATCH, or DELETE risks repeating it, e.g. a duplicate announcement email.
func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if method, _ := ctx.Value(methodContextKey{}).(string); method != http.MethodGet {
		return false, nil
	}

	retry, err := retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	if err != nil {
		return retry, fmt.Errorf("default retry policy: %w", err)
	}

	return retry, nil
}
