# Omniflix Block Indexer

This repository contains the block indexer for the Omniflixhub blockchain. The indexer is responsible for storing block details in a Postgres database and providing an API to query the data.

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
    git clone https://github.com/username/omniFlix.git
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

The API provides several endpoints to interact with the block data. Below is an example of how to fetch block details:

- `GET /api/block/{block_id}`: Fetch details of a specific block.

Example request:
```bash
curl http://localhost:8080/block/14041989
```

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
