Below is the **README.md** file for the provided service setup.

---

# Image Analysis Platform Lite

This service provides APIs for image management, including **uploading images**, **retrieving image metadata**, and **basic image analysis**. The system is designed to be modular, efficient, and scalable.

## Overview

The service offers core CRUD (Create, Read, Update, Delete) APIs for managing image metadata. It allows users to:
- Upload images.
- Retrieve image details.
- Perform basic analysis and fetch results.

This repository also includes tools for **database migrations** and a structured project setup to ensure clean maintainability.

---

## Pre-requisites

Before setting up the project locally, ensure you have the following:

1. **Golang** (version 1.20 or later)
2. **PostgreSQL** (for database management)
3. **Make** (to run the Makefile scripts)

---

## Directory Structure

```plaintext
.
├── bin/                        # Output binaries (gitignored)
├── cmd/                        # Entry points for the service
│   ├── server/                 # Main API server code
│   │   └── main.go             # API server entry point
│   └── migration/              # Database migration logic
│       └── main.go
├── config/                     # App config
├── pkg/                        # depedencies
├── internal/                   # Service logic and core components
├── Makefile                    # Build automation
└── README.md                   # Project setup documentation
```

---

## Setup Instructions

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone <repository-url>
cd <repository-name>
```

---

### 2. Install Dependencies

Ensure Go modules are initialized and dependencies are downloaded:

```bash
go mod tidy
```

---

### 3. Set Up PostgreSQL

Create a PostgreSQL database and update the environment variables:

1. Run PostgreSQL on your machine (e.g., using Docker or directly on your OS).
2. Create a database:

```sql
CREATE DATABASE image_service;
```

3. Set the following environment variables:

```bash
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="<your-db-user>"
export DB_PASSWORD="<your-db-password>"
export DB_NAME="image_service"
```

---

### 4. Build the Binaries

The `Makefile` automates building the service binaries.

#### Build the API Server:

```bash
make go-build-api
```

The output binary will be located in the `bin/` directory:

```plaintext
bin/api
```

#### Build the Migration Tool:

```bash
make go-build-migration
```

The migration binary will also be located in the `bin/` directory:

```plaintext
bin/migration
```

---

### 5. Run Database Migrations

Execute database migrations using the migration binary:

```bash
./bin/migration up
```

Ensure your database schema is initialized properly.

---

### 6. Start the API Server

Run the API server:

```bash
./bin/api
```

The server will start and be accessible at:

```plaintext
http://localhost:8081
```

---