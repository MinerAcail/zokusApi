#!/bin/sh
# Wait for MinIO to be available
until curl -s http://minio:9000; do
  echo "Waiting for MinIO..."
  sleep 5
done

# Configure MinIO client
mc alias set minio http://minio:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD}

# Run the application
exec air ./cmd/main.go -b 0.0.0.0
    