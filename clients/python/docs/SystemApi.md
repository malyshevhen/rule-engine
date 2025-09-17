# rule_engine_client.SystemApi

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**health_get**](SystemApi.md#health_get) | **GET** /health | Health check


# **health_get**
> Dict[str, str] health_get()

Health check

Get the health status of the service

### Example


```python
import rule_engine_client
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
    api_instance = rule_engine_client.SystemApi(api_client)

    try:
        # Health check
        api_response = api_instance.health_get()
        print("The response of SystemApi->health_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling SystemApi->health_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

**Dict[str, str]**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

