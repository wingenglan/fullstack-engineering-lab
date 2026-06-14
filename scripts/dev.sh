#!/bin/bash
set -e

echo "🚀 Starting local development..."

cleanup() {
    echo ""
    echo "🛑 Stopping all services..."
    kill 0
    wait
    echo "✅ All services stopped"
}

trap cleanup EXIT INT TERM

# Start Go backend
echo "🔧 Starting Go backend..."
cd apps/server
go run ./cmd/server/ &
SERVER_PID=$!
cd ../..

# Wait for server to be ready
sleep 2

# Start Vue frontend
echo "🎨 Starting Vue frontend..."
cd apps/web
npm run dev &
WEB_PID=$!
cd ../..

# Start docs (optional, uncomment if needed)
# echo "📚 Starting docs..."
# cd apps/docs
# npm run dev &
# DOCS_PID=$!
# cd ../..

echo ""
echo "✅ Development services started!"
echo "   Frontend: http://localhost:5173"
echo "   Backend:  http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop all services"

wait
