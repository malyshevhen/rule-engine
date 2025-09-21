# Rule Engine API Client

A production-ready Go client library for the Rule Engine API. This client provides a convenient interface to interact with all endpoints of the Rule Engine microservice.

## Features

- **Complete API Coverage**: Supports all endpoints defined in the OpenAPI specification
- **Authentication**: Supports both Bearer token and API key authentication
- **Type Safety**: Strongly typed request/response structures
- **Error Handling**: Comprehensive error handling with custom error types
- **Context Support**: Full context support for request cancellation and timeouts
- **Pagination**: Built-in support for paginated responses
- **JSON Patch**: Support for partial updates using JSON Patch operations

## Installation

```bash
go get github.com/malyshevhen/rule-engine/client
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/malyshevhen/rule-engine/client"
)

func main() {
    // Create client with API key authentication
    c := client.NewClient("http://localhost:8080", client.AuthConfig{
        APIKey: "your-api-key-here",
    })

    // Or with Bearer token authentication
    c := client.NewClient("http://localhost:8080", client.AuthConfig{
        BearerToken: "your-jwt-token-here",
    })

    ctx := context.Background()

    // Check service health
    health, err := c.Health(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Database: %s, Redis: %s", health.Database, health.Redis)
}
```

## Authentication

The client supports two authentication methods:

### API Key Authentication

```go
auth := client.AuthConfig{
    APIKey: "your-api-key",
}
c := client.NewClient("http://localhost:8080", auth)
```

### Bearer Token Authentication

```go
auth := client.AuthConfig{
    BearerToken: "your-jwt-token",
}
c := client.NewClient("http://localhost:8080", auth)
```

You can also use both methods simultaneously if required.

## API Methods

### Health & Monitoring

#### Health Check

```go
health, err := client.Health(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Database: %s, Redis: %s\n", health.Database, health.Redis)
```

#### Metrics

```go
metrics, err := client.Metrics(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println(metrics) // Prometheus format metrics
```

### Rules

#### Create a Rule

```go
rule, err := client.CreateRule(ctx, client.CreateRuleRequest{
    Name:      "Temperature Alert",
    LuaScript: "if event.temperature > 25 then return true end",
    Priority:  &[]int{1}[0], // optional
    Enabled:   &[]bool{true}[0], // optional
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created rule: %s\n", rule.ID)
```

#### Get a Rule

```go
rule, err := client.GetRule(ctx, ruleID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Rule: %+v\n", rule)
```

#### List Rules

```go
rules, err := client.ListRules(ctx, 50, 0) // limit=50, offset=0
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d rules out of %d total\n", len(rules.Rules), rules.Total)
for _, rule := range rules.Rules {
    fmt.Printf("- %s: %s\n", rule.ID, rule.Name)
}
```

#### Update a Rule (JSON Patch)

```go
patches := client.PatchRequest{
    {Op: "replace", Path: "/name", Value: "Updated Rule Name"},
    {Op: "replace", Path: "/enabled", Value: false},
}

updatedRule, err := client.UpdateRule(ctx, ruleID, patches)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Updated rule: %s\n", updatedRule.Name)
```

#### Delete a Rule

```go
err := client.DeleteRule(ctx, ruleID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Rule deleted")
```

#### Add Action to Rule

```go
err := client.AddActionToRule(ctx, ruleID, client.AddActionToRuleRequest{
    ActionID: actionID,
})
if err != nil {
    log.Fatal(err)
}
fmt.Println("Action added to rule")
```

### Triggers

#### Create a Trigger

```go
trigger, err := client.CreateTrigger(ctx, client.CreateTriggerRequest{
    RuleID:          ruleID,
    Type:            "CONDITIONAL",
    ConditionScript: "if event.device_id == 'sensor_1' then return true end",
    Enabled:         &[]bool{true}[0], // optional
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created trigger: %s\n", trigger.ID)
```

#### Get a Trigger

```go
trigger, err := client.GetTrigger(ctx, triggerID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Trigger: %+v\n", trigger)
```

#### List Triggers

```go
triggers, err := client.ListTriggers(ctx, 50, 0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d triggers\n", len(triggers.Triggers))
```

#### Delete a Trigger

```go
err := client.DeleteTrigger(ctx, triggerID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Trigger deleted")
```

### Actions

#### Create an Action

```go
action, err := client.CreateAction(ctx, client.CreateActionRequest{
    Name:      "Send Alert",
    LuaScript: "log_message('info', 'Temperature alert triggered')",
    Enabled:   &[]bool{true}[0], // optional
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created action: %s\n", action.ID)
```

#### Get an Action

```go
action, err := client.GetAction(ctx, actionID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Action: %+v\n", action)
```

#### List Actions

```go
actions, err := client.ListActions(ctx, 50, 0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d actions\n", len(actions.Actions))
```

#### Delete an Action

```go
err := client.DeleteAction(ctx, actionID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Action deleted")
```

### Script Evaluation

#### Evaluate Lua Script

```go
result, err := client.EvaluateScript(ctx, client.EvaluateScriptRequest{
    Script: "return 2 + 3",
    Context: map[string]interface{}{
        "temperature": 25,
        "device_id":   "sensor_1",
    },
})
if err != nil {
    log.Fatal(err)
}

if result.Success {
    fmt.Printf("Result: %v\n", result.Result)
    fmt.Printf("Output: %v\n", result.Output)
    fmt.Printf("Duration: %s\n", result.Duration)
} else {
    fmt.Printf("Error: %s\n", result.Error)
}
```

## Error Handling

The client returns custom error types for API errors:

```go
rule, err := client.GetRule(ctx, ruleID)
if err != nil {
    var apiErr *client.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error: %s (%d) - %s\n",
            apiErr.Code, apiErr.StatusCode, apiErr.Message)
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

## Custom HTTP Client

You can provide your own HTTP client for custom configuration:

```go
httpClient := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    },
}

c := client.NewClientWithHTTPClient("https://api.example.com", auth, httpClient)
```

## Context Support

All methods accept a context for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

rule, err := client.GetRule(ctx, ruleID)
// Request will timeout after 5 seconds
```

## Examples

See the `examples/` directory for complete usage examples.

## License

This project is licensed under the MIT License.