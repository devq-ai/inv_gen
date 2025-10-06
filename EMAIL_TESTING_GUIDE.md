# 📧 Email Invoice System - Testing Guide

**Complete guide for testing the invoice email system with Nu shell**

---

## 🎯 Overview

The invoice email system is the **MOST CRITICAL** component of the invoicing process. Without email delivery, invoices cannot reach InfoObjects for payment processing.

### What Gets Emailed:

1. **Approval Email** (to dion@devq.ai)
   - Invoice PDF attached
   - Summary of hours and amount
   - Link to approve endpoint

2. **Final Invoice** (to InfoObjects)
   - TO: `infoobjects@bill.com`, `timesheets@infoobjects.com`
   - CC: `dion@devq.ai` (for verification)
   - Invoice PDF attached
   - Professional formatted email

---

## ⚙️ Setup Requirements

### 1. Gmail App Password

You need a Gmail App Password (NOT your regular Gmail password).

**Get it here:** https://myaccount.google.com/apppasswords

```bash
# Create .env file
cat > .env << 'EOF'
GMAIL_ADDRESS=dion@devq.ai
GMAIL_APP_PASSWORD=your_16_char_app_password_here
EOF
```

### 2. Python Dependencies

```bash
pip install -r requirements.txt
```

Required packages:
- `fastapi` - Web framework
- `uvicorn` - ASGI server
- `aiosmtplib` - Async SMTP client
- `pydantic-settings` - Configuration management

---

## 🚀 Quick Start (Nu Shell)

### Start the Server

```nu
cd ~/devqai/inv_gen
./start_server.nu
```

This will:
- ✅ Check for `.env` file
- ✅ Verify dependencies
- ✅ Start FastAPI on port 8000
- ✅ Enable email functionality

### Run Email Tests

In a **new terminal**:

```nu
cd ~/devqai/inv_gen
./test_email.nu
```

This will:
1. Generate a new invoice
2. Send approval email (with PDF attached)
3. Approve the invoice
4. Send to InfoObjects (with PDF attached, CC to dion@devq.ai)
5. Verify both emails were sent

---

## 📋 Manual Testing (Nu Shell)

### Test 1: Generate Invoice & Send for Approval

```nu
# Generate invoice
let response = (http post http://localhost:8000/invoice/generate)

# View response
print ($response | to json --indent 2)

# Expected output:
# {
#   "invoice_number": "N001",
#   "invoice_path": "./invoices/invoice_N001.pdf",
#   "total_hours": 40.0,
#   "total_amount": 3200.0,
#   "status": "pending_approval",
#   "message": "Invoice generated and sent to dion@devq.ai for approval"
# }
```

**Check your email:** Look for approval email at `dion@devq.ai` with PDF attached.

### Test 2: Approve & Send to InfoObjects

```nu
# Get invoice number from previous response
let invoice_num = "N001"

# Approve and send
let approval = (http post $"http://localhost:8000/invoice/approve/($invoice_num)")

# View response
print ($approval | to json --indent 2)

# Expected output:
# {
#   "invoice_number": "N001",
#   "invoice_path": "./invoices/invoice_N001.pdf",
#   "total_hours": 40.0,
#   "total_amount": 3200.0,
#   "status": "approved",
#   "message": "Invoice sent to infoobjects@bill.com and timesheets@infoobjects.com"
# }
```

**Check your email:**
1. InfoObjects should receive the invoice
2. You should be CC'd at `dion@devq.ai`
3. PDF should be attached to both emails

---

## 🔍 Verification Checklist

After running tests, verify:

### Email 1: Approval Email
- [ ] Received at `dion@devq.ai`
- [ ] Subject: "Invoice N001 - Pending Approval"
- [ ] PDF attached: `invoice_N001.pdf`
- [ ] Body shows invoice summary
- [ ] Body shows approval instructions

### Email 2: InfoObjects Invoice
- [ ] Received at `infoobjects@bill.com`
- [ ] Received at `timesheets@infoobjects.com`
- [ ] CC'd to `dion@devq.ai` ✅ **VERIFICATION COPY**
- [ ] Subject: "Weekly Invoice N001 - Dion Edge"
- [ ] PDF attached: `invoice_N001.pdf`
- [ ] Professional formatted body
- [ ] Shows work summary and payment terms

---

## 🐛 Troubleshooting

### Server Won't Start

```nu
# Check if .env exists
ls .env

# Check if port 8000 is available
lsof -i :8000

# View server logs
python3 main.py
```

### Email Not Sending

```nu
# Test Gmail credentials
python3 -c "
import os
from dotenv import load_dotenv
load_dotenv()
print(f'Gmail: {os.getenv(\"GMAIL_ADDRESS\")}')
print(f'Password set: {bool(os.getenv(\"GMAIL_APP_PASSWORD\"))}')
"
```

