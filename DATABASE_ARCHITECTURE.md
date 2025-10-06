# Database-Driven Invoice Architecture

## Overview

This invoice system is built on a **database-first architecture** where the SQLite database serves as the single source of truth for all invoice operations. Unlike traditional systems that generate invoices from templates and store PDFs, this system stores structured data and generates PDFs on-demand.

## Core Architecture Principle

```
┌─────────────────────────────────────────────────────────────┐
│                    SQLite Database                          │
│                  (Single Source of Truth)                   │
│                                                             │
│  ┌───────────────────────────────────────────────────┐    │
│  │  Invoice Records (27 weeks pre-populated)        │    │
│  │  - All invoice data                               │    │
│  │  - Work dates and times                           │    │
│  │  - Payment information                            │    │
│  │  - Status tracking                                │    │
│  └───────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┼───────────┐
                │           │           │
                ▼           ▼           ▼
         ┌──────────┐ ┌──────────┐ ┌──────────┐
         │   PDF    │ │  Status  │ │ Reports  │
         │Generator │ │ Tracking │ │& Stats   │
         └──────────┘ └──────────┘ └──────────┘
```

## Why Database-Driven?

### 1. **Data Consistency**
- All invoice data exists in one canonical location
- No duplicate information across files
- Changes update once, affect everywhere

### 2. **Flexibility**
- Generate PDFs any time without data entry
- Regenerate invoices if format changes
- Query data in unlimited ways

### 3. **Scalability**
- Handle thousands of invoices efficiently
- Fast queries with proper indexing
- Minimal storage requirements

### 4. **Integration**
- Easy to connect to other systems
- Standard SQL interface
- Export to any format needed

### 5. **Reporting**
- Real-time statistics
- Complex queries for analysis
- Historical data tracking

## Database Schema Design

### Core Fields

```sql
CREATE TABLE invoices (
    -- Identity
    pk INTEGER PRIMARY KEY AUTOINCREMENT,
    invoice_number TEXT NOT NULL UNIQUE,
    
    -- Dates
    invoice_create_date TEXT NOT NULL,
    due_date TEXT NOT NULL,
    payment_terms INTEGER NOT NULL DEFAULT 15,
    
    -- Parties
    payee TEXT NOT NULL,
    payee_address TEXT NOT NULL,
    payor TEXT NOT NULL,
    payor_address TEXT NOT NULL,
    payor_phone TEXT NOT NULL,
    
    -- Work Details (5 days × 6 fields = 30 fields)
    monday_date TEXT,
    monday_in TEXT,
    monday_out TEXT,
    monday_hours_worked REAL,
    monday_unit_price REAL,
    monday_line_total REAL,
    -- ... repeated for tuesday, wednesday, thursday, friday
    
    -- Totals
    total_hours REAL NOT NULL,
    line_total REAL NOT NULL,
    
    -- Status Tracking
    submitted INTEGER NOT NULL DEFAULT 0,
    paid INTEGER NOT NULL DEFAULT 0,
    
    -- Metadata
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
)
```

### Design Decisions

#### **Denormalized Structure**
- Each invoice is a complete record
- Monday through Friday fields stored separately
- Trade-off: Some redundancy for query simplicity

**Why?**
- Simpler queries (no joins needed)
- Faster PDF generation
- Easier to understand and maintain
- Typical use case: read entire invoice at once

#### **Text-Based Dates**
- Dates stored as TEXT (e.g., "10/05/2025")
- Not using SQLite DATE type

**Why?**
- Matches PDF output format exactly
- No conversion needed for display
- Human-readable in database
- Sufficient for this use case

#### **Boolean as INTEGER**
- `submitted` and `paid` stored as 0/1
- SQLite doesn't have native BOOLEAN type

**Why?**
- Standard SQLite practice
- Easy to query: `WHERE paid = 1`
- Compatible with all tools

## Data Flow

### 1. Database Population

