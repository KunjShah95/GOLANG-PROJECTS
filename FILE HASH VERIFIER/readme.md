# File Hash Verifier

![Go](https://img.shields.io/badge/Go-1.17-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![Contributions](https://img.shields.io/badge/Contributions-Welcome-brightgreen)

🔍 **File Hash Verifier** is a simple tool written in Go to verify the integrity of your files by comparing their hash values.

## 🚀 Features

- Supports multiple hash algorithms: MD5, SHA-1, SHA-256, and more.
- Easy to use command-line interface.
- Fast and efficient hashing.

## 📦 Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/file-hash-verifier.git
   ```
2. Navigate to the project directory:
   ```sh
   cd file-hash-verifier
   ```
3. Build the project:
   ```sh
   go build
   ```

## 🛠 Usage

To verify a file's hash, use the following command:

```sh
./file-hash-verifier -file <path_to_file> -hash <expected_hash> -algo <hash_algorithm>
```

Example:

```sh
./file-hash-verifier -file example.txt -hash d41d8cd98f00b204e9800998ecf8427e -algo md5
```

## 🤝 Contributing

Contributions are welcome! Please fork this repository and submit a pull request.

## 📄 License

This project is licensed under the MIT License.

## 📧 Contact

For any questions or suggestions, feel free to open an issue or contact me at [your-email@example.com](mailto:your-email@example.com).

Happy hashing! 😊
