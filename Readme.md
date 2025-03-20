# Students API

[Postman Collection](https://documenter.getpostman.com/view/34455053/2sAYkGKKDR)

A simple RESTful API built with Go to manage student records. It supports CRUD operations (Create, Read, Update, Delete) for student data, stored in a SQLite database.

## Features

- Create a new student (`POST /api/students`)
- Retrieve a student by ID (`GET /api/students/{id}`)
- List all students (`GET /api/students`)
- Update a student by ID (`PUT /api/students/{id}`)
- Delete a student by ID (`DELETE /api/students/{id}`)

## Tech Stack

- **Language**: Go (Golang)
- **Database**: SQLite
- **Configuration**: YAML with `cleanenv` for environment management
- **Validation**: `go-playground/validator` for request validation
- **Logging**: `log/slog` for structured logging

## Prerequisites

- Go 1.18 or higher
- SQLite (included via `github.com/mattn/go-sqlite3`)

## Setup and Installation

1. Clone the repository:
2. Install dependencies:
3. Configure the application:
   Copy `local.yaml` to a new file or set the `CONFIG_PATH` environment variable:
4. Run the application: Using command : go run cmd/students-api/main.go -config config/local.yaml

The API will be available at [http://localhost:8082](http://localhost:8082).

## API Endpoints

| Method | Endpoint           | Description            |
| ------ | ------------------ | ---------------------- |
| POST   | /api/students      | Create a new student   |
| GET    | /api/students/{id} | Get a student by ID    |
| GET    | /api/students      | Get all students       |
| PUT    | /api/students/{id} | Update a student by ID |
| DELETE | /api/students/{id} | Delete a student by ID |

## License

This project is licensed under the MIT License.

---
