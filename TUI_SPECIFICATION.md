# Invoice TUI Specification - Charm.sh Implementation

## Overview

A beautiful, interactive Terminal User Interface (TUI) for managing invoices using Charm.sh tools.

## Technology Stack

### Charm.sh Libraries

1. **Bubble Tea** - Main TUI framework (like React for terminals)
   - Event-driven architecture
   - Component-based UI
   - State management

2. **Bubbles** - Pre-built UI components
   - Lists
   - Tables
   - Text inputs
   - Spinners
   - Progress bars

3. **Lip Gloss** - Styling and layout
   - Colors and themes
   - Borders and padding
   - Layout management
   - Responsive design

4. **Huh** - Interactive forms
   - Form validation
   - Multi-step wizards
   - Input fields
   - Confirmations

## Integration Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Go TUI Application                    │
│                  (Bubble Tea + Charm.sh)                │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              SQLite Database (invoices.db)              │
│                  Direct database access                 │
│                  using go-sqlite3 driver                │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│          Python PDF Generator (subprocess)              │
│         Call: db_invoice_generator.py                   │
└─────────────────────────────────────────────────────────┘
```

## TUI Features & Screens

### 1. Main Dashboard (Home Screen)

**Purpose**: Overview of invoice system status

**Layout**:
```
╔══════════════════════════════════════════════════════════╗
║               📋 INVOICE MANAGEMENT SYSTEM               ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║   Financial Overview                                     ║
║   ┌────────────────────────────────────────────────┐   ║
║   │  Total Invoices:        27                     │   ║
║   │  Total Value:           $86,400.00             │   ║
║   │                                                 │   ║
║   │  📤 Submitted:     3    $9,600.00             │   ║
║   │  ✅ Paid:          1    $3,200.00             │   ║
║   │  ⏳ Pending:       24   $76,800.00            │   ║
║   └────────────────────────────────────────────────┘   ║
║                                                          ║
║   Quick Actions                                          ║
║   ┌────────────────────────────────────────────────┐   ║
║   │  → View All Invoices                           │   ║
║   │  → Generate Invoice                            │   ║
║   │  → Approve Invoices                            │   ║
║   │  → View Reports                                │   ║
║   │  → Settings                                    │   ║
║   │  → Exit                                        │   ║
║   └────────────────────────────────────────────────┘   ║
║                                                          ║
║  [↑/↓] Navigate  [Enter] Select  [q] Quit              ║
╚══════════════════════════════════════════════════════════╝
```

**Components**:
- Lip Gloss: Styled panels, colors, borders
- Bubbles: List component for quick actions
- Live data from SQLite database

---

### 2. Invoice List View

**Purpose**: Browse and filter all invoices

**Layout**:
```
╔══════════════════════════════════════════════════════════╗
║                    📋 ALL INVOICES                       ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  Filter: [All ▼] Search: [_________]  [27 items]       ║
║                                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │ Invoice  Created     Due Date    Amount    Status  │ ║
║  ├────────────────────────────────────────────────────┤ ║
║  │ N001    10/05/2025  10/20/2025  $3,200  ✓ Submit │ ║
║  │ N002    10/12/2025  10/27/2025  $3,200  ✓ Submit │ ║
║  │ N003    10/19/2025  11/03/2025  $3,200  ✓✓ Paid  │ ║
║  │ N004    10/26/2025  11/10/2025  $3,200  ⏳ Pending│ ║
║  │ N005    11/02/2025  11/17/2025  $3,200  ⏳ Pending│ ║
║  │ N006    11/09/2025  11/24/2025  $3,200  ⏳ Pending│ ║
║  │ ...                                                │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  [↑/↓] Navigate  [Enter] View  [/] Search  [f] Filter  ║
║  [g] Generate PDF  [a] Approve  [Esc] Back             ║
╚══════════════════════════════════════════════════════════╝
```

**Components**:
- Bubbles: Table component with sorting
- Bubbles: Text input for search
- Bubbles: Dropdown for filter
- Color-coded status indicators (Lip Gloss)

**Features**:
- Real-time search filtering
- Sort by any column
- Filter by status (All, Pending, Submitted, Paid)
- Keyboard navigation
- Select invoice to view details

---

### 3. Invoice Detail View

**Purpose**: View complete invoice information

**Layout**:
```
╔══════════════════════════════════════════════════════════╗
║              📄 INVOICE DETAILS - N001                   ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  Invoice Information                                     ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Number:      N001                                 │ ║
║  │  Created:     10/05/2025                           │ ║
║  │  Due Date:    10/20/2025                           │ ║
║  │  Terms:       Net 15                               │ ║
║  │  Status:      ✓ Submitted                          │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  From: Dion Edge                                         ║
║        10705 Pinehurst Drive, Austin, TX 78747          ║
║                                                          ║
║  To:   InfoObjects, Inc.                                 ║
║        2041 Mission College Blvd, Ste 280               ║
║        Santa Clara, CA 95054                            ║
║        (408) 988-2000                                   ║
║                                                          ║
║  Work Details                                            ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │ Day       Date          In     Out    Hrs   Total  │ ║
║  │ Monday    09/29/2025  10:00  18:00   8.0   $640   │ ║
║  │ Tuesday   09/30/2025  10:00  18:00   8.0   $640   │ ║
║  │ Wednesday 10/01/2025  10:00  18:00   8.0   $640   │ ║
║  │ Thursday  10/02/2025  10:00  18:00   8.0   $640   │ ║
║  │ Friday    10/03/2025  10:00  18:00   8.0   $640   │ ║
║  │                                                    │ ║
║  │ TOTAL                              40.0  $3,200   │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  [g] Generate PDF  [a] Approve  [p] Mark Paid          ║
║  [Esc] Back                                             ║
╚══════════════════════════════════════════════════════════╝
```

**Components**:
- Lip Gloss: Styled panels and tables
- Custom rendering for invoice data
- Action buttons at bottom

---

### 4. Invoice Approval Workflow (KEY FEATURE)

**Purpose**: Interactive approval process for invoices

**Step 1: Select Invoices to Approve**

```
╔══════════════════════════════════════════════════════════╗
║              ✅ APPROVE INVOICES - Select                ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  Select invoices to approve and send:                    ║
║                                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │ [x] N001  10/05/2025  $3,200  Week: 09/29-10/03  │ ║
║  │ [x] N002  10/12/2025  $3,200  Week: 10/06-10/10  │ ║
║  │ [ ] N003  10/19/2025  $3,200  Week: 10/13-10/17  │ ║
║  │ [ ] N004  10/26/2025  $3,200  Week: 10/20-10/24  │ ║
║  │ [ ] N005  11/02/2025  $3,200  Week: 10/27-10/31  │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Selected: 2 invoices  Total: $6,400.00                ║
║                                                          ║
║  [Space] Toggle  [a] Select All  [n] Select None       ║
║  [Enter] Continue  [Esc] Cancel                         ║
╚══════════════════════════════════════════════════════════╝
```

**Step 2: Preview & Confirm**

```
╔══════════════════════════════════════════════════════════╗
║           ✅ APPROVE INVOICES - Confirm                  ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  You are about to approve:                               ║
║                                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Invoice N001                                      │ ║
║  │  • Created: 10/05/2025                            │ ║
║  │  • Amount: $3,200.00                              │ ║
║  │  • Week: Mon 09/29/2025 - Fri 10/03/2025         │ ║
║  │                                                    │ ║
║  │  Invoice N002                                      │ ║
║  │  • Created: 10/12/2025                            │ ║
║  │  • Amount: $3,200.00                              │ ║
║  │  • Week: Mon 10/06/2025 - Fri 10/10/2025         │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Total: 2 invoices  Amount: $6,400.00                  ║
║                                                          ║
║  Actions to perform:                                     ║
║  [✓] Generate PDF for each invoice                      ║
║  [✓] Mark as submitted in database                      ║
║  [ ] Email to client (optional)                         ║
║                                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  ⚠️  Are you sure you want to approve these?      │ ║
║  │                                                    │ ║
║  │     [Yes, Approve]    [No, Go Back]               │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
```

**Step 3: Processing**

```
╔══════════════════════════════════════════════════════════╗
║           ✅ APPROVE INVOICES - Processing               ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  Approving invoices...                                   ║
║                                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  ✓ Invoice N001 PDF generated                     │ ║
║  │  ✓ Invoice N001 marked as submitted               │ ║
║  │                                                    │ ║
║  │  ⏳ Invoice N002 generating PDF...                 │ ║
║  │    ████████░░░░░░░░░░░░░░░░░░░░░░░  30%           │ ║
║  │                                                    │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Progress: 1 of 2 complete                              ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
```

**Step 4: Success**

```
╔══════════════════════════════════════════════════════════╗
║            ✅ APPROVE INVOICES - Complete                ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  🎉 Successfully approved 2 invoices!                   ║
║                                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  ✓ Invoice N001 - $3,200.00                       │ ║
║  │    PDF: invoices/invoice_N001.pdf                 │ ║
║  │                                                    │ ║
║  │  ✓ Invoice N002 - $3,200.00                       │ ║
║  │    PDF: invoices/invoice_N002.pdf                 │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Total approved: $6,400.00                              ║
║                                                          ║
║  PDFs saved to: ./invoices/                             ║
║  Database updated successfully                           ║
║                                                          ║
║  [Enter] Return to Dashboard                            ║
╚══════════════════════════════════════════════════════════╝
```

**Components**:
- Bubbles: Multi-select list
- Bubbles: Progress bar
- Huh: Confirmation dialog
- Bubbles: Spinner for processing
- Lip Gloss: Success/error styling

---

### 5. Create/Edit Invoice (Huh Forms)

**Purpose**: Create new invoice or edit existing one

**Form Steps** (using Huh):

**Step 1: Basic Information**
```
┌──────────────────────────────────────────────────────────┐
│ Create New Invoice - Basic Information                   │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  Invoice Number: N028___                                 │
│                                                           │
│  Invoice Date: [10/05/2025]                              │
│  Select: [ Today ] [ Custom Date ]                       │
│                                                           │
│  Payment Terms: [15▼] days                               │
│                                                           │
│  Due Date: 10/20/2025 (calculated)                       │
│                                                           │
│                        [Next →]                           │
└──────────────────────────────────────────────────────────┘
```

**Step 2: Work Week Details**
```
┌──────────────────────────────────────────────────────────┐
│ Create New Invoice - Work Details                        │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  Week Starting: [10/13/2025] (Monday)                    │
│                                                           │
│  Default Hours/Day: [8.0]                                │
│  Default Rate/Hour: [$80.00]                             │
│                                                           │
│  Apply to all days? [Yes] [No, customize]                │
│                                                           │
│  If customizing:                                          │
│  Monday:    [8.0] hrs × [$80] = $640                     │
│  Tuesday:   [8.0] hrs × [$80] = $640                     │
│  Wednesday: [8.0] hrs × [$80] = $640                     │
│  Thursday:  [8.0] hrs × [$80] = $640                     │
│  Friday:    [8.0] hrs × [$80] = $640                     │
│                                                           │
│  Total: 40.0 hrs  $3,200.00                              │
│                                                           │
│                 [← Back]  [Next →]                        │
└──────────────────────────────────────────────────────────┘
```

**Step 3: Review & Create**
```
┌──────────────────────────────────────────────────────────┐
│ Create New Invoice - Review                              │
├──────────────────────────────────────────────────────────┤
│                                                           │
│  Invoice Number: N028                                    │
│  Date: 10/05/2025                                        │
│  Due: 10/20/2025 (Net 15)                               │
│                                                           │
│  Work Week: 10/13/2025 - 10/17/2025                     │
│  Total Hours: 40.0                                       │
│  Total Amount: $3,200.00                                 │
│                                                           │
│  ✓ All required fields completed                         │
│                                                           │
│  Create this invoice?                                    │
│                                                           │
│          [← Back]  [Create Invoice]                      │
└──────────────────────────────────────────────────────────┘
```

**Components**:
- Huh: Multi-step form
- Huh: Input validation
- Huh: Date picker
- Huh: Confirmation
- Real-time calculation display

---

### 6. Reports & Statistics

**Purpose**: Financial reporting and analytics

```
╔══════════════════════════════════════════════════════════╗
║                📊 FINANCIAL REPORTS                      ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  Overall Statistics                                      ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Total Invoices:       27                          │ ║
║  │  Total Value:          $86,400.00                  │ ║
║  │                                                    │ ║
║  │  Submitted:      3     $9,600.00    (11.1%)      │ ║
║  │  Paid:           1     $3,200.00    (3.7%)       │ ║
║  │  Pending:        24    $76,800.00   (88.9%)      │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Monthly Breakdown                                       ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  October 2025:    4 invoices   $12,800            │ ║
║  │  November 2025:   5 invoices   $16,000            │ ║
║  │  December 2025:   4 invoices   $12,800            │ ║
║  │  January 2026:    5 invoices   $16,000            │ ║
║  │  February 2026:   4 invoices   $12,800            │ ║
║  │  March 2026:      4 invoices   $12,800            │ ║
║  │  April 2026:      1 invoice    $3,200             │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Revenue Chart                                           ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Oct ████████ 14.8%                                │ ║
║  │  Nov ██████████ 18.5%                              │ ║
║  │  Dec ████████ 14.8%                                │ ║
║  │  Jan ██████████ 18.5%                              │ ║
║  │  Feb ████████ 14.8%                                │ ║
║  │  Mar ████████ 14.8%                                │ ║
║  │  Apr ██ 3.7%                                       │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  [Esc] Back                                             ║
╚══════════════════════════════════════════════════════════╝
```

**Components**:
- Lip Gloss: Styled panels
- Custom bar charts
- Percentage calculations
- Color-coded status

---

### 7. Settings

**Purpose**: Configure application settings

```
╔══════════════════════════════════════════════════════════╗
║                    ⚙️  SETTINGS                          ║
╠══════════════════════════════════════════════════════════╣
║                                                          ║
║  Default Values                                          ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Hourly Rate:         [$80.00]                     │ ║
║  │  Payment Terms:       [15] days                    │ ║
║  │  Hours per Day:       [8.0]                        │ ║
║  │  Work Start Time:     [10:00]                      │ ║
║  │  Work End Time:       [18:00]                      │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  PDF Output                                              ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Output Directory:    [./invoices/]                │ ║
║  │  Auto-open PDFs:      [Yes] [No]                   │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  Email Settings (Optional)                               ║
║  ┌────────────────────────────────────────────────────┐ ║
║  │  Enable Email:        [Yes] [No]                   │ ║
║  │  SMTP Server:         [smtp.gmail.com]             │ ║
║  │  From Email:          [dion@devq.ai]               │ ║
║  │  To Email:            [billing@client.com]         │ ║
║  └────────────────────────────────────────────────────┘ ║
║                                                          ║
║  [Save]  [Cancel]  [Reset to Defaults]                 ║
╚══════════════════════════════════════════════════════════╝
```

**Components**:
- Huh: Form inputs
- Input validation
- Save/cancel actions

---

## Technical Implementation

### File Structure

```
inv_gen/
├── tui/
│   ├── main.go                    # Entry point
│   ├── models/
│   │   ├── invoice.go            # Invoice data model
│   │   ├── database.go           # SQLite operations
│   │   └── pdf.go                # PDF generation wrapper
│   ├── views/
│   │   ├── dashboard.go          # Main dashboard
│   │   ├── invoice_list.go       # Invoice list view
│   │   ├── invoice_detail.go     # Invoice detail view
│   │   ├── approval.go           # Approval workflow
│   │   ├── create_invoice.go     # Create/edit form
│   │   ├── reports.go            # Reports view
│   │   └── settings.go           # Settings view
│   ├── components/
│   │   ├── table.go              # Custom table component
│   │   ├── panel.go              # Styled panels
│   │   ├── statusbar.go          # Status bar
│   │   └── charts.go             # Simple charts
│   ├── styles/
│   │   └── theme.go              # Lip Gloss styles
│   └── utils/
│       ├── formatting.go         # Number/date formatting
│       └── validation.go         # Input validation
├── go.mod
├── go.sum
└── README_TUI.md
```

### Dependencies (go.mod)

```go
module github.com/devqai/invoice-tui

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/charmbracelet/huh v0.3.0
    github.com/mattn/go-sqlite3 v1.14.18
)
```

### Main Application Flow

```go
type sessionState int

