# 🐚 Nu Shell Commands Reference - Invoice System

**Complete Nu shell command reference for the DevQ.ai Invoice Management System**

---

## 🚀 Quick Start Commands

### Start the FastAPI Server
```nu
cd ~/devqai/inv_gen
./start_server.nu
```

### Run Email Tests
```nu
cd ~/devqai/inv_gen
./test_email.nu
```

### Complete Workflow Test
```nu
cd ~/devqai/inv_gen
./test_complete_workflow.nu
```

---

## 📊 Database Queries (Nu Shell)

### View All Invoices
```nu
sqlite3 invoices.db "SELECT invoice_number, line_total, total_hours, submitted, paid FROM invoices" | from csv
```

### Count by Status
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

### Financial Summary
```nu
sqlite3 invoices.db "
SELECT
  printf('$%.2f', SUM(line_total)) as total_invoiced,
  printf('$%.2f', SUM(CASE WHEN paid = 1 THEN line_total ELSE 0 END)) as total_paid,
  printf('$%.2f', SUM(CASE WHEN submitted = 1 AND paid = 0 THEN line_total ELSE 0 END)) as awaiting_payment,
  printf('$%.2f', SUM(CASE WHEN submitted = 0 AND paid = 0 THEN line_total ELSE 0 END)) as pending
FROM invoices
" | from csv
```

### Get Specific Invoice Details
```nu
let invoice_num = "N001"
sqlite3 invoices.db $"
SELECT * FROM invoices 
WHERE invoice_number = '($invoice_num)'
" | from csv
```

### List Pending Invoices
```nu
sqlite3 invoices.db "
SELECT invoice_number, line_total, total_hours 
FROM invoices 
WHERE submitted = 0 AND paid = 0
" | from csv
```

### List Submitted But Unpaid
```nu
sqlite3 invoices.db "
SELECT invoice_number, line_total, total_hours, invoice_create_date
FROM invoices 
WHERE submitted = 1 AND paid = 0
" | from csv
```

---

## 🌐 API Requests (Nu Shell)

### Health Check
```nu
http get http://localhost:8000/health
```

### Generate New Invoice
```nu
http post http://localhost:8000/invoice/generate
```

### Generate for Specific Date
```nu
http post http://localhost:8000/invoice/generate {
    start_date: "2025-10-07"
}
```

### Approve Invoice
```nu
let invoice_num = "N001"
http post $"http://localhost:8000/invoice/approve/($invoice_num)"
```

### List All Invoices (API)
```nu
http get http://localhost:8000/invoice/list
```

### Download Invoice PDF
```nu
let invoice_num = "N001"
http get $"http://localhost:8000/invoice/download/($invoice_num)" | save $"invoice_($invoice_num).pdf"
```

---

## 🐍 Python CLI Commands (Nu Shell)

### List All Invoices
```nu
python3 invoice_cli.py list
```

### Show Invoice Details
```nu
let invoice_num = "N001"
python3 invoice_cli.py show $invoice_num
```

### Generate PDF
```nu
let invoice_num = "N001"
python3 invoice_cli.py generate $invoice_num
```

### Mark as Submitted
```nu
let invoice_num = "N001"
python3 invoice_cli.py submit $invoice_num
```

### Mark as Paid
```nu
let invoice_num = "N001"
python3 invoice_cli.py paid $invoice_num
```

### View Statistics
```nu
python3 invoice_cli.py stats
```

---

## 📧 Email Testing Commands

### Test Email Configuration
```nu
# Check if .env exists
ls .env

# Verify environment variables
cat .env | lines | where $it =~ "GMAIL"
```

### Send Test Email (via API)
```nu
# Start server first
./start_server.nu

# In another terminal, generate and send
http post http://localhost:8000/invoice/generate
```

### Check Email Delivery
```nu
# Monitor server logs
tail -f /tmp/fastapi.log

# Or run server with verbose output
python3 main.py
```

---

## 🖥️ TUI Commands (Go Interface)

### Start TUI
```nu
cd ~/devqai/inv_gen/tui
./bin/invoice-tui
```

### Rebuild TUI
```nu
cd ~/devqai/inv_gen/tui
go build -o bin/invoice-tui
```

### TUI Keyboard Shortcuts
- `q` - Quit (from dashboard)
- `Esc` - Back to dashboard
- `Ctrl+C` - Force quit
- `↑/↓` - Navigate
- `Enter` - Select
- `f` - Cycle filters (List view)
- `g` - Generate PDF
- `s` - Mark as submitted
- `p` - Mark as paid

