# rule_engine_client.ActionsApi

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**actions_get**](ActionsApi.md#actions_get) | **GET** /actions | List all actions
[**actions_id_get**](ActionsApi.md#actions_id_get) | **GET** /actions/{id} | Get action by ID
[**actions_post**](ActionsApi.md#actions_post) | **POST** /actions | Create a new action


# **actions_get**
> List[ApiActionDTO] actions_get()

List all actions

Get a list of all actions

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_action_dto import ApiActionDTO
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
    api_instance = rule_engine_client.ActionsApi(api_client)

    try:
        # List all actions
        api_response = api_instance.actions_get()
        print("The response of ActionsApi->actions_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ActionsApi->actions_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[ApiActionDTO]**](ApiActionDTO.md)

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

# **actions_id_get**
> ApiActionDTO actions_id_get(id)

Get action by ID

Get a specific action by its ID

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_action_dto import ApiActionDTO
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
    api_instance = rule_engine_client.ActionsApi(api_client)
    id = 'id_example' # str | Action ID

    try:
        # Get action by ID
        api_response = api_instance.actions_id_get(id)
        print("The response of ActionsApi->actions_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ActionsApi->actions_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Action ID | 

### Return type

[**ApiActionDTO**](ApiActionDTO.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **actions_post**
> ApiActionDTO actions_post(action)

Create a new action

Create a new action for rule execution

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_action_dto import ApiActionDTO
from rule_engine_client.models.api_create_action_request import ApiCreateActionRequest
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
    api_instance = rule_engine_client.ActionsApi(api_client)
    action = rule_engine_client.ApiCreateActionRequest() # ApiCreateActionRequest | Action data

    try:
        # Create a new action
        api_response = api_instance.actions_post(action)
        print("The response of ActionsApi->actions_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling ActionsApi->actions_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **action** | [**ApiCreateActionRequest**](ApiCreateActionRequest.md)| Action data | 

### Return type

[**ApiActionDTO**](ApiActionDTO.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**200** | OK |  -  |
**400** | Bad Request |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

