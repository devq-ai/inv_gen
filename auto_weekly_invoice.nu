#!/usr/bin/env nu

# Automated Weekly Invoice Generation
# Runs every Monday at 9 AM to generate and send last week's invoice
# Designed for cron job execution

const LOG_FILE = "/Users/dionedge/devqai/inv_gen/logs/auto_invoice.log"
const INVOICE_DIR = "/Users/dionedge/devqai/inv_gen"

# Ensure log directory exists
mkdir ($LOG_FILE | path dirname)

# Log function
def log [message: string] {
    let timestamp = (date now | format date "%Y-%m-%d %H:%M:%S")
    print $"[($timestamp)] ($message)"
    echo $"[($timestamp)] ($message)" | save --append $LOG_FILE
}

log "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
log "🚀 Automated Weekly Invoice Generation - Starting"
log "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Change to invoice directory
cd $INVOICE_DIR

# Check if .env exists
if not (".env" | path exists) {
    log "❌ ERROR: .env file not found"
    log "   Email credentials required for invoice sending"
    exit 1
}

log "✅ Found .env configuration"

# Check if FastAPI server is already running
let server_check = (do {
    http get http://localhost:8000/health
} | complete)

let server_was_running = ($server_check.exit_code == 0)

if $server_was_running {
    log "✅ FastAPI server already running"
} else {
    log "⚙️  Starting FastAPI server..."

    # Start server in background
    bash -c "cd /Users/dionedge/devqai/inv_gen && python3 main.py > /tmp/fastapi_auto.log 2>&1 &"

    # Wait for server to start
    sleep 5sec

    let startup_check = (do {
        http get http://localhost:8000/health
    } | complete)

    if $startup_check.exit_code != 0 {
        log "❌ ERROR: Failed to start FastAPI server"
        log "   Check /tmp/fastapi_auto.log for details"
        exit 1
    }

    log "✅ FastAPI server started successfully"
}

# Generate invoice for last week
log "📄 Generating invoice for last week..."

let generate_result = (do {
    http post http://localhost:8000/invoice/generate
} | complete)

if $generate_result.exit_code != 0 {
    log "❌ ERROR: Failed to generate invoice"
    log $"   Response: ($generate_result.stderr)"

    # Stop server if we started it
    if not $server_was_running {
        log "🛑 Stopping FastAPI server..."
        bash -c "pkill -f 'python3 main.py'"
    }

    exit 1
}

let response = ($generate_result.stdout | from json)
let invoice_number = ($response | get invoice_number)
let total_amount = ($response | get total_amount)
let total_hours = ($response | get total_hours)
let status = ($response | get status)

log "✅ Invoice generated successfully"
log $"   Invoice Number: ($invoice_number)"
log $"   Total Amount: $($total_amount)"
log $"   Total Hours: ($total_hours)"
log $"   Status: ($status)"

# Wait a moment for email to send
log "⏳ Waiting for approval email to send..."
sleep 3sec

# Check email was sent successfully
if ($response | get message | str contains "sent to") {
    log "✅ Approval email sent to dion@devq.ai"
    log "   📧 Email includes PDF attachment"
    log "   📧 Subject: Invoice ($invoice_number) - Pending Approval"
} else {
    log "⚠️  WARNING: Email may not have been sent"
    log $"   Message: ($response | get message)"
}

# Update database to mark as submitted
log "📤 Marking invoice as submitted in database..."

let submit_result = (do {
    python3 invoice_cli.py submit $invoice_number
} | complete)

if $submit_result.exit_code == 0 {
    log "✅ Invoice marked as submitted in database"
} else {
    log "⚠️  WARNING: Failed to update database status"
    log $"   Error: ($submit_result.stderr)"
}

# Generate statistics
log "📊 Generating statistics..."

let stats_result = (do {
    python3 invoice_cli.py stats
} | complete)

if $stats_result.exit_code == 0 {
    log "📊 Current Statistics:"
    log $"($stats_result.stdout)"
} else {
    log "⚠️  Could not generate statistics"
}

# Stop server if we started it
if not $server_was_running {
    log "🛑 Stopping FastAPI server..."
    bash -c "pkill -f 'python3 main.py'"
    sleep 2sec
    log "✅ Server stopped"
}

log "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
log "✅ Automated Invoice Generation - COMPLETED"
log "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
log ""
log "📧 NEXT STEPS:"
log "   1. Check email at dion@devq.ai"
log "   2. Review invoice PDF attachment"
log "   3. Approve via: http post http://localhost:8000/invoice/approve/($invoice_number)"
log "   4. Or use: python3 invoice_cli.py paid ($invoice_number) after payment received"
log ""

exit 0
