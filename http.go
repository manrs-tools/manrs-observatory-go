package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/cast"
)

const (
	ContentTypeJSON  = "application/json"
	EncodingTypeGZIP = "gzip"
)

type (
	// GenericResponse is the API's problem-detail payload for a non-2xx response.
	GenericResponse struct {
		Title    string        `json:"title"`
		Status   int           `json:"status"`
		Detail   string        `json:"detail"`
		Instance string        `json:"instance,omitempty"`
		Type     string        `json:"type,omitempty"`
		Errors   []ErrorDetail `json:"errors,omitempty"`
	}

	// ErrorDetail is one field-level failure within a GenericResponse. Location names the
	// offending field as a dotted path from the request root, e.g. "body.favorite.manrs",
	// and Value carries what was rejected.
	ErrorDetail struct {
		Location string `json:"location,omitempty"`
		Message  string `json:"message"`
		Value    any    `json:"value,omitempty"`
	}

	// ResponseError is the error every request method returns for a non-2xx response. A caller
	// wanting the structured detail behind the message recovers it with errors.As.
	ResponseError struct {
		GenericResponse
	}
)

func (e ResponseError) Error() string {
	return e.ToErrorString()
}

func (d ErrorDetail) String() string {
	if d.Location == "" {
		return d.Message
	}

	return d.Location + ": " + d.Message
}

func (r GenericResponse) ToErrorString() string {
	var sb strings.Builder

	if r.Detail == "" {
		r.Detail = r.Title
	}

	fmt.Fprintf(&sb, "%s (HTTP %d)", r.Detail, r.Status)

	if len(r.Errors) > 0 {
		sb.WriteString(": ")
		for i, e := range r.Errors {
			sb.WriteString(e.String())
			if i < len(r.Errors)-1 {
				sb.WriteString("; ")
			}
		}
	}

	return sb.String()
}

func (c *Client) DELETE(ctx context.Context, endpoint string, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "DELETE", nil, obj)
}

func (c *Client) GET(ctx context.Context, endpoint string, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "GET", nil, obj)
}

func (c *Client) PATCH(ctx context.Context, endpoint string, body []byte, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "PATCH", body, obj)
}

func (c *Client) POST(ctx context.Context, endpoint string, body []byte, obj any) ([]byte, error) {
	return c.makeRequest(ctx, endpoint, "POST", body, obj)
}