const (
    dashboardView sessionState = iota
    invoiceListView
    invoiceDetailView
    approvalView
    createInvoiceView
    reportsView
    settingsView
)

type model struct {
    state          sessionState
    db             *Database
    invoices       []Invoice
    selectedInvoice *Invoice
    // ... view-specific state
}

func (m model) Init() tea.Cmd {
    // Load initial data from SQLite
    return loadInvoicesCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle keyboard events
    // Route to appropriate view
    // Update state
}

func (m model) View() string {
    // Render current view
    switch m.state {
    case dashboardView:
        return renderDashboard(m)
    case invoiceListView:
        return renderInvoiceList(m)
    // ...
    }
}
```

### Database Integration (Go)

```go
type Database struct {
    db *sql.DB
}

func OpenDatabase(path string) (*Database, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    return &Database{db: db}, nil
}

func (d *Database) GetAllInvoices() ([]Invoice, error) {
    rows, err := d.db.Query("SELECT * FROM invoices ORDER BY pk")
    // Parse and return invoices
}

func (d *Database) MarkAsSubmitted(invoiceNumber string) error {
    _, err := d.db.Exec(
        "UPDATE invoices SET submitted = 1 WHERE invoice_number = ?",
        invoiceNumber,
    )
    return err
}

// More database operations...
```

### PDF Generation Wrapper

```go
func GeneratePDF(invoiceNumber string) error {
    cmd := exec.Command(
        "python3",
        "db_invoice_generator.py",
        "--invoice", invoiceNumber,
    )
    return cmd.Run()
}

