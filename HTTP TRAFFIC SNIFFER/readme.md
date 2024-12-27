# HTTP Traffic Sniffer

This project is a simple HTTP traffic sniffer written in Go. It logs incoming HTTP requests based on specified filters and rate limits. The log files are rotated based on size to ensure manageability.

## Features

- Filters requests based on IP, URL, and HTTP method.
- Rate limits the logging to prevent excessive log entries.
- Rotates log files when they exceed a specified size (5MB).

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/http-traffic-sniffer.git
   cd http-traffic-sniffer
   ```

2. Build the project:
   ```sh
   go build -o http-traffic-sniffer
   ```

## Usage

1. Run the executable:

   ```sh
   ./http-traffic-sniffer
   ```

2. The server will start on port `8080`. You can send HTTP requests to `http://localhost:8080`.

## Screenshots

![HTTP Traffic Sniffer](https://via.placeholder.com/800x400.png?text=HTTP+Traffic+Sniffer)

## Badges

![Go](https://img.shields.io/badge/Go-1.17-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)

## Stickers

![Sticker](https://via.placeholder.com/150x150.png?text=Sticker+1)
![Sticker](https://via.placeholder.com/150x150.png?text=Sticker+2)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Contact

For any inquiries, please contact [yourname@example.com](mailto:yourname@example.com).

🚀 Happy Coding!
