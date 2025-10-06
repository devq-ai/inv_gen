package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Invoice represents a complete invoice record from the database
type Invoice struct {
	PK                 int       `db:"pk"`
	InvoiceNumber      string    `db:"invoice_number"`
	InvoiceCreateDate  string    `db:"invoice_create_date"`
	PaymentTerms       int       `db:"payment_terms"`
	DueDate            string    `db:"due_date"`
	Payee              string    `db:"payee"`
	PayeeAddress       string    `db:"payee_address"`
	Payor              string    `db:"payor"`
	PayorAddress       string    `db:"payor_address"`
	PayorPhone         string    `db:"payor_phone"`

	// Monday
	MondayDate         string    `db:"monday_date"`
	MondayIn           string    `db:"monday_in"`
	MondayOut          string    `db:"monday_out"`
	MondayHours        float64   `db:"monday_hours_worked"`
	MondayRate         float64   `db:"monday_unit_price"`
	MondayTotal        float64   `db:"monday_line_total"`

	// Tuesday
	TuesdayDate        string    `db:"tuesday_date"`
	TuesdayIn          string    `db:"tuesday_in"`
	TuesdayOut         string    `db:"tuesday_out"`
	TuesdayHours       float64   `db:"tuesday_hours_worked"`
	TuesdayRate        float64   `db:"tuesday_unit_price"`
	TuesdayTotal       float64   `db:"tuesday_line_total"`

	// Wednesday
	WednesdayDate      string    `db:"wednesday_date"`
	WednesdayIn        string    `db:"wednesday_in"`
	WednesdayOut       string    `db:"wednesday_out"`
	WednesdayHours     float64   `db:"wednesday_hours_worked"`
	WednesdayRate      float64   `db:"wednesday_unit_price"`
	WednesdayTotal     float64   `db:"wednesday_line_total"`

	// Thursday
	ThursdayDate       string    `db:"thursday_date"`
	ThursdayIn         string    `db:"thursday_in"`
	ThursdayOut        string    `db:"thursday_out"`
	ThursdayHours      float64   `db:"thursday_hours_worked"`
	ThursdayRate       float64   `db:"thursday_unit_price"`
	ThursdayTotal      float64   `db:"thursday_line_total"`

	// Friday
	FridayDate         string    `db:"friday_date"`
	FridayIn           string    `db:"friday_in"`
	FridayOut          string    `db:"friday_out"`
	FridayHours        float64   `db:"friday_hours_worked"`
	FridayRate         float64   `db:"friday_unit_price"`
	FridayTotal        float64   `db:"friday_line_total"`

	// Totals
	TotalHours         float64   `db:"total_hours"`
	LineTotal          float64   `db:"line_total"`

	// Status
	Submitted          bool      `db:"submitted"`
	Paid               bool      `db:"paid"`
	CreatedAt          time.Time `db:"created_at"`
}

// InvoiceSummary contains aggregate statistics for all invoices
type InvoiceSummary struct {
	TotalCount       int
	SubmittedCount   int
	PaidCount        int
	PendingCount     int
	UnpaidCount      int
	TotalAmount      float64
	SubmittedAmount  float64
	PaidAmount       float64
	PendingAmount    float64
	UnpaidAmount     float64
}

// StatusString returns a human-readable status string
func (i *Invoice) StatusString() string {
	if i.Paid {
		return "Paid"
	}
	if i.Submitted {
		return "Submitted"
	}
	return "Pending"
}

// StatusColor returns the status for color coding
func (i *Invoice) StatusColor() string {
	if i.Paid {
		return "paid"
	}
	if i.Submitted {
		return "submitted"
	}
	return "pending"
}

// FormatAmount formats a float as currency
func FormatAmount(amount float64) string {
	return fmt.Sprintf("$%8.2f", amount)
}

