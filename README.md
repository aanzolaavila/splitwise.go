# splitwise.go
[![Go tests](https://github.com/aanzolaavila/splitwise.go/actions/workflows/go.yml/badge.svg)](https://github.com/aanzolaavila/splitwise.go/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/aanzolaavila/splitwise.go/branch/main/graph/badge.svg?token=YB1C70XRP7)](https://codecov.io/gh/aanzolaavila/splitwise.go)

Yet another community driven Golang SDK for [Splitwise](https://splitwise.com/) 3rd-party APIs made in Go. Inspired from the work done by [anvari1313](https://github.com/anvari1313/splitwise.go/tree/main).

## How to use start using it?
1. Get this package into your project's dependencies
```bash
$ go get -u github.com/aanzolaavila/splitwise.go
```

2. Get you application credentials
   To get them, register you app in [here](https://secure.splitwise.com/apps).

   You can either use [`ApiKeyAuth`](https://dev.splitwise.com/#section/Authentication/ApiKeyAuth), recommended for testing.

   Or, you can use the much more secure option with [OAuth](https://dev.splitwise.com/#section/Authentication/OAuth), in [this example](https://github.com/aanzolaavila/splitwise.go/blob/main/examples/run.go#L25) you can see how you can set it up.

## Examples
Get your app credentials, to use the examples, you can run them like this

### With API key
```bash
$ TOKEN="<your-api-token>" make examples
```

### With OAuth 2.0
For this method, make sure that the Callback URI is the same one as this [line](https://github.com/aanzolaavila/splitwise.go/blob/main/examples/run.go#L39)
```bash
$ CLIENT_ID="<client-id>" CLIENT_SECRET="<client-secret>" make examples
# The app will ask you to open a link in the web browser,
# it will redirect you to a page with a URI that looks
# like https://localhost:17000/oauth/redirect/?code=YOUR_CODE&state=state

# copy YOUR_CODE into the terminal and it will execute the examples
```

### With third party JSON library
By default, the used library for marshaling and unmarshaling JSON is the native one, but in case you want to include another such as [go-json](https://github.com/goccy/go-json), [jsoniter](https://github.com/json-iterator/go), [sonic](https://github.com/bytedance/sonic), among others that might exist, you can specify another one.

For this you just need to point out the marshaling and unmarshaling methods like:
```go
import (
   	gojson "github.com/goccy/go-json"
)

// ...

client.JsonMarshaler = gojson.Marshal
client.JsonUnmarshaler = gojson.Unmarshal
```

Both fields use different interfaces, which are the same as the ones in the standard Go library.
```go
type jsonMarshaler func(interface{}) ([]byte, error)
type jsonUnmarshaler func([]byte, interface{}) error

type Client struct {
   // ...
   JsonMarshaler   jsonMarshaler
   JsonUnmarshaler jsonUnmarshaler
}

```

See a concrete example [here](https://github.com/aanzolaavila/splitwise.go/blob/main/examples/run.go#L76-L78).

### With custom logger
You can specify your desired logger implementation, it just has to fulfill the following interface.
```go
type logger interface {
	Printf(string, ...interface{})
}
```
You can find an example in [here](https://github.com/aanzolaavila/splitwise.go/blob/main/examples/run.go#L81-L82).

You can also use another logging solution such as [logrus](https://github.com/sirupsen/logrus). See this example for reference:
```go
import (
   "log"

   logrus "github.com/sirupsen/logrus"
)

// ...

logger := logrus.New()
logger.Formatter = &logrus.JSONFormatter{}

// Use logrus for standard log output
// Note that `log` here references stdlib's log
// Not logrus imported under the name `log`.
stdlog := log.New(os.Stdout, "", log.Lshortfile)
stdlog.SetOutput(logger.Writer())

splitwiseClient.Logger = stdlog
```
*Based from a logrus [example](https://github.com/sirupsen/logrus#logger-as-an-iowriter).*

#### Disable logging
If you wish to disable logging, you can do this
```go
import "log"

// ...

nopLogger := log.New(io.Discard, "", log.LstdFlags)
splitwiseClient.Logger = nopLogger
```

Or you can implement a `nopLogger` struct that fulfills the logger interface
```go
type nopLogger struct {}

func (l nopLogger) Printf(s string, args ...interface{}) {
   return // nothing to do here
}

// ...

splitwiseClient.Logger = nopLogger
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
