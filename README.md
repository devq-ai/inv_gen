# Invoice Management System

A streamlined invoice generation and email delivery system for weekly contractor invoicing to InfoObjects.

## Quick Start

### Installation
```bash
# Install Python dependencies
pip install -r requirements.txt

# Build Go TUI (optional)
cd tui && make build
```

### Configure Email
```bash
# Create .env file
cp .env.example .env

# Add your Gmail app password
# Get it from: https://myaccount.google.com/apppasswords
GMAIL_ADDRESS=your-email@gmail.com
GMAIL_APP_PASSWORD=your_16_char_app_password
```

### Run
```bash
# Python CLI
python3 invoice_cli.py list

# Go TUI
cd tui && ./bin/invoice-tui
```

## System Components

### 1. Python CLI (`invoice_cli.py`)
Command-line interface for invoice management.

**Commands:**
```bash
python3 invoice_cli.py list                    # List all invoices
python3 invoice_cli.py show N001               # Show invoice details
python3 invoice_cli.py generate N001           # Generate PDF
python3 invoice_cli.py submit N001             # Mark as submitted
python3 invoice_cli.py paid N001               # Mark as paid
python3 invoice_cli.py stats                   # Show statistics
```

### 2. Go TUI (`tui/`)
Terminal user interface with native email support.

**Features:**
- Dashboard with financial overview
- Invoice list with filtering (all/pending/submitted/paid)
- PDF generation
- Email delivery via native Go SMTP
- Status tracking

**Keyboard Shortcuts:**
- `↑/↓` - Navigate
- `Enter` - Select
- `f` - Cycle filters (List view)
- `g` - Generate PDF
- `s` - Send email & mark submitted
- `p` - Mark as paid
- `q` - Quit/Back
- `Esc` - Back to dashboard

### 3. Email Service (`email_service.py`)
SMTP email delivery with PDF attachments.

**Recipients:**
- **TO:** InfoObjects billing/timesheet emails
- **CC:** dion@wrench.chat (for verification)

### 4. Cron Automation (`cron_weekly_invoice.py`)
Automated weekly invoice generation and email delivery.

**Schedule:** Every Saturday at 2:00 AM CST

**Cron Entry:**
```cron
0 2 * * 6 cd /path/to/inv_gen && /usr/bin/python3 cron_weekly_invoice.py
```

## Database

**Location:** `invoices.db` (SQLite)
**Current State:** 27 invoices from Oct 2025 - Apr 2026

**Schema:**
```sql
CREATE TABLE invoices (
    invoice_number TEXT PRIMARY KEY,
    invoice_create_date TEXT,
    due_date TEXT,
    work_week_start TEXT,
    work_week_end TEXT,
    line_total REAL,
    submitted INTEGER DEFAULT 0,
    paid INTEGER DEFAULT 0,
    -- Daily hours: monday_hours through sunday_hours
    -- Each day: REAL DEFAULT 8.0
);
```

**Quick Stats:**
```bash
# View statistics
python3 invoice_cli.py stats

# Or via SQL
sqlite3 invoices.db "SELECT 
  COUNT(*) as total,
  SUM(CASE WHEN paid = 1 THEN 1 ELSE 0 END) as paid,
  SUM(CASE WHEN submitted = 1 AND paid = 0 THEN 1 ELSE 0 END) as submitted,
  SUM(CASE WHEN submitted = 0 AND paid = 0 THEN 1 ELSE 0 END) as pending,
  printf('$%.2f', SUM(line_total)) as total_value
FROM invoices;"
```

## Email Configuration

### Gmail Setup

1. **Enable 2FA** on your Google account
2. **Create App Password:**
   - Visit: https://myaccount.google.com/apppasswords
   - Select: Mail > Other (Custom name)
   - Copy the 16-character password
3. **Update .env:**
   ```bash
   GMAIL_ADDRESS=your-email@gmail.com
   GMAIL_APP_PASSWORD=abcdefghijklmnop
   ```

