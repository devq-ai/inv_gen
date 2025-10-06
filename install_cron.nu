#!/usr/bin/env nu

# Install Automated Weekly Invoice Generation Cron Job
# One-command setup for automated invoice generation every Monday at 9 AM

print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print "⏰ Installing Automated Weekly Invoice Cron Job"
print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

# Check prerequisites
print "📋 Checking Prerequisites...\n"

# Check if we're in the right directory
if not ("invoices.db" | path exists) {
    print "❌ Error: Must run from ~/devqai/inv_gen directory"
    print "   Run: cd ~/devqai/inv_gen && ./install_cron.nu\n"
    exit 1
}

# Check .env exists
if not (".env" | path exists) {
    print "❌ Error: .env file not found"
    print "   Email credentials required for automated invoicing"
    print "   See GMAIL_SETUP.md for configuration\n"
    exit 1
}

print "✅ Found .env configuration"

# Check if auto_weekly_invoice.nu exists
if not ("auto_weekly_invoice.nu" | path exists) {
    print "❌ Error: auto_weekly_invoice.nu script not found"
    exit 1
}

print "✅ Found automation script"

# Create logs directory
mkdir logs
print "✅ Created logs directory"

# Make scripts executable
chmod +x auto_weekly_invoice.nu
print "✅ Made automation script executable\n"

# Test the script first
print "🧪 Testing automation script...\n"
print "   This will generate a test invoice and send email"
print "   Press Ctrl+C to cancel, or wait 5 seconds to continue..."

sleep 5sec

let test_result = (do {
    ./auto_weekly_invoice.nu
} | complete)

if $test_result.exit_code != 0 {
    print "\n❌ Test failed! Not installing cron job."
    print "   Check logs/auto_invoice.log for details"
    print $"   Error: ($test_result.stderr)\n"
    exit 1
}

print "\n✅ Test successful! Script works correctly.\n"

# Get current crontab
let current_crontab = (do { crontab -l } | complete)
let existing_jobs = if $current_crontab.exit_code == 0 {
    $current_crontab.stdout
} else {
    ""
}

# Check if already installed
if ($existing_jobs | str contains "auto_weekly_invoice.nu") {
    print "⚠️  Cron job already installed!"
    print "\n📅 Current invoice automation entry:"
    print ($existing_jobs | lines | where $it =~ "auto_weekly_invoice")
    print "\n❓ Do you want to reinstall? (yes/no)"

    let response = (input)
    if $response != "yes" {
        print "   Cancelled. No changes made.\n"
        exit 0
    }

    # Remove old entry
    let cleaned = ($existing_jobs | lines | where $it !~ "auto_weekly_invoice" | str join "\n")
    echo $cleaned | crontab -
    print "   Removed old cron entry"
}

# Build new crontab with invoice automation
let invoice_dir = $env.PWD
let cron_line = $"0 9 * * 1 /usr/bin/env nu ($invoice_dir)/auto_weekly_invoice.nu >> ($invoice_dir)/logs/cron.log 2>&1"

let new_crontab = if ($existing_jobs | str trim | is-empty) {
    # No existing crontab, create new one
    $"# Automated Weekly Invoice Generation
# Runs every Monday at 9:00 AM Central Time
# Generated: (date now | format date '%Y-%m-%d %H:%M:%S')
($cron_line)"
} else {
    # Append to existing crontab
    $"($existing_jobs)

# Automated Weekly Invoice Generation
# Runs every Monday at 9:00 AM Central Time
# Generated: (date now | format date '%Y-%m-%d %H:%M:%S')
($cron_line)"
}

# Install new crontab
echo $new_crontab | crontab -

print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
print "✅ CRON JOB INSTALLED SUCCESSFULLY"
print "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"

print "📅 Schedule: Every Monday at 9:00 AM Central Time\n"

print "📋 What will happen automatically:"
print "   1. Generate invoice for previous week"
print "   2. Send approval email to dion@devq.ai"
print "   3. Include PDF attachment"
print "   4. Update database status"
print "   5. Log all activities\n"

print "📊 Monitoring:"
print $"   • Activity Log: ($invoice_dir)/logs/auto_invoice.log"
print $"   • Cron Log: ($invoice_dir)/logs/cron.log"
print $"   • View logs: cat ($invoice_dir)/logs/auto_invoice.log\n"

print "🔍 Verification:"
print "   • View crontab: crontab -l"
print "   • Test manually: ./auto_weekly_invoice.nu"
print "   • Check logs: tail -f logs/auto_invoice.log\n"

print "📧 After each run:"
print "   1. Check email at dion@devq.ai"
print "   2. Review invoice PDF"
print "   3. Approve via API or mark as paid when payment received\n"

print "🛑 To disable:"
print "   • Edit crontab: crontab -e"
print "   • Comment out or delete the invoice line\n"

print "✨ Setup complete! Your invoices will be generated automatically every Monday.\n"
