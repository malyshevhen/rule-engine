package modules

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

// TimeModule provides functions that Lua scripts can call
type TimeModule struct {
}

// NewTimeModule creates a new TimeModule
func NewTimeModule() *TimeModule {
	return &TimeModule{}
}

// Name returns the name of the module
func (s *TimeModule) Name() string {
	return "time"
}

// GetCurrentTime accepts the optional Go-sryle formatter, and returns the current time string
// If no formatter is provided, the default Go time format is used
func (s *TimeModule) GetCurrentTime(L *lua.LState) int {
	format := L.ToString(1)
	now := time.Now()
	if format == "" {
		L.Push(lua.LString(now.String()))
	} else {
		L.Push(lua.LString(now.Format(format)))
	}
	return 1
}

// Loader loads the Time module into the Lua state
func (s *TimeModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"now": s.GetCurrentTime,
	})

	fields := map[string]lua.LValue{
		"ANSIC":       lua.LString(time.ANSIC),
		"UNIX":        lua.LString(time.UnixDate),
		"RubyDate":    lua.LString(time.RubyDate),
		"RFC822":      lua.LString(time.RFC822),
		"RFC822Z":     lua.LString(time.RFC822Z),
		"RFC850":      lua.LString(time.RFC850),
		"RFC1123":     lua.LString(time.RFC1123),
		"RFC1123Z":    lua.LString(time.RFC1123Z),
		"RFC3339":     lua.LString(time.RFC3339),
		"RFC3339Nano": lua.LString(time.RFC3339Nano),
		"Kitchen":     lua.LString(time.Kitchen),
		"Stamp":       lua.LString(time.Stamp),
		"StampMilli":  lua.LString(time.StampMilli),
		"StampMicro":  lua.LString(time.StampMicro),
		"StampNano":   lua.LString(time.StampNano),
		"DateTime":    lua.LString(time.DateTime),
		"DateOnly":    lua.LString(time.DateOnly),
		"TimeOnly":    lua.LString(time.TimeOnly),
	}

	for name, field := range fields {
		L.SetField(mod, name, field)
	}

	L.Push(mod)
	return 1
}
