# File Encryptor

This project is a simple file encryption and decryption tool written in Go. It allows you to securely encrypt and decrypt files using a specified key.

## Features

- Encrypt files using AES encryption.
- Decrypt files using the same key used for encryption.
- Simple command-line interface.

## Requirements

- Go 1.16 or higher

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/file-encryptor.git
   ```
2. Navigate to the project directory:
   ```sh
   cd file-encryptor
   ```
3. Build the project:
   ```sh
   go build -o file-encryptor
   ```

## Usage

### Encrypt a file

```sh
./file-encryptor encrypt -key your-encryption-key -input /path/to/inputfile -output /path/to/outputfile
```

### Decrypt a file

```sh
./file-encryptor decrypt -key your-encryption-key -input /path/to/inputfile -output /path/to/outputfile
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
