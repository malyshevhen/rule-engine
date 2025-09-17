# ApiTriggerDTO


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**condition_script** | **str** |  | [optional] 
**created_at** | **str** |  | [optional] 
**enabled** | **bool** |  | [optional] 
**id** | **str** |  | [optional] 
**type** | **str** |  | [optional] 
**updated_at** | **str** |  | [optional] 

## Example

```python
from rule_engine_client.models.api_trigger_dto import ApiTriggerDTO

# TODO update the JSON string below
json = "{}"
# create an instance of ApiTriggerDTO from a JSON string
api_trigger_dto_instance = ApiTriggerDTO.from_json(json)
# print the JSON string representation of the object
print(ApiTriggerDTO.to_json())

# convert the object into a dict
api_trigger_dto_dict = api_trigger_dto_instance.to_dict()
# create an instance of ApiTriggerDTO from a dict
api_trigger_dto_from_dict = ApiTriggerDTO.from_dict(api_trigger_dto_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


