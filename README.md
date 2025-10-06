# 📧 DevQ.ai Invoice Management System

**Complete invoice generation, management, and email delivery system**

---

## 🎯 Critical Feature: Email Delivery

**EMAIL IS THE MOST IMPORTANT FEATURE** - Without email delivery, invoices cannot reach InfoObjects for payment processing.

### Email Flow:
1. **Generate Invoice** → Sends approval email to `dion@devq.ai` with PDF
2. **Approve Invoice** → Sends to InfoObjects with CC to `dion@devq.ai` for verification
3. **PDF Attached** → Every email includes the invoice PDF attachment

---

## 🚀 Quick Start

### 1. Configure Email (Required!)

```nu
# Create .env file with Gmail App Password
cat > .env << 'EOF'
GMAIL_ADDRESS=dion@devq.ai
GMAIL_APP_PASSWORD=your_16_char_app_password_here
EOF
```

**Get Gmail App Password:** https://myaccount.google.com/apppasswords

### 2. Install Dependencies

```nu
pip install -r requirements.txt
```

### 3. Start Email Server

```nu
./start_server.nu
```

### 4. Test Email System

```nu
# In another terminal
./test_email.nu
```

---

## 📦 System Components

### 1. **FastAPI Server** (Email Enabled)
```nu
./start_server.nu
```
- **Port:** 8000
- **Email:** Sends invoices via Gmail SMTP
- **API Docs:** http://localhost:8000/docs
- **Critical:** Handles email delivery with PDF attachments

### 2. **Python CLI** (Local Management)
```nu
python3 invoice_cli.py list              # View all invoices
python3 invoice_cli.py show N001         # Invoice details
python3 invoice_cli.py generate N001     # Generate PDF
python3 invoice_cli.py submit N001       # Mark as submitted
python3 invoice_cli.py paid N001         # Mark as paid
python3 invoice_cli.py stats             # Financial statistics
```

### 3. **Go TUI** (Terminal Interface)
```nu
cd tui && ./bin/invoice-tui
```
- Dashboard with financial overview
- Invoice list with filtering
- PDF generation
- Status updates
- **Note:** TUI does NOT send emails (use FastAPI for emails)

---

## 📧 Email System (The Critical Part)

### How It Works:

#### Step 1: Generate & Send for Approval
```nu
http post http://localhost:8000/invoice/generate
```

**Email Sent:**
- **TO:** dion@devq.ai
- **CC:** dion@devq.ai (verification)
- **Subject:** "Invoice N001 - Pending Approval"
- **Attachment:** invoice_N001.pdf

#### Step 2: Approve & Send to InfoObjects
```nu
http post http://localhost:8000/invoice/approve/N001
```

**Email Sent:**
- **TO:** infoobjects@bill.com, timesheets@infoobjects.com
- **CC:** dion@devq.ai (verification)
- **Subject:** "Weekly Invoice N001 - Dion Edge"
- **Attachment:** invoice_N001.pdf

### Email Verification Checklist:
- [ ] Approval email received at dion@devq.ai
- [ ] PDF attached to approval email
- [ ] InfoObjects email sent to both addresses
- [ ] CC copy received at dion@devq.ai
- [ ] PDF attached to InfoObjects email
- [ ] Professional email formatting

---

## 📊 Database

### Quick Stats (Nu Shell)
```nu
sqlite3 invoices.db "
SELECT 
  CASE 
    WHEN paid = 1 THEN 'Paid'
    WHEN submitted = 1 THEN 'Submitted'
    ELSE 'Pending'
  END as status,
  COUNT(*) as count,
  printf('$%.2f', SUM(line_total)) as total
FROM invoices 
GROUP BY status
" | from csv
```

### Current State:
- **Total Invoices:** 27
- **Date Range:** Oct 2025 - Apr 2026
- **Weekly Rate:** $3,200 (40 hours × $80/hour)
- **Total Value:** $86,400

---

## 🧪 Testing

### Complete Workflow Test
```nu
./test_complete_workflow.nu
```

Tests entire lifecycle:
1. Generate PDF
2. Send approval email
3. Update database
4. Send to InfoObjects
5. Verify all steps

### Quick Email Test
```nu
./test_email.nu
```

Tests email delivery:
1. Generate invoice
2. Send approval email
3. Approve invoice
4. Send to InfoObjects
5. Check both emails sent

---

## 📁 Project Structure

```
inv_gen/
├── .env                          # Email credentials (REQUIRED)
├── .env.example                  # Template
├── invoices.db                   # SQLite database (27 invoices)
├── main.py                       # FastAPI server (EMAIL ENABLED)
├── invoice_cli.py                # Python CLI
├── db_invoice_generator.py       # PDF generator
├── email_service.py              # Email delivery (CRITICAL)
├── requirements.txt              # Python dependencies
│
├── invoices/                     # Generated PDFs
│   └── invoice_N*.pdf           # 27 invoice PDFs
│
├── tui/                          # Go TUI
│   ├── main.go
│   ├── models/
│   ├── views/
│   └── bin/invoice-tui          # Compiled binary
│
└── Nu Shell Scripts:
    ├── start_server.nu           # Start FastAPI
    ├── test_email.nu             # Test email system
    └── test_complete_workflow.nu # Full test
```

