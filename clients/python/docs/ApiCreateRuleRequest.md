# ApiCreateRuleRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **bool** |  | [optional] 
**lua_script** | **str** |  | 
**name** | **str** |  | 
**priority** | **int** |  | [optional] 

## Example

```python
from rule_engine_client.models.api_create_rule_request import ApiCreateRuleRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ApiCreateRuleRequest from a JSON string
api_create_rule_request_instance = ApiCreateRuleRequest.from_json(json)
# print the JSON string representation of the object
print(ApiCreateRuleRequest.to_json())

# convert the object into a dict
api_create_rule_request_dict = api_create_rule_request_instance.to_dict()
# create an instance of ApiCreateRuleRequest from a dict
api_create_rule_request_from_dict = ApiCreateRuleRequest.from_dict(api_create_rule_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