---

## 🔧 System Maintenance (Nu Shell)

### Check Python Dependencies
```nu
pip list | grep -E "fastapi|uvicorn|aiosmtplib|pydantic"
```

### Install/Update Dependencies
```nu
pip install -r requirements.txt
```

### Backup Database
```nu
let timestamp = (date now | format date "%Y%m%d_%H%M%S")
cp invoices.db $"invoices_backup_($timestamp).db"
```

### Reset All Invoices to Pending
```nu
sqlite3 invoices.db "UPDATE invoices SET submitted = 0, paid = 0"
```

### Reset Specific Invoice
```nu
let invoice_num = "N001"
sqlite3 invoices.db $"UPDATE invoices SET submitted = 0, paid = 0 WHERE invoice_number = '($invoice_num)'"
```

### View Invoice PDFs
```nu
ls invoices/ | where type == file and name =~ "pdf"
```

### Count Generated PDFs
```nu
ls invoices/*.pdf | length
```

---

## 🧪 Testing Workflows

### Full Workflow Test
```nu
# 1. Start server
./start_server.nu

# 2. In new terminal, run complete test
./test_complete_workflow.nu

# 3. Verify email delivery
# Check dion@devq.ai inbox
```

### Quick Email Test
```nu
# With server running
./test_email.nu
```

### Manual Step-by-Step Test
```nu
# Step 1: Check database status
sqlite3 invoices.db "SELECT COUNT(*), SUM(line_total) FROM invoices WHERE submitted = 0" | from csv

# Step 2: Generate invoice via API
let response = (http post http://localhost:8000/invoice/generate)
print ($response | to json --indent 2)

# Step 3: Get invoice number
let invoice_num = ($response | get invoice_number)

# Step 4: Verify PDF exists
ls $"invoices/invoice_($invoice_num).pdf"

# Step 5: Approve and send
http post $"http://localhost:8000/invoice/approve/($invoice_num)"

# Step 6: Check email
print "Check dion@devq.ai for emails"
```

---

## 📊 Analytics & Reporting

### Monthly Revenue Summary
```nu
sqlite3 invoices.db "
SELECT 
  strftime('%Y-%m', invoice_create_date) as month,
  COUNT(*) as invoices,
  printf('$%.2f', SUM(line_total)) as revenue
FROM invoices 
GROUP BY month 
ORDER BY month
" | from csv
```

### Unpaid Invoice Report
```nu
sqlite3 invoices.db "
SELECT 
  invoice_number,
  invoice_create_date,
  due_date,
  printf('$%.2f', line_total) as amount,
  CASE 
    WHEN submitted = 1 THEN 'Awaiting Payment'
    ELSE 'Not Submitted'
  END as status
FROM invoices 
WHERE paid = 0
ORDER BY due_date
" | from csv
```

### Hours Worked Summary
```nu
sqlite3 invoices.db "
SELECT 
  SUM(total_hours) as total_hours,
  AVG(total_hours) as avg_hours_per_invoice,
  COUNT(*) as invoice_count
FROM invoices
" | from csv
```

---

## 🔍 Debugging Commands

### Check Server Status
```nu
# Test connection
http get http://localhost:8000/health | complete

# Check if port is in use
lsof -i :8000
```

### View Recent Logs
```nu
# Python CLI operations
python3 invoice_cli.py list | tail -n 20

# FastAPI server logs (if logging to file)
tail -f /tmp/fastapi.log
```

### Test Email Credentials
```nu
python3 -c "
import os
from dotenv import load_dotenv
load_dotenv()
print(f'Gmail: {os.getenv(\"GMAIL_ADDRESS\")}')
print(f'Password set: {bool(os.getenv(\"GMAIL_APP_PASSWORD\"))}')
print(f'Password length: {len(os.getenv(\"GMAIL_APP_PASSWORD\", \"\"))}')
"
```

### Validate Database Schema
```nu
sqlite3 invoices.db ".schema invoices"
```

### Check File Permissions
```nu
ls -la | where name =~ ".nu"
```

---

## 🎯 Production Workflows

### Weekly Invoice Generation
```nu
# Monday morning workflow
cd ~/devqai/inv_gen

# 1. Check last week's data
sqlite3 invoices.db "SELECT * FROM invoices ORDER BY pk DESC LIMIT 1" | from csv

# 2. Generate new invoice
http post http://localhost:8000/invoice/generate

# 3. Review and approve (get invoice number from response)
let invoice_num = "N027"
http post $"http://localhost:8000/invoice/approve/($invoice_num)"

# 4. Verify email sent
print "✅ Check email for confirmation"

# 5. Update database
python3 invoice_cli.py submit $invoice_num
```

