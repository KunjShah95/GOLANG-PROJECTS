# Blockchain From Scratch

This repository contains a Go-based implementation of a basic blockchain built from scratch. The project is designed to demonstrate foundational blockchain concepts, such as block creation, hashing, and verification, and is an excellent starting point for developers exploring blockchain technology.

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Custom Blockchain Implementation**: Built without relying on third-party libraries to understand core blockchain principles.
- **Block and Chain Management**: Supports adding blocks to the chain and maintaining a valid blockchain.
- **Hashing Mechanism**: Ensures block integrity using cryptographic hash functions.
- **Proof of Work (Optional)**: Demonstrates basic consensus mechanisms for blockchain validation.
- **Efficient Data Structure**: Designed for simplicity and clarity.

---

## Project Structure

```plaintext
blockchain-from-scratch/
|-- main.go           # Entry point of the application
|-- blockchain.go     # Core blockchain logic
|-- block.go          # Implementation of individual blocks
|-- utils.go          # Utility functions (e.g., hashing, timestamping)
|-- README.md         # Documentation (You are here!)
```

---

## Getting Started

### Prerequisites

To run this project, ensure you have the following installed:

- [Go](https://golang.org/) (version 1.18 or later)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/KunjShah95/GOLANG-PROJECTS.git
   cd GOLANG-PROJECTS/BLOCKCHAIN/blockchain%20from%20scratch
   ```

2. Install dependencies (if any):

   ```bash
   go mod tidy
   ```

3. Run the application:

   ```bash
   go run main.go
   ```

---

## Usage

The blockchain implementation allows you to:

- Add blocks with custom data to the blockchain.
- Validate the chain to ensure data integrity.

Modify the `main.go` file to add blocks or experiment with blockchain functionalities. Example code snippets are included in the [Examples](#examples) section.

---

## Examples

### Adding a New Block

```go
package main

func main() {
    blockchain := NewBlockchain()

    blockchain.AddBlock("First block after genesis")
    blockchain.AddBlock("Second block after genesis")

    for _, block := range blockchain.Blocks {
        fmt.Printf("Block Data: %s\n", block.Data)
        fmt.Printf("Hash: %x\n", block.Hash)
        fmt.Printf("Previous Hash: %x\n\n", block.PrevHash)
    }
}
```

### Output

```plaintext
Block Data: Genesis Block
Hash: <hash_value>
Previous Hash: <previous_hash_value>

Block Data: First block after genesis
Hash: <hash_value>
Previous Hash: <genesis_hash>

Block Data: Second block after genesis
Hash: <hash_value>
Previous Hash: <first_block_hash>
```

---

## Contributing

Contributions are welcome! If you have ideas for improving this project or would like to fix any issues, feel free to submit a pull request.

### Steps to Contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m "Add new feature"`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/KunjShah95/GOLANG-PROJECTS/blob/main/LICENSE) file for details.

---

## Contact

For any queries or discussions, feel free to connect via:

- **GitHub**: [KunjShah95](https://github.com/KunjShah95)
- **LinkedIn**:https://www.linkedin.com/in/kunj-shah15957477/