### Email Flow

1. **Generate Invoice** → Creates PDF in `invoices/` directory
2. **Send Email** (TUI: press `s` or Cron: automated)
   - Attaches PDF
   - Sends to InfoObjects
   - CCs you for verification
   - Marks invoice as submitted in database

### Test Email
```bash
# Test email configuration
python3 -c "from email_service import send_invoice; \
    send_invoice('N001', 'test@example.com')"
```

## Cron Setup

### Install Cron Job
```bash
# Edit crontab
crontab -e

# Add this line (adjust path):
0 2 * * 6 cd /Users/yourname/devqai/inv_gen && /usr/bin/python3 cron_weekly_invoice.py >> logs/cron.log 2>&1
```

### Verify Cron
```bash
# List active cron jobs
crontab -l

# Check logs
tail -f logs/cron.log

# Test manually
python3 cron_weekly_invoice.py
```

### What the Cron Job Does
1. Finds next pending invoice
2. Generates PDF
3. Sends email with PDF attached
4. Marks invoice as submitted
5. Logs results

## Development

### Code Analysis
```bash
# Count lines by file type
./count_code.nu

# Analyze codebase
./analyze_codebase.nu
```

### File Structure
```
inv_gen/
├── invoice_cli.py           # Python CLI (368 lines)
├── email_service.py         # Email functionality (227 lines)
├── cron_weekly_invoice.py   # Automation (178 lines)
├── invoices.db              # SQLite database
├── requirements.txt         # Python dependencies
│
├── tui/                     # Go TUI (2,230 lines)
│   ├── main.go              # Main application
│   ├── email.go             # Native Go SMTP
│   ├── models/              # Data models
│   ├── views/               # UI views
│   └── styles/              # Color theme
│
├── invoices/                # Generated PDFs
├── logs/                    # Log files
│
├── analyze_codebase.nu      # Code analysis tool
├── count_code.nu            # Line counter
└── run.sh                   # Quick start script
```

### Building TUI
```bash
cd tui
make build    # Build binary
make install  # Install to bin/
make clean    # Clean build artifacts
```

## Troubleshooting

### Email Not Sending

**Check Gmail credentials:**
```bash
# Verify .env exists
cat .env | grep GMAIL

# Test SMTP connection
python3 -c "import smtplib; \
    smtp = smtplib.SMTP('smtp.gmail.com', 587); \
    smtp.starttls(); \
    print('SMTP connection successful')"
```

**Common issues:**
- App password not enabled (requires 2FA)
- Incorrect app password (no spaces, 16 chars)
- Gmail blocking "less secure apps" (use app password)

### PDF Not Generating

**Check WeasyPrint:**
```bash
# Test PDF generation
python3 invoice_cli.py generate N001

# Check output
ls -lh invoices/invoice_N001.pdf
```

**Common issues:**
- Missing WeasyPrint dependencies
- Template file not found
- Database connection failed

### Database Issues

**Check database:**
```bash
# Verify database exists
ls -lh invoices.db

# Check invoice count
sqlite3 invoices.db "SELECT COUNT(*) FROM invoices;"

# View schema
sqlite3 invoices.db ".schema invoices"
```

### TUI Not Starting

**Check binary:**
```bash
# Verify binary exists
ls -lh tui/bin/invoice-tui

# Check database path
cd tui && ./bin/invoice-tui ../invoices.db

# Build from source
cd tui && make clean && make build
```

## Configuration

### Environment Variables (.env)

