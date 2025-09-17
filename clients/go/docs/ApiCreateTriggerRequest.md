# ApiCreateTriggerRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConditionScript** | **string** |  | 
**Enabled** | Pointer to **bool** |  | [optional] 
**RuleId** | **string** |  | 
**Type** | **string** |  | 

## Methods

### NewApiCreateTriggerRequest

`func NewApiCreateTriggerRequest(conditionScript string, ruleId string, type_ string, ) *ApiCreateTriggerRequest`

NewApiCreateTriggerRequest instantiates a new ApiCreateTriggerRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiCreateTriggerRequestWithDefaults

`func NewApiCreateTriggerRequestWithDefaults() *ApiCreateTriggerRequest`

NewApiCreateTriggerRequestWithDefaults instantiates a new ApiCreateTriggerRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConditionScript

`func (o *ApiCreateTriggerRequest) GetConditionScript() string`

GetConditionScript returns the ConditionScript field if non-nil, zero value otherwise.

### GetConditionScriptOk

`func (o *ApiCreateTriggerRequest) GetConditionScriptOk() (*string, bool)`

GetConditionScriptOk returns a tuple with the ConditionScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConditionScript

`func (o *ApiCreateTriggerRequest) SetConditionScript(v string)`

SetConditionScript sets ConditionScript field to given value.


### GetEnabled

`func (o *ApiCreateTriggerRequest) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ApiCreateTriggerRequest) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ApiCreateTriggerRequest) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ApiCreateTriggerRequest) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetRuleId

`func (o *ApiCreateTriggerRequest) GetRuleId() string`

GetRuleId returns the RuleId field if non-nil, zero value otherwise.

### GetRuleIdOk

`func (o *ApiCreateTriggerRequest) GetRuleIdOk() (*string, bool)`

GetRuleIdOk returns a tuple with the RuleId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleId

`func (o *ApiCreateTriggerRequest) SetRuleId(v string)`

SetRuleId sets RuleId field to given value.


### GetType

`func (o *ApiCreateTriggerRequest) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ApiCreateTriggerRequest) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ApiCreateTriggerRequest) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