// FormatAmountCompact formats a float as currency without padding
func FormatAmountCompact(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// GetWeekEnding returns the last working day (Friday) of the invoice week
func (i *Invoice) GetWeekEnding() string {
	return i.FridayDate
}

// Database represents a connection to the SQLite database
type Database struct {
	db *sql.DB
}

// OpenDatabase opens a connection to the invoice database
func OpenDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetAllInvoices retrieves all invoices from the database
func (d *Database) GetAllInvoices() ([]Invoice, error) {
	query := `
		SELECT
			pk, invoice_number, invoice_create_date, payment_terms, due_date,
			payee, payee_address, payor, payor_address, payor_phone,
			monday_date, monday_in, monday_out, monday_hours_worked, monday_unit_price, monday_line_total,
			tuesday_date, tuesday_in, tuesday_out, tuesday_hours_worked, tuesday_unit_price, tuesday_line_total,
			wednesday_date, wednesday_in, wednesday_out, wednesday_hours_worked, wednesday_unit_price, wednesday_line_total,
			thursday_date, thursday_in, thursday_out, thursday_hours_worked, thursday_unit_price, thursday_line_total,
			friday_date, friday_in, friday_out, friday_hours_worked, friday_unit_price, friday_line_total,
			total_hours, line_total, submitted, paid, created_at
		FROM invoices
		ORDER BY pk ASC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query invoices: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		var submitted, paid int
		var createdAt string

		err := rows.Scan(
			&inv.PK, &inv.InvoiceNumber, &inv.InvoiceCreateDate, &inv.PaymentTerms, &inv.DueDate,
			&inv.Payee, &inv.PayeeAddress, &inv.Payor, &inv.PayorAddress, &inv.PayorPhone,
			&inv.MondayDate, &inv.MondayIn, &inv.MondayOut, &inv.MondayHours, &inv.MondayRate, &inv.MondayTotal,
			&inv.TuesdayDate, &inv.TuesdayIn, &inv.TuesdayOut, &inv.TuesdayHours, &inv.TuesdayRate, &inv.TuesdayTotal,
			&inv.WednesdayDate, &inv.WednesdayIn, &inv.WednesdayOut, &inv.WednesdayHours, &inv.WednesdayRate, &inv.WednesdayTotal,
			&inv.ThursdayDate, &inv.ThursdayIn, &inv.ThursdayOut, &inv.ThursdayHours, &inv.ThursdayRate, &inv.ThursdayTotal,
			&inv.FridayDate, &inv.FridayIn, &inv.FridayOut, &inv.FridayHours, &inv.FridayRate, &inv.FridayTotal,
			&inv.TotalHours, &inv.LineTotal, &submitted, &paid, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}

		inv.Submitted = submitted == 1
		inv.Paid = paid == 1
		inv.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		invoices = append(invoices, inv)
	}

	return invoices, nil
}

// GetInvoiceByNumber retrieves a specific invoice by its number
func (d *Database) GetInvoiceByNumber(invoiceNumber string) (*Invoice, error) {
	query := `
		SELECT
			pk, invoice_number, invoice_create_date, payment_terms, due_date,
			payee, payee_address, payor, payor_address, payor_phone,
			monday_date, monday_in, monday_out, monday_hours_worked, monday_unit_price, monday_line_total,
			tuesday_date, tuesday_in, tuesday_out, tuesday_hours_worked, tuesday_unit_price, tuesday_line_total,
			wednesday_date, wednesday_in, wednesday_out, wednesday_hours_worked, wednesday_unit_price, wednesday_line_total,
			thursday_date, thursday_in, thursday_out, thursday_hours_worked, thursday_unit_price, thursday_line_total,
			friday_date, friday_in, friday_out, friday_hours_worked, friday_unit_price, friday_line_total,
			total_hours, line_total, submitted, paid, created_at
		FROM invoices
		WHERE invoice_number = ?
	`

	var inv Invoice
	var submitted, paid int
	var createdAt string

	err := d.db.QueryRow(query, invoiceNumber).Scan(
		&inv.PK, &inv.InvoiceNumber, &inv.InvoiceCreateDate, &inv.PaymentTerms, &inv.DueDate,
		&inv.Payee, &inv.PayeeAddress, &inv.Payor, &inv.PayorAddress, &inv.PayorPhone,
		&inv.MondayDate, &inv.MondayIn, &inv.MondayOut, &inv.MondayHours, &inv.MondayRate, &inv.MondayTotal,
		&inv.TuesdayDate, &inv.TuesdayIn, &inv.TuesdayOut, &inv.TuesdayHours, &inv.TuesdayRate, &inv.TuesdayTotal,
		&inv.WednesdayDate, &inv.WednesdayIn, &inv.WednesdayOut, &inv.WednesdayHours, &inv.WednesdayRate, &inv.WednesdayTotal,
		&inv.ThursdayDate, &inv.ThursdayIn, &inv.ThursdayOut, &inv.ThursdayHours, &inv.ThursdayRate, &inv.ThursdayTotal,
		&inv.FridayDate, &inv.FridayIn, &inv.FridayOut, &inv.FridayHours, &inv.FridayRate, &inv.FridayTotal,
		&inv.TotalHours, &inv.LineTotal, &submitted, &paid, &createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invoice %s not found", invoiceNumber)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query invoice: %w", err)
	}

	inv.Submitted = submitted == 1
	inv.Paid = paid == 1
	inv.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

	return &inv, nil
}

// GetInvoicesByStatus retrieves invoices filtered by status
func (d *Database) GetInvoicesByStatus(submitted, paid *bool) ([]Invoice, error) {
	query := `
		SELECT
			pk, invoice_number, invoice_create_date, payment_terms, due_date,
			payee, payee_address, payor, payor_address, payor_phone,
			monday_date, monday_in, monday_out, monday_hours_worked, monday_unit_price, monday_line_total,
			tuesday_date, tuesday_in, tuesday_out, tuesday_hours_worked, tuesday_unit_price, tuesday_line_total,
			wednesday_date, wednesday_in, wednesday_out, wednesday_hours_worked, wednesday_unit_price, wednesday_line_total,
			thursday_date, thursday_in, thursday_out, thursday_hours_worked, thursday_unit_price, thursday_line_total,
			friday_date, friday_in, friday_out, friday_hours_worked, friday_unit_price, friday_line_total,
			total_hours, line_total, submitted, paid, created_at
		FROM invoices
		WHERE 1=1
	`

	args := []interface{}{}
	if submitted != nil {
		query += " AND submitted = ?"
		if *submitted {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}
	if paid != nil {
		query += " AND paid = ?"
		if *paid {
			args = append(args, 1)
		} else {
			args = append(args, 0)
		}
	}

	query += " ORDER BY pk ASC"

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query invoices: %w", err)
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var inv Invoice
		var submittedInt, paidInt int
		var createdAt string

		err := rows.Scan(
			&inv.PK, &inv.InvoiceNumber, &inv.InvoiceCreateDate, &inv.PaymentTerms, &inv.DueDate,
			&inv.Payee, &inv.PayeeAddress, &inv.Payor, &inv.PayorAddress, &inv.PayorPhone,
			&inv.MondayDate, &inv.MondayIn, &inv.MondayOut, &inv.MondayHours, &inv.MondayRate, &inv.MondayTotal,
			&inv.TuesdayDate, &inv.TuesdayIn, &inv.TuesdayOut, &inv.TuesdayHours, &inv.TuesdayRate, &inv.TuesdayTotal,
			&inv.WednesdayDate, &inv.WednesdayIn, &inv.WednesdayOut, &inv.WednesdayHours, &inv.WednesdayRate, &inv.WednesdayTotal,
			&inv.ThursdayDate, &inv.ThursdayIn, &inv.ThursdayOut, &inv.ThursdayHours, &inv.ThursdayRate, &inv.ThursdayTotal,
			&inv.FridayDate, &inv.FridayIn, &inv.FridayOut, &inv.FridayHours, &inv.FridayRate, &inv.FridayTotal,
			&inv.TotalHours, &inv.LineTotal, &submittedInt, &paidInt, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan invoice: %w", err)
		}

		inv.Submitted = submittedInt == 1
		inv.Paid = paidInt == 1
		inv.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)

		invoices = append(invoices, inv)
	}

	return invoices, nil
}

// MarkAsSubmitted marks an invoice as submitted
func (d *Database) MarkAsSubmitted(invoiceNumber string) error {
	query := `UPDATE invoices SET submitted = 1 WHERE invoice_number = ?`
	result, err := d.db.Exec(query, invoiceNumber)
	if err != nil {
		return fmt.Errorf("failed to mark invoice as submitted: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("invoice %s not found", invoiceNumber)
	}

	return nil
}

// MarkAsPaid marks an invoice as paid (and submitted)
func (d *Database) MarkAsPaid(invoiceNumber string) error {
	query := `UPDATE invoices SET submitted = 1, paid = 1 WHERE invoice_number = ?`
	result, err := d.db.Exec(query, invoiceNumber)
	if err != nil {
		return fmt.Errorf("failed to mark invoice as paid: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("invoice %s not found", invoiceNumber)
	}

	return nil
}

// GetSummaryStats calculates aggregate statistics for all invoices
func (d *Database) GetSummaryStats() (*InvoiceSummary, error) {
	query := `
		SELECT
			COUNT(*) as total_count,
			SUM(CASE WHEN submitted = 1 THEN 1 ELSE 0 END) as submitted_count,
			SUM(CASE WHEN paid = 1 THEN 1 ELSE 0 END) as paid_count,
			SUM(CASE WHEN submitted = 0 THEN 1 ELSE 0 END) as pending_count,
			SUM(CASE WHEN paid = 0 THEN 1 ELSE 0 END) as unpaid_count,
			SUM(line_total) as total_amount,
			SUM(CASE WHEN submitted = 1 THEN line_total ELSE 0 END) as submitted_amount,
			SUM(CASE WHEN paid = 1 THEN line_total ELSE 0 END) as paid_amount,
			SUM(CASE WHEN submitted = 0 THEN line_total ELSE 0 END) as pending_amount,
			SUM(CASE WHEN paid = 0 THEN line_total ELSE 0 END) as unpaid_amount
		FROM invoices
	`

	var summary InvoiceSummary
	err := d.db.QueryRow(query).Scan(
		&summary.TotalCount,
		&summary.SubmittedCount,
		&summary.PaidCount,
		&summary.PendingCount,
		&summary.UnpaidCount,
		&summary.TotalAmount,
		&summary.SubmittedAmount,
		&summary.PaidAmount,
		&summary.PendingAmount,
		&summary.UnpaidAmount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get summary stats: %w", err)
	}

	return &summary, nil
}
