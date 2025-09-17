# \ActionsAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ActionsGet**](ActionsAPI.md#ActionsGet) | **Get** /actions | List all actions
[**ActionsIdGet**](ActionsAPI.md#ActionsIdGet) | **Get** /actions/{id} | Get action by ID
[**ActionsPost**](ActionsAPI.md#ActionsPost) | **Post** /actions | Create a new action



## ActionsGet

> []ApiActionDTO ActionsGet(ctx).Execute()

List all actions



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
	resp, r, err := apiClient.ActionsAPI.ActionsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ActionsAPI.ActionsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ActionsGet`: []ApiActionDTO
	fmt.Fprintf(os.Stdout, "Response from `ActionsAPI.ActionsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiActionsGetRequest struct via the builder pattern


### Return type

[**[]ApiActionDTO**](ApiActionDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ActionsIdGet

> ApiActionDTO ActionsIdGet(ctx, id).Execute()

Get action by ID



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
	id := "id_example" // string | Action ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ActionsAPI.ActionsIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ActionsAPI.ActionsIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ActionsIdGet`: ApiActionDTO
	fmt.Fprintf(os.Stdout, "Response from `ActionsAPI.ActionsIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Action ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiActionsIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ApiActionDTO**](ApiActionDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ActionsPost

> ApiActionDTO ActionsPost(ctx).Action(action).Execute()

Create a new action



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
	action := *openapiclient.NewApiCreateActionRequest("log_message('info', 'Temperature alert triggered')") // ApiCreateActionRequest | Action data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ActionsAPI.ActionsPost(context.Background()).Action(action).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ActionsAPI.ActionsPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ActionsPost`: ApiActionDTO
	fmt.Fprintf(os.Stdout, "Response from `ActionsAPI.ActionsPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiActionsPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **action** | [**ApiCreateActionRequest**](ApiCreateActionRequest.md) | Action data | 

### Return type

[**ApiActionDTO**](ApiActionDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