```bash
# Required
GMAIL_ADDRESS=your-email@gmail.com
GMAIL_APP_PASSWORD=your_app_password

# Invoice Details
CONTRACTOR_NAME="Dion Edge"
CONTRACTOR_ADDRESS="123 Main St"
CONTRACTOR_CITY="Pinehurst"
CONTRACTOR_STATE="NC"
CONTRACTOR_ZIP="28374"
CONTRACTOR_PHONE="910-988-2000"

# Client Details
CLIENT_NAME="Info Objects, Inc."
CLIENT_ADDRESS="2033 Gateway Pl"
CLIENT_CITY="San Jose"
CLIENT_STATE="CA"
CLIENT_ZIP="95110"

# Email Recipients
BILLING_EMAIL=infoobjects@bill.com
TIMESHEET_EMAIL=timesheets@infoobjects.com
CC_EMAIL=dion@wrench.chat

# Invoice Settings
HOURLY_RATE=80.0
PAYMENT_TERMS="Net 15"
```

## Testing

### Manual Test Workflow
```bash
# 1. Generate PDF
python3 invoice_cli.py generate N001

# 2. Verify PDF created
ls -lh invoices/invoice_N001.pdf

# 3. Test email (won't actually send)
python3 invoice_cli.py submit N001

# 4. Check database updated
python3 invoice_cli.py show N001

# 5. Mark as paid
python3 invoice_cli.py paid N001
```

### TUI Test
```bash
cd tui
./bin/invoice-tui

# Test workflow:
# 1. Navigate to invoice list
# 2. Press 'g' to generate PDF
# 3. Press 's' to send email
# 4. Press 'p' to mark paid
# 5. Press 'q' to quit
```

### Cron Test
```bash
# Run cron script manually
python3 cron_weekly_invoice.py

# Check output
cat logs/cron.log
```

## Performance

**Current System:**
- Total Lines: ~2,800 (refactored from 11,361)
- Code Files: 10 active files
- Database: 27 invoices, ~86KB
- PDF Generation: ~2 seconds per invoice
- Email Delivery: ~3 seconds per email

## Invoice Details

**Standard Invoice:**
- Weekly timesheet (Monday-Sunday)
- 40 hours per week (8 hours/day)
- $80/hour rate
- Total: $3,200 per invoice
- Payment terms: Net 15

**Invoice Numbering:**
- Format: N### (N001, N002, etc.)
- Sequential by week
- Generated every Monday for previous week

## Security

**Best Practices:**
- `.env` is gitignored (never committed)
- App passwords instead of main password
- TLS encryption for SMTP
- Local SQLite database (not exposed)
- PDF files gitignored

## Backup

**Cron Backup (Daily 2 AM):**
```cron
0 2 * * * cd /Users/yourname/devqai && tar -czf backups/inv_gen_$(date +\%Y\%m\%d).tar.gz inv_gen/
```

**Manual Backup:**
```bash
# Backup database
cp invoices.db invoices.db.backup

# Backup entire project
tar -czf inv_gen_backup.tar.gz inv_gen/
```

## Statistics (Current)

- **Total Invoices:** 27
- **Date Range:** Oct 2025 - Apr 2026  
- **Total Value:** $86,400
- **Average Invoice:** $3,200
- **Hourly Rate:** $80
- **Hours/Week:** 40
- **Payment Terms:** Net 15

## Support

**Logs:**
```bash
# Cron logs
tail -f logs/cron.log

# Email logs (if configured)
tail -f logs/email.log
```

**Common Commands:**
```bash
# Check system status
python3 invoice_cli.py stats

# View next pending invoice
sqlite3 invoices.db "SELECT * FROM invoices WHERE submitted=0 ORDER BY invoice_create_date LIMIT 1;"

# Regenerate all PDFs
for inv in $(sqlite3 invoices.db "SELECT invoice_number FROM invoices;"); do
    python3 invoice_cli.py generate $inv
done
```

## License

MIT License - See LICENSE file for details

## Contact

**Developer:** Dion Edge  
**Email:** dion@devq.ai  
**Project:** DevQ.ai Invoice Management System  
**GitHub:** https://github.com/devq-ai/inv_gen

---

**Quick Reference:**
- Generate: `python3 invoice_cli.py generate N###`
- Send Email: TUI press `s` or Cron automated
- View Stats: `python3 invoice_cli.py stats`
- TUI: `cd tui && ./bin/invoice-tui`
