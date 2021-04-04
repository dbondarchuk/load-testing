# Compare variable

Makes an assertion on the varaible's value and fails the test if the assertion is wrong

## Type name

`CompareVariable`

## Properties

|Property|Type|Is Required|Supports interpolation|Description|
|---|---|---|---|---|
|method|ValueComparisonMethodEnum|yes|no|Method to compare too values|
|value|string|yes|no|Name of the variable for the left part of the comparison|
|to|string|yes|yes|Right part of the comparison|
|ignoreCase|boolean|yes|no|Indicates whether to make comparison case insensetive|

## ValueComparisonMethodEnum

Please look [here](compare-values.md)