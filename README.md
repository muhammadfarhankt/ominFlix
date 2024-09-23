# Omniflix Block Indexer

This project implements a block indexer for the Omniflixhub blockchain using GoLang. The indexer fetches block data from the Omniflixhub RPC API, stores it in a PostgreSQL database, and provides an API endpoint to access the indexed data.

## Features
- Fetches block data concurrently using goroutines.
- Stores block height, block ID, proposer address, number of transactions, and other relevant details in the database.
- Provides an API endpoint (/block/:height) to fetch block details by block height.
- Handles errors gracefully and includes basic error handling for API requests and database interactions.
- Uses a semaphore to limit concurrent API requests and prevent overloading the blockchain nodes.
- Includes timestamps (created_at, updated_at) for tracking changes in the database.

## Table of Contents
- [Project Overview](#project-overview)
- [Installation](#installation)
- [Project Structure](#project-structure)
- [Usage](#usage)
- [Configuration](#configuration)
- [Docker Setup](#docker-setup)
- [Makefile Commands](#makefile-commands)
- [API Documentation](#api-documentation)

## Project Overview

The indexer processes block data from the Omniflixhub blockchain and stores the block details in a PostgreSQL database. It also provides an API that allows querying of these details.

## Installation

### Prerequisites
- Go 1.19+
- Docker and Docker Compose
- PostgreSQL
- Git

### Steps

1. Clone the repository:
    ```bash
    git clone https://github.com/muhammadfarhankt/omniFlix.git
    cd omniFlix
    ```

2. Install dependencies:
    ```bash
    go mod download
    ```

3. Configure environment variables by editing the `.env` file:
    ```bash
    cp .env.example .env
    # Edit the `.env` file with your configuration
    ```

4. Run the application:
    ```bash
    go run main.go
    ```

## Project Structure

```plaintext
omniFlix/
├── api/                # API implementation files
├── db/                 # Database migration and setup files
├── indexer/            # Core indexer logic
├── main.go             # Entry point of the application
├── Dockerfile          # Docker configuration for the application
├── docker-compose.yml  # Docker Compose setup for multi-container deployment
├── .env                # Environment variables configuration
├── Makefile            # Contains various build commands
├── go.mod              # Go module dependencies
├── go.sum              # Checksums for dependencies
├── README.md           # Project documentation
└── .gitignore          # Git ignore file
```

## Usage

1. Start the application:
    ```bash
    go run main.go
    ```

2. Access the API to query block data. Example:
    ```bash
    curl http://localhost:8080/api/block/1
    ```

## Configuration

- `.env`: Store your environment variables here.
    - `POSTGRES_USER`: Username for PostgreSQL
    - `POSTGRES_PASSWORD`: Password for PostgreSQL
    - `POSTGRES_DB`: Database name for PostgreSQL
    - `BLOCKCHAIN_API_URL`: URL for accessing the Omniflixhub blockchain

## Docker Setup

1. Build and run the services using Docker Compose:
    ```bash
    docker-compose up --build
    ```

2. To stop the services:
    ```bash
    docker-compose down
    ```

## Makefile Commands

The `Makefile` contains useful commands to manage the project. Here are some examples:

- Build the project:
    ```bash
    make build
    ```

- Run tests:
    ```bash
    make test
    ```

- Clean up:
    ```bash
    make clean
    ```

## API Documentation

The API provides a single endpoint:

*   **`GET /block/:height`**

    Fetches the details of the block with the specified height.

    **Parameters:**

    *   `height`: The height of the block (integer).

    **Response:**

    *   If the block is found in the database or can be fetched from the blockchain, the API returns a JSON object with the block details (height, block ID, proposer address, number of transactions, timestamps, and other details).
    *   If there's an error, the API returns a JSON object with an error message.
 
    Response:
```plaintext
{
  "height": 14041989,
  "block_id": "E1677CD5F68547CF2A4E0781C26A0D30E48887291CFE3CD0883E2B790FC03B6A",
  "num_transactions": 0,
  "proposer": "032B564B7C99BB9C127F8CDE514C54F167D84979",
  "created_at": "2024-09-23T15:01:50.44084+05:30",
  "updated_at": "2024-09-23T16:17:52.44333+05:30",
  "deleted_at": {
    "Time": "0001-01-01T00:00:00Z",
    "Valid": false
  },
  "details": null
}
```


## Code Structure

The project is organized into the following packages:

*   `api`: Handles the API endpoint and request handling.
*   `db`: Manages the database connection and table creation.
*   `indexer`: Contains the core logic for fetching and indexing block data.

## Further Improvements

*   More robust error handling: Implement retry mechanisms, exponential backoff, and more informative error responses.
*   API documentation: Generate detailed API documentation using a tool like Swagger.
*   Data validation: Add validation for the data fetched from the blockchain APIs.
*   Caching: Implement caching to improve performance for frequently accessed blocks.
*   Metrics: Add metrics to monitor the indexer's performance and health.
*   Testing: Write unit and integration tests to ensure code quality and reliability.
