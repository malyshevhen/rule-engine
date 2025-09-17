# AnalyticsRuleStats


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**average_latency_ms** | **float** |  | [optional] 
**failed_executions** | **int** |  | [optional] 
**last_executed** | **str** |  | [optional] 
**rule_id** | **str** |  | [optional] 
**rule_name** | **str** |  | [optional] 
**success_rate** | **float** |  | [optional] 
**successful_executions** | **int** |  | [optional] 
**total_executions** | **int** |  | [optional] 

## Example

```python
from rule_engine_client.models.analytics_rule_stats import AnalyticsRuleStats

# TODO update the JSON string below
json = "{}"
# create an instance of AnalyticsRuleStats from a JSON string
analytics_rule_stats_instance = AnalyticsRuleStats.from_json(json)
# print the JSON string representation of the object
print(AnalyticsRuleStats.to_json())

# convert the object into a dict
analytics_rule_stats_dict = analytics_rule_stats_instance.to_dict()
# create an instance of AnalyticsRuleStats from a dict
analytics_rule_stats_from_dict = AnalyticsRuleStats.from_dict(analytics_rule_stats_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


