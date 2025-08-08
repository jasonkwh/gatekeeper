# Gatekeeper

[![Go Reference](https://pkg.go.dev/badge/github.com/jasonkwh/gatekeeper.svg)](https://pkg.go.dev/github.com/jasonkwh/gatekeeper)

Gatekeeper is a Go library that provides gRPC interceptors for validating incoming requests. It allows you to easily add validation to your gRPC services by checking for a `Validate()` method on your request messages.

## Features

-   **Unary and Streaming Support:** Provides interceptors for both unary and streaming gRPC services.
-   **Simple Validation Interface:** Any request message that implements the `Validate() error` method will be automatically validated.
-   **Easy to Use:** Simply register your request types and add the interceptors to your gRPC server.

## Installation

```bash
go get github.com/jasonkwh/gatekeeper
```

## Usage

1.  **Define your protobuf messages**

    Define your request messages in your `.proto` file as you normally would.

2.  **Implement the `Validate()` method**

    For each request message that you want to validate, implement the `incomingRequest` interface. This interface has one method: `Validate() error`.

    ```go
    func (r *YourRequestMessage) Validate() error {
        if r.SomeField == "" {
            return errors.New("some_field is required")
        }
        return nil
    }
    ```

3.  **Register your request types**

    In your application's `init()` function, register the request types that you want to validate.

    ```go
    import "github.com/jasonkwh/gatekeeper"

    func init() {
        gatekeeper.RegisterRequest(&YourRequestMessage{})
        // You can also register multiple requests at once
        // gatekeeper.RegisterRequests(&Request1{}, &Request2{})
    }
    ```

4.  **Add the interceptors to your gRPC server**

    When creating your gRPC server, add the `UnaryServerInterceptor` and `StreamServerInterceptor` from the `gatekeeper` package.

    ```go
    import (
        "google.golang.org/grpc"
        "github.com/jasonkwh/gatekeeper"
    )

    func main() {
        server := grpc.NewServer(
            grpc.UnaryInterceptor(gatekeeper.UnaryServerInterceptor()),
            grpc.StreamInterceptor(gatekeeper.StreamServerInterceptor()),
        )

        // ... register your services and start the server
    }
    ```

## How it works

The `gatekeeper` interceptors check if an incoming request message has been registered. If it has, the interceptor will call the `Validate()` method on the request message. If the `Validate()` method returns an error, the interceptor will return a gRPC error with the code `InvalidArgument`. Otherwise, the request will be passed to the next handler.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
