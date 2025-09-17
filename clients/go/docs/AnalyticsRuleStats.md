# AnalyticsRuleStats

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AverageLatencyMs** | Pointer to **float32** |  | [optional] 
**FailedExecutions** | Pointer to **int32** |  | [optional] 
**LastExecuted** | Pointer to **string** |  | [optional] 
**RuleId** | Pointer to **string** |  | [optional] 
**RuleName** | Pointer to **string** |  | [optional] 
**SuccessRate** | Pointer to **float32** |  | [optional] 
**SuccessfulExecutions** | Pointer to **int32** |  | [optional] 
**TotalExecutions** | Pointer to **int32** |  | [optional] 

## Methods

### NewAnalyticsRuleStats

`func NewAnalyticsRuleStats() *AnalyticsRuleStats`

NewAnalyticsRuleStats instantiates a new AnalyticsRuleStats object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAnalyticsRuleStatsWithDefaults

`func NewAnalyticsRuleStatsWithDefaults() *AnalyticsRuleStats`

NewAnalyticsRuleStatsWithDefaults instantiates a new AnalyticsRuleStats object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAverageLatencyMs

`func (o *AnalyticsRuleStats) GetAverageLatencyMs() float32`

GetAverageLatencyMs returns the AverageLatencyMs field if non-nil, zero value otherwise.

### GetAverageLatencyMsOk

`func (o *AnalyticsRuleStats) GetAverageLatencyMsOk() (*float32, bool)`

GetAverageLatencyMsOk returns a tuple with the AverageLatencyMs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAverageLatencyMs

`func (o *AnalyticsRuleStats) SetAverageLatencyMs(v float32)`

SetAverageLatencyMs sets AverageLatencyMs field to given value.

### HasAverageLatencyMs

`func (o *AnalyticsRuleStats) HasAverageLatencyMs() bool`

HasAverageLatencyMs returns a boolean if a field has been set.

### GetFailedExecutions

`func (o *AnalyticsRuleStats) GetFailedExecutions() int32`

GetFailedExecutions returns the FailedExecutions field if non-nil, zero value otherwise.

### GetFailedExecutionsOk

`func (o *AnalyticsRuleStats) GetFailedExecutionsOk() (*int32, bool)`

GetFailedExecutionsOk returns a tuple with the FailedExecutions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailedExecutions

`func (o *AnalyticsRuleStats) SetFailedExecutions(v int32)`

SetFailedExecutions sets FailedExecutions field to given value.

### HasFailedExecutions

`func (o *AnalyticsRuleStats) HasFailedExecutions() bool`

HasFailedExecutions returns a boolean if a field has been set.

### GetLastExecuted

`func (o *AnalyticsRuleStats) GetLastExecuted() string`

GetLastExecuted returns the LastExecuted field if non-nil, zero value otherwise.

### GetLastExecutedOk

`func (o *AnalyticsRuleStats) GetLastExecutedOk() (*string, bool)`

GetLastExecutedOk returns a tuple with the LastExecuted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastExecuted

`func (o *AnalyticsRuleStats) SetLastExecuted(v string)`

SetLastExecuted sets LastExecuted field to given value.

### HasLastExecuted

`func (o *AnalyticsRuleStats) HasLastExecuted() bool`

HasLastExecuted returns a boolean if a field has been set.

### GetRuleId

`func (o *AnalyticsRuleStats) GetRuleId() string`

GetRuleId returns the RuleId field if non-nil, zero value otherwise.

### GetRuleIdOk

`func (o *AnalyticsRuleStats) GetRuleIdOk() (*string, bool)`

GetRuleIdOk returns a tuple with the RuleId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleId

`func (o *AnalyticsRuleStats) SetRuleId(v string)`

SetRuleId sets RuleId field to given value.

### HasRuleId

`func (o *AnalyticsRuleStats) HasRuleId() bool`

HasRuleId returns a boolean if a field has been set.

### GetRuleName

`func (o *AnalyticsRuleStats) GetRuleName() string`

GetRuleName returns the RuleName field if non-nil, zero value otherwise.

### GetRuleNameOk

`func (o *AnalyticsRuleStats) GetRuleNameOk() (*string, bool)`

GetRuleNameOk returns a tuple with the RuleName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleName

`func (o *AnalyticsRuleStats) SetRuleName(v string)`

SetRuleName sets RuleName field to given value.

### HasRuleName

`func (o *AnalyticsRuleStats) HasRuleName() bool`

HasRuleName returns a boolean if a field has been set.

### GetSuccessRate

`func (o *AnalyticsRuleStats) GetSuccessRate() float32`

GetSuccessRate returns the SuccessRate field if non-nil, zero value otherwise.

### GetSuccessRateOk

`func (o *AnalyticsRuleStats) GetSuccessRateOk() (*float32, bool)`

GetSuccessRateOk returns a tuple with the SuccessRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccessRate

`func (o *AnalyticsRuleStats) SetSuccessRate(v float32)`

SetSuccessRate sets SuccessRate field to given value.

### HasSuccessRate

`func (o *AnalyticsRuleStats) HasSuccessRate() bool`

HasSuccessRate returns a boolean if a field has been set.

### GetSuccessfulExecutions

`func (o *AnalyticsRuleStats) GetSuccessfulExecutions() int32`

GetSuccessfulExecutions returns the SuccessfulExecutions field if non-nil, zero value otherwise.

### GetSuccessfulExecutionsOk

`func (o *AnalyticsRuleStats) GetSuccessfulExecutionsOk() (*int32, bool)`

GetSuccessfulExecutionsOk returns a tuple with the SuccessfulExecutions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccessfulExecutions

`func (o *AnalyticsRuleStats) SetSuccessfulExecutions(v int32)`

SetSuccessfulExecutions sets SuccessfulExecutions field to given value.

### HasSuccessfulExecutions

`func (o *AnalyticsRuleStats) HasSuccessfulExecutions() bool`

HasSuccessfulExecutions returns a boolean if a field has been set.

### GetTotalExecutions

`func (o *AnalyticsRuleStats) GetTotalExecutions() int32`

GetTotalExecutions returns the TotalExecutions field if non-nil, zero value otherwise.

### GetTotalExecutionsOk

`func (o *AnalyticsRuleStats) GetTotalExecutionsOk() (*int32, bool)`

GetTotalExecutionsOk returns a tuple with the TotalExecutions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalExecutions

`func (o *AnalyticsRuleStats) SetTotalExecutions(v int32)`

SetTotalExecutions sets TotalExecutions field to given value.

### HasTotalExecutions

`func (o *AnalyticsRuleStats) HasTotalExecutions() bool`

HasTotalExecutions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


