# Vending_api

## Overview

Vending_api is a Go project that implements a backend server for a candy vending machine system based on a Swagger API specification. The server handles candy purchase requests, validates input, calculates change, and returns responses in JSON format. Additionally, it supports mutual TLS authentication and integrates a fun ASCII cow message in the response.

---

## Table of Contents

- [Introduction](#introduction)  
- [Features](#features)  
- [Getting Started](#getting-started)  
- [Usage](#usage)  
  - [Exercise 00: Catching the Fortune](#exercise-00-catching-the-fortune)  
  - [Exercise 01: Law and Order](#exercise-01-law-and-order)  
  - [Exercise 02: Old Cow](#exercise-02-old-cow)  
- [Project Structure](#project-structure)  

---

## Introduction

Mister Rogers runs a candy vending business with five unique candy types. After a server theft, the backend must be recreated quickly to keep the business alive.

This project implements:

- A REST API for candy purchases following a Swagger 2.0 specification.
- Validation of candy type abbreviations and candy count.
- Proper HTTP status codes and JSON responses for success, insufficient funds, and invalid input.
- Mutual TLS authentication with self-signed certificates for secure communication.
- Integration of an ASCII art "thank you" cow message in the purchase response, reusing a provided C function via interoperability or embedding.

---

## Features

- Validates candy purchase requests ensuring:
  - Candy type is one of the predefined abbreviations: CE, AA, NT, DE, YR.
  - Candy count is non-negative.
- Calculates total cost and returns:
  - HTTP 201 with JSON containing a "thanks" message and "change" if enough money is provided.
  - HTTP 402 with an error message if money is insufficient.
  - HTTP 400 with an error message for invalid input.
- Supports mutual TLS authentication using self-signed certificates and a local certificate authority (minica).
- Provides a test client supporting flags to specify candy type, count, and money.
- Returns an ASCII art cow with the thank you message in the JSON response.

---

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git
- Minica (for generating self-signed certificates)
- A C compiler (optional, if integrating the original C function directly)

### Installation

```bash
git clone https://github.com/Desolitto/Go_Vending_api
cd Go_Vending_api
go mod tidy
```

---

## Usage

### Exercise 00: Catching the Fortune

Run the server implementing the `/buy_candy` endpoint according to the Swagger spec:

- **Request JSON fields:**
  - `money` (integer): amount inserted
  - `candyType` (string): one of CE, AA, NT, DE, YR
  - `candyCount` (integer): number of candies to buy

- **Responses:**
  - 201 Created:  
    ```json
    {
      "thanks": "Thank you!",
      "change": <integer>
    }
    ```
  - 400 Bad Request (invalid input):  
    ```json
    {
      "error": "<description>"
    }
    ```
  - 402 Payment Required (not enough money):  
    ```json
    {
      "error": "You need <amount> more money!"
    }
    ```

Example curl command:

```bash
curl -XPOST -H "Content-Type: application/json" -d '{"money": 20, "candyType": "AA", "candyCount": 1}' http://127.0.0.1:3333/buy_candy
```

---

### Exercise 01: Law and Order

Enhance the server with mutual TLS authentication:

- Use self-signed certificates generated via minica.
- Server and client authenticate each other.
- The client supports flags:
  - `-k` candy type (two-letter abbreviation)
  - `-c` candy count
  - `-m` money amount

Example client usage:

```bash
./candy-client -k AA -c 2 -m 50
Thank you! Your change is 20
```

---

### Exercise 02: Old Cow

Modify the server to include an ASCII art cow in the "thanks" message, reusing the provided C function `ask_cow()`.

Example JSON response:

```json
{
  "change": 0,
  "thanks": " ____________\n< Thank you! >\n ------------\n        \\   ^__^\n         \\  (oo)\\_______\n            (__)\\       )\\/\\\n                ||----w |\n                ||     ||\n"
}
```

Example curl command with TLS:

```bash
curl -s --key cert/client/key.pem --cert cert/client/cert.pem --cacert cert/minica.pem -XPOST -H "Content-Type: application/json" -d '{"candyType": "NT", "candyCount": 2, "money": 34}' "https://candy.tld:3333/buy_candy"
```

---

## Project Structure

```
src/
├── ex00/
│   ├── client/
│   ├── cmd/candy-server-server/
│   ├── restapi/
│   ├── Makefile
│   ├── go.mod
│   ├── go.sum
│   ├── swagger.yaml
│   └── task.md
├── ex01/
│   ├── candy.tld/
│   ├── cmd/candy-server-server/
│   ├── localhost/
│   ├── restapi/
│   ├── Makefile
│   ├── go.mod
│   ├── go.sum
│   ├── main_client.go
│   ├── minica-key.pem
│   ├── minica.pem
│   ├── swagger.yaml
│   └── task.md
├── ex02/
│   ├── cmd/candy-server-server/
│   ├── localhost/
│   ├── restapi/
│   ├── Makefile
│   ├── go.mod
│   ├── go.sum
│   ├── main_client
│   ├── main_client.go
│   ├── minica-key.pem
│   ├── minica.pem
│   ├── swagger.yaml
│   └── task.md
└── README.md
```

---
