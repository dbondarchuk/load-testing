# Reading results

Each test run will produce two result files: `{{ID}}.http.log` and `{ID}.result.log`.

## `result.log` file

The result file is CSV file with `|` symbol as a delimeter and has the following structure:

`USER|LOOP|IS_SUCCESSFUL|EXECUTION_TIME|ERROR_LIST`

|Property|Type|Description|
|---|---|---|
|USER|int|User's number starting from 1|
|LOOP|int|Number of the loop for this user starting from 1|
|IS_SUCCESSFUL|boolean|Indicates whether the loop was succesful (no errors)|
|EXECUTION_TIME|int|Loop's executation time in micro seconds|
|ERROR_LIST|Array of strings as JSON|Contains erros that were caught during the loop execution serialized as JSON|

## `http.log` file

The result file is CSV file with `|` symbol as a delimeter and has the following structure:

`TIME|NUMBER_OF_REQUESTS|AVG_LATENCY|STATUS_BREAKDOWN`

|Property|Type|Description|
|---|---|---|
|TIME|int|Number of seconds passed|
|NUMBER_OF_REQUESTS|int|Number of requests that were done in this second|
|AVG_LATENCY|int|Average response time of requests that were done in this second|
|STATUS_BREAKDOWN|JSON object|JSON object with urls as keys and object `{statusCode: numberOfRequestsWithThisStatusCode}` as value|

## Console log
Console log contains same information in a simplified way:

```
Second: {{SECOND}}. Loops done: {{Loops done in this second}}. Total loops done: {{Total loops done till this second}}
Second: {Second}. Avg latency: {{Average response time of requests that finnished this second in micro seconds}} μs ({{Average response time of requests that finnished this second in milliseconds}} ms). req/s: {{Number of requests that are done in that second}}
```