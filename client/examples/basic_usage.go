package main

import (
	"context"
	"fmt"
	"log"

	"github.com/malyshevhen/rule-engine/client"
)

func main() {
	// Create a client with API key authentication
	c := client.NewClient("http://localhost:8080", client.AuthConfig{
		APIKey: "your-api-key-here", // Replace with your actual API key
	})

	ctx := context.Background()

	// Check service health
	fmt.Println("Checking service health...")
	health, err := c.Health(ctx)
	if err != nil {
		log.Fatalf("Failed to check health: %v", err)
	}
	fmt.Printf("Service health - Database: %s, Redis: %s\n\n", health.Database, health.Redis)

	// Evaluate a simple Lua script
	fmt.Println("Evaluating Lua script...")
	result, err := c.EvaluateScript(ctx, client.EvaluateScriptRequest{
		Script: "return 2 + 3",
		Context: map[string]interface{}{
			"temperature": 25,
			"device_id":   "sensor_1",
		},
	})
	if err != nil {
		log.Fatalf("Failed to evaluate script: %v", err)
	}

	if result.Success {
		fmt.Printf("Script executed successfully!\n")
		fmt.Printf("Result: %v\n", result.Result)
		fmt.Printf("Output: %v\n", result.Output)
		fmt.Printf("Duration: %s\n\n", result.Duration)
	} else {
		fmt.Printf("Script execution failed: %s\n\n", result.Error)
	}

	// Create a rule
	fmt.Println("Creating a rule...")
	rule, err := c.CreateRule(ctx, client.CreateRuleRequest{
		Name:      "Temperature Alert Rule",
		LuaScript: "if event.temperature > 25 then return true end",
		Priority:  &[]int{1}[0],
		Enabled:   &[]bool{true}[0],
	})
	if err != nil {
		log.Fatalf("Failed to create rule: %v", err)
	}
	fmt.Printf("Created rule: %s (%s)\n\n", rule.Name, rule.ID)

	// List rules
	fmt.Println("Listing rules...")
	rules, err := c.ListRules(ctx, 10, 0)
	if err != nil {
		log.Fatalf("Failed to list rules: %v", err)
	}
	fmt.Printf("Found %d rules out of %d total:\n", len(rules.Rules), rules.Total)
	for _, r := range rules.Rules {
		fmt.Printf("- %s: %s\n", r.ID, r.Name)
	}
	fmt.Println()

	// Create an action
	fmt.Println("Creating an action...")
	action, err := c.CreateAction(ctx, client.CreateActionRequest{
		Name:      "Send Temperature Alert",
		LuaScript: "log_message('info', 'Temperature alert triggered for device: ' .. event.device_id)",
		Enabled:   &[]bool{true}[0],
	})
	if err != nil {
		log.Fatalf("Failed to create action: %v", err)
	}
	fmt.Printf("Created action: %s (%s)\n\n", action.Name, action.ID)

	// Add action to rule
	fmt.Println("Adding action to rule...")
	err = c.AddActionToRule(ctx, rule.ID, client.AddActionToRuleRequest{
		ActionID: action.ID,
	})
	if err != nil {
		log.Fatalf("Failed to add action to rule: %v", err)
	}
	fmt.Println("Action added to rule successfully!\n")

	// Create a trigger
	fmt.Println("Creating a trigger...")
	trigger, err := c.CreateTrigger(ctx, client.CreateTriggerRequest{
		RuleID:          rule.ID,
		Type:            "CONDITIONAL",
		ConditionScript: "if event.device_id == 'sensor_1' then return true end",
		Enabled:         &[]bool{true}[0],
	})
	if err != nil {
		log.Fatalf("Failed to create trigger: %v", err)
	}
	fmt.Printf("Created trigger: %s (Type: %s)\n\n", trigger.ID, trigger.Type)

	// Get the complete rule with triggers and actions
	fmt.Println("Getting complete rule details...")
	completeRule, err := c.GetRule(ctx, rule.ID)
	if err != nil {
		log.Fatalf("Failed to get rule: %v", err)
	}
	fmt.Printf("Rule: %s\n", completeRule.Name)
	fmt.Printf("Script: %s\n", completeRule.LuaScript)
	fmt.Printf("Triggers: %d\n", len(completeRule.Triggers))
	fmt.Printf("Actions: %d\n", len(completeRule.Actions))
	fmt.Println()

	// Clean up - delete created resources
	fmt.Println("Cleaning up...")

	// Delete trigger
	err = c.DeleteTrigger(ctx, trigger.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete trigger: %v", err)
	} else {
		fmt.Println("Deleted trigger")
	}

	// Delete action
	err = c.DeleteAction(ctx, action.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete action: %v", err)
	} else {
		fmt.Println("Deleted action")
	}

	// Delete rule
	err = c.DeleteRule(ctx, rule.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete rule: %v", err)
	} else {
		fmt.Println("Deleted rule")
	}

	fmt.Println("\nExample completed successfully!")
}
