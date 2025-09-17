# ApiCreateActionRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **bool** |  | [optional] 
**lua_script** | **str** |  | 

## Example

```python
from rule_engine_client.models.api_create_action_request import ApiCreateActionRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ApiCreateActionRequest from a JSON string
api_create_action_request_instance = ApiCreateActionRequest.from_json(json)
# print the JSON string representation of the object
print(ApiCreateActionRequest.to_json())

# convert the object into a dict
api_create_action_request_dict = api_create_action_request_instance.to_dict()
# create an instance of ApiCreateActionRequest from a dict
api_create_action_request_from_dict = ApiCreateActionRequest.from_dict(api_create_action_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


