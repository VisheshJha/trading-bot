#!/bin/bash

set -e

echo "🔧 Building Docker images..."
docker compose build

echo "🛑 Stopping existing services..."
docker compose down || true

echo "🚀 Starting new deployment..."
docker compose up -d --force-recreate

echo "✅ Deployment completed successfully!"