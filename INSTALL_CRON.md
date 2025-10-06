# ⏰ Install Invoice Automation Cron Jobs

**Complete automation: Generate → Draft → Remind → Auto-Send**

---

## 🎯 What Gets Installed

Three cron jobs that ensure invoices are ALWAYS sent on time:

1. **Saturday 2:00 AM CST** - Generate PDF, create Gmail draft, send approval request
2. **Sunday 12:00 PM CST** - Send reminder if not approved
3. **Sunday 9:00 PM CST** - **AUTO-SEND if not approved (CRITICAL!)**

---

## 📋 Prerequisites

### 1. Gmail API Setup (REQUIRED)

You need OAuth2 credentials (NOT just app password):

1. Go to https://console.cloud.google.com/apis/credentials
2. Create OAuth 2.0 Client ID (Desktop application)
3. Download JSON and save as `~/devqai/inv_gen/credentials.json`
4. First run will open browser for authorization
5. Token saved to `token.pickle` for future use

### 2. Install Python Dependencies

```bash
cd ~/devqai/inv_gen
pip install -r requirements.txt
```

This installs:
- `google-auth-oauthlib`
- `google-auth-httplib2`
- `google-api-python-client`

### 3. Test Gmail API Authentication

```bash
cd ~/devqai/inv_gen
python3 -c "from gmail_draft_service import GmailDraftService; g = GmailDraftService(); print('✅ Gmail API working')"
```

This will open browser for first-time authorization.

---

## 🚀 Installation

### Option 1: One Command Install

```bash
cd ~/devqai/inv_gen
(crontab -l 2>/dev/null; cat << 'CRON'

# DevQ.ai Automated Invoice System
# Saturday 2 AM: Generate PDF, create draft, send approval
0 2 * * 6 cd /Users/dionedge/devqai/inv_gen && /usr/bin/python3 cron_saturday_generate.py >> logs/cron_saturday.log 2>&1

# Sunday 12 PM: Send reminder if not approved
0 12 * * 0 cd /Users/dionedge/devqai/inv_gen && /usr/bin/python3 cron_sunday_reminder.py >> logs/cron_sunday_reminder.log 2>&1

# Sunday 9 PM: AUTO-SEND if not approved (CRITICAL!)
0 21 * * 0 cd /Users/dionedge/devqai/inv_gen && /usr/bin/python3 cron_sunday_autosend.py >> logs/cron_sunday_autosend.log 2>&1

CRON
) | crontab -
```

### Option 2: Manual Install

```bash
crontab -e
```

Add these lines:

```cron
# DevQ.ai Automated Invoice System
# Saturday 2 AM: Generate PDF, create draft, send approval
0 2 * * 6 cd /Users/dionedge/devqai/inv_gen && /usr/bin/python3 cron_saturday_generate.py >> logs/cron_saturday.log 2>&1

# Sunday 12 PM: Send reminder if not approved
0 12 * * 0 cd /Users/dionedge/devqai/inv_gen && /usr/bin/python3 cron_sunday_reminder.py >> logs/cron_sunday_reminder.log 2>&1

# Sunday 9 PM: AUTO-SEND if not approved (CRITICAL!)
0 21 * * 0 cd /Users/dionedge/devqai/inv_gen && /usr/bin/python3 cron_sunday_autosend.py >> logs/cron_sunday_autosend.log 2>&1
```

---

## ✅ Verify Installation

```bash
# View crontab
crontab -l

# Create logs directory
mkdir -p ~/devqai/inv_gen/logs

# Make scripts executable
chmod +x ~/devqai/inv_gen/cron_saturday_generate.py
chmod +x ~/devqai/inv_gen/cron_sunday_reminder.py
chmod +x ~/devqai/inv_gen/cron_sunday_autosend.py
```

---

## 🧪 Test Before Production

### Test Saturday Generation

```bash
cd ~/devqai/inv_gen
python3 cron_saturday_generate.py
```

**Should:**
- Generate PDF for most recent pending invoice
- Create Gmail draft with PDF attached
- Send approval request email to dion@devq.ai
- Update database with draft_id

**Check:**
- Email in inbox
- Draft in Gmail
- Log file: `cat logs/cron_saturday.log`

### Test Sunday Reminder

```bash
cd ~/devqai/inv_gen
python3 cron_sunday_reminder.py
```

**Should:**
- Check for pending drafts
- Send reminder email if draft not sent
- Update database with reminder timestamp

### Test Sunday Auto-Send

```bash
cd ~/devqai/inv_gen
python3 cron_sunday_autosend.py
```

**Should:**
- Check for pending drafts
- **SEND THE DRAFT** if not already sent
- Mark as submitted in database
- Send confirmation email

⚠️ **WARNING:** This will actually send the invoice! Only test with test data.

---

## 📧 Email Flow

### Saturday 2 AM
**Subject:** "⚠️ Invoice N001 - Approval Required"
**To:** dion@devq.ai
**Content:** Invoice details, approval instructions, deadlines

### Sunday 12 PM (if not approved)
**Subject:** "🚨 REMINDER: Invoice N001 - Action Required"
**To:** dion@devq.ai
**Content:** Reminder with 9-hour warning, quick approve command

