# HTTP Server From Scratch in Go

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A minimalist HTTP/1.1 web server built from scratch in Go, using only the standard library. This project is a deep dive into the fundamentals of network programming and the HTTP protocol.

## Overview

This project implements a fully concurrent, static file web server. It is built without any external frameworks or libraries to demonstrate a foundational understanding of how web servers operate at a low level. The server is capable of parsing raw HTTP requests, routing based on URL paths, handling multiple connections simultaneously, and serving HTML content from the local filesystem.

## Key Features

- **Built From Scratch:** Uses only Go's standard `net` and `os` packages. No external dependencies.
- **Concurrent Architecture:** Leverages Go's powerful concurrency model, using a separate goroutine for each client connection to handle thousands of simultaneous requests efficiently.
- **HTTP/1.1 Request Parsing:** Manually parses the request line to extract the HTTP method, path, and version.
- **Static File Serving:** Reads and serves HTML files from a `/static` directory on the local filesystem.
- **Basic Routing:** Implements a simple router to serve different content based on the requested path (`/`, `/about`, etc.) and returns a custom 404 page for unknown routes.
- **JSON API Endpoint:** Includes a simple API at `/api/time` that returns the current server time in JSON format.
- **Clean Code Structure:** The project is organized following Go's best practices, with a clear separation between the application's entry point (`/cmd`) and its core logic (`/internal`).

## Core Concepts Explored

This project was a practical exercise in understanding:
- **TCP Sockets:** The fundamentals of `net.Listen`, `listener.Accept`, and managing `net.Conn` objects.
- **HTTP Protocol:** The structure of requests and responses, including status lines, headers (`Content-Type`, `Content-Length`), and the importance of CRLF (`\r\n`).
- **Concurrency in Go:** The practical application of goroutines to solve the challenge of handling multiple clients without blocking.
- **I/O Operations:** Reading from network connections and local files using buffers (`[]byte`).
- **API Development:** Understanding how to serve dynamic data as JSON and setting the correct `application/json` content type.
- **Code Refactoring:** The importance of breaking down complex functions into smaller, single-responsibility components for better readability and maintenance.

## Getting Started

### Prerequisites

- Go (version 1.18 or higher)

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/Oblutack/HTTP-Server.git
    ```
2.  **Navigate to the project directory:**
    ```sh
    cd HTTP-Server
    ```
3.  **Run the server:**
    ```sh
    go run ./cmd/webserver/main.go
    ```
4.  The server will start listening on `localhost:8080`. You can access it from your browser or using `curl`:
    ```sh
    # Access the home page
    curl http://localhost:8080/

    # Access the about page
    curl http://localhost:8080/about

    # Access the time API
    curl http://localhost:8080/api/time
    ```
