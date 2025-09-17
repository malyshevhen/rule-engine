# \AnalyticsAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AnalyticsDashboardGet**](AnalyticsAPI.md#AnalyticsDashboardGet) | **Get** /analytics/dashboard | Get analytics dashboard data



## AnalyticsDashboardGet

> AnalyticsDashboardData AnalyticsDashboardGet(ctx).TimeRange(timeRange).Execute()

Get analytics dashboard data



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
	timeRange := "timeRange_example" // string | Time range (1h, 24h, 7d, 30d) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnalyticsAPI.AnalyticsDashboardGet(context.Background()).TimeRange(timeRange).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnalyticsAPI.AnalyticsDashboardGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AnalyticsDashboardGet`: AnalyticsDashboardData
	fmt.Fprintf(os.Stdout, "Response from `AnalyticsAPI.AnalyticsDashboardGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAnalyticsDashboardGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **timeRange** | **string** | Time range (1h, 24h, 7d, 30d) | 

### Return type

[**AnalyticsDashboardData**](AnalyticsDashboardData.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

