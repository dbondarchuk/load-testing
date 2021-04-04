# Writing tests

Test in load tester are defined as *.json files. You can find an example of such file in [examples/input.json](../examples/input.json).

## JSON file properties

Each test file contains such properties:

|Property|Type|Is Required|Description|
|---|---|---|---|
|usersCount|int|yes|Number of parallel processes (users)|
|loopCount|int|yes|Number of loops which each user should run|
|rampup|int|yes|Number of seconds during which all users should start. Should be 0 or larger|
|steps|Array of *Step*|yes|List of steps to execute by each user|
|variables|Key-Value pair|no|List of predefined variables for the tests: `"name": "value"`|

## Step common properties

Each step should contain following properties:

|Property|Type|Is Required|Description|
|---|---|---|---|
|typeName|string|yes|Type of the action to execute|
|name|string|yes|Display name of the step|
|enabled|boolean|yes|Determines whether this step should be executed|
|ignoreError|boolean|yes|Determines whether test execution should ignore any error produced by this step|
|runOnFailure|boolean|yes|Determines whether this step should run even if the test is already failed|
|propertyValues|Key-Value pair of string-object|yes|Contains all properties for the step described as map of property name and property value|

## Supported Steps

1. [RestRequest](steps/rest-request.md) - Sends HTTP request.
2. [Sleep](steps/sleep.md) - Delay
3. [CompareNumberValues](steps/compare-number-values.md) - Compares numeric values
4. [CompareValues](steps/compare-values.md) - Compares two values
5. [CompareVariable](steps/compare-variable.md) - Compares variable's value to another value
6. [RandomNumberValue](steps/random-number-value.md) - Generates random number value
7. [RandomStringValue](steps/random-string-value.md) - Generates random string

## Variable interpolation

When the step has a property which supports variable interpolation you can use it to run an expression to be evaluated and passed into the step.

The basic syntax for this is `$(expression)`. i.e: `"value": "$(posts.Body[0].title)"`

You can use any JavaScript expression insid `$()`.

You also can use the expression along side with regular text. i.e: `"value": "Title is $(posts.Body[0].title)"`