### Month-End Review
```nu
# Current month summary
let current_month = (date now | format date "%Y-%m")

sqlite3 invoices.db $"
SELECT 
  COUNT(*) as invoices,
  SUM(total_hours) as hours,
  printf('$%.2f', SUM(line_total)) as revenue,
  SUM(CASE WHEN paid = 1 THEN 1 ELSE 0 END) as paid_count,
  printf('$%.2f', SUM(CASE WHEN paid = 1 THEN line_total ELSE 0 END)) as paid_amount
FROM invoices 
WHERE strftime('%Y-%m', invoice_create_date) = '($current_month)'
" | from csv
```

---

## 🛠️ Advanced Nu Shell Patterns

### Loop Through Pending Invoices
```nu
let pending = (sqlite3 invoices.db "
SELECT invoice_number FROM invoices 
WHERE submitted = 0 AND paid = 0
" | from csv)

for invoice in $pending {
    let num = ($invoice | get invoice_number)
    print $"Processing ($num)..."
    python3 invoice_cli.py generate $num
}
```

### Batch Generate PDFs
```nu
sqlite3 invoices.db "SELECT invoice_number FROM invoices" 
| from csv 
| each { |row| 
    python3 invoice_cli.py generate ($row | get invoice_number)
}
```

### Custom Report Function
```nu
def invoice-summary [] {
    let data = (sqlite3 invoices.db "
    SELECT 
      COUNT(*) as total,
      SUM(CASE WHEN paid = 1 THEN 1 ELSE 0 END) as paid,
      printf('$%.2f', SUM(line_total)) as total_amount
    FROM invoices
    " | from csv | first)
    
    print "📊 Invoice Summary"
    print $"   Total: ($data.total)"
    print $"   Paid: ($data.paid)"
    print $"   Amount: ($data.total_amount)"
}

# Use it:
invoice-summary
```

---

## 📚 Useful Aliases (Add to ~/.config/nushell/config.nu)

```nu
# Quick navigation
alias inv = cd ~/devqai/inv_gen

# Start services
alias inv-server = cd ~/devqai/inv_gen; ./start_server.nu
alias inv-tui = cd ~/devqai/inv_gen/tui; ./bin/invoice-tui

# Quick commands
alias inv-list = python3 ~/devqai/inv_gen/invoice_cli.py list
alias inv-stats = python3 ~/devqai/inv_gen/invoice_cli.py stats

# Database shortcuts
alias inv-db = sqlite3 ~/devqai/inv_gen/invoices.db
alias inv-status = sqlite3 ~/devqai/inv_gen/invoices.db "SELECT CASE WHEN paid = 1 THEN 'Paid' WHEN submitted = 1 THEN 'Submitted' ELSE 'Pending' END as status, COUNT(*) FROM invoices GROUP BY status" | from csv
```

---

## 🎓 Learning Resources

### Nu Shell Basics
```nu
# Variables
let name = "value"

# Command substitution
let result = (command arg1 arg2)

# Piping
command1 | command2 | command3

# JSON parsing
http get url | from json | select field1 field2

# CSV parsing
command | from csv

# Tables
ls | where size > 1mb | sort-by size
```

### Common Patterns
```nu
# Check if file exists
if ("file.txt" | path exists) { print "exists" }

# String interpolation
let name = "World"
print $"Hello ($name)!"

# Error handling
let result = (do { risky-command } | complete)
if $result.exit_code != 0 {
    print $"Error: ($result.stderr)"
}
```

---

## 🆘 Troubleshooting Guide

### Server Won't Start
```nu
# Check if port is in use
lsof -i :8000

# Kill existing process
kill -9 (lsof -ti:8000)

# Check .env
cat .env | lines | where $it =~ "GMAIL"
```

### Database Issues
```nu
# Verify database exists
ls invoices.db

# Check table structure
sqlite3 invoices.db ".tables"

# Validate data
sqlite3 invoices.db "SELECT COUNT(*) FROM invoices" | from csv
```

### Email Problems
```nu
# Test credentials
python3 -c "from email_service import EmailService; import os; from dotenv import load_dotenv; load_dotenv(); print('Credentials loaded')"

# Check SMTP connection
openssl s_client -starttls smtp -connect smtp.gmail.com:587
```

---

**Remember: All Nu shell scripts use `.nu` extension and are executable!** 🐚✨