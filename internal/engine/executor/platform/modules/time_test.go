package modules

import (
	"fmt"
	"strings"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestNewTimeModule(t *testing.T) {
	mod := NewTimeModule()
	if mod == nil {
		t.Fatal("NewTimeModule returned nil")
	}
}

func TestGetCurrentTime_NoFormat(t *testing.T) {
	mod := NewTimeModule()
	L := lua.NewState()
	defer L.Close()

	result := mod.GetCurrentTime(L)
	if result != 1 {
		t.Errorf("GetCurrentTime returned %d, expected 1", result)
	}

	if L.GetTop() != 1 {
		t.Fatalf("Stack top is %d, expected 1", L.GetTop())
	}

	timeStr := L.ToString(1)
	if timeStr == "" {
		t.Error("Returned empty time string")
	}

	// Since default format is time.Time.String(), just check it's not empty
	// and contains current year
	now := time.Now()
	if !strings.Contains(timeStr, fmt.Sprintf("%d", now.Year())) {
		t.Errorf("Time string '%s' does not contain current year %d", timeStr, now.Year())
	}
}

func TestGetCurrentTime_WithFormat(t *testing.T) {
	mod := NewTimeModule()
	L := lua.NewState()
	defer L.Close()

	// Push format
	L.Push(lua.LString(time.RFC3339))

	result := mod.GetCurrentTime(L)
	if result != 1 {
		t.Errorf("GetCurrentTime returned %d, expected 1", result)
	}

	timeStr := L.ToString(-1)
	if timeStr == "" {
		t.Error("Returned empty time string")
	}

	// Parse with the format
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Errorf("Failed to parse time string '%s' with RFC3339: %v", timeStr, err)
	}

	// Check it's recent
	now := time.Now()
	if parsed.Sub(now) > time.Minute || now.Sub(parsed) > time.Minute {
		t.Errorf("Parsed time %v is not within 1 minute of now %v", parsed, now)
	}
}

func TestTimeModuleLoader(t *testing.T) {
	mod := NewTimeModule()
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

	// Check now function
	nowFunc := L.RawGet(table, lua.LString("now"))
	if nowFunc.Type() != lua.LTFunction {
		t.Error("now function not found in table")
	}

	// Check constants
	constants := map[string]string{
		"ANSIC":       time.ANSIC,
		"UNIX":        time.UnixDate,
		"RubyDate":    time.RubyDate,
		"RFC822":      time.RFC822,
		"RFC822Z":     time.RFC822Z,
		"RFC850":      time.RFC850,
		"RFC1123":     time.RFC1123,
		"RFC1123Z":    time.RFC1123Z,
		"RFC3339":     time.RFC3339,
		"RFC3339Nano": time.RFC3339Nano,
		"Kitchen":     time.Kitchen,
		"Stamp":       time.Stamp,
		"StampMilli":  time.StampMilli,
		"StampMicro":  time.StampMicro,
		"StampNano":   time.StampNano,
		"DateTime":    time.DateTime,
		"DateOnly":    time.DateOnly,
		"TimeOnly":    time.TimeOnly,
	}

	for name, expected := range constants {
		val := L.RawGet(table, lua.LString(name))
		if val.Type() != lua.LTString {
			t.Errorf("Constant %s not found or not a string", name)
			continue
		}
		if string(val.(lua.LString)) != expected {
			t.Errorf("Constant %s: expected '%s', got '%s'", name, expected, val)
		}
	}
}

func TestScriptExecution_Now(t *testing.T) {
	mod := NewTimeModule()
	L := lua.NewState()
	defer L.Close()

	// Load the module
	mod.Loader(L)
	timeTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("time", timeTable)

	// Run script
	script := `
		local current = time.now()
		return current
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if L.GetTop() != 1 {
		t.Fatalf("Expected 1 return value, got %d", L.GetTop())
	}

	timeStr := L.ToString(1)
	if timeStr == "" {
		t.Error("Script returned empty time string")
	}

	now := time.Now()
	if !strings.Contains(timeStr, fmt.Sprintf("%d", now.Year())) {
		t.Errorf("Time string '%s' does not contain current year %d", timeStr, now.Year())
	}
}

func TestScriptExecution_NowWithFormat(t *testing.T) {
	mod := NewTimeModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	timeTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("time", timeTable)

	// Run script with format
	script := `
		local formatted = time.now(time.RFC3339)
		return formatted
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	timeStr := L.ToString(1)
	if timeStr == "" {
		t.Error("Script returned empty time string")
	}

	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Errorf("Failed to parse RFC3339 time string '%s': %v", timeStr, err)
	}

	now := time.Now()
	if parsed.Sub(now) > time.Second || now.Sub(parsed) > time.Second {
		t.Errorf("Parsed time %v is not within 1 second of now %v", parsed, now)
	}
}

func TestScriptExecution_Constants(t *testing.T) {
	mod := NewTimeModule()
	L := lua.NewState()
	defer L.Close()

	mod.Loader(L)
	timeTable := L.ToTable(-1)
	L.Pop(1)
	L.SetGlobal("time", timeTable)

	// Run script using constants
	script := `
		local formatted = time.now(time.Kitchen)
		return formatted, time.Kitchen
	`

	err := L.DoString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if L.GetTop() != 2 {
		t.Fatalf("Expected 2 return values, got %d", L.GetTop())
	}

	timeStr := L.ToString(1)
	formatStr := L.ToString(2)

	if timeStr == "" {
		t.Error("Script returned empty time string")
	}

	if formatStr != time.Kitchen {
		t.Errorf("Expected format '%s', got '%s'", time.Kitchen, formatStr)
	}

	_, err = time.Parse(time.Kitchen, timeStr)
	if err != nil {
		t.Errorf("Failed to parse Kitchen time string '%s': %v", timeStr, err)
	}

	// Since Kitchen is time-only format, we can't check the date
	// Just ensure parsing succeeds
}