func (c *Client) makeRequest(ctx context.Context, endpoint string, method string, requestBody []byte, object any) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/%s", EndpointURL, endpoint)

	req, err := retryablehttp.NewRequestWithContext(withMethod(ctx, method), method, endpointURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Add("Authorization", "Bearer "+c.apiKey)
	}

	req.Header.Add("Accept", ContentTypeJSON)
	req.Header.Add("Accept-Encoding", EncodingTypeGZIP)
	req.Header.Add("Content-Type", ContentTypeJSON)
	req.Header.Add("Manrs-Client-Version", Version)

	if c.debug {
		fmt.Fprintf(os.Stderr, "> %s %s\n", method, redactURL(endpointURL))
		if len(requestBody) > 0 {
			fmt.Fprintf(os.Stderr, "> body: %s\n", redactSecrets(requestBody))
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	if res == nil {
		return nil, errors.New("read response")
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	resBody, err = unmarshalResBody(resBody, res.Header.Get("Content-Encoding"))
	if err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	if c.debug {
		fmt.Fprintf(os.Stderr, "< %s %s %d\n", method, redactURL(endpointURL), res.StatusCode)
		if len(resBody) > 0 {
			fmt.Fprintln(os.Stderr, string(redactSecrets(resBody)))
		}
	}

	if res.StatusCode >= http.StatusBadRequest {
		if len(resBody) > 0 {
			resp := GenericResponse{}

			if err = json.Unmarshal(resBody, &resp); err == nil {
				return nil, ResponseError{GenericResponse: resp}
			}
		}

		return resBody, errors.New("request status code " + cast.ToString(res.StatusCode))
	}

	if len(resBody) > 0 && object != nil {
		if err = json.Unmarshal(resBody, object); err != nil {
			return nil, fmt.Errorf("unmarshal response body: %w", err)
		}
	}

	return resBody, nil
}

// download issues a GET that stops at the API's redirect instead of following it, then retrieves
// the object from the storage URL that redirect names. The storage URL authorises itself, and the
// JSON content negotiation makeRequest performs does not apply to a binary object, so neither the
// API's Authorization header nor its Accept headers are carried across.
func (c *Client) download(ctx context.Context, endpoint string) ([]byte, error) {
	endpointURL := fmt.Sprintf("%s/%s", EndpointURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	req.Header.Set("Manrs-Client-Version", Version)

	httpClient := *c.client.HTTPClient
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	if c.debug {
		fmt.Fprintf(os.Stderr, "> GET %s\n", redactURL(endpointURL))
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if c.debug {
		fmt.Fprintf(os.Stderr, "< GET %s %d\n", redactURL(endpointURL), res.StatusCode)
	}

	if res.StatusCode >= http.StatusMultipleChoices && res.StatusCode < http.StatusBadRequest {
		location := res.Header.Get("Location")
		if location == "" {
			return nil, errors.New("redirect response carried no location")
		}

		target, parseErr := req.URL.Parse(location)
		if parseErr != nil {
			return nil, fmt.Errorf("parse redirect location: %w", parseErr)
		}

		return c.fetchObject(ctx, target.String())
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if res.StatusCode >= http.StatusBadRequest {
		resp := GenericResponse{}

		if len(resBody) > 0 && json.Unmarshal(resBody, &resp) == nil {
			return nil, ResponseError{GenericResponse: resp}
		}

		return nil, errors.New("request status code " + cast.ToString(res.StatusCode))
	}

	return resBody, nil
}

// fetchObject retrieves rawURL with none of the API request headers attached.
func (c *Client) fetchObject(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("new storage request: %w", err)
	}

	res, err := c.client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch object: %w", err)
	}
	defer res.Body.Close()

	if c.debug {
		fmt.Fprintf(os.Stderr, "< GET %s %d\n", redactStorageURL(rawURL), res.StatusCode)
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, errors.New("fetch object: status code " + cast.ToString(res.StatusCode))
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read object body: %w", err)
	}

	return resBody, nil
}

const redactedValue = "***REDACTED***"

var sensitiveBodyKeys = map[string]struct{}{
	"password": {}, //nolint:goconst // Redaction match target, unrelated to the request body field name in auth.go and invite.go.
	"apikey":   {},
	"key":      {}, //nolint:goconst // Redaction match target, unrelated to the request body field name in api_key.go.
	"token":    {},
	"secret":   {},
}

func redactValue(v any) any {
	switch val := v.(type) {
	case []any:
		out := make([]any, len(val))
		for i, item := range val {
			out[i] = redactValue(item)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, item := range val {
			if _, sensitive := sensitiveBodyKeys[strings.ToLower(k)]; sensitive {
				out[k] = redactedValue
			} else {
				out[k] = redactValue(item)
			}
		}
		return out
	default:
		return val
	}
}

var sensitivePathParents = map[string]struct{}{
	"invites":        {},
	"reset-password": {},
	"verify":         {},
}

// redactStorageURL redacts every query value, because the query of a signed storage URL carries the
// credential and signature that grant access to the object.
func redactStorageURL(rawURL string) string {
	base, query, hasQuery := strings.Cut(rawURL, "?")
	if !hasQuery {
		return base
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return base
	}
	for k := range values {
		values.Set(k, redactedValue)
	}

	return base + "?" + values.Encode()
}

// redactURL redacts sensitive query values, and the path segment that directly follows an invites,
// reset-password or verify segment, before a URL reaches the debug log.
func redactURL(rawURL string) string {
	base, query, hasQuery := strings.Cut(rawURL, "?")

	segments := strings.Split(base, "/")
	for i := 1; i < len(segments); i++ {
		if _, sensitive := sensitivePathParents[segments[i-1]]; sensitive && segments[i] != "" {
			segments[i] = redactedValue
		}
	}
	base = strings.Join(segments, "/")

	if !hasQuery {
		return base
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return base
	}
	for k := range values {
		if _, sensitive := sensitiveBodyKeys[strings.ToLower(k)]; sensitive {
			values.Set(k, redactedValue)
		}
	}

	return base + "?" + values.Encode()
}

// RedactSecrets redacts password, apiKey, key, token, and secret fields before a body reaches the
// debug log. The search recurses because every secret this SDK sends is nested at least one level
// deep, e.g. {login: {password}}. Non-JSON bodies are returned unchanged.
func redactSecrets(body []byte) []byte {
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return body
	}

	redacted, err := json.Marshal(redactValue(parsed))
	if err != nil {
		return body
	}

	return redacted
}

// ParseTag splits a tag body like `name,omitempty,other` into its name and its option set.
func parseTag(tag string) (string, map[string]bool) {
	name := ""
	opts := map[string]bool{}
	if tag == "" {
		return name, opts
	}

	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, p := range parts[1:] {
		if p != "" {
			opts[p] = true
		}
	}

	return name, opts
}

// WithPageToken appends pageToken to a query string that may or may not already have parameters.
func withPageToken(query string, pageToken string) string {
	if pageToken == "" {
		return query
	}
	if query != "" {
		query += "&"
	}
	return query + "pageToken=" + url.QueryEscape(pageToken)
}

func structToQueryParams(data any) string {
	values := url.Values{}
	addStructQueryParams(values, data)

	return values.Encode()
}

// AddStructQueryParams adds the query-tagged fields of data to values. Embedded structs add into
// the same set, so a field encodes identically whether declared directly or promoted. A nil
// pointer or non-struct adds nothing.
func addStructQueryParams(values url.Values, data any) {
	v, ok := structValue(data)
	if !ok {
		return
	}

	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		value := v.Field(i)

		if field.Anonymous && value.Kind() == reflect.Struct {
			addStructQueryParams(values, value.Interface())
			continue
		}

		qName, _ := parseTag(field.Tag.Get("query"))
		if qName == "" || qName == "-" {
			continue
		}

		setQueryValue(values, qName, value, fieldOmitsEmpty(field))
	}
}

// StructValue dereferences data to a struct, reporting false for a nil pointer or a non-struct.
func structValue(data any) (reflect.Value, bool) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}, false
		}

		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	return v, true
}

