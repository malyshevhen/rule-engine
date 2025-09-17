# AnalyticsTimeSeriesData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | Pointer to [**[]AnalyticsTimeSeriesPoint**](AnalyticsTimeSeriesPoint.md) |  | [optional] 
**Metric** | Pointer to **string** |  | [optional] 

## Methods

### NewAnalyticsTimeSeriesData

`func NewAnalyticsTimeSeriesData() *AnalyticsTimeSeriesData`

NewAnalyticsTimeSeriesData instantiates a new AnalyticsTimeSeriesData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAnalyticsTimeSeriesDataWithDefaults

`func NewAnalyticsTimeSeriesDataWithDefaults() *AnalyticsTimeSeriesData`

NewAnalyticsTimeSeriesDataWithDefaults instantiates a new AnalyticsTimeSeriesData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *AnalyticsTimeSeriesData) GetData() []AnalyticsTimeSeriesPoint`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *AnalyticsTimeSeriesData) GetDataOk() (*[]AnalyticsTimeSeriesPoint, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *AnalyticsTimeSeriesData) SetData(v []AnalyticsTimeSeriesPoint)`

SetData sets Data field to given value.

### HasData

`func (o *AnalyticsTimeSeriesData) HasData() bool`

HasData returns a boolean if a field has been set.

### GetMetric

`func (o *AnalyticsTimeSeriesData) GetMetric() string`

GetMetric returns the Metric field if non-nil, zero value otherwise.

### GetMetricOk

`func (o *AnalyticsTimeSeriesData) GetMetricOk() (*string, bool)`

GetMetricOk returns a tuple with the Metric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetric

`func (o *AnalyticsTimeSeriesData) SetMetric(v string)`

SetMetric sets Metric field to given value.

### HasMetric

`func (o *AnalyticsTimeSeriesData) HasMetric() bool`

HasMetric returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


