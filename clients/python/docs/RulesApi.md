# rule_engine_client.RulesApi

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**rules_get**](RulesApi.md#rules_get) | **GET** /rules | List all rules
[**rules_id_delete**](RulesApi.md#rules_id_delete) | **DELETE** /rules/{id} | Delete rule
[**rules_id_get**](RulesApi.md#rules_id_get) | **GET** /rules/{id} | Get rule by ID
[**rules_id_put**](RulesApi.md#rules_id_put) | **PUT** /rules/{id} | Update rule
[**rules_post**](RulesApi.md#rules_post) | **POST** /rules | Create a new rule


# **rules_get**
> List[ApiRuleDTO] rules_get()

List all rules

Get a list of all automation rules

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_rule_dto import ApiRuleDTO
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
    api_instance = rule_engine_client.RulesApi(api_client)

    try:
        # List all rules
        api_response = api_instance.rules_get()
        print("The response of RulesApi->rules_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling RulesApi->rules_get: %s\n" % e)
```



### Parameters

This endpoint does not need any parameter.

### Return type

[**List[ApiRuleDTO]**](ApiRuleDTO.md)

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

# **rules_id_delete**
> rules_id_delete(id)

Delete rule

Delete an automation rule by its ID

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
    api_instance = rule_engine_client.RulesApi(api_client)
    id = 'id_example' # str | Rule ID

    try:
        # Delete rule
        api_instance.rules_id_delete(id)
    except Exception as e:
        print("Exception when calling RulesApi->rules_id_delete: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Rule ID | 

### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

### HTTP response details

| Status code | Description | Response headers |
|-------------|-------------|------------------|
**204** | No Content |  -  |
**400** | Bad Request |  -  |
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **rules_id_get**
> ApiRuleDTO rules_id_get(id)

Get rule by ID

Get a specific automation rule by its ID

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_rule_dto import ApiRuleDTO
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
    api_instance = rule_engine_client.RulesApi(api_client)
    id = 'id_example' # str | Rule ID

    try:
        # Get rule by ID
        api_response = api_instance.rules_id_get(id)
        print("The response of RulesApi->rules_id_get:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling RulesApi->rules_id_get: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Rule ID | 

### Return type

[**ApiRuleDTO**](ApiRuleDTO.md)

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

# **rules_id_put**
> ApiRuleDTO rules_id_put(id, rule)

Update rule

Update an existing automation rule

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_rule_dto import ApiRuleDTO
from rule_engine_client.models.api_update_rule_request import ApiUpdateRuleRequest
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
    api_instance = rule_engine_client.RulesApi(api_client)
    id = 'id_example' # str | Rule ID
    rule = rule_engine_client.ApiUpdateRuleRequest() # ApiUpdateRuleRequest | Updated rule data

    try:
        # Update rule
        api_response = api_instance.rules_id_put(id, rule)
        print("The response of RulesApi->rules_id_put:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling RulesApi->rules_id_put: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **str**| Rule ID | 
 **rule** | [**ApiUpdateRuleRequest**](ApiUpdateRuleRequest.md)| Updated rule data | 

### Return type

[**ApiRuleDTO**](ApiRuleDTO.md)

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
**404** | Not Found |  -  |
**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **rules_post**
> ApiRuleDTO rules_post(rule)

Create a new rule

Create a new automation rule with Lua script

### Example


```python
import rule_engine_client
from rule_engine_client.models.api_create_rule_request import ApiCreateRuleRequest
from rule_engine_client.models.api_rule_dto import ApiRuleDTO
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
    api_instance = rule_engine_client.RulesApi(api_client)
    rule = rule_engine_client.ApiCreateRuleRequest() # ApiCreateRuleRequest | Rule data

    try:
        # Create a new rule
        api_response = api_instance.rules_post(rule)
        print("The response of RulesApi->rules_post:\n")
        pprint(api_response)
    except Exception as e:
        print("Exception when calling RulesApi->rules_post: %s\n" % e)
```



### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **rule** | [**ApiCreateRuleRequest**](ApiCreateRuleRequest.md)| Rule data | 

### Return type

[**ApiRuleDTO**](ApiRuleDTO.md)

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

