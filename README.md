# instruqt-go

[![Go Report Card](https://goreportcard.com/badge/github.com/isovalent/instruqt-go)](https://goreportcard.com/report/github.com/isovalent/instruqt-go)  
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

`instruqt-go` is a Go client library for interacting with the Instruqt platform. It provides a simple and convenient way to programmatically access Instruqt's APIs, manage content, retrieve user data and track information.

## Features

- **Manage Instruqt Teams and Challenges**: Retrieve team information, challenges, and user progress.

## Installation

To install the `instruqt-go` library, run:

```shell
go get github.com/isovalent/instruqt-go
```


## Example Usage


```go
package main

import (
    "github.com/isovalent/instruqt-go/instruqt"
    "cloud.google.com/go/logging"
)

func main() {
    // Initialize the Instruqt client
    client := instruqt.NewClient("your-api-token", "your-team-slug")

    // Get all tracks
    tracks, err := client.GetTracks()

    // Add context to calls
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    clientWithTimeout := client.WithContext(ctx)
    userInfo, err := clientWithTimeout.GetUserInfo("user-id")

    // Attach a logger
    logClient, err := logging.NewClient(ctx, "some-gcp-project")
    if err != nil {
    	panic("failed to create log client")
    }
    stdLogger := logClient.Logger("my-std-logger")
    client.InfoLogger = stdLogger.StandardLogger(logging.Info)
}
```

## Contributing

We welcome contributions! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch with your feature or bug fix.
3. Make your changes and add tests.
4. Submit a pull request with a detailed description of your changes.

## Running Tests

To run the tests, use:

```shell
go test ./...
```


Make sure to write tests for any new functionality and ensure that all existing tests pass.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help, feel free to open an issue in the [GitHub repository](https://github.com/isovalent/instruqt-go/issues).

## Acknowledgments

- [Instruqt](https://instruqt.com/) for their amazing platform.