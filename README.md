# Log Aggregator with MongoDB Backend

## Overview

This project is a simple log aggregator that allows users to ingest, retrieve, and filter logs via HTTP APIs. The logs are stored in a MongoDB database with a structured format, supporting metadata and filtering capabilities.

## Features

- **Log Ingestion**: Push logs with structured data, including levels, messages, and metadata.
- **Retrieve Logs by ID**: Fetch logs using their unique identifier.
- **Filter Logs**: Query logs based on various parameters like level, message, timestamps, and metadata.
- **Delete Logs**: Remove logs before a specific timestamp or by ID.
- **Ingestion Script**: Run the ingest.sh file in scripts to populate your mongodb

## API Endpoints

### 1. Ingest Logs

**Endpoint:**

```
POST /v1.0/logs
```

**Request Example:**

```sh
curl --location 'http://localhost:6060/v1.0/logs' \
--header 'Content-Type: application/json' \
--data '{
    "level": "INFO",
    "message": "Application started successfully",
    "metadata": {
      "service": "api",
      "version": "1.0.0",
      "environment": "PROD",
      "user_id": "12345",
      "ip_address": "192.168.1.100",
      "session_id": "abcde12345"
    }
  }'
```

### 2. Retrieve Log by ID

**Endpoint:**

```
GET /v1.0/logs/{log_id}
```

**Request Example:**

```sh
curl --location 'http://localhost:6060/v1.0/logs/67e8fa498aea23c72b9908da'
```

### 3. Filter Logs

**Endpoint:**

```
GET /v1.0/logs
```

**Query Parameters:**

- `level` - Filter logs by level (e.g., INFO, ERROR, DEBUG).
- `message` - Search for logs containing a specific message.
- `starttime` - Start timestamp (epoch) to filter logs.
- `endtime` - End timestamp (epoch) to filter logs.
- `recent` - Number of recent logs to fetch.
- `metadata` - Filter logs based on metadata fields.

**Request Example:**

```sh
curl --location 'http://localhost:6060/v1.0/logs?level=ERROR&starttime=1743321000&endtime=1743322000'
```

### 4. Delete Logs

**Endpoint:**

```
DELETE /v1.0/logs
```

**Query Parameters:**

- `before` - Delete logs before a specific epoch timestamp.
- `id` - Delete a specific log by ID.

**Request Example:**

```sh
curl --location --request DELETE 'http://localhost:6060/v1.0/logs?before=1743321727'
```

## Setup Instructions

### Prerequisites

- Go 1.18+
- MongoDB (running on `localhost:27017` by default)

### Running the Service

1. Clone the repository:
   ```sh
   git clone https://github.com/bhuvankumar123/klg.git
   cd klg
   ```
2. Install dependencies:
   ```sh
   make goensure
   ```
3. Run the service:
   ```sh
   go run cmd/klg/main.go cmd/klg/flags.go start
   ```

## Environment Variables

| Variable             | Default Value               | Description            |
| -------------------- | --------------------------- | ---------------------- |
| `APP_MONGO_URI`      | `mongodb://localhost:27017` | MongoDB connection URI |
| `APP_MONGO_DATABASE` | `logs`                      | MongoDB database name  |

## License

This project is open-source and available under the MIT License.

## Contributions

Feel free to submit issues or pull requests to improve the project!


