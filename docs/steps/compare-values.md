# Compare values

Makes comparison of two string values and fails the test if the assertion is wrong

## Type name

`CompareValues`

## Properties

|Property|Type|Is Required|Supports interpolation|Description|
|---|---|---|---|---|
|method|ValueComparisonMethodEnum|yes|no|Method to compare too values|
|value|string|yes|yes|Left part of the comparison|
|to|string|yes|yes|Right part of the comparison|
|ignoreCase|boolean|yes|no|Indicates whether to make comparison case insensetive|

## ValueComparisonMethodEnum

|Name|Description|
|---|---|
|equals|Both parts are equal|
|notEquals|Both parts aren't equal|
|startsWith|Left part starts with right one|
|notStartsWith|Left part doesn't start with right one|
|endsWith|Left part ends with right one|
|notEndsWith|Left part doesn't end with right one|
|contains|Left part contais right one|
|notEndsWith|Left part doesn't contain right one|
|regex|Checks RegEx on the left part. Right part is the pattern|