```
create_invoice_db.py
        │
        ▼
┌────────────────────┐
│ Calculate Dates    │  Start: Oct 5, 2025
│ - Sundays (create) │  End: Apr 3, 2026
│ - Mon-Fri (work)   │  Total: 27 invoices
└────────────────────┘
        │
        ▼
┌────────────────────┐
│ Generate Data      │  For each week:
│ - Invoice number   │  - N001, N002, etc.
│ - Work days        │  - Mon-Fri dates
│ - Calculations     │  - Hours, rates, totals
└────────────────────┘
        │
        ▼
┌────────────────────┐
│ Insert Records     │  INSERT INTO invoices
│ - 27 week records  │  VALUES (...)
│ - All fields       │
└────────────────────┘
        │
        ▼
    invoices.db
```

### 2. PDF Generation

```
invoice_cli.py generate N001
        │
        ▼
┌──────────────────────────┐
│ Query Database           │  SELECT * FROM invoices
│ - Get invoice record     │  WHERE invoice_number = 'N001'
└──────────────────────────┘
        │
        ▼
┌──────────────────────────┐
│ DatabaseInvoiceGenerator │
│ - Parse record data      │
│ - Format for PDF         │
└──────────────────────────┘
        │
        ▼
┌──────────────────────────┐
│ ReportLab PDF Canvas     │
│ - Draw invoice header    │
│ - Add table rows         │
│ - Calculate layout       │
└──────────────────────────┘
        │
        ▼
  invoice_N001.pdf
```

### 3. Status Updates

```
invoice_cli.py submit N001
        │
        ▼
┌──────────────────────────┐
│ Update Database          │  UPDATE invoices
│ - Set submitted = 1      │  SET submitted = 1
│                          │  WHERE invoice_number = 'N001'
└──────────────────────────┘
        │
        ▼
    invoices.db
    (status updated)
```

### 4. Reporting

```
invoice_cli.py stats
        │
        ▼
┌──────────────────────────┐
│ Aggregate Queries        │  SELECT COUNT(*), SUM(line_total)
│ - Count invoices         │  FROM invoices
│ - Sum amounts            │  GROUP BY status
│ - Group by status        │
└──────────────────────────┘
        │
        ▼
┌──────────────────────────┐
│ Format Results           │
│ - Display statistics     │
│ - Calculate percentages  │
└──────────────────────────┘
```

## Key Operations

### Query Examples

#### 1. Get Single Invoice
```sql
SELECT * FROM invoices 
WHERE invoice_number = 'N001';
```

#### 2. List Pending Invoices
```sql
SELECT invoice_number, invoice_create_date, line_total
FROM invoices 
WHERE submitted = 0
ORDER BY invoice_create_date;
```

#### 3. Calculate Total Revenue
```sql
SELECT 
    COUNT(*) as total_invoices,
    SUM(line_total) as total_amount,
    SUM(CASE WHEN paid = 1 THEN line_total ELSE 0 END) as paid_amount,
    SUM(CASE WHEN paid = 0 THEN line_total ELSE 0 END) as unpaid_amount
FROM invoices;
```

#### 4. Find Overdue Invoices
```sql
SELECT invoice_number, due_date, line_total
FROM invoices
WHERE paid = 0
  AND date(substr(due_date, 7, 4) || '-' || 
           substr(due_date, 1, 2) || '-' || 
           substr(due_date, 4, 2)) < date('now')
ORDER BY due_date;
```

#### 5. Monthly Revenue Report
```sql
SELECT 
    substr(invoice_create_date, 1, 2) as month,
    COUNT(*) as invoice_count,
    SUM(total_hours) as total_hours,
    SUM(line_total) as revenue
FROM invoices
GROUP BY month
ORDER BY invoice_create_date;
```

## Advantages Over File-Based Systems

| Aspect | Database-Driven | File-Based |
|--------|----------------|------------|
| **Data Entry** | Once, centralized | Multiple times |
| **Consistency** | Guaranteed | Manual sync needed |
| **Querying** | SQL queries | Parse files |
| **Updates** | Update one record | Edit multiple files |
| **Reporting** | Real-time | Must aggregate files |
| **Backup** | Single DB file | Many files to backup |
| **Scalability** | Handles thousands | File system limits |
| **Integration** | Standard SQL | Custom parsing |

