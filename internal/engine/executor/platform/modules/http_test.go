package modules

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestNewHTTPModule(t *testing.T) {
	mod := NewHTTPModule()
	if mod == nil {
		t.Fatal("NewHTTPModule returned nil")
	}
	if mod.client == nil {
		t.Fatal("HTTP client not set")
	}
}

func TestNewHTTPModuleWithCustomClient(t *testing.T) {
	client := &http.Client{}
	mod := NewHTTPModule(WithHTTPClient(client))
	if mod.client != client {
		t.Fatal("Custom client not set")
	}
}

func TestMakeHTTPRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected /test, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Test") != "value" {
			t.Errorf("Expected X-Test header, got %s", r.Header.Get("X-Test"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	mod := NewHTTPModule()

	result, err := mod.MakeHTTPRequest(context.Background(), HTTPMethodGet, server.URL+"/test", map[string]string{"X-Test": "value"}, "")
	if err != nil {
		t.Fatalf("MakeHTTPRequest failed: %v", err)
	}

	if result["status"] != 200 {
		t.Errorf("Expected status 200, got %v", result["status"])
	}

	if result["body"] != `{"message": "success"}` {
		t.Errorf("Expected body '{\"message\": \"success\"}', got %v", result["body"])
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response body"))
	}))
	defer server.Close()

	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	// Load the module
	mod.Loader(L)
	httpTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("http", httpTable)

	// Run script
	script := fmt.Sprintf(`
		local res, err = http.get("%s")
		return res, err
	`, server.URL)

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	// Check results
	if L.GetTop() != 2 {
		t.Fatalf("Expected 2 return values, got %d", L.GetTop())
	}

	resultTable := L.ToTable(1)
	if resultTable == nil {
		t.Fatal("First return value is not a table")
	}

	status := L.RawGet(resultTable, lua.LString("status"))
	if status.Type() != lua.LTNumber || int(status.(lua.LNumber)) != 200 {
		t.Errorf("Expected status 200, got %v", status)
	}

	body := L.RawGet(resultTable, lua.LString("body"))
	if body.Type() != lua.LTString || string(body.(lua.LString)) != "response body" {
		t.Errorf("Expected body 'response body', got %v", body)
	}

	errValue := L.ToString(2)
	if errValue != "" {
		t.Errorf("Expected no error, got %s", errValue)
	}
}

func TestPost(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	}))
	defer server.Close()

	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	httpTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("http", httpTable)

	script := `
		local headers = {["Content-Type"] = "application/json"}
		local res, err = http.post("` + server.URL + `", headers, '{"key": "value"}')
		return res, err
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if receivedBody != `{"key": "value"}` {
		t.Errorf("Expected body '{\"key\": \"value\"}', got '%s'", receivedBody)
	}

	resultTable := L.ToTable(1)
	status := L.RawGet(resultTable, lua.LString("status"))
	if int(status.(lua.LNumber)) != 201 {
		t.Errorf("Expected status 201, got %v", status)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	httpTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("http", httpTable)

	script := `
		local res, err = http.delete("` + server.URL + `")
		return res, err
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	resultTable := L.ToTable(1)
	status := L.RawGet(resultTable, lua.LString("status"))
	if int(status.(lua.LNumber)) != 204 {
		t.Errorf("Expected status 204, got %v", status)
	}
}

func TestPut(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("updated"))
	}))
	defer server.Close()

	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	httpTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("http", httpTable)

	script := `
		local headers = {}
		local res, err = http.put("` + server.URL + `", headers, "new content")
		return res, err
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if receivedBody != "new content" {
		t.Errorf("Expected body 'new content', got '%s'", receivedBody)
	}
}

func TestPatch(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH, got %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("patched"))
	}))
	defer server.Close()

	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	httpTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("http", httpTable)

	script := `
		local res, err = http.patch("` + server.URL + `", {}, "patch data")
		return res, err
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if receivedBody != "patch data" {
		t.Errorf("Expected body 'patch data', got '%s'", receivedBody)
	}
}

func TestHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	httpTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("http", httpTable)

	script := `
		local res, err = http.get("` + server.URL + `")
		return res, err
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	resultTable := L.ToTable(1)
	if resultTable == nil {
		t.Fatal("Expected table even on error")
	}

	status := L.RawGet(resultTable, lua.LString("status"))
	if int(status.(lua.LNumber)) != 500 {
		t.Errorf("Expected status 500, got %v", status)
	}

	body := L.RawGet(resultTable, lua.LString("body"))
	if string(body.(lua.LString)) != "server error" {
		t.Errorf("Expected body 'server error', got %v", body)
	}

	errValue := L.ToString(2)
	if errValue != "" {
		t.Errorf("Expected no error string, got %s", errValue)
	}
}

func TestHTTPModuleLoader(t *testing.T) {
	mod := NewHTTPModule()
	L := lua.NewState()
	defer L.Close()

	result := mod.Loader(L)
	if result != 1 {
		t.Errorf("Loader returned %d, expected 1", result)
	}

	if L.GetTop() != 1 {
		t.Fatalf("Stack top is %d, expected 1", L.GetTop())
	}

	table := L.ToTable(1)
	if table == nil {
		t.Fatal("Loader did not push a table")
	}

	expectedFuncs := []string{"get", "post", "delete", "put", "patch"}
	for _, funcName := range expectedFuncs {
		val := L.RawGet(table, lua.LString(funcName))
		if val.Type() != lua.LTFunction {
			t.Errorf("Function %s not found in table or not a function", funcName)
		}
	}
}
