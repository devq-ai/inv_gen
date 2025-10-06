# ⏰ Automated Weekly Invoice Generation - Cron Setup

**Automated invoice generation and email delivery every Monday at 9 AM**

---

## 🎯 Overview

This cron job will automatically:
1. Generate an invoice for the previous week (Monday-Friday)
2. Send approval email to `dion@devq.ai` with PDF attached
3. Update database to mark as submitted
4. Log all activities for review

**Schedule:** Every Monday at 9:00 AM Central Time

---

## 📋 Prerequisites

### 1. Email Configuration (Required!)
```bash
# Verify .env exists with Gmail credentials
cat ~/devqai/inv_gen/.env
```

Must contain:
```bash
GMAIL_ADDRESS=dion@devq.ai
GMAIL_APP_PASSWORD=your_16_char_app_password
```

### 2. Python Dependencies
```bash
cd ~/devqai/inv_gen
pip install -r requirements.txt
```

### 3. Test Script Manually
```bash
cd ~/devqai/inv_gen
./auto_weekly_invoice.nu
```

**Verify:**
- ✅ Invoice generates
- ✅ Email sends to dion@devq.ai
- ✅ PDF is attached
- ✅ Database updates
- ✅ Log file created in `logs/auto_invoice.log`

---

## 🚀 Cron Installation

### Method 1: Nu Shell Cron (Recommended)

Edit your crontab:
```bash
crontab -e
```

Add this line for **9 AM Central Time every Monday**:
```cron
# Automated Weekly Invoice Generation
# Runs every Monday at 9:00 AM Central Time (CDT/CST)
# Timezone: America/Chicago
0 9 * * 1 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu >> /Users/dionedge/devqai/inv_gen/logs/cron.log 2>&1
```

### Method 2: Bash Wrapper (Alternative)

If Nu shell path issues, create wrapper:
```bash
cat > ~/devqai/inv_gen/auto_weekly_invoice_wrapper.sh << 'EOF'
#!/bin/bash
export PATH="/usr/local/bin:/usr/bin:/bin:/opt/homebrew/bin:$PATH"
cd /Users/dionedge/devqai/inv_gen
/usr/bin/env nu auto_weekly_invoice.nu
EOF

chmod +x ~/devqai/inv_gen/auto_weekly_invoice_wrapper.sh
```

Then add to crontab:
```cron
0 9 * * 1 /Users/dionedge/devqai/inv_gen/auto_weekly_invoice_wrapper.sh >> /Users/dionedge/devqai/inv_gen/logs/cron.log 2>&1
```

---

## 🕐 Schedule Options

### Every Monday at 9 AM
```cron
0 9 * * 1 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu
```

### Every Monday at 8 AM
```cron
0 8 * * 1 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu
```

### Every Friday at 5 PM (end of week)
```cron
0 17 * * 5 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu
```

### First Monday of Each Month at 9 AM
```cron
0 9 1-7 * 1 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu
```

---

## 📊 Monitoring & Logs

### View Auto-Invoice Log
```nu
cat ~/devqai/inv_gen/logs/auto_invoice.log | tail -n 50
```

### View Cron Execution Log
```nu
cat ~/devqai/inv_gen/logs/cron.log | tail -n 50
```

### Watch Live (during scheduled run)
```nu
tail -f ~/devqai/inv_gen/logs/auto_invoice.log
```

### Check Last Run Status
```nu
cat ~/devqai/inv_gen/logs/auto_invoice.log | lines | last 20
```

---

## ✅ Verification After First Run

### 1. Check Cron Executed
```bash
# View system cron logs (macOS)
log show --predicate 'process == "cron"' --last 1h

# Or check if our log was updated
ls -la ~/devqai/inv_gen/logs/auto_invoice.log
```

### 2. Verify Email Sent
- Check `dion@devq.ai` inbox
- Subject: "Invoice N### - Pending Approval"
- PDF attachment should be present

### 3. Check Database Updated
```nu
cd ~/devqai/inv_gen
python3 invoice_cli.py list | tail -n 1
```

Should show latest invoice with `submitted=1`

### 4. Review Log File
```nu
cat ~/devqai/inv_gen/logs/auto_invoice.log | tail -n 30
```

Look for:
- ✅ "Invoice generated successfully"
- ✅ "Approval email sent to dion@devq.ai"
- ✅ "Invoice marked as submitted in database"
- ✅ "COMPLETED"

---

## 🔧 Troubleshooting

### Cron Job Not Running

**Check if cron service is active:**
```bash
# macOS
sudo launchctl list | grep cron

# If not running, start it
sudo launchctl load -w /System/Library/LaunchDaemons/com.vix.cron.plist
```

**Verify crontab entry:**
```bash
crontab -l | grep invoice
```

**Test manually at any time:**
```bash
cd ~/devqai/inv_gen
./auto_weekly_invoice.nu
```

### Email Not Sending

**Check .env credentials:**
```bash
cat ~/devqai/inv_gen/.env | grep GMAIL
```

**Test email manually:**
```bash
cd ~/devqai/inv_gen
./start_server.nu
# In another terminal:
./test_email.nu
```

**Check log for errors:**
```bash
grep -i "error\|fail" ~/devqai/inv_gen/logs/auto_invoice.log
```

### Permission Issues

**Make script executable:**
```bash
chmod +x ~/devqai/inv_gen/auto_weekly_invoice.nu
```

