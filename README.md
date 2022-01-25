[![Go Report Card](https://goreportcard.com/badge/github.com/reecerussell/httpy)](https://goreportcard.com/badge/github.com/reecerussell/httpy)
[![codecov](https://codecov.io/gh/reecerussell/httpy/branch/master/graph/badge.svg)](https://codecov.io/gh/reecerussell/httpy)
[![Go Docs](https://godoc.org/github.com/reecerussell/httpy?status.svg)](https://godoc.org/github.com/reecerussell/httpy)
![Actions](https://github.com/reecerussell/httpy/actions/workflows/test.yaml/badge.svg)

# httpy

A lightweight package used to make sending HTTP requests easier, with the fluent API of Requests.

## Get started

Using Go modules, you can install httpy into your code by running:

```
> go get -u github.com/reecerussell/httpy
```

Once installed, you can start using `Requests`.

```go
import (
    "context"
    "fmt"

    "github.com/reecerussell/httpy"
)

...

func SendRequest() {
    ctx := context.Background()

    // Configure your request, then call `.Do(ctx)` at the end.
    resp, err := httpy.Get("https://jsonplaceholder.typicode.com/todos/1").Do(ctx)
    if err != nil {
        panic(err)
    }

    var data map[string]interface{}

    // You can use the helpers on the Response object to decode JSON data.
    err = resp.DecodeJSON(&data)
    if err != nil {
        panic(err)
    }

    // Do something with data.
    fmt.Print(data)
}
```

## Fluent API

Requests can be configured really easily by using the fluent-like API that Requests provide.

### Data

Setting the body of a Request is as easy as calling `.SetBody`. This is a generic function which just sets the body as what is passed to the function. For example:

```go
var myBody io.Reader

// Creates a new POST requests with the raw data.
httpy.Post("https://my-app.io").
    SetBody(myBody)
```

#### JSON

To help setting the body of a request to JSON, there's a function which can be used to encode data as JSON for the request body. This also set's the content type to `application/json`. For example:

```go
data := map[string]string{
    "foo": "bar",
}

// Creates a new POST requests with data as the JSON body.
httpy.Post("https://my-app.io").
    WithJSON(data)
```

#### Plain text

As well as having a helper function to set JSON data, there's a function that can be used to just set the body of a request to plain text, then set the content type to `text/plain`. For example:

```go
// Creates a new POST requests with data as the plain text body.
httpy.Post("https://my-app.io").
    WithPlainText("foobar")
```

### Headers

Requests provide quite a few different functions you can use to configure headers. The first one being the most generic, which takes a name and a set of values, which looks like the following:

```go
httpy.Get("https://my-app.io").
    // Adds the "my-header" header with the values "foo" and "bar".
    // This function can be passed an array of values. However, if
    // no values are given, the header will not be set.
    SetHeader("my-header", "foo", "bar")
```

#### Content-Type

As `Content-Type` is a fairly common header, there is a function which can be used to just set the content type.

```go
var myBody io.Reader

httpy.Post("https://my-app.io").
    SetBody(myBody).
    // Sets the "Content-Type" header to "application/json".
    SetContentType("application/json")
```

#### Authorization

There are two functions which can be used to set the `Authorization` header on a request. These are `WithBearer` which sets a bearer token and `WithBasicAuth` which sets the header value, after given a username and password.

```go
httpy.Get("https://my-app.io").
    // Sets the "Authorization" header to "Bearer my-token".
    WithBearer("my-token").

    // Sets the "Authorization" header to "Basic ...".
    WithBasicAuth("myUser", "myPass")
```

### Sending requests

To make sending requests just as simple as making them, with any `Request`, you can call `.Do(...)` on the end. This will send the request, using the default `Client`, and return a `Response`. An example of this is in the first code block in the "Get started" section.

When calling the `.Do` function, aside from passing it a `context.Context`, you can optionally pass it an instance of `Client`, which it will use to make the request. `Client` is custom interface within httpy, which is used to send requests - more on this further down.

## Responses

In httpy, we have custom `Request` and `Response` objects. Requests are what can be seen above - they are used to build and make requests. Responses, however, are what the Requests return after being sent. The `Response` object is actually just a wrapper around `http.Response`, but provides more functionality.

For example, here is the .`DecodeJSON(...)` function in action:

```go
ctx := context.Background()
resp, err := httpy.Get("https://jsonplaceholder.typicode.com/todos/1").Do(ctx)
if err != nil {
    panic(err)
}

var data map[string]interface{}

// You can use the DecodeJSON object to decode a JSON response to an object.
err = resp.DecodeJSON(&data)
if err != nil {
    panic(err)
}
```

## Clients

httpy provides a `Client` interface, which is used to send HTTP requests. This can be used to create custom abstractions of a HTTP client, and also provides a mockable interface used for testing. Clients can be used to configure aspects about requests. This can include the timeout of the requests and also a base URL.

For more about what is available on the `Client` interface, see the docs.

### Default

By default, there is a client configured at `httpy.DefaultClient` which is what the `Requests` used to send requests. This can be overrwritten be set with your own, or can be configured to how you'd like it.