## Extension Points

### Future Capabilities Enabled by Database

#### 1. **Multi-Client Support**
```sql
ALTER TABLE invoices ADD COLUMN client_id INTEGER;
CREATE TABLE clients (
    id INTEGER PRIMARY KEY,
    name TEXT,
    address TEXT,
    payment_terms INTEGER
);
```

#### 2. **Payment Tracking**
```sql
CREATE TABLE payments (
    id INTEGER PRIMARY KEY,
    invoice_id INTEGER,
    payment_date TEXT,
    amount REAL,
    payment_method TEXT,
    FOREIGN KEY (invoice_id) REFERENCES invoices(pk)
);
```

#### 3. **Time Entry Tracking**
```sql
CREATE TABLE time_entries (
    id INTEGER PRIMARY KEY,
    invoice_id INTEGER,
    work_date TEXT,
    time_in TEXT,
    time_out TEXT,
    hours REAL,
    description TEXT,
    FOREIGN KEY (invoice_id) REFERENCES invoices(pk)
);
```

#### 4. **Email History**
```sql
CREATE TABLE email_log (
    id INTEGER PRIMARY KEY,
    invoice_id INTEGER,
    sent_date TEXT,
    recipient TEXT,
    subject TEXT,
    status TEXT,
    FOREIGN KEY (invoice_id) REFERENCES invoices(pk)
);
```

## Performance Considerations

### Current Scale
- **27 invoices**: Instantaneous queries
- **Database size**: ~100KB
- **PDF generation**: <1 second per invoice

### Optimization for Growth

#### Indexes for Common Queries
```sql
CREATE INDEX idx_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_status ON invoices(submitted, paid);
CREATE INDEX idx_dates ON invoices(invoice_create_date, due_date);
```

#### Partitioning Strategy (Future)
```sql
-- For thousands of invoices, partition by year
CREATE TABLE invoices_2025 AS SELECT * FROM invoices WHERE ...;
CREATE TABLE invoices_2026 AS SELECT * FROM invoices WHERE ...;
```

## Database Maintenance

### Backup Strategy
```bash
# Daily backup
sqlite3 invoices.db ".backup invoices_$(date +%Y%m%d).db"

# Export to SQL
sqlite3 invoices.db .dump > invoices_backup.sql

# Restore from backup
sqlite3 invoices_new.db < invoices_backup.sql
```

### Data Integrity Checks
```sql
-- Verify totals match
SELECT invoice_number 
FROM invoices
WHERE line_total != (
    monday_line_total + tuesday_line_total + 
    wednesday_line_total + thursday_line_total + friday_line_total
);

-- Check for missing data
SELECT invoice_number
FROM invoices
WHERE monday_date IS NULL OR friday_date IS NULL;
```

## Best Practices

### 1. **Always Query First**
- Never generate data outside the database
- Query database for invoice data
- Generate PDFs from query results

### 2. **Update Status in Database**
- Mark submitted when email sent
- Mark paid when payment received
- Never track status in files or spreadsheets

### 3. **Use Database for Reports**
- Run SQL queries for analytics
- Don't parse PDFs for data
- Database is the source of truth

### 4. **Backup Regularly**
- Automated daily backups
- Test restore procedure
- Keep multiple versions

### 5. **Validate Data**
- Check totals after updates
- Verify date calculations
- Run integrity checks

## Conclusion

The database-driven architecture provides:

✅ **Single Source of Truth**: All data in one place  
✅ **Data Consistency**: No duplication or sync issues  
✅ **Flexibility**: Generate PDFs anytime, any format  
✅ **Scalability**: Handle thousands of invoices  
✅ **Reporting**: Real-time business intelligence  
✅ **Integration**: Easy to connect other systems  
✅ **Maintenance**: Simple backup and recovery  
✅ **Future-Proof**: Easy to extend and enhance  

This architecture ensures that the **database drives the process**, not files or manual data entry. The result is a robust, maintainable, and scalable invoice management system.