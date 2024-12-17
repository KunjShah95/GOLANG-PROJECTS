# DNS Server Project

Welcome to the DNS Server Project! 🎉

## Overview

This project is a simple DNS server written in Go. It is designed to handle DNS queries and provide responses based on predefined rules.

## Features

- 📡 Handles DNS queries
- ⚡ Fast and efficient
- 🔧 Easy to configure

## Installation

To install the DNS server, follow these steps:

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/dns_server.git
   ```
2. Navigate to the project directory:
   ```sh
   cd dns_server
   ```
3. Build the project:
   ```sh
   go build
   ```

## Usage

To start the DNS server, run the following command:

```sh
./dns_server
```

## Configuration

You can configure the DNS server by editing the `config.json` file. Here is an example configuration:

```json
{
  "port": 53,
  "rules": [
    {
      "domain": "example.com",
      "ip": "192.168.1.1"
    }
  ]
}
```

## Contributing

We welcome contributions! Please read our [contributing guidelines](CONTRIBUTING.md) before submitting a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
