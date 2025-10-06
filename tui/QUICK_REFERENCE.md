# Invoice TUI - Quick Reference Card

## 🚀 Quick Start

```bash
cd inv_gen/tui
make build
./bin/invoice-tui
```

## ⌨️ Keyboard Shortcuts

### Global Navigation
```
q        Quit (dashboard) / Back to dashboard (other views)
Esc      Return to dashboard
Ctrl+C   Force quit
↑/↓      Navigate up/down
Enter    Select/Confirm
```

### Dashboard
```
↑/↓      Navigate menu
Enter    Select action
q        Quit application
```

### Invoice List
```
↑/↓      Navigate invoices
Enter    View details
f        Cycle filters (All→Pending→Submitted→Paid)
g        Generate PDF
s        Mark as submitted
p        Mark as paid
Esc      Back to dashboard
```

### Invoice Detail
```
g        Generate PDF
s        Mark as submitted
p        Mark as paid
Esc      Back to invoice list
```

### Reports & Settings
```
Esc      Back to dashboard
q        Back to dashboard
```

## 📋 Common Tasks

### View All Invoices
```
Dashboard → "View All Invoices" → Enter
```

### Filter Pending Invoices
```
Invoice List → Press 'f' until "Pending" selected
```

### Generate PDF
```
1. Navigate to invoice (List or Detail)
2. Press 'g'
3. Success message appears
4. PDF saved to ../invoices/
```

### Mark Invoice as Submitted
```
1. Navigate to invoice (List or Detail)
2. Press 's'
3. Confirmation message appears
```

### Mark Invoice as Paid
```
1. Navigate to invoice (List or Detail)
2. Press 'p'
3. Confirmation message appears
```

### View Financial Reports
```
Dashboard → "View Reports" → Enter
```

## 🎨 Status Colors

```
⏳ Pending    - Blue (#A4C2F4)
📤 Submitted  - Pink (#F4A6C0)
✅ Paid       - Purple (#B5A0E3)
```

## 🔧 Build Commands

```bash
make help           # Show all commands
make build          # Build binary
make run            # Build and run
make run-dev        # Run without building
make clean          # Remove build files
make install        # Install system-wide
make test           # Run tests
make fmt            # Format code
```

## 📁 File Locations

```
Database:    ../invoices.db
PDFs:        ../invoices/invoice_*.pdf
Python CLI:  ../invoice_cli.py
Binary:      bin/invoice-tui
```

## 🐛 Troubleshooting

### Database not found
```bash
cd .. && python3 create_invoice_db.py
```

### Build fails
```bash
make clean
go mod download
make build
```

### PDF generation fails
```bash
# Check Python CLI
python3 ../invoice_cli.py list

# Make executable
chmod +x ../invoice_cli.py
```

## 📊 Dashboard Overview

```
Financial Overview:
  Total Invoices:    27
  Total Value:       $86,400.00
  
  📤 Submitted:      3    $9,600.00
  ✅ Paid:           1    $3,200.00
  ⏳ Pending:        24   $76,800.00

Quick Actions:
  → 📋 View All Invoices
  → 📄 Generate Invoice
  → ✅ Approve Invoices
  → 📊 View Reports
  → ⚙️  Settings
  → 🚪 Exit
```

## 🔄 Workflow Examples

### Weekly Invoice Submission
```
1. Dashboard → "View All Invoices"
2. Press 'f' → Filter "Pending"
3. Select this week's invoice
4. Press 'g' → Generate PDF
5. Review PDF in ../invoices/
6. Press 's' → Mark as submitted
7. Email PDF to client
8. When paid: Press 'p' → Mark as paid
```

### Monthly Review
```
1. Dashboard → "View Reports"
2. Review statistics
3. Check submitted percentages
4. Esc → Back to dashboard
5. "View All Invoices" → Filter "Paid"
6. Verify all paid invoices
```

### Batch Processing
```
1. Dashboard → "Approve Invoices"
2. Lists all pending invoices
3. Navigate with ↑/↓
4. Press 'g' for each to generate PDF
5. Press 's' for each to mark submitted
```

## 💡 Tips & Tricks

- **Fast Navigation**: Use 'f' to quickly cycle filters
- **Escape Key**: Always returns you to dashboard
- **Status Updates**: 's' and 'p' work in both List and Detail views
- **PDF Location**: All PDFs saved to ../invoices/ directory
- **Database Safety**: All operations are atomic and safe

## 📦 System Requirements

```
Go:        1.21+
CGo:       Enabled
Python:    3.12+
SQLite:    3.x
Database:  invoices.db (required)
Terminal:  80x24 minimum
```

## 🆘 Help

```bash
# In-app help
Press '?' (if implemented) or see help text at bottom of each screen

# Command help
make help

# Documentation
cat README.md
cat SETUP.md
```

## 📞 Contact

**Developer**: Dion Edge  
**Email**: dion@devq.ai  
**Company**: DevQ.ai

---

**Version**: 1.0.0  
**Last Updated**: January 2025  
**License**: DevQ.ai Internal Tool

---

**🎨 Happy invoice managing! ✨**