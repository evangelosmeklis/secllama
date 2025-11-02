#!/bin/bash
# SecLlama Desktop App Launcher

echo "🚀 Starting SecLlama Desktop App..."
echo ""

# Navigate to macapp directory
cd "$(dirname "$0")/macapp" || exit 1

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
    echo ""
fi

# Start the app
echo "🔨 Building and launching Electron app..."
echo "   This may take 20-30 seconds on first run..."
echo ""

npm start

echo ""
echo "❌ App closed or failed to start"

