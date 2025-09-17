# rule_engine_client.TriggersApi

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**triggers_get**](TriggersApi.md#triggers_get) | **GET** /triggers | List all triggers
[**triggers_id_get**](TriggersApi.md#triggers_id_get) | **GET** /triggers/{id} | Get trigger by ID
[**triggers_post**](TriggersApi.md#triggers_post) | **POST** /triggers | Create a new trigger


# **triggers_get**
> List[ApiTriggerDTO] triggers_get()

List all triggers

Get a list of all triggers

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_trigger_dto import ApiTriggerDTO
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
    api_instance = rule_engine_client.TriggersApi(api_client)

    try:
        # List all triggers
        api_response = api_instance.triggers_get()
        print("The response of TriggersApi->triggers_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TriggersApi->triggers_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[ApiTriggerDTO]**](ApiTriggerDTO.md)

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

# **triggers_id_get**
> ApiTriggerDTO triggers_id_get(id)

Get trigger by ID

Get a specific trigger by its ID

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_trigger_dto import ApiTriggerDTO
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
    api_instance = rule_engine_client.TriggersApi(api_client)
    id = 'id_example' # str | Trigger ID

    try:
        # Get trigger by ID
        api_response = api_instance.triggers_id_get(id)
        print("The response of TriggersApi->triggers_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TriggersApi->triggers_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Trigger ID | 

### Return type

[**ApiTriggerDTO**](ApiTriggerDTO.md)

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

# **triggers_post**
> ApiTriggerDTO triggers_post(trigger)

Create a new trigger

Create a new trigger for rule execution

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_create_trigger_request import ApiCreateTriggerRequest
from rule_engine_client.models.api_trigger_dto import ApiTriggerDTO
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
    api_instance = rule_engine_client.TriggersApi(api_client)
    trigger = rule_engine_client.ApiCreateTriggerRequest() # ApiCreateTriggerRequest | Trigger data

    try:
        # Create a new trigger
        api_response = api_instance.triggers_post(trigger)
        print("The response of TriggersApi->triggers_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling TriggersApi->triggers_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **trigger** | [**ApiCreateTriggerRequest**](ApiCreateTriggerRequest.md)| Trigger data | 

### Return type

[**ApiTriggerDTO**](ApiTriggerDTO.md)

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

