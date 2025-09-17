# ApiCreateTriggerRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**condition_script** | **str** |  | 
**enabled** | **bool** |  | [optional] 
**rule_id** | **str** |  | 
**type** | **str** |  | 

## Example

```python
from rule_engine_client.models.api_create_trigger_request import ApiCreateTriggerRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ApiCreateTriggerRequest from a JSON string
api_create_trigger_request_instance = ApiCreateTriggerRequest.from_json(json)
# print the JSON string representation of the object
print(ApiCreateTriggerRequest.to_json())

# convert the object into a dict
api_create_trigger_request_dict = api_create_trigger_request_instance.to_dict()
# create an instance of ApiCreateTriggerRequest from a dict
api_create_trigger_request_from_dict = ApiCreateTriggerRequest.from_dict(api_create_trigger_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


