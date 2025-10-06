#!/usr/bin/env nu

# Test Email Invoice System
# Tests sending invoices via FastAPI with CC to dion@devq.ai

print "🧪 Invoice Email System Test"
print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

# Check if .env file exists
if not (".env" | path exists) {
    print "❌ Error: .env file not found"
    print "Please create .env with Gmail credentials"
    exit 1
}

print "✅ Found .env configuration\n"

# Check if FastAPI server is running
print "🔍 Checking if FastAPI server is running..."
let server_check = (do { http get http://localhost:8000/health } | complete)

if $server_check.exit_code != 0 {
    print "❌ FastAPI server not running"
    print "\n📝 Starting FastAPI server..."
    print "Run in another terminal: cd ~/devqai/inv_gen && python3 main.py\n"
    exit 1
}

print "✅ FastAPI server is running\n"

# Test 1: Generate invoice and send for approval
print "📧 Test 1: Generate invoice and send for approval"
print "   TO: dion@devq.ai (approval)"
print "   CC: dion@devq.ai (verification)"
print "   Attachment: invoice PDF\n"

let generate_response = (http post http://localhost:8000/invoice/generate)

if ($generate_response | get status) == "pending_approval" {
    print "✅ Invoice generated successfully"
    print $"   Invoice Number: ($generate_response | get invoice_number)"
    print $"   Total Amount: $($generate_response | get total_amount)"
    print $"   Total Hours: ($generate_response | get total_hours)"
    print $"   Status: ($generate_response | get status)"
    print $"   Message: ($generate_response | get message)\n"

    let invoice_number = ($generate_response | get invoice_number)

    # Wait a moment for email to send
    print "⏳ Waiting 3 seconds for email to send..."
    sleep 3sec

    # Test 2: Approve and send to InfoObjects
    print "\n📧 Test 2: Approve and send to InfoObjects"
    print "   TO: infoobjects@bill.com, timesheets@infoobjects.com"
    print "   CC: dion@devq.ai (verification)"
    print "   Attachment: invoice PDF\n"

    let approve_response = (http post $"http://localhost:8000/invoice/approve/($invoice_number)")

    if ($approve_response | get status) == "approved" {
        print "✅ Invoice approved and sent to InfoObjects"
        print $"   Invoice Number: ($approve_response | get invoice_number)"
        print $"   Total Amount: $($approve_response | get total_amount)"
        print $"   Status: ($approve_response | get status)"
        print $"   Message: ($approve_response | get message)\n"

        print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        print "✅ ALL TESTS PASSED"
        print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

        print "📬 Check your email at dion@devq.ai:"
        print "   1. Approval email (TO field)"
        print "   2. Approval email (CC field - verification)"
        print "   3. InfoObjects email (CC field - verification)"
        print "\n📄 Both emails should have the invoice PDF attached\n"

    } else {
        print "❌ Failed to approve invoice"
        print ($approve_response | to json)
        exit 1
    }

} else {
    print "❌ Failed to generate invoice"
    print ($generate_response | to json)
    exit 1
}

# Show invoice list
print "\n📋 Current invoice list:"
let invoice_list = (http get http://localhost:8000/invoice/list)
print ($invoice_list | get invoices | to json --indent 2)

print "\n✨ Test complete!"
