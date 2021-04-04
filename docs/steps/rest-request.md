# REST Request

Makes a rest request, measures response time, captures response headers, status code and parses responses body (only JSON)

## Type name

`RestRequest`

## Properties

|Property|Type|Is Required|Supports interpolation|Description|
|---|---|---|---|---|
|method|HttpMethod|yes|no|HTTP method|
|endPoint|string|yes|yes|URL of the endpoint to reach|
|bodyType|HttpRequestBodyType|no|no|Type of the HTTP Request body. If ommitted - no body will be sent|
|rawData|string|no|yes|Raw string of data. Used when `bodyType` is `raw`|
|formData|Array of Key-Value pairs|no|yes|Array of key-values of the form data. Used when `bodyType` is `formdata` or `multipart`|
|files|Array of Key-Value pairs|no|yes|Array of key-values of file property name - path to file to upload. Used when `bodyType` is `multipart`|
|headers|Array of Key-Value pairs|no|yes|Array of key-values of HTTP request headers|
|variableName|string|no|yes|Name of the variable where to save the result of the execution of type `HttpResponseObject`|
|timeOut|int|no|no|Max allowed time in seconds for the request to be executed|

## HttpMethod

HTTP standard methods: `get`, `post` etc

## HttpRequestBodyType

`raw` - raw string of data. i.e JSON or XML. `rawData` property should be specified.
`dataform` - Form data
`multipart` - MultiPart form data. Can include files

## HttpResponseObject

|Property Name|Property type|Description|
|---|---|---|
|StatusCode|int|HTTP response status code|
|Headers|Array of Key-Value pairs|HTTP response headers|
|ResponseTime|int|Response time in milliseconds|
|Body|object|If the response result is JSON, `Body` property will contain parsed result|
