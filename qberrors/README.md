# qberrors

The `qberrors` package provides consistent error handling for Quickbase applications written in Go.

## Goals

* Simplify code by sticking to [common Go conventions](https://blog.golang.org/go1.13-errors) with transparent abstractions for conveninence.
* Improve security by separating internal errors from those that are safe for the user to see.
* Improve user experience by returning appropriate HTTP status codes with detailed error messages.
* Simplify internal logging by maintaining the entire error chain through the edge of the app.
* Promote resiliency by recommending whether the operation that caused the error should be retried.

## Usage

* `Client` errors are the result of user input and should not be retried until the input is changed.
* `Internal` errors are internal to the application and should not be retried until code is fixed.
* `Service` errors are temporary problems that should be retried using a backoff algorithm.

Standard errors are treated as `Internal` and unsafe for users to see.

```go
err := errors.New("something bad happened")

fmt.Printf("%t\n", qberrors.IsSafe(err))
// Output: false

fmt.Println(qberrors.SafeMessage(err))
// Output: internal error

fmt.Println(qberrors.Upstream(err))
// Output: something bad happened

fmt.Println(qberrors.StatusCode(err))
// Output: 500
```

`Service` errors imply that the operation should be retried. The example below displays a helpful message to the user while maintaining the internal error chain for logging:

```go
connect := func() error {
	return errors.New("timeout connecting to service")
}

cerr := connect()
werr := fmt.Errorf("additional context: %w", cerr)

time := time.Now().Add(5 * time.Minute).Format("3:04 PM") // A time 5 minutes from now.
err := qberrors.Service(werr).Safef(qberrors.ServceUnavailable, "please retry at %s", time)

fmt.Println(err)
// Output: please retry at 11:19 AM: service unavailable

fmt.Println(qberrors.SafeMessage(err))
// Output: service unavailable

fmt.Println(qberrors.SafeDetail(err))
// Output: please retry at 11:19 AM

fmt.Println(qberrors.Upstream(err))
// Output: additional context: timeout connecting to service

fmt.Println(qberrors.StatusCode(err))
// Output: 503
```

Handling "not found" errors is common in applications. This library treats them as `Client` errors so that developers can use Go's error handling capabilities to control the logic of the application and show users a helpful message with an appropriate status code.

```go
id := "123"
err := qberrors.NotFoundError("item %q", id)

fmt.Printf("%t\n", errors.Is(err, qberrors.NotFound))
// Output: true

fmt.Println(err)
// Output: item "123": not found

fmt.Println(qberrors.Upstream(err))
// Output: <nil>

fmt.Println(qberrors.StatusCode(err))
// Output: 404
```