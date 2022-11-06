# splitwise.go
[![Go tests](https://github.com/aanzolaavila/splitwise.go/actions/workflows/go.yml/badge.svg)](https://github.com/aanzolaavila/splitwise.go/actions/workflows/go.yml)

Yet another community driven Golang SDK for [Splitwise](https://splitwise.com/) 3rd-party APIs made in Go. Inspired from the work done by [anvari1313](https://github.com/anvari1313/splitwise.go/tree/main).

## How to use start using it?
1. Get this package into your project's dependencies
```bash
$ go get -u github.com/aanzolaavila/splitwise.go
```

2. Get you application credentials
   To get them, register you app in [here](https://secure.splitwise.com/apps).

   You can either use [`ApiKeyAuth`](https://dev.splitwise.com/#section/Authentication/ApiKeyAuth), recommended for testing.

   Or, you can use the much more secure option with [OAuth](https://dev.splitwise.com/#section/Authentication/OAuth), in [this example](https://github.com/aanzolaavila/splitwise.go/blob/main/examples/run.go#L23) you can see how you can set it up.

## Examples
Get your app credentials, to use the examples, you can run them like this

### With API key
```bash
$ TOKEN="<your-api-token>" make examples
```

### With OAuth 2.0
For this method, make sure that the Callback URI is the same one as this [line](https://github.com/aanzolaavila/splitwise.go/blob/main/examples/run.go#L37)
```bash
$ CLIENT_ID="<client-id>" CLIENT_SECRET="<client-secret>" make examples
# The app will ask you to open a link in the web browser,
# it will redirect you to a page with a URI that looks
# like https://localhost:17000/oauth/redirect/?code=YOUR_CODE&state=state

# copy YOUR_CODE into the terminal and it will execute the examples
```

## How to collaborate?
You can create a PR to this repo for corrections or improvements. One of the main guidelines for approving them is that all tests must pass.

### Execute the tests
Every test uses either a mock or a httptest server for responding with custom responses.
``` bash
$ make test
# or
$ go test .
```

### Show coverage information
To generate and open coverage reports you can do
```bash
$ make coverage
# or
$ go test -v -cover -coverprofile=coverage.out .
$ go tool cover -html=coverage.out
```

This will open the web browser with the coverage report.
