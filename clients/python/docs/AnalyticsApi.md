# rule_engine_client.AnalyticsApi

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**analytics_dashboard_get**](AnalyticsApi.md#analytics_dashboard_get) | **GET** /analytics/dashboard | Get analytics dashboard data


# **analytics_dashboard_get**
> AnalyticsDashboardData analytics_dashboard_get(time_range=time_range)

Get analytics dashboard data

Get aggregated analytics data for the dashboard

### Example


```python
import rule_engine_client
from rule_engine_client.models.analytics_dashboard_data import AnalyticsDashboardData
from rule_engine_client.rest import ApiException
from pprint import pprint

# Defining the host is optional and defaults to http://localhost:8080/api/v1
# See configuration.py for a list of all supported configuration parameters.
configuration = rule_engine_client.Configuration(
    host = "http://localhost:8080/api/v1"
)


# Enter a context with an instance of the API client
with rule_engine_client.ApiClient(configuration) as api_client:
    # Create an instance of the API class
    api_instance = rule_engine_client.AnalyticsApi(api_client)
    time_range = 'time_range_example' # str | Time range (1h, 24h, 7d, 30d) (optional)

    try:
        # Get analytics dashboard data
        api_response = api_instance.analytics_dashboard_get(time_range=time_range)
        print("The response of AnalyticsApi->analytics_dashboard_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling AnalyticsApi->analytics_dashboard_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **time_range** | **str**| Time range (1h, 24h, 7d, 30d) | [optional] 

### Return type

[**AnalyticsDashboardData**](AnalyticsDashboardData.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

