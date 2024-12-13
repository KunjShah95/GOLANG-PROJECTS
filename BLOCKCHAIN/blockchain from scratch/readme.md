# 📚 Book Checkout Blockchain

This project implements a simple blockchain to manage book checkouts. Each block in the blockchain represents a book checkout transaction, containing details about the book, the user, and the checkout date. The blockchain ensures the integrity and immutability of the checkout records. 🔒

## 📋 Table of Contents

- [✨ Features](#features)
- [🛠️ Technologies Used](#technologies-used)
- [🚀 Getting Started](#getting-started)
- [📡 API Endpoints](#api-endpoints)
- [🧪 How to Test](#how-to-test)
- [🤝 Contributing](#contributing)
- [📄 License](#license)

## ✨ Features

- Create a new book checkout record. 📖
- Retrieve the entire blockchain of book checkouts. 🔍
- Generate a unique ID for each book based on its ISBN and publish date. 🆔
- Ensure data integrity through hashing. 🔐

## 🛠️ Technologies Used

- Go (Golang) 🦙
- Gorilla Mux (for routing) 🚦
- JSON (for data interchange) 📄
- SHA-256 and MD5 (for hashing) 🔑

## 🚀 Getting Started

To run this project locally, follow these steps:

1. **Clone the repository:**

   git clone https://github.com/yourusername/book-checkout-blockchain.git

   cd book-checkout-blockchain

3. Install dependencies:

Make sure you have Go installed. You can download it from golang.org. 🌐

3. Run the application:
   go run main.go
   The server will start listening on port 3000. 🎉

📡 API Endpoints

1. Get Blockchain
   Endpoint: GET /
   Description: Retrieves the entire blockchain of book checkouts. 📜
   Response:
   Returns a JSON array of blocks.


2. Write Block (Checkout a Book)
   Endpoint: POST /
   Description: Creates a new book checkout record. 📝
   Request Body:
   json
   {
   "book_id": "12345",
   "user": "John Doe",
   "checkout_date": "2023-10-01",
   "is_genesis": false
   }
   Response:
   Returns the created checkout record in JSON format.


3. Create New Book
   Endpoint: POST /new
   Description: Creates a new book record. 📚
   Request Body:
   json
   {
   "title": "The Great Gatsby",
   "author": "F. Scott Fitzgerald",
   "publish_date": "1925-04-10",
   "isbn": "9780743273565"
   }


    Response:
   Returns the created book record with a unique ID in JSON format.
   🧪 How to Test
   You can use tools like Postman or curl to test the API endpoints. 🛠️

Example using curl:

1. Create a new book:
   curl -X POST http://localhost:3000/new -H "Content-Type: application/json" -d '{"title": "The Great Gatsby", "author": "F. Scott Fitzgerald", "publish_date": "1925-04-10", "isbn": "9780743273565"}'

2. Checkout a book:
   curl -X POST http://localhost:3000/ -H "Content-Type: application/json" -d '{"book_id": "12345", "user": "John Doe", "checkout_date": "2023-10-01", "is_genesis": false}'

3. Get the blockchain:
   curl -X GET http://localhost:3000/

🤝 Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements. 💡

📄 License
This project is licensed under the MIT License. See the LICENSE file for details. 📝
