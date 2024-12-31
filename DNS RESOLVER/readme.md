# DNS Resolver

![Go](https://img.shields.io/badge/Go-1.17-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)

🚀 A simple and efficient DNS resolver written in Go.

## Features

- 🔍 Fast and reliable DNS resolution
- 📦 Lightweight and easy to use
- 🔧 Configurable and extensible

## Installation

```bash
go get -u github.com/KunjShah95/dns-resolver
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/KunjShah95/dns-resolver"
)

func main() {
    resolver := dnsresolver.New()
    ip, err := resolver.Resolve("example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("IP Address:", ip)
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

## Contact

Feel free to reach out if you have any questions or suggestions! 😊
