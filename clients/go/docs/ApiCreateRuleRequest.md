# ApiCreateRuleRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**LuaScript** | **string** |  | 
**Name** | **string** |  | 
**Priority** | Pointer to **int32** |  | [optional] 

## Methods

### NewApiCreateRuleRequest

`func NewApiCreateRuleRequest(luaScript string, name string, ) *ApiCreateRuleRequest`

NewApiCreateRuleRequest instantiates a new ApiCreateRuleRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiCreateRuleRequestWithDefaults

`func NewApiCreateRuleRequestWithDefaults() *ApiCreateRuleRequest`

NewApiCreateRuleRequestWithDefaults instantiates a new ApiCreateRuleRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *ApiCreateRuleRequest) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ApiCreateRuleRequest) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ApiCreateRuleRequest) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ApiCreateRuleRequest) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetLuaScript

`func (o *ApiCreateRuleRequest) GetLuaScript() string`

GetLuaScript returns the LuaScript field if non-nil, zero value otherwise.

### GetLuaScriptOk

`func (o *ApiCreateRuleRequest) GetLuaScriptOk() (*string, bool)`

GetLuaScriptOk returns a tuple with the LuaScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLuaScript

`func (o *ApiCreateRuleRequest) SetLuaScript(v string)`

SetLuaScript sets LuaScript field to given value.


### GetName

`func (o *ApiCreateRuleRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ApiCreateRuleRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ApiCreateRuleRequest) SetName(v string)`

SetName sets Name field to given value.


### GetPriority

`func (o *ApiCreateRuleRequest) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ApiCreateRuleRequest) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ApiCreateRuleRequest) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ApiCreateRuleRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