**Check log directory permissions:**
```bash
mkdir -p ~/devqai/inv_gen/logs
chmod 755 ~/devqai/inv_gen/logs
```

### Path Issues

**Full path cron entry (most reliable):**
```cron
0 9 * * 1 /bin/bash -c 'export PATH="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:$PATH" && cd /Users/dionedge/devqai/inv_gen && /opt/homebrew/bin/nu auto_weekly_invoice.nu' >> /Users/dionedge/devqai/inv_gen/logs/cron.log 2>&1
```

---

## 📧 Manual Approval Workflow

After the cron job runs, you'll receive an approval email. To complete the workflow:

### Option 1: Via API
```nu
# Start server if not running
cd ~/devqai/inv_gen
./start_server.nu

# In another terminal, approve
http post http://localhost:8000/invoice/approve/N027
```

### Option 2: Via CLI (when payment received)
```nu
cd ~/devqai/inv_gen
python3 invoice_cli.py paid N027
```

---

## 🔐 Security Considerations

### Protect Credentials
```bash
# .env should not be readable by others
chmod 600 ~/devqai/inv_gen/.env

# Verify
ls -la ~/devqai/inv_gen/.env
# Should show: -rw-------
```

### Log Rotation
```bash
# Create log rotation config
cat > ~/devqai/inv_gen/rotate_logs.nu << 'EOF'
#!/usr/bin/env nu

# Keep last 10 log files
let log_dir = "/Users/dionedge/devqai/inv_gen/logs"
let max_logs = 10

cd $log_dir

# Archive current log if it exists and is large
if ("auto_invoice.log" | path exists) {
    let size = (ls auto_invoice.log | get size | first)
    if $size > 1mb {
        let timestamp = (date now | format date "%Y%m%d_%H%M%S")
        mv auto_invoice.log $"auto_invoice_($timestamp).log"
    }
}

# Remove old logs
ls auto_invoice_*.log 
| sort-by modified 
| reverse 
| skip $max_logs 
| each { |file| rm $file.name }
EOF

chmod +x ~/devqai/inv_gen/rotate_logs.nu
```

Add to crontab (run monthly):
```cron
0 0 1 * * /usr/bin/env nu /Users/dionedge/devqai/inv_gen/rotate_logs.nu
```

---

## 📊 Complete Crontab Example

```cron
# DevQ.ai Automated Systems
SHELL=/bin/bash
PATH=/usr/local/bin:/usr/bin:/bin:/opt/homebrew/bin
MAILTO=dion@devq.ai

# Daily backup (2 AM)
0 2 * * * /Users/dionedge/backups/backup_devqai.sh >> /Users/dionedge/backups/cron.log 2>&1

# Weekly invoice generation (Monday 9 AM)
0 9 * * 1 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu >> /Users/dionedge/devqai/inv_gen/logs/cron.log 2>&1

# Monthly log rotation (1st of month, midnight)
0 0 1 * * /usr/bin/env nu /Users/dionedge/devqai/inv_gen/rotate_logs.nu >> /Users/dionedge/devqai/inv_gen/logs/cron.log 2>&1
```

---

## 🧪 Testing Before Production

### Dry Run Test
```nu
cd ~/devqai/inv_gen

# Run the script manually
./auto_weekly_invoice.nu

# Check output
cat logs/auto_invoice.log | tail -n 50
```

### Test at Custom Time (Today)
```cron
# Temporarily add to crontab for testing
# Run 5 minutes from now
[CURRENT_TIME + 5 minutes] * * * /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu
```

Example: If it's 2:30 PM, add:
```cron
35 14 * * * /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu
```

Wait 5 minutes, then check logs.

---

## 📈 Success Metrics

After cron runs successfully, you should see:

1. ✅ Log file updated with timestamp
2. ✅ Email in dion@devq.ai inbox
3. ✅ PDF attachment in email
4. ✅ Database shows new invoice as submitted
5. ✅ No errors in log file

---

## 🆘 Emergency Disable

### Temporarily Disable Cron
```bash
# Comment out the line in crontab
crontab -e
# Add # at start of invoice line

# Or remove entirely
crontab -e
# Delete the invoice line
```

### Re-enable
```bash
crontab -e
# Remove # or re-add the line
```

---

## 📞 Support Checklist

If automated invoicing fails:

1. **Check cron is running:** `sudo launchctl list | grep cron`
2. **Check crontab entry:** `crontab -l | grep invoice`
3. **Check logs:** `cat ~/devqai/inv_gen/logs/auto_invoice.log`
4. **Test manually:** `cd ~/devqai/inv_gen && ./auto_weekly_invoice.nu`
5. **Check email config:** `cat .env | grep GMAIL`
6. **Verify database:** `python3 invoice_cli.py list`

---

## ✨ Final Setup Command

**Complete setup in one command:**
```bash
cd ~/devqai/inv_gen && \
chmod +x auto_weekly_invoice.nu && \
mkdir -p logs && \
(crontab -l 2>/dev/null; echo "# Automated Weekly Invoice Generation"; echo "0 9 * * 1 /usr/bin/env nu /Users/dionedge/devqai/inv_gen/auto_weekly_invoice.nu >> /Users/dionedge/devqai/inv_gen/logs/cron.log 2>&1") | crontab - && \
echo "✅ Cron job installed!" && \
echo "📅 Will run every Monday at 9 AM" && \
echo "📋 View schedule: crontab -l"
```

---

**Ready for automated weekly invoice generation!** ⏰✅