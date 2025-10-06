#!/bin/bash

# Auto-Invoice Generator Startup Script

echo "=========================================="
echo "Auto-Invoice Generator"
echo "=========================================="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo ""
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo ""
    echo "⚠️  IMPORTANT: Edit .env and add your Gmail App Password"
    echo ""
    echo "Get Gmail App Password:"
    echo "  1. Go to: https://myaccount.google.com/security"
    echo "  2. Enable 2-Step Verification"
    echo "  3. Go to 'App passwords'"
    echo "  4. Generate password for 'Mail'"
    echo "  5. Copy to .env file: GMAIL_APP_PASSWORD=your_password"
    echo ""
    read -p "Press Enter after updating .env to continue..."
fi

# Check if Gmail password is set
if grep -q "GET_THIS_FROM_GOOGLE_ACCOUNT_SECURITY" .env; then
    echo "⚠️  Gmail App Password not configured in .env"
    echo ""
    echo "Email functionality will not work until you:"
    echo "  1. Get App Password from: https://myaccount.google.com/security"
    echo "  2. Update .env: GMAIL_APP_PASSWORD=your_password"
    echo ""
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Create invoices directory
mkdir -p invoices

echo ""
echo "Starting FastAPI server..."
echo ""
echo "API Documentation: http://localhost:8000/docs"
echo "Health Check: http://localhost:8000/health"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Start the server
python3 main.py