**Common Issues:**

1. **"Authentication failed"**
   - Wrong app password
   - Using regular Gmail password instead of app password
   - App password has spaces (remove them)

2. **"Connection refused"**
   - Check internet connection
   - Gmail SMTP might be blocked by firewall
   - Try port 465 instead of 587

3. **"Attachment too large"**
   - Gmail limit is 25MB
   - Invoice PDFs should be ~100KB

### Check Email Logs

```nu
# View FastAPI logs
tail -f /tmp/fastapi.log

# Or run server with verbose output
python3 main.py --log-level debug
```

---

## 📊 API Endpoints

### Health Check
```nu
http get http://localhost:8000/health
```

### Generate Invoice
```nu
http post http://localhost:8000/invoice/generate
```

### Approve Invoice
```nu
http post http://localhost:8000/invoice/approve/N001
```

### List All Invoices
```nu
http get http://localhost:8000/invoice/list
```

### Download Invoice PDF
```nu
http get http://localhost:8000/invoice/download/N001 | save invoice_N001.pdf
```

---

## 🎨 Email Templates

### Approval Email Template

```
Subject: Invoice N001 - Pending Approval

Invoice Ready for Approval

Invoice Number: N001
Total Hours: 40.0 hours
Total Amount: $3,200.00

Please review the attached invoice. If approved, click the button below to send it to InfoObjects.

Next Steps:
1. Review the attached invoice PDF
2. If approved, use the API endpoint to send: POST /invoice/approve/N001
3. The invoice will be automatically sent to:
   • infoobjects@bill.com
   • timesheets@infoobjects.com

[PDF Attachment: invoice_N001.pdf]
```

### InfoObjects Email Template

```
Subject: Weekly Invoice N001 - Dion Edge

Weekly Timesheet & Invoice

Dear InfoObjects Team,

Please find attached my weekly invoice for your review and processing.

Invoice Number: N001
Total Hours Worked: 40.0 hours
Total Amount Due: $3,200.00

Work Summary:
• Monday through Friday: 8 hours per day
• Total: 40.0 hours @ $80/hour

Please process this invoice according to our agreed payment terms (Net 15).

If you have any questions or require additional information, please don't hesitate to contact me.

Best regards,
Dion Edge
dion@devq.ai

[PDF Attachment: invoice_N001.pdf]
```

---

## 🔐 Security Notes

### Environment Variables

**NEVER commit `.env` to Git!**

The `.gitignore` file includes `.env`, but verify:

```nu
cat .gitignore | grep ".env"
```

### App Password Safety

- Store in `.env` only
- Don't share or commit
- Rotate periodically
- Use only for this application

### Email Content

- No sensitive data beyond invoice details
- PDFs are password-protected (optional)
- TLS encryption for SMTP

---

## 📈 Production Workflow

### Weekly Invoice Process

```nu
# Monday morning: Generate last week's invoice
cd ~/devqai/inv_gen
./start_server.nu

# In another terminal
cd ~/devqai/inv_gen

# Generate invoice for last week
http post http://localhost:8000/invoice/generate

# Check email for approval
# Review PDF attachment

# Approve and send to InfoObjects
http post http://localhost:8000/invoice/approve/N027

# Verify CC email received
# Confirm InfoObjects received invoice
```

### Automation Option

Create a cron job (using Nu shell):

```nu
# Add to crontab
# Every Monday at 9 AM
0 9 * * 1 cd ~/devqai/inv_gen && ./auto_invoice.nu
```

---

## 🎯 Success Criteria

An invoice email is successful when:

1. ✅ PDF is generated correctly
2. ✅ Approval email sent to dion@devq.ai
3. ✅ PDF is attached to approval email
4. ✅ Invoice sent to both InfoObjects emails
5. ✅ CC copy received at dion@devq.ai
6. ✅ PDF is attached to InfoObjects email
7. ✅ Email formatting is professional
8. ✅ All invoice details are accurate

**Without email delivery, the invoice system is incomplete!**

---

## 📚 Additional Resources

- **Gmail App Passwords**: https://support.google.com/accounts/answer/185833
- **FastAPI Docs**: https://fastapi.tiangolo.com
- **Nu Shell Docs**: https://nushell.sh
- **SMTP Debugging**: Use `openssl s_client -starttls smtp -connect smtp.gmail.com:587`

---

## 🆘 Support

If you encounter issues:

1. Check this guide's troubleshooting section
2. Review FastAPI logs
3. Verify `.env` configuration
4. Test Gmail credentials manually
5. Check firewall/network settings

---

**Remember: Email delivery is THE most important feature of the invoice system!** 📧✅