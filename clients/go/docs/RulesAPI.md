# \RulesAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RulesGet**](RulesAPI.md#RulesGet) | **Get** /rules | List all rules
[**RulesIdDelete**](RulesAPI.md#RulesIdDelete) | **Delete** /rules/{id} | Delete rule
[**RulesIdGet**](RulesAPI.md#RulesIdGet) | **Get** /rules/{id} | Get rule by ID
[**RulesIdPut**](RulesAPI.md#RulesIdPut) | **Put** /rules/{id} | Update rule
[**RulesPost**](RulesAPI.md#RulesPost) | **Post** /rules | Create a new rule



## RulesGet

> []ApiRuleDTO RulesGet(ctx).Execute()

List all rules



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
	resp, r, err := apiClient.RulesAPI.RulesGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RulesAPI.RulesGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RulesGet`: []ApiRuleDTO
	fmt.Fprintf(os.Stdout, "Response from `RulesAPI.RulesGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiRulesGetRequest struct via the builder pattern


### Return type

[**[]ApiRuleDTO**](ApiRuleDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RulesIdDelete

> RulesIdDelete(ctx, id).Execute()

Delete rule



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
	id := "id_example" // string | Rule ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RulesAPI.RulesIdDelete(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RulesAPI.RulesIdDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Rule ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRulesIdDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RulesIdGet

> ApiRuleDTO RulesIdGet(ctx, id).Execute()

Get rule by ID



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
	id := "id_example" // string | Rule ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RulesAPI.RulesIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RulesAPI.RulesIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RulesIdGet`: ApiRuleDTO
	fmt.Fprintf(os.Stdout, "Response from `RulesAPI.RulesIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Rule ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRulesIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**ApiRuleDTO**](ApiRuleDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RulesIdPut

> ApiRuleDTO RulesIdPut(ctx, id).Rule(rule).Execute()

Update rule



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
	id := "id_example" // string | Rule ID
	rule := *openapiclient.NewApiUpdateRuleRequest() // ApiUpdateRuleRequest | Updated rule data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RulesAPI.RulesIdPut(context.Background(), id).Rule(rule).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RulesAPI.RulesIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RulesIdPut`: ApiRuleDTO
	fmt.Fprintf(os.Stdout, "Response from `RulesAPI.RulesIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Rule ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiRulesIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **rule** | [**ApiUpdateRuleRequest**](ApiUpdateRuleRequest.md) | Updated rule data | 

### Return type

[**ApiRuleDTO**](ApiRuleDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RulesPost

> ApiRuleDTO RulesPost(ctx).Rule(rule).Execute()

Create a new rule



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
	rule := *openapiclient.NewApiCreateRuleRequest("if event.temperature > 25 then return true end", "Temperature Alert Rule") // ApiCreateRuleRequest | Rule data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RulesAPI.RulesPost(context.Background()).Rule(rule).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RulesAPI.RulesPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RulesPost`: ApiRuleDTO
	fmt.Fprintf(os.Stdout, "Response from `RulesAPI.RulesPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRulesPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **rule** | [**ApiCreateRuleRequest**](ApiCreateRuleRequest.md) | Rule data | 

### Return type

[**ApiRuleDTO**](ApiRuleDTO.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

