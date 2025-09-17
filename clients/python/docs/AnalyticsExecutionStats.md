# AnalyticsExecutionStats


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**average_latency_ms** | **float** |  | [optional] 
**failed_executions** | **int** |  | [optional] 
**success_rate** | **float** |  | [optional] 
**successful_executions** | **int** |  | [optional] 
**total_executions** | **int** |  | [optional] 

## Example

```python
from rule_engine_client.models.analytics_execution_stats import AnalyticsExecutionStats

# TODO update the JSON string below
json = "{}"
# create an instance of AnalyticsExecutionStats from a JSON string
analytics_execution_stats_instance = AnalyticsExecutionStats.from_json(json)
# print the JSON string representation of the object
print(AnalyticsExecutionStats.to_json())

# convert the object into a dict
analytics_execution_stats_dict = analytics_execution_stats_instance.to_dict()
# create an instance of AnalyticsExecutionStats from a dict
analytics_execution_stats_from_dict = AnalyticsExecutionStats.from_dict(analytics_execution_stats_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


