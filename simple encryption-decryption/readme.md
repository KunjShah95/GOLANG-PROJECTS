# Simple Encryption-Decryption

🔒 A simple Go project for encryption and decryption of text.

## Features

- Encrypt text using a symmetric key
- Decrypt text using the same symmetric key
- Easy to use command-line interface

## Installation

```bash
go get github.com/KunjShah95/simple-encryption-decryption
```

## Usage

### Encrypting Text

```bash
go run main.go encrypt -k yourkey -t "your text to encrypt"
```

### Decrypting Text

```bash
go run main.go decrypt -k yourkey -t "your encrypted text"
```

## Example

```bash
go run main.go encrypt -k mysecretkey -t "Hello, World!"
# Output: Encrypted text

go run main.go decrypt -k mysecretkey -t "Encrypted text"
# Output: Hello, World!
```

## Contributing

🤝 Contributions are welcome! Please submit a pull request or open an issue.

## License

📄 This project is licensed under the MIT License.

## Contact

📧 For any inquiries, please contact [your email].

---

Made with ❤️ by [Your Name]
