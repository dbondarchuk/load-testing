# Random string value

Generates a random string value in specified bounds and stores it into the variable.

It uses upper and lower characters of the English alphabet, digits and can include special characters from the list: `!@#$%^&*()-_+=[]\/,.<>;:'"`

## Type name

`RandomStringValue`

## Properties

|Property|Type|Is Required|Supports interpolation|Description|
|---|---|---|---|---|
|minLength|int as string|yes|yes|The minimal length of string to generate|
|maxLength|int as string|yes|yes|The maximal length of string to generate|
|specialCharacters|bool|yes|no|Indicates whether to include special characters|
|excludedCharacters|string|yes|no|Comma separated string of characters to exclude from the generation. If you need to exclude comma, use double comma|
|variableName|string|yes|no|Name of the variable where the generated string will be stored|