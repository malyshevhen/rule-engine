# ApiUpdateRuleRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**LuaScript** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Priority** | Pointer to **int32** |  | [optional] 

## Methods

### NewApiUpdateRuleRequest

`func NewApiUpdateRuleRequest() *ApiUpdateRuleRequest`

NewApiUpdateRuleRequest instantiates a new ApiUpdateRuleRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiUpdateRuleRequestWithDefaults

`func NewApiUpdateRuleRequestWithDefaults() *ApiUpdateRuleRequest`

NewApiUpdateRuleRequestWithDefaults instantiates a new ApiUpdateRuleRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *ApiUpdateRuleRequest) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ApiUpdateRuleRequest) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ApiUpdateRuleRequest) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ApiUpdateRuleRequest) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetLuaScript

`func (o *ApiUpdateRuleRequest) GetLuaScript() string`

GetLuaScript returns the LuaScript field if non-nil, zero value otherwise.

### GetLuaScriptOk

`func (o *ApiUpdateRuleRequest) GetLuaScriptOk() (*string, bool)`

GetLuaScriptOk returns a tuple with the LuaScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLuaScript

`func (o *ApiUpdateRuleRequest) SetLuaScript(v string)`

SetLuaScript sets LuaScript field to given value.

### HasLuaScript

`func (o *ApiUpdateRuleRequest) HasLuaScript() bool`

HasLuaScript returns a boolean if a field has been set.

### GetName

`func (o *ApiUpdateRuleRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ApiUpdateRuleRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ApiUpdateRuleRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ApiUpdateRuleRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPriority

`func (o *ApiUpdateRuleRequest) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ApiUpdateRuleRequest) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ApiUpdateRuleRequest) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ApiUpdateRuleRequest) HasPriority() bool`

HasPriority returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


