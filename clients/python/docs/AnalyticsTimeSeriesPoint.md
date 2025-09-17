# AnalyticsTimeSeriesPoint


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**timestamp** | **str** |  | [optional] 
**value** | **float** |  | [optional] 

## Example

```python
from rule_engine_client.models.analytics_time_series_point import AnalyticsTimeSeriesPoint

# TODO update the JSON string below
json = "{}"
# create an instance of AnalyticsTimeSeriesPoint from a JSON string
analytics_time_series_point_instance = AnalyticsTimeSeriesPoint.from_json(json)
# print the JSON string representation of the object
print(AnalyticsTimeSeriesPoint.to_json())

# convert the object into a dict
analytics_time_series_point_dict = analytics_time_series_point_instance.to_dict()
# create an instance of AnalyticsTimeSeriesPoint from a dict
analytics_time_series_point_from_dict = AnalyticsTimeSeriesPoint.from_dict(analytics_time_series_point_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


