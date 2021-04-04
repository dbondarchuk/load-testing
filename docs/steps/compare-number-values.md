# Compare nummber values

Makes comparison of the the numeric values and fails the test if the assertion is wrong

## Type name

`CompareNumberValues`

## Properties

|Property|Type|Is Required|Supports interpolation|Description|
|---|---|---|---|---|
|method|NumericComparisonMethodEnum|yes|no|Method to compare too values|
|value|string|yes|yes|Left part of the comparison|
|to|string|yes|yes|Right part of the comparison|

## NumericComparisonMethodEnum

|Name|Description|
|---|---|
|lt|Left part is strictly less than right one|
|le|Left part is less or eaual than right one|
|gt|Left part is strictly greater than right one|
|ge|Left part is greater or eaual than right one|
|eq|Both parts are equal|
|ne|Both parts aren't equal|