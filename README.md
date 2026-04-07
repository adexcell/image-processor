# Image Processor

A background image processing service built with Go and Apache Kafka.

## Features

- **HTTP API**: Upload images, get processed results, and delete images.
- **Background Processing**: Resize and thumbnail generation using Kafka for task queueing.
- **Web UI**: Simple browser-based interface to interact with the service.
- **Robustness**: Uses the `wbf` framework for logging, configuration, and Kafka integration.

## Architecture

1. **API Service**: Receives image uploads, saves them to local storage, and publishes a task to Kafka.
2. **Worker Service**: Consumes tasks from Kafka, processes images using the `imaging` library, and updates the status.
3. **Kafka**: Acts as the message broker between the API and the Worker.
4. **Local Storage**: Stores original images, processed results, and status information.

## API Endpoints

- `POST /api/upload`: Upload an image (form-data: `image`).
- `GET /api/image/:id`: Get the processed image.
- `GET /api/status/:id`: Get the transformation status.
- `DELETE /api/image/:id`: Delete an image and its results.
- `GET /api/list`: List all uploaded images and their statuses.

## Web Interface

Access the web interface at `http://localhost:8080/`.

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Make (optional)

### Running with Docker

```bash
make up
# or
docker-compose up --build
```

The API will be available at `http://localhost:8080`.

### Running Locally

You'll need a running Kafka instance. Update `config.yaml` or use environment variables to point to your Kafka brokers.

```bash
# Start API
make run-api

# Start Worker
make run-worker
```

## Configuration

Configuration can be managed via `config.yaml` or environment variables:

| Env Var | Description | Default |
|---------|-------------|---------|
| `HTTP_PORT` | Port for the API server | `8080` |
| `KAFKA_BROKERS` | List of Kafka brokers | `localhost:9092` |
| `KAFKA_TOPIC` | Kafka topic for tasks | `image-tasks` |
| `STORAGE_UPLOADS_DIR` | Directory for original images | `./uploads` |
| `STORAGE_PROCESSED_DIR` | Directory for processed images | `./processed` |

## License

MIT
