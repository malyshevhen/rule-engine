# ApiUpdateRuleRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **bool** |  | [optional] 
**lua_script** | **str** |  | [optional] 
**name** | **str** |  | [optional] 
**priority** | **int** |  | [optional] 

## Example

```python
from rule_engine_client.models.api_update_rule_request import ApiUpdateRuleRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ApiUpdateRuleRequest from a JSON string
api_update_rule_request_instance = ApiUpdateRuleRequest.from_json(json)
# print the JSON string representation of the object
print(ApiUpdateRuleRequest.to_json())

# convert the object into a dict
api_update_rule_request_dict = api_update_rule_request_instance.to_dict()
# create an instance of ApiUpdateRuleRequest from a dict
api_update_rule_request_from_dict = ApiUpdateRuleRequest.from_dict(api_update_rule_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


