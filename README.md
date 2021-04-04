# Introduction 

**Load tester** is a simple tool created for easy load testing HTTP endpoints focused on response time and verifying expected results.

What it can do:

1. Launch specific number of users (streams) and loops for each user over specific amount of time (ramp up).
2. Supports writing general results & http results.
3. Supports different tasks: http_request, sleep, comparisons.
4. Supports variables: both predefined and parsed from http results (status, headers, response time, body (json)).
5. Supports comparison actions,  which can be used to check if you got correct response code, got response in satisfying amount of time or if you got expected result -> number of items returned, property equals/contains/starts(ends)With/lesser/greater/regex - or their ‘not’ counterparts.
6. Accepts json as test definition.

# Getting Started

To start using the tool, download the latest binaries for your OS from the 'Releases' section on GitHub. This tools supports Windows (x86/x64), Linux (x86/x64) and MacOS (x64).
To run the test, execute following command in your terminal:
`load-tester -i TEST_FILE -id TEST_ID -o OUTPUT_FOLDER`

Learn about writing a test file definition [here](docs/writing-test.md) or have a look into the [example file](examples/input.json)

Learn about results produced by this application [here](docs/reading-results.md)

# Build

If you would like to build the application yourself, you can run included build script:
`build.sh PACKAGE_NAME [PLATFORM]`

Where `PACKAGE_NAME` would be the base name of output execution file (i.e if the package name is **load-tester** it will produce **load-tester-x64** or **load-tester-x86** files).

`PLATFORM` argument is optional and can be one of the following: `"windows/amd64" "windows/386" "darwin/amd64" "linux/amd64" "linux/386"`.
If this argument is ommited, than the build script will produce packages for all the mentioned above platforms.

# Contribute

Please feel free to create a pull request with your suggestions and fixes.