# ApiAPIErrorResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**error** | **str** |  | [optional] 

## Example

```python
from rule_engine_client.models.api_api_error_response import ApiAPIErrorResponse

# TODO update the JSON string below
json = "{}"
# create an instance of ApiAPIErrorResponse from a JSON string
api_api_error_response_instance = ApiAPIErrorResponse.from_json(json)
# print the JSON string representation of the object
print(ApiAPIErrorResponse.to_json())

# convert the object into a dict
api_api_error_response_dict = api_api_error_response_instance.to_dict()
# create an instance of ApiAPIErrorResponse from a dict
api_api_error_response_from_dict = ApiAPIErrorResponse.from_dict(api_api_error_response_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


