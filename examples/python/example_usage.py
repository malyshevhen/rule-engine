#!/usr/bin/env python3
"""
Example usage of the Rule Engine Python client.

This example demonstrates how to use the generated Python client to interact
with the Rule Engine API for creating rules, triggers, and actions.
"""

import sys
import os

# Add the client to the Python path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '../../clients/python'))

from rule_engine_client import ApiClient, Configuration
from rule_engine_client.api import RulesApi, ActionsApi, TriggersApi, AnalyticsApi
from rule_engine_client.models import (
    ApiCreateRuleRequest, ApiCreateActionRequest, ApiCreateTriggerRequest,
    ApiUpdateRuleRequest
)


def main():
    # Create configuration
    config = Configuration()
    config.host = "http://localhost:8080/api/v1"

    # Set API key for authentication
    config.api_key['Authorization'] = 'ApiKey your-api-key-here'
    config.api_key_prefix['Authorization'] = ''

    # Create API client
    api_client = ApiClient(config)

    # Create API instances
    rules_api = RulesApi(api_client)
    actions_api = ActionsApi(api_client)
    triggers_api = TriggersApi(api_client)
    analytics_api = AnalyticsApi(api_client)

    try:
        print("🚀 Rule Engine Python Client Example")
        print("=" * 40)

        # Example 1: Create a new rule
        print("\n📝 Creating a new rule...")
        create_rule_request = ApiCreateRuleRequest(
            name="Temperature Alert Rule",
            lua_script="if event.temperature > 25 then return true end",
            priority=0,
            enabled=True
        )

        rule = rules_api.rules_post(create_rule_request)
        print(f"✅ Created rule: {rule.name} (ID: {rule.id})")

        # Example 2: List all rules
        print("\n📋 Listing all rules...")
        rules = rules_api.rules_get()
        print(f"📊 Found {len(rules)} rules:")
        for r in rules:
            print(f"   • {r.name} (ID: {r.id})")

        # Example 3: Get analytics dashboard
        print("\n📊 Getting analytics dashboard...")
        dashboard = analytics_api.analytics_dashboard_get(time_range="24h")
        print(f"⏰ Dashboard time range: {dashboard.time_range}")
        print(f"📈 Total executions: {dashboard.overall_stats.total_executions}")
        print(".2f"
        # Example 4: Create an action
        print("\n⚡ Creating an action...")
        create_action_request = ApiCreateActionRequest(
            lua_script="log_message('info', 'Temperature alert triggered')",
            enabled=True
        )

        action = actions_api.actions_post(create_action_request)
        print(f"✅ Created action (ID: {action.id})")

        # Example 5: Create a trigger
        print("\n🎯 Creating a trigger...")
        create_trigger_request = ApiCreateTriggerRequest(
            rule_id=rule.id,
            type="CONDITIONAL",
            condition_script="if event.device_id == 'sensor_1' then return true end",
            enabled=True
        )

        trigger = triggers_api.triggers_post(create_trigger_request)
        print(f"✅ Created trigger (ID: {trigger.id})")

        # Example 6: Update a rule
        print("\n✏️  Updating the rule...")
        update_rule_request = ApiUpdateRuleRequest(
            name="Updated Temperature Alert Rule",
            priority=5
        )

        updated_rule = rules_api.rules_id_put(rule.id, update_rule_request)
        print(f"✅ Updated rule name: {updated_rule.name}")

        # Example 7: Get specific rule with details
        print("\n🔍 Getting rule details...")
        detailed_rule = rules_api.rules_id_get(rule.id)
        print(f"📋 Rule: {detailed_rule.name}")
        print(f"   Script: {detailed_rule.lua_script}")
        print(f"   Priority: {detailed_rule.priority}")
        print(f"   Enabled: {detailed_rule.enabled}")
        print(f"   Actions: {len(detailed_rule.actions) if detailed_rule.actions else 0}")
        print(f"   Triggers: {len(detailed_rule.triggers) if detailed_rule.triggers else 0}")

        # Example 8: Delete the rule (cleanup)
        print("\n🗑️  Deleting the rule...")
        rules_api.rules_id_delete(rule.id)
        print("✅ Rule deleted successfully")

        print("\n🎉 All examples completed successfully!")

    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()