#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# ✅ Load environment variables from .env
if [ -f .env ]; then
  echo "📦 Loading environment variables from .env..."
  set -o allexport
  source .env
  set +o allexport
else
  echo "❌ .env file not found. Aborting."
  exit 1
fi

# echo "🧹 Removing old containers, volumes, and images..."
# docker-compose down --rmi all --volumes --remove-orphans

echo "🚀 Starting services..."
docker-compose up --build -d

echo "⏳ Waiting for MySQL to be ready..."
# Wait for MySQL to become healthy
for i in {1..30}; do
  STATUS=$(docker inspect --format='{{.State.Health.Status}}' task-manager-mysql-1 2>/dev/null || echo "not_found")
  if [[ "$STATUS" == "healthy" ]]; then
    echo "✅ MySQL is healthy and ready."
    break
  fi
  echo "⏳ Waiting... attempt $i"
  sleep 2
done

if [[ "$STATUS" != "healthy" ]]; then
  echo "❌ MySQL did not become ready in time. Aborting."
  exit 1
fi

echo "📥 Importing init.sql into MySQL database..."
docker exec -i task-manager-mysql-1 mysql -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < ./init.sql

echo "✅ All set. Application started successfully."
echo ""
echo "🌐 Frontend available at: http://localhost:5173"
echo "🛠️  Backend API at:       http://localhost:8080"
echo "📦 MySQL Database on:     port 3306"
echo ""
echo "🪵 View logs with: docker-compose logs -f"
