# AnalyticsDashboardData


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**execution_trend** | [**AnalyticsTimeSeriesData**](AnalyticsTimeSeriesData.md) |  | [optional] 
**latency_trend** | [**AnalyticsTimeSeriesData**](AnalyticsTimeSeriesData.md) |  | [optional] 
**overall_stats** | [**AnalyticsExecutionStats**](AnalyticsExecutionStats.md) |  | [optional] 
**rule_stats** | [**List[AnalyticsRuleStats]**](AnalyticsRuleStats.md) |  | [optional] 
**success_rate_trend** | [**AnalyticsTimeSeriesData**](AnalyticsTimeSeriesData.md) |  | [optional] 
**time_range** | **str** |  | [optional] 

## Example

```python
from rule_engine_client.models.analytics_dashboard_data import AnalyticsDashboardData

# TODO update the JSON string below
json = "{}"
# create an instance of AnalyticsDashboardData from a JSON string
analytics_dashboard_data_instance = AnalyticsDashboardData.from_json(json)
# print the JSON string representation of the object
print(AnalyticsDashboardData.to_json())

# convert the object into a dict
analytics_dashboard_data_dict = analytics_dashboard_data_instance.to_dict()
# create an instance of AnalyticsDashboardData from a dict
analytics_dashboard_data_from_dict = AnalyticsDashboardData.from_dict(analytics_dashboard_data_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