### Sunday 9 PM (if STILL not approved)
**Action:** Draft automatically sent to InfoObjects
**Subject:** "✅ Invoice N001 AUTO-SENT"
**To:** dion@devq.ai
**Content:** Confirmation that invoice was auto-sent

---

## 🎮 Manual Control

### Approve Invoice (Before Auto-Send)

```bash
cd ~/devqai/inv_gen
python3 invoice_cli.py approve N001
```

This will:
- Send the Gmail draft to InfoObjects
- Mark as submitted in database
- Prevent auto-send from triggering

### Check Current Status

```bash
cd ~/devqai/inv_gen
python3 invoice_cli.py show N001
```

### View Pending Drafts

```bash
sqlite3 ~/devqai/inv_gen/invoices.db "
SELECT invoice_number, draft_id, approval_sent_at, submitted 
FROM invoices 
WHERE draft_id IS NOT NULL 
ORDER BY approval_sent_at DESC
"
```

---

## 📊 Monitoring

### View Logs

```bash
# Saturday generation log
tail -f ~/devqai/inv_gen/logs/cron_saturday.log

# Sunday reminder log
tail -f ~/devqai/inv_gen/logs/cron_sunday_reminder.log

# Sunday auto-send log (CRITICAL)
tail -f ~/devqai/inv_gen/logs/cron_sunday_autosend.log
```

### Check Last Run

```bash
# Saturday
cat ~/devqai/inv_gen/logs/cron_saturday.log | tail -n 30

# Sunday reminder
cat ~/devqai/inv_gen/logs/cron_sunday_reminder.log | tail -n 30

# Sunday auto-send
cat ~/devqai/inv_gen/logs/cron_sunday_autosend.log | tail -n 30
```

---

## 🚨 Critical Alerts

### If Saturday Generation Fails

**Symptoms:** No approval email received Saturday morning

**Fix:**
```bash
cd ~/devqai/inv_gen
cat logs/cron_saturday.log
# Fix error, then run manually:
python3 cron_saturday_generate.py
```

### If Sunday Auto-Send Fails

**Symptoms:** No confirmation email, invoice not in "Sent" folder

**IMMEDIATE ACTION REQUIRED:**
```bash
# Check log
cat ~/devqai/inv_gen/logs/cron_sunday_autosend.log

# Manually send from Gmail:
# 1. Go to Gmail Drafts
# 2. Find invoice draft
# 3. Send immediately

# Then mark as submitted:
cd ~/devqai/inv_gen
python3 invoice_cli.py submit N001
```

---

## 🛠️ Troubleshooting

### Cron Job Not Running

```bash
# Check if cron is active (macOS)
sudo launchctl list | grep cron

# View system cron logs
log show --predicate 'process == "cron"' --last 1h

# Check crontab syntax
crontab -l
```

### Gmail API Issues

```bash
# Re-authenticate
rm ~/devqai/inv_gen/token.pickle
python3 -c "from gmail_draft_service import GmailDraftService; GmailDraftService()"
```

### Permission Errors

```bash
# Ensure scripts are executable
chmod +x ~/devqai/inv_gen/cron_*.py

# Ensure log directory exists
mkdir -p ~/devqai/inv_gen/logs
chmod 755 ~/devqai/inv_gen/logs
```

---

## 📅 Timeline Example

**Saturday 2:00 AM:**
- ✅ Invoice PDF generated
- ✅ Gmail draft created with PDF attached
- ✅ Approval email sent to dion@devq.ai

**Sunday 12:00 PM:**
- ❓ Check if draft sent
- ⚠️ If NOT sent → Send reminder email
- ✅ If already sent → Skip

**Sunday 9:00 PM:**
- ❓ Check if draft sent
- 🚨 If NOT sent → **AUTO-SEND TO INFOOBJECTS**
- ✅ If already sent → Skip
- 📧 Send confirmation email

---

## 🎯 Success Criteria

The system is working correctly when:

1. ✅ Saturday 2 AM: Approval email received
2. ✅ Gmail draft exists with PDF attached
3. ✅ Sunday 12 PM: Reminder received (if not approved)
4. ✅ Sunday 9 PM: Either manually approved OR auto-sent
5. ✅ Invoice marked as submitted in database
6. ✅ No errors in log files
7. ✅ InfoObjects receives invoice by Sunday 9 PM

**THE CRITICAL METRIC: Invoice sent by Sunday 9 PM → WE GET PAID!**

---

## 🆘 Emergency Disable

If something is broken and you need to stop automation:

```bash
# Comment out cron jobs
crontab -e
# Add # at start of invoice lines

# Or remove completely
crontab -l | grep -v "cron_" | crontab -
```

**Re-enable:**
```bash
# Re-run installation command from above
```

---

## 📞 Quick Reference

**Install:** Run one-command install above
**Check Status:** `crontab -l`
**View Logs:** `tail -f ~/devqai/inv_gen/logs/cron_*.log`
**Manual Approve:** `python3 invoice_cli.py approve N001`
**Emergency Send:** Go to Gmail → Drafts → Send manually

**REMEMBER: If auto-send fails, we don't get paid! Monitor Sunday 9 PM closely.**

---

**Installation complete! Invoice automation now active.** ⏰✅