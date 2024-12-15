# Go Rate Limit

This project is a simple rate limiting library for Go.

## Features

- Easy to use
- Lightweight
- Configurable

## Installation

To install the package, run:

```sh
go get github.com/yourusername/go-ratelimit
```

## Usage

Here's a basic example of how to use the rate limiter:

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/go-ratelimit"
)

func main() {
    limiter := ratelimit.New(1, time.Second) // 1 request per second

    for i := 0; i < 5; i++ {
        if limiter.Allow() {
            fmt.Println("Request allowed")
        } else {
            fmt.Println("Rate limit exceeded")
        }
        time.Sleep(500 * time.Millisecond)
    }
}
```

## Configuration

You can configure the rate limiter by specifying the number of requests and the time window:

```go
limiter := ratelimit.New(requests int, duration time.Duration)
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For any questions or suggestions, please contact at kkshah2005@gmail.com.
