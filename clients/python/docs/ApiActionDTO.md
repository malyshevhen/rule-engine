# ApiActionDTO


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **str** |  | [optional] 
**enabled** | **bool** |  | [optional] 
**id** | **str** |  | [optional] 
**lua_script** | **str** |  | [optional] 
**updated_at** | **str** |  | [optional] 

## Example

```python
from rule_engine_client.models.api_action_dto import ApiActionDTO

# TODO update the JSON string below
json = "{}"
# create an instance of ApiActionDTO from a JSON string
api_action_dto_instance = ApiActionDTO.from_json(json)
# print the JSON string representation of the object
print(ApiActionDTO.to_json())

# convert the object into a dict
api_action_dto_dict = api_action_dto_instance.to_dict()
# create an instance of ApiActionDTO from a dict
api_action_dto_from_dict = ApiActionDTO.from_dict(api_action_dto_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


