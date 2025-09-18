package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// HTTPMethod represents an HTTP method
type HTTPMethod string

// HTTP methods enum
const (
	HTTPMethodGet    HTTPMethod = "GET"
	HTTPMethodPost   HTTPMethod = "POST"
	HTTPMethodDelete HTTPMethod = "DELETE"
	HTTPMethodPatch  HTTPMethod = "PATCH"
	HTTPMethodPut    HTTPMethod = "PUT"
)

// HTTPModule provides functions that Lua scripts can call
type HTTPModule struct {
	client *http.Client
}

// HTTPModuleOption allows to configure the HTTP module
type HTTPModuleOption func(hm *HTTPModule) *HTTPModule

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) HTTPModuleOption {
	return func(hm *HTTPModule) *HTTPModule {
		hm.client = client
		return hm
	}
}

// NewHTTPModule creates a new HTTP API service
func NewHTTPModule(opts ...HTTPModuleOption) *HTTPModule {
	hm := &HTTPModule{client: &http.Client{Timeout: 5 * time.Second}}
	for _, opt := range opts {
		opt(hm)
	}
	return hm
}

// Name returns the name of the module
func (s *HTTPModule) Name() string {
	return "http"
}

// Get makes an HTTP GET request
func (s *HTTPModule) Get(L *lua.LState) int {
	url := L.ToString(1)
	headersTable := L.ToTable(2)

	headers := make(map[string]string)
	if headersTable != nil {
		headersTable.ForEach(func(k, v lua.LValue) {
			headers[k.String()] = v.String()
		})
	}

	ctx := context.Background()

	result, err := s.MakeHTTPRequest(ctx, HTTPMethodGet, url, headers, "")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for k, v := range result {
		L.SetField(table, k, luaValueFromGo(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// Post makes an HTTP POST request
func (s *HTTPModule) Post(L *lua.LState) int {
	url := L.ToString(1)
	headersTable := L.ToTable(2)
	body := L.ToString(3)

	headers := make(map[string]string)
	if headersTable != nil {
		headersTable.ForEach(func(k, v lua.LValue) {
			headers[k.String()] = v.String()
		})
	}

	ctx := context.Background()

	result, err := s.MakeHTTPRequest(ctx, HTTPMethodPost, url, headers, body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for k, v := range result {
		L.SetField(table, k, luaValueFromGo(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// Delete makes an HTTP DELETE request
func (s *HTTPModule) Delete(L *lua.LState) int {
	url := L.ToString(1)
	headersTable := L.ToTable(2)

	headers := make(map[string]string)
	if headersTable != nil {
		headersTable.ForEach(func(k, v lua.LValue) {
			headers[k.String()] = v.String()
		})
	}

	ctx := context.Background()

	result, err := s.MakeHTTPRequest(ctx, HTTPMethodDelete, url, headers, "")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for k, v := range result {
		L.SetField(table, k, luaValueFromGo(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// Put makes an HTTP PUT request
func (s *HTTPModule) Put(L *lua.LState) int {
	url := L.ToString(1)
	headersTable := L.ToTable(2)
	body := L.ToString(3)

	headers := make(map[string]string)
	if headersTable != nil {
		headersTable.ForEach(func(k, v lua.LValue) {
			headers[k.String()] = v.String()
		})
	}

	ctx := context.Background()

	result, err := s.MakeHTTPRequest(ctx, HTTPMethodPut, url, headers, body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for k, v := range result {
		L.SetField(table, k, luaValueFromGo(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// Patch makes an HTTP PATCH request
func (s *HTTPModule) Patch(L *lua.LState) int {
	url := L.ToString(1)
	headersTable := L.ToTable(2)
	body := L.ToString(3)

	headers := make(map[string]string)
	if headersTable != nil {
		headersTable.ForEach(func(k, v lua.LValue) {
			headers[k.String()] = v.String()
		})
	}

	ctx := context.Background()

	result, err := s.MakeHTTPRequest(ctx, HTTPMethodPatch, url, headers, body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	table := L.NewTable()
	for k, v := range result {
		L.SetField(table, k, luaValueFromGo(L, v))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

func (s *HTTPModule) MakeHTTPRequest(
	ctx context.Context,
	method HTTPMethod,
	url string,
	headers map[string]string,
	body string,
) (map[string]any, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, string(method), url, bodyReader)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"status": resp.StatusCode,
		"body":   string(respBody),
	}, nil
}

// Loader loads the HTTP module into the Lua state
func (s *HTTPModule) Loader(L *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"get":    s.Get,
		"post":   s.Post,
		"delete": s.Delete,
		"put":    s.Put,
		"patch":  s.Patch,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

// luaValueFromGo converts a Go value to a Lua value
func luaValueFromGo(L *lua.LState, v any) lua.LValue {
	switch val := v.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(val)
	case int:
		return lua.LNumber(val)
	case int32:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case float32:
		return lua.LNumber(val)
	case float64:
		return lua.LNumber(val)
	case string:
		return lua.LString(val)
	case map[string]any:
		table := L.NewTable()
		for k, v := range val {
			L.SetField(table, k, luaValueFromGo(L, v))
		}
		return table
	case []any:
		table := L.NewTable()
		for i, v := range val {
			table.Insert(i+1, luaValueFromGo(L, v))
		}
		return table
	default:
		// For complex types, convert to JSON string
		if jsonBytes, err := json.Marshal(v); err == nil {
			return lua.LString(jsonBytes)
		}
		return lua.LString(fmt.Sprintf("%v", v))
	}
}
