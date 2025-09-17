# ApiRuleDTO


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**actions** | [**List[ApiActionDTO]**](ApiActionDTO.md) |  | [optional] 
**created_at** | **str** |  | [optional] 
**enabled** | **bool** |  | [optional] 
**id** | **str** |  | [optional] 
**lua_script** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**priority** | **int** |  | [optional] 
**triggers** | [**List[ApiTriggerDTO]**](ApiTriggerDTO.md) |  | [optional] 
**updated_at** | **str** |  | [optional] 

## Example

```python
from rule_engine_client.models.api_rule_dto import ApiRuleDTO

# TODO update the JSON string below
json = "{}"
# create an instance of ApiRuleDTO from a JSON string
api_rule_dto_instance = ApiRuleDTO.from_json(json)
# print the JSON string representation of the object
print(ApiRuleDTO.to_json())

# convert the object into a dict
api_rule_dto_dict = api_rule_dto_instance.to_dict()
# create an instance of ApiRuleDTO from a dict
api_rule_dto_from_dict = ApiRuleDTO.from_dict(api_rule_dto_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


