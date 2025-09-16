package executor

import (
	"context"
	"testing"
	"time"

	execCtx "github.com/malyshevhen/rule-engine/internal/engine/executor/context"
	"github.com/stretchr/testify/assert"
)

func TestExecutorService_ExecuteScript_Success(t *testing.T) {
	ctxSvc := execCtx.NewService()
	svc := NewService(ctxSvc)

	ctx := ctxSvc.CreateContext("test-rule", "test-trigger")

	result := svc.ExecuteScript(context.Background(), "return 'success'", ctx)

	assert.True(t, result.Success)
	assert.NotEmpty(t, result.Output)
	assert.Empty(t, result.Error)
	assert.Greater(t, result.Duration, time.Duration(0))
}

func TestExecutorService_ExecuteScript_Error(t *testing.T) {
	ctxSvc := execCtx.NewService()
	svc := NewService(ctxSvc)

	ctx := ctxSvc.CreateContext("test-rule", "test-trigger")

	result := svc.ExecuteScript(context.Background(), "invalid lua syntax", ctx)

	assert.False(t, result.Success)
	assert.NotEmpty(t, result.Error)
	assert.Nil(t, result.Output)
	assert.Greater(t, result.Duration, time.Duration(0))
}

func TestExecutorService_ExecuteScript_Sandboxing(t *testing.T) {
	ctxSvc := execCtx.NewService()
	svc := NewService(ctxSvc)

	ctx := ctxSvc.CreateContext("test-rule", "test-trigger")

	// Test that io is disabled
	result := svc.ExecuteScript(context.Background(), "io.write('test')", ctx)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "attempt to index")

	// Test that os is disabled
	result = svc.ExecuteScript(context.Background(), "os.exit()", ctx)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "attempt to index")
}

func TestExecutorService_ExecuteScript_WithContext(t *testing.T) {
	ctxSvc := execCtx.NewService()
	svc := NewService(ctxSvc)

	ctx := ctxSvc.CreateContext("rule-123", "trigger-456")

	result := svc.ExecuteScript(context.Background(), "return rule_id", ctx)

	assert.True(t, result.Success)
	assert.Equal(t, "rule-123", result.Output[len(result.Output)-1])
}
