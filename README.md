# VastCap Go

[![Go Reference](https://pkg.go.dev/badge/github.com/justyuuto/vastcap-go.svg)](https://pkg.go.dev/github.com/justyuuto/vastcap-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/justyuuto/vastcap-go)](https://goreportcard.com/report/github.com/justyuuto/vastcap-go)

A Go wrapper for the [VastCap](https://captcha.vast.sh) API.

## Installation

```bash
go get github.com/justyuuto/vastcap-go
```

## Usage

```go
package main

import "github.com/justyuuto/vastcap-go"

func main() {
    // Initialize the VastCap client with your API key
    client := vastcap.New("your_api_key")

    // Create a new task
    taskId, err := client.HCaptcha(vastcap.HCaptchaTask{
        WebsiteURL: "https://example.com",
        WebsiteKey: "website_key",
        // ...
    })
    if err != nil {
        panic(err)
    }

    // Get its result
    result, err := client.GetResult(taskId)
    if err != nil {
        panic(err)
    }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.