---

## 🔑 Email Configuration

### Gmail App Password Setup:

1. Go to https://myaccount.google.com/apppasswords
2. Select "Mail" and device type
3. Copy the 16-character password
4. Add to `.env`:
   ```
   GMAIL_APP_PASSWORD=abcd efgh ijkl mnop
   ```
   (spaces will be removed automatically)

### Security:
- ✅ `.env` is in `.gitignore`
- ✅ Never commit credentials
- ✅ Use app passwords, not main password
- ✅ TLS encryption for SMTP

---

## 📚 Documentation

- **[EMAIL_TESTING_GUIDE.md](EMAIL_TESTING_GUIDE.md)** - Complete email testing guide
- **[NU_COMMANDS.md](NU_COMMANDS.md)** - Nu shell command reference
- **[GMAIL_SETUP.md](GMAIL_SETUP.md)** - Email configuration
- **[DATABASE_ARCHITECTURE.md](DATABASE_ARCHITECTURE.md)** - Database schema
- **[TUI_SPECIFICATION.md](TUI_SPECIFICATION.md)** - TUI design

---

## 🎯 Common Workflows

### Weekly Invoice (Production)
```nu
# Monday morning
cd ~/devqai/inv_gen
./start_server.nu

# In another terminal
cd ~/devqai/inv_gen

# Generate invoice for last week
http post http://localhost:8000/invoice/generate

# Check email, review PDF
# Then approve and send
http post http://localhost:8000/invoice/approve/N027

# Verify CC email received at dion@devq.ai
```

### Check Invoice Status
```nu
python3 invoice_cli.py stats
```

### Generate Missing PDFs
```nu
python3 invoice_cli.py generate N001
```

### View in TUI
```nu
cd tui && ./bin/invoice-tui
```

---

## 🐛 Troubleshooting

### Email Not Sending?
```nu
# Check .env exists
ls .env

# Verify credentials loaded
cat .env | grep GMAIL

# Test server
http get http://localhost:8000/health
```

### Server Won't Start?
```nu
# Check if port 8000 is in use
lsof -i :8000

# Kill existing process
kill -9 (lsof -ti:8000)

# Start fresh
./start_server.nu
```

### Database Issues?
```nu
# Check database
sqlite3 invoices.db "SELECT COUNT(*) FROM invoices"

# Recreate if needed
python3 create_invoice_db.py
```

---

## 🎨 Features

### Email System:
- ✅ Gmail SMTP integration
- ✅ PDF attachments
- ✅ Professional HTML formatting
- ✅ CC for verification
- ✅ Approval workflow
- ✅ Error handling

### Python CLI:
- ✅ List all invoices
- ✅ Show details
- ✅ Generate PDFs
- ✅ Update status
- ✅ Financial statistics

### Go TUI:
- ✅ Beautiful dashboard
- ✅ Invoice filtering
- ✅ PDF generation
- ✅ Status updates
- ✅ Neon color scheme

### Database:
- ✅ SQLite (portable)
- ✅ 27 pre-populated invoices
- ✅ Time tracking per day
- ✅ Status tracking

---

## 🔐 Environment Variables

Required in `.env`:
```bash
# Email (REQUIRED for sending invoices)
GMAIL_ADDRESS=dion@devq.ai
GMAIL_APP_PASSWORD=your_app_password

# Optional overrides
BILLING_EMAIL=infoobjects@bill.com
TIMESHEET_EMAIL=timesheets@infoobjects.com
HOURLY_RATE=80.0
PAYMENT_TERMS=Net 15
```

---

## 📈 Statistics

### Current System State:
- **Total Invoices:** 27
- **Total Value:** $86,400
- **Average per Invoice:** $3,200
- **Hours per Week:** 40
- **Hourly Rate:** $80
- **Payment Terms:** Net 15

---

## ✨ Success Criteria

An invoice is successfully processed when:

1. ✅ PDF generated correctly
2. ✅ Approval email sent to dion@devq.ai
3. ✅ PDF attached to approval email
4. ✅ Approval confirmed
5. ✅ Invoice sent to InfoObjects (both emails)
6. ✅ CC copy received at dion@devq.ai
7. ✅ PDF attached to InfoObjects email
8. ✅ Database updated to "submitted"
9. ✅ Payment received (eventually)
10. ✅ Database updated to "paid"

**Without email delivery (steps 2-7), the invoice system is incomplete!**

---

## 🆘 Support

### Check Logs:
```nu
# FastAPI logs
tail -f /tmp/fastapi.log

# Or run with verbose output
python3 main.py
```

### Test Components:
```nu
# Test database
python3 invoice_cli.py list

# Test API
http get http://localhost:8000/health

# Test TUI
cd tui && ./bin/invoice-tui

# Test email
./test_email.nu
```

---

## 📞 Contact

**Developer:** Dion Edge  
**Email:** dion@devq.ai  
**Project:** DevQ.ai Invoice Management System  
**Location:** ~/devqai/inv_gen

---

**Remember: Email delivery is the MOST IMPORTANT feature!** 📧✅