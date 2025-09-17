# ApiCreateActionRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**LuaScript** | **string** |  | 

## Methods

### NewApiCreateActionRequest

`func NewApiCreateActionRequest(luaScript string, ) *ApiCreateActionRequest`

NewApiCreateActionRequest instantiates a new ApiCreateActionRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiCreateActionRequestWithDefaults

`func NewApiCreateActionRequestWithDefaults() *ApiCreateActionRequest`

NewApiCreateActionRequestWithDefaults instantiates a new ApiCreateActionRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *ApiCreateActionRequest) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ApiCreateActionRequest) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ApiCreateActionRequest) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ApiCreateActionRequest) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetLuaScript

`func (o *ApiCreateActionRequest) GetLuaScript() string`

GetLuaScript returns the LuaScript field if non-nil, zero value otherwise.

### GetLuaScriptOk

`func (o *ApiCreateActionRequest) GetLuaScriptOk() (*string, bool)`

GetLuaScriptOk returns a tuple with the LuaScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLuaScript

`func (o *ApiCreateActionRequest) SetLuaScript(v string)`

SetLuaScript sets LuaScript field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


