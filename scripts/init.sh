#!/bin/bash
set -e

echo "🚀 Initializing FullStack Engineering Lab..."

# Copy .env if not exists
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env from .env.example"
else
    echo "⏭️  .env already exists, skipping"
fi

# Check Docker
if command -v docker &> /dev/null; then
    echo "✅ Docker found: $(docker --version)"
else
    echo "❌ Docker not found. Please install Docker first."
    exit 1
fi

# Check Docker Compose
if docker compose version &> /dev/null; then
    echo "✅ Docker Compose found"
else
    echo "❌ Docker Compose not found. Please install Docker Compose first."
    exit 1
fi

# Check Go (optional for local dev)
if command -v go &> /dev/null; then
    echo "✅ Go found: $(go version)"
else
    echo "⚠️  Go not found (optional for Docker mode)"
fi

# Check Node (optional for local dev)
if command -v node &> /dev/null; then
    echo "✅ Node found: $(node --version)"
else
    echo "⚠️  Node not found (optional for Docker mode)"
fi

echo ""
echo "🎉 Initialization complete!"
echo ""
echo "Next steps:"
echo "  make up    - Start all services with Docker"
echo "  make dev   - Start local development"
