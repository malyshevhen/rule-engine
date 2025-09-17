package main

import (
	"context"
	"fmt"
	"log"

	ruleengine "github.com/malyshevhen/rule-engine/clients/go"
)

func main() {
	// Create a new API client
	config := ruleengine.NewConfiguration()
	config.Host = "localhost:8080"
	config.Scheme = "http"

	// Set API key for authentication
	config.AddDefaultHeader("Authorization", "ApiKey your-api-key-here")

	client := ruleengine.NewAPIClient(config)

	ctx := context.Background()

	// Example 1: Create a new rule
	fmt.Println("Creating a new rule...")
	createRuleRequest := ruleengine.ApiCreateRuleRequest{
		Name:      "Temperature Alert Rule",
		LuaScript: "if event.temperature > 25 then return true end",
		Priority:  &[]int32{0}[0],
		Enabled:   &[]bool{true}[0],
	}

	rule, resp, err := client.RulesAPI.RulesPost(ctx).ApiCreateRuleRequest(createRuleRequest).Execute()
	if err != nil {
		log.Fatalf("Error creating rule: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Created rule: %s (ID: %s)\n", *rule.Name, rule.Id)

	// Example 2: List all rules
	fmt.Println("\nListing all rules...")
	rules, resp, err := client.RulesAPI.RulesGet(ctx).Execute()
	if err != nil {
		log.Fatalf("Error listing rules: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Found %d rules:\n", len(rules))
	for _, r := range rules {
		fmt.Printf("  - %s (ID: %s)\n", *r.Name, r.Id)
	}

	// Example 3: Get analytics dashboard
	fmt.Println("\nGetting analytics dashboard...")
	dashboard, resp, err := client.AnalyticsAPI.AnalyticsDashboardGet(ctx).TimeRange("24h").Execute()
	if err != nil {
		log.Fatalf("Error getting dashboard: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Dashboard time range: %s\n", *dashboard.TimeRange)
	fmt.Printf("Total executions: %d\n", *dashboard.OverallStats.TotalExecutions)
	fmt.Printf("Success rate: %.2f%%\n", *dashboard.OverallStats.SuccessRate)

	// Example 4: Create an action
	fmt.Println("\nCreating an action...")
	createActionRequest := ruleengine.ApiCreateActionRequest{
		LuaScript: "log_message('info', 'Temperature alert triggered')",
		Enabled:   &[]bool{true}[0],
	}

	action, resp, err := client.ActionsAPI.ActionsPost(ctx).ApiCreateActionRequest(createActionRequest).Execute()
	if err != nil {
		log.Fatalf("Error creating action: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Created action (ID: %s)\n", action.Id)

	// Example 5: Create a trigger
	fmt.Println("\nCreating a trigger...")
	createTriggerRequest := ruleengine.ApiCreateTriggerRequest{
		RuleId:          rule.Id,
		Type:            "CONDITIONAL",
		ConditionScript: "if event.device_id == 'sensor_1' then return true end",
		Enabled:         &[]bool{true}[0],
	}

	trigger, resp, err := client.TriggersAPI.TriggersPost(ctx).ApiCreateTriggerRequest(createTriggerRequest).Execute()
	if err != nil {
		log.Fatalf("Error creating trigger: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Created trigger (ID: %s)\n", trigger.Id)

	// Example 6: Update a rule
	fmt.Println("\nUpdating the rule...")
	updateRuleRequest := ruleengine.ApiUpdateRuleRequest{
		Name:     &[]string{"Updated Temperature Alert Rule"}[0],
		Priority: &[]int32{5}[0],
	}

	updatedRule, resp, err := client.RulesAPI.RulesIdPut(ctx, rule.Id).ApiUpdateRuleRequest(updateRuleRequest).Execute()
	if err != nil {
		log.Fatalf("Error updating rule: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Updated rule name: %s\n", *updatedRule.Name)

	// Example 7: Delete the rule (cleanup)
	fmt.Println("\nDeleting the rule...")
	resp, err = client.RulesAPI.RulesIdDelete(ctx, rule.Id).Execute()
	if err != nil {
		log.Fatalf("Error deleting rule: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Rule deleted successfully")
}
