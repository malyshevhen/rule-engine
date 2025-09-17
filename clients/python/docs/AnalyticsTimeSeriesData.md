# AnalyticsTimeSeriesData


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**data** | [**List[AnalyticsTimeSeriesPoint]**](AnalyticsTimeSeriesPoint.md) |  | [optional] 
**metric** | **str** |  | [optional] 

## Example

```python
from rule_engine_client.models.analytics_time_series_data import AnalyticsTimeSeriesData

# TODO update the JSON string below
json = "{}"
# create an instance of AnalyticsTimeSeriesData from a JSON string
analytics_time_series_data_instance = AnalyticsTimeSeriesData.from_json(json)
# print the JSON string representation of the object
print(AnalyticsTimeSeriesData.to_json())

# convert the object into a dict
analytics_time_series_data_dict = analytics_time_series_data_instance.to_dict()
# create an instance of AnalyticsTimeSeriesData from a dict
analytics_time_series_data_from_dict = AnalyticsTimeSeriesData.from_dict(analytics_time_series_data_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


