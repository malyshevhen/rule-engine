# AnalyticsDashboardData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExecutionTrend** | Pointer to [**AnalyticsTimeSeriesData**](AnalyticsTimeSeriesData.md) |  | [optional] 
**LatencyTrend** | Pointer to [**AnalyticsTimeSeriesData**](AnalyticsTimeSeriesData.md) |  | [optional] 
**OverallStats** | Pointer to [**AnalyticsExecutionStats**](AnalyticsExecutionStats.md) |  | [optional] 
**RuleStats** | Pointer to [**[]AnalyticsRuleStats**](AnalyticsRuleStats.md) |  | [optional] 
**SuccessRateTrend** | Pointer to [**AnalyticsTimeSeriesData**](AnalyticsTimeSeriesData.md) |  | [optional] 
**TimeRange** | Pointer to **string** |  | [optional] 

## Methods

### NewAnalyticsDashboardData

`func NewAnalyticsDashboardData() *AnalyticsDashboardData`

NewAnalyticsDashboardData instantiates a new AnalyticsDashboardData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAnalyticsDashboardDataWithDefaults

`func NewAnalyticsDashboardDataWithDefaults() *AnalyticsDashboardData`

NewAnalyticsDashboardDataWithDefaults instantiates a new AnalyticsDashboardData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExecutionTrend

`func (o *AnalyticsDashboardData) GetExecutionTrend() AnalyticsTimeSeriesData`

GetExecutionTrend returns the ExecutionTrend field if non-nil, zero value otherwise.

### GetExecutionTrendOk

`func (o *AnalyticsDashboardData) GetExecutionTrendOk() (*AnalyticsTimeSeriesData, bool)`

GetExecutionTrendOk returns a tuple with the ExecutionTrend field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExecutionTrend

`func (o *AnalyticsDashboardData) SetExecutionTrend(v AnalyticsTimeSeriesData)`

SetExecutionTrend sets ExecutionTrend field to given value.

### HasExecutionTrend

`func (o *AnalyticsDashboardData) HasExecutionTrend() bool`

HasExecutionTrend returns a boolean if a field has been set.

### GetLatencyTrend

`func (o *AnalyticsDashboardData) GetLatencyTrend() AnalyticsTimeSeriesData`

GetLatencyTrend returns the LatencyTrend field if non-nil, zero value otherwise.

### GetLatencyTrendOk

`func (o *AnalyticsDashboardData) GetLatencyTrendOk() (*AnalyticsTimeSeriesData, bool)`

GetLatencyTrendOk returns a tuple with the LatencyTrend field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLatencyTrend

`func (o *AnalyticsDashboardData) SetLatencyTrend(v AnalyticsTimeSeriesData)`

SetLatencyTrend sets LatencyTrend field to given value.

### HasLatencyTrend

`func (o *AnalyticsDashboardData) HasLatencyTrend() bool`

HasLatencyTrend returns a boolean if a field has been set.

### GetOverallStats

`func (o *AnalyticsDashboardData) GetOverallStats() AnalyticsExecutionStats`

GetOverallStats returns the OverallStats field if non-nil, zero value otherwise.

### GetOverallStatsOk

`func (o *AnalyticsDashboardData) GetOverallStatsOk() (*AnalyticsExecutionStats, bool)`

GetOverallStatsOk returns a tuple with the OverallStats field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOverallStats

`func (o *AnalyticsDashboardData) SetOverallStats(v AnalyticsExecutionStats)`

SetOverallStats sets OverallStats field to given value.

### HasOverallStats

`func (o *AnalyticsDashboardData) HasOverallStats() bool`

HasOverallStats returns a boolean if a field has been set.

### GetRuleStats

`func (o *AnalyticsDashboardData) GetRuleStats() []AnalyticsRuleStats`

GetRuleStats returns the RuleStats field if non-nil, zero value otherwise.

### GetRuleStatsOk

`func (o *AnalyticsDashboardData) GetRuleStatsOk() (*[]AnalyticsRuleStats, bool)`

GetRuleStatsOk returns a tuple with the RuleStats field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleStats

`func (o *AnalyticsDashboardData) SetRuleStats(v []AnalyticsRuleStats)`

SetRuleStats sets RuleStats field to given value.

### HasRuleStats

`func (o *AnalyticsDashboardData) HasRuleStats() bool`

HasRuleStats returns a boolean if a field has been set.

### GetSuccessRateTrend

`func (o *AnalyticsDashboardData) GetSuccessRateTrend() AnalyticsTimeSeriesData`

GetSuccessRateTrend returns the SuccessRateTrend field if non-nil, zero value otherwise.

### GetSuccessRateTrendOk

`func (o *AnalyticsDashboardData) GetSuccessRateTrendOk() (*AnalyticsTimeSeriesData, bool)`

GetSuccessRateTrendOk returns a tuple with the SuccessRateTrend field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccessRateTrend

`func (o *AnalyticsDashboardData) SetSuccessRateTrend(v AnalyticsTimeSeriesData)`

SetSuccessRateTrend sets SuccessRateTrend field to given value.

### HasSuccessRateTrend

`func (o *AnalyticsDashboardData) HasSuccessRateTrend() bool`

HasSuccessRateTrend returns a boolean if a field has been set.

### GetTimeRange

`func (o *AnalyticsDashboardData) GetTimeRange() string`

GetTimeRange returns the TimeRange field if non-nil, zero value otherwise.

### GetTimeRangeOk

`func (o *AnalyticsDashboardData) GetTimeRangeOk() (*string, bool)`

GetTimeRangeOk returns a tuple with the TimeRange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeRange

`func (o *AnalyticsDashboardData) SetTimeRange(v string)`

SetTimeRange sets TimeRange field to given value.

### HasTimeRange

`func (o *AnalyticsDashboardData) HasTimeRange() bool`

HasTimeRange returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