func GenerateBatchPDFs(invoiceNumbers []string) error {
    // Generate multiple PDFs
    // Show progress in TUI
}
```

## Key Bindings

### Global
- `q` / `Ctrl+C` - Quit
- `Esc` - Go back / Cancel
- `?` - Help
- `Tab` - Switch focus

### Dashboard
- `↑/↓` - Navigate menu
- `Enter` - Select action
- `r` - Refresh data

### Invoice List
- `↑/↓` - Navigate invoices
- `Enter` - View details
- `/` - Search
- `f` - Filter
- `g` - Generate PDF
- `a` - Approve selected
- `Space` - Multi-select

### Invoice Detail
- `g` - Generate PDF
- `a` - Approve
- `p` - Mark as paid
- `e` - Edit

### Approval Workflow
- `Space` - Toggle selection
- `a` - Select all
- `n` - Deselect all
- `Enter` - Continue/Confirm
- `Esc` - Cancel

## Color Scheme (Lip Gloss)

### Status Colors
```go
var (
    PendingColor   = lipgloss.Color("#FFE599")  // Yellow
    SubmittedColor = lipgloss.Color("#A4C2F4")  // Blue
    PaidColor      = lipgloss.Color("#A1D9A0")  // Green
    ErrorColor     = lipgloss.Color("#FF6961")  // Red
    SuccessColor   = lipgloss.Color("#77DD77")  // Green
)
```

### UI Elements
```go
var (
    BorderColor    = lipgloss.Color("#9D00FF")  // Neon Purple
    AccentColor    = lipgloss.Color("#FF10F0")  // Neon Pink
    TextColor      = lipgloss.Color("#E3E3E3")  // Soft White
    DimTextColor   = lipgloss.Color("#7D8B99")  // Cool Grey
    BackgroundColor = lipgloss.Color("#0A0A0A") // Pure Black
)
```

## Performance Considerations

1. **Lazy Loading**: Load invoice details on demand
2. **Caching**: Cache frequently accessed data
3. **Efficient Rendering**: Only re-render changed components
4. **Background Operations**: Run PDF generation in goroutines
5. **Database Optimization**: Use prepared statements

## Error Handling

1. **Database Errors**: Show user-friendly messages
2. **PDF Generation Errors**: Display detailed error info
3. **Validation Errors**: Inline form validation
4. **Network Errors**: Retry logic for email sending

## Future Enhancements

1. **Themes**: Light/dark mode toggle
2. **Keyboard Shortcuts**: Customizable key bindings
3. **Export**: CSV/Excel export functionality
4. **Email Integration**: Direct email sending from TUI
5. **Multi-user**: User authentication
6. **Backup**: Database backup/restore from TUI
7. **Templates