// FieldOmitsEmpty reports whether either the json or query tag of field carries
// omitempty or omitzero.
func fieldOmitsEmpty(field reflect.StructField) bool {
	_, jsonOpts := parseTag(field.Tag.Get("json"))
	_, queryOpts := parseTag(field.Tag.Get("query"))

	return jsonOpts["omitempty"] || jsonOpts["omitzero"] || queryOpts["omitempty"] || queryOpts["omitzero"]
}

// SetQueryValue encodes value under qName, skipping it when omitEmpty applies. Slices join with
// commas rather than repeating the key, matching explode: false semantics, and empty ones are
// always skipped.
func setQueryValue(values url.Values, qName string, value reflect.Value, omitEmpty bool) {
	if value.Kind() == reflect.Slice {
		if value.IsNil() || value.Len() == 0 {
			return
		}

		parts := make([]string, value.Len())
		for j := range value.Len() {
			parts[j] = fmt.Sprintf("%v", value.Index(j).Interface())
		}

		values.Set(qName, strings.Join(parts, ","))

		return
	}

	switch val := value.Interface().(type) {
	case time.Time:
		if omitEmpty && val.IsZero() {
			return
		}

		values.Set(qName, val.Format(time.RFC3339))
	default:
		if omitEmpty && value.IsZero() {
			return
		}

		values.Set(qName, fmt.Sprintf("%v", val))
	}
}

func unmarshalResBody(body []byte, encodingType string) ([]byte, error) {
	if len(body) == 0 || encodingType == "" {
		return body, nil
	}

	if encodingType != EncodingTypeGZIP {
		return body, nil
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("decode gzip content: %w", err)
	}
	defer gzipReader.Close()

	resBody, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("read decompressed gzip body: %w", err)
	}

	return resBody, nil
}
