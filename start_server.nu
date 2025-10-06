#!/usr/bin/env nu

# Start FastAPI Invoice Server
# Starts the invoice generation API with email capabilities

print "🚀 Starting FastAPI Invoice Server\n"
print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

# Check if .env file exists
if not (".env" | path exists) {
    print "❌ Error: .env file not found"
    print "\n📝 Please create .env with the following:"
    print "   GMAIL_ADDRESS=dion@devq.ai"
    print "   GMAIL_APP_PASSWORD=your_app_password_here\n"
    print "📖 See GMAIL_SETUP.md for instructions\n"
    exit 1
}

print "✅ Found .env configuration"

# Check if requirements are installed
print "🔍 Checking Python dependencies..."
let check_deps = (do {
    python3 -c "import fastapi, aiosmtplib, pydantic_settings"
} | complete)

if $check_deps.exit_code != 0 {
    print "❌ Missing dependencies"
    print "\n📦 Installing requirements..."
    pip install -r requirements.txt
    print ""
}

print "✅ Dependencies installed\n"

# Show server info
print "📡 Server Configuration:"
print "   • Host: 0.0.0.0"
print "   • Port: 8000"
print "   • API Docs: http://localhost:8000/docs"
print "   • Health: http://localhost:8000/health\n"

print "📧 Email Features:"
print "   • Generate & send for approval"
print "   • Approve & send to InfoObjects"
print "   • CC to dion@devq.ai for verification"
print "   • PDF attachments included\n"

print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print "🎯 Starting server...\n"

# Start the FastAPI server
python3 main.py
