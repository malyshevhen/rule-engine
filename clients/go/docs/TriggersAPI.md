# \TriggersAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TriggersGet**](TriggersAPI.md#TriggersGet) | **Get** /triggers | List all triggers
[**TriggersIdGet**](TriggersAPI.md#TriggersIdGet) | **Get** /triggers/{id} | Get trigger by ID
[**TriggersPost**](TriggersAPI.md#TriggersPost) | **Post** /triggers | Create a new trigger



## TriggersGet

> []ApiTriggerDTO TriggersGet(ctx).Execute()

List all triggers



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/ruleengine"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TriggersAPI.TriggersGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TriggersAPI.TriggersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TriggersGet`: []ApiTriggerDTO
	fmt.Fprintf(os.Stdout, "Response from `TriggersAPI.TriggersGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiTriggersGetRequest struct via the builder pattern


### Return type

[**[]ApiTriggerDTO**](ApiTriggerDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TriggersIdGet

> ApiTriggerDTO TriggersIdGet(ctx, id).Execute()

Get trigger by ID



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/ruleengine"
)

func main() {
	id := "id_example" // string | Trigger ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TriggersAPI.TriggersIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TriggersAPI.TriggersIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TriggersIdGet`: ApiTriggerDTO
	fmt.Fprintf(os.Stdout, "Response from `TriggersAPI.TriggersIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Trigger ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiTriggersIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ApiTriggerDTO**](ApiTriggerDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TriggersPost

> ApiTriggerDTO TriggersPost(ctx).Trigger(trigger).Execute()

Create a new trigger



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/ruleengine"
)

func main() {
	trigger := *openapiclient.NewApiCreateTriggerRequest("if event.device_id == 'sensor_1' then return true end", "550e8400-e29b-41d4-a716-446655440000", "CONDITIONAL") // ApiCreateTriggerRequest | Trigger data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TriggersAPI.TriggersPost(context.Background()).Trigger(trigger).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TriggersAPI.TriggersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TriggersPost`: ApiTriggerDTO
	fmt.Fprintf(os.Stdout, "Response from `TriggersAPI.TriggersPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTriggersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **trigger** | [**ApiCreateTriggerRequest**](ApiCreateTriggerRequest.md) | Trigger data | 

### Return type

[**ApiTriggerDTO**](ApiTriggerDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

