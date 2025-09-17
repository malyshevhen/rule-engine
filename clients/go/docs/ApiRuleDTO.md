# ApiRuleDTO

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Actions** | Pointer to [**[]ApiActionDTO**](ApiActionDTO.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**Enabled** | Pointer to **bool** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**LuaScript** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Priority** | Pointer to **int32** |  | [optional] 
**Triggers** | Pointer to [**[]ApiTriggerDTO**](ApiTriggerDTO.md) |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 

## Methods

### NewApiRuleDTO

`func NewApiRuleDTO() *ApiRuleDTO`

NewApiRuleDTO instantiates a new ApiRuleDTO object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiRuleDTOWithDefaults

`func NewApiRuleDTOWithDefaults() *ApiRuleDTO`

NewApiRuleDTOWithDefaults instantiates a new ApiRuleDTO object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActions

`func (o *ApiRuleDTO) GetActions() []ApiActionDTO`

GetActions returns the Actions field if non-nil, zero value otherwise.

### GetActionsOk

`func (o *ApiRuleDTO) GetActionsOk() (*[]ApiActionDTO, bool)`

GetActionsOk returns a tuple with the Actions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActions

`func (o *ApiRuleDTO) SetActions(v []ApiActionDTO)`

SetActions sets Actions field to given value.

### HasActions

`func (o *ApiRuleDTO) HasActions() bool`

HasActions returns a boolean if a field has been set.

### GetCreatedAt

`func (o *ApiRuleDTO) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ApiRuleDTO) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ApiRuleDTO) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *ApiRuleDTO) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetEnabled

`func (o *ApiRuleDTO) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *ApiRuleDTO) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *ApiRuleDTO) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *ApiRuleDTO) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetId

`func (o *ApiRuleDTO) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ApiRuleDTO) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ApiRuleDTO) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ApiRuleDTO) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLuaScript

`func (o *ApiRuleDTO) GetLuaScript() string`

GetLuaScript returns the LuaScript field if non-nil, zero value otherwise.

### GetLuaScriptOk

`func (o *ApiRuleDTO) GetLuaScriptOk() (*string, bool)`

GetLuaScriptOk returns a tuple with the LuaScript field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLuaScript

`func (o *ApiRuleDTO) SetLuaScript(v string)`

SetLuaScript sets LuaScript field to given value.

### HasLuaScript

`func (o *ApiRuleDTO) HasLuaScript() bool`

HasLuaScript returns a boolean if a field has been set.

### GetName

`func (o *ApiRuleDTO) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ApiRuleDTO) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ApiRuleDTO) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ApiRuleDTO) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPriority

`func (o *ApiRuleDTO) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ApiRuleDTO) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ApiRuleDTO) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ApiRuleDTO) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetTriggers

`func (o *ApiRuleDTO) GetTriggers() []ApiTriggerDTO`

GetTriggers returns the Triggers field if non-nil, zero value otherwise.

### GetTriggersOk

`func (o *ApiRuleDTO) GetTriggersOk() (*[]ApiTriggerDTO, bool)`

GetTriggersOk returns a tuple with the Triggers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTriggers

`func (o *ApiRuleDTO) SetTriggers(v []ApiTriggerDTO)`

SetTriggers sets Triggers field to given value.

### HasTriggers

`func (o *ApiRuleDTO) HasTriggers() bool`

HasTriggers returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ApiRuleDTO) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ApiRuleDTO) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ApiRuleDTO) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ApiRuleDTO) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


