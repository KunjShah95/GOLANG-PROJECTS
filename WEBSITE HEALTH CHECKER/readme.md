# Health Checker

![Go](https://img.shields.io/badge/Go-1.17-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)

## 🚀 Introduction

Health Checker is a simple Go application that checks the health of various services and endpoints. It is designed to be lightweight and easy to use.

## 📋 Features

- Check the health of HTTP endpoints
- Configurable check intervals
- Lightweight and fast

## 🛠️ Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/health-checker.git
   ```
2. Navigate to the project directory:
   ```sh
   cd health-checker
   ```
3. Build the application:
   ```sh
   go build
   ```

## 🚦 Usage

Run the application with the following command:

```sh
./health-checker -config=config.yaml
```

## 📝 Configuration

The application uses a YAML file for configuration. Below is an example configuration file:

```yaml
endpoints:
  - url: "https://example.com/health"
    interval: 60
  - url: "https://another-service.com/health"
    interval: 120
```

## 🤝 Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
