package executor

import (
	"context"
	"testing"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/malyshevhen/rule-engine/internal/engine/executor/platform"
)

// BenchmarkLuaExecution_SimpleArithmetic benchmarks simple arithmetic operations
func BenchmarkLuaExecution_SimpleArithmetic(b *testing.B) {
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := NewService(contextSvc, platformSvc)

	execContext := &execCtx.ExecutionContext{
		RuleID:    "benchmark-rule",
		TriggerID: "benchmark-trigger",
		Data: map[string]any{
			"x": 10,
			"y": 20,
		},
	}

	script := "return x + y"

	for b.Loop() {
		result := executorSvc.ExecuteScript(context.Background(), script, execContext)
		if !result.Success {
			b.Fatalf("Script execution failed: %s", result.Error)
		}
		if result.Output[0] != 30.0 {
			b.Fatalf("Expected 30.0, got %v", result.Output[0])
		}
	}
}

// BenchmarkLuaExecution_ComplexLogic benchmarks complex conditional logic
func BenchmarkLuaExecution_ComplexLogic(b *testing.B) {
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := NewService(contextSvc, platformSvc)

	execContext := &execCtx.ExecutionContext{
		RuleID:    "benchmark-rule",
		TriggerID: "benchmark-trigger",
		Data: map[string]any{
			"temperature": 25.5,
			"humidity":    60,
			"threshold":   20,
		},
	}

	script := `
		if temperature > threshold then
			if humidity > 50 then
				return temperature * 2
			else
				return temperature + 10
			end
		else
			return temperature
		end
	`

	for b.Loop() {
		result := executorSvc.ExecuteScript(context.Background(), script, execContext)
		if !result.Success {
			b.Fatalf("Script execution failed: %s", result.Error)
		}
		if result.Output[0] != 51.0 {
			b.Fatalf("Expected 51.0, got %v", result.Output[0])
		}
	}
}

// BenchmarkLuaExecution_PlatformAPICalls benchmarks scripts with platform API calls
func BenchmarkLuaExecution_PlatformAPICalls(b *testing.B) {
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := NewService(contextSvc, platformSvc)

	execContext := &execCtx.ExecutionContext{
		RuleID:    "benchmark-rule",
		TriggerID: "benchmark-trigger",
		Data: map[string]any{
			"device_id": "device123",
		},
	}

	script := `
		log_message("Benchmarking platform API")
		local device_data, err = get_device_state(device_id)
		if err ~= nil then
			return false
		end
		store_data("benchmark_key", device_data)
		return device_data.online
	`

	for b.Loop() {
		result := executorSvc.ExecuteScript(context.Background(), script, execContext)
		if !result.Success {
			b.Fatalf("Script execution failed: %s", result.Error)
		}
		if result.Output[0] != true {
			b.Fatalf("Expected true, got %v", result.Output[0])
		}
	}
}

// BenchmarkLuaExecution_LargeScript benchmarks execution of larger scripts
func BenchmarkLuaExecution_LargeScript(b *testing.B) {
	contextSvc := execCtx.NewService()
	platformSvc := platform.NewService()
	executorSvc := NewService(contextSvc, platformSvc)

	execContext := &execCtx.ExecutionContext{
		RuleID:    "benchmark-rule",
		TriggerID: "benchmark-trigger",
		Data: map[string]any{
			"values": []any{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
	}

	script := `
		local sum = 0
		local max_val = 0
		for i, v in ipairs(values) do
			sum = sum + v
			if v > max_val then
				max_val = v
			end
		end
		local avg = sum / #values
		return {sum = sum, max = max_val, avg = avg}
	`

	for b.Loop() {
		result := executorSvc.ExecuteScript(context.Background(), script, execContext)
		if !result.Success {
			b.Fatalf("Script execution failed: %s", result.Error)
		}
		// Verify the result structure
		resultMap := result.Output[0].(map[string]any)
		if resultMap["sum"] != 55 {
			b.Fatalf("Expected sum 55, got %v", resultMap["sum"])
		}
	}
}

// BenchmarkLuaExecution_Concurrent benchmarks concurrent script execution
func BenchmarkLuaExecution_Concurrent(b *testing.B) {
	script := "return math.sin(x) + math.cos(y)"

	b.RunParallel(func(pb *testing.PB) {
		localExecutor := NewService(execCtx.NewService(), platform.NewService())
		execContext := &execCtx.ExecutionContext{
			RuleID:    "concurrent-rule",
			TriggerID: "concurrent-trigger",
			Data: map[string]any{
				"x": 1.5,
				"y": 2.3,
			},
		}

		for pb.Next() {
			result := localExecutor.ExecuteScript(context.Background(), script, execContext)
			if !result.Success {
				b.Fatalf("Concurrent script execution failed: %s", result.Error)
			}
		}
	})
}
