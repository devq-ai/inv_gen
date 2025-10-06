package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	"github.com/devqai/invoice-tui/models"
	"github.com/devqai/invoice-tui/styles"
	"github.com/devqai/invoice-tui/views"
)

// sessionState represents the current view/screen
type sessionState int

const (
	dashboardView sessionState = iota
	invoiceListView
	invoiceDetailView
	reportsView
	settingsView
)

// mainModel is the top-level application model
type mainModel struct {
	state          sessionState
	db             *models.Database
	invoices       []models.Invoice
	summary        *models.InvoiceSummary
	selectedInvoice *models.Invoice
	width          int
	height         int
	err            error

	// View models
	dashboard     views.DashboardModel
	invoiceList   views.InvoiceListModel
	invoiceDetail views.InvoiceDetailModel

	// Status
	message string
	loading bool
	quitting bool
}

// Message types for async operations
type dataLoadedMsg struct {
	invoices []models.Invoice
	summary  *models.InvoiceSummary
	err      error
}

type invoiceUpdatedMsg struct {
	invoiceNumber string
	success       bool
	message       string
	err           error
}

type pdfGeneratedMsg struct {
	invoiceNumber string
	success       bool
	message       string
	err           error
}

func main() {
	// Check for database file
	dbPath := "../invoices.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	// Verify database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("Error: Database file not found at %s\n", dbPath)
		fmt.Println("\nUsage: invoice-tui [path/to/invoices.db]")
		fmt.Println("\nDefault: ../invoices.db")
		os.Exit(1)
	}

	// Open database connection
	db, err := models.OpenDatabase(dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize the application
	p := tea.NewProgram(
		initialModel(db),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

// initialModel creates the initial application model
func initialModel(db *models.Database) mainModel {
	m := mainModel{
		state:       dashboardView,
		db:          db,
		loading:     true,
		dashboard:   views.NewDashboard(80, 24),
		invoiceList: views.NewInvoiceList(80, 24),
		invoiceDetail: views.NewInvoiceDetail(80, 24),
	}

	return m
}

// Init initializes the application
func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		loadDataCmd(m.db),
	)
}

// Update handles all application events
func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboard = m.dashboard.SetSize(m.width, m.height)
		m.invoiceList = m.invoiceList.SetSize(m.width, m.height)
		m.invoiceDetail = m.invoiceDetail.SetSize(m.width, m.height)
		return m, nil

	case dataLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.invoices = msg.invoices
		m.summary = msg.summary
		m.dashboard = m.dashboard.SetSummary(msg.summary)
		m.invoiceList = m.invoiceList.SetInvoices(msg.invoices)
		return m, nil

	case invoiceUpdatedMsg:
		if msg.success {
			m.message = msg.message
			// Reload data to reflect changes
			return m, loadDataCmd(m.db)
		} else {
			m.err = msg.err
			return m, nil
		}

	case pdfGeneratedMsg:
		if msg.success {
			m.message = msg.message
			m.invoiceDetail = m.invoiceDetail.SetMessage(msg.message)
		} else {
			m.err = msg.err
		}
		return m, nil
	}

	// Update active view
	return m.updateActiveView(msg)
}

// handleKeyPress processes keyboard input
func (m mainModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == dashboardView {
			m.quitting = true
			return m, tea.Quit
		}
		// q goes back to dashboard from other views
		m.state = dashboardView
		m.message = ""
		m.err = nil
		return m, nil

	case "esc":
		// Always go back to dashboard
		m.state = dashboardView
		m.message = ""
		m.err = nil
		return m, nil
	}

	// State-specific key handling
	switch m.state {
	case dashboardView:
		return m.handleDashboardKeys(msg)
	case invoiceListView:
		return m.handleInvoiceListKeys(msg)
	case invoiceDetailView:
		return m.handleInvoiceDetailKeys(msg)
	}

	return m, nil
}

// handleDashboardKeys handles dashboard-specific keys
func (m mainModel) handleDashboardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		action := m.dashboard.SelectedAction()
		switch action {
		case "📋 View All Invoices":
			m.state = invoiceListView
			m.invoiceList = m.invoiceList.SetFilter("all")
			return m, nil
		case "📄 Generate Invoice":
			m.state = invoiceListView
			m.invoiceList = m.invoiceList.SetFilter("pending")
			m.message = "Select an invoice to generate PDF"
			return m, nil
		case "✅ Approve Invoices":
			m.state = invoiceListView
			m.invoiceList = m.invoiceList.SetFilter("pending")
			m.message = "Select invoices to approve"
			return m, nil
		case "🚪 Exit":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Update dashboard
	var cmd tea.Cmd
	m.dashboard, cmd = m.dashboard.Update(msg)
	return m, cmd
}

// handleInvoiceListKeys handles invoice list-specific keys
func (m mainModel) handleInvoiceListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// View invoice details
		selected := m.invoiceList.SelectedInvoice()
		if selected != nil {
			m.selectedInvoice = selected
			m.invoiceDetail = views.NewInvoiceDetail(m.width, m.height).SetInvoice(selected)
			m.state = invoiceDetailView
		}
		return m, nil

	case "f":
		// Cycle through filters
		currentFilter := m.invoiceList.GetFilter()
		nextFilter := "all"
		switch currentFilter {
		case "all":
			nextFilter = "pending"
		case "pending":
			nextFilter = "submitted"
		case "submitted":
			nextFilter = "paid"
		case "paid":
			nextFilter = "all"
		}
		m.invoiceList = m.invoiceList.SetFilter(nextFilter)
		return m, nil

	case "g":
		// Generate PDF for selected invoice
		selected := m.invoiceList.SelectedInvoice()
		if selected != nil {
			return m, generatePDFCmd(selected.InvoiceNumber)
		}
		return m, nil

	case "s":
		// Mark selected invoice as submitted
		selected := m.invoiceList.SelectedInvoice()
		if selected != nil && !selected.Submitted {
			return m, markSubmittedCmd(m.db, selected.InvoiceNumber)
		}
		return m, nil

	case "p":
		// Mark selected invoice as paid
		selected := m.invoiceList.SelectedInvoice()
		if selected != nil && !selected.Paid {
			return m, markPaidCmd(m.db, selected.InvoiceNumber)
		}
		return m, nil
	}

	// Update invoice list
	var cmd tea.Cmd
	m.invoiceList, cmd = m.invoiceList.Update(msg)
	return m, cmd
}

// handleInvoiceDetailKeys handles invoice detail-specific keys
func (m mainModel) handleInvoiceDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "g":
		// Generate PDF
		if m.selectedInvoice != nil {
			return m, generatePDFCmd(m.selectedInvoice.InvoiceNumber)
		}
		return m, nil

	case "s":
		// Mark as submitted
		if m.selectedInvoice != nil && !m.selectedInvoice.Submitted {
			return m, sendInvoiceEmailCmdGo(m.db, m.selectedInvoice)
		}
		return m, nil

	case "p":
		// Mark as paid
		if m.selectedInvoice != nil && !m.selectedInvoice.Paid {
			return m, markPaidCmd(m.db, m.selectedInvoice.InvoiceNumber)
		}
		return m, nil

	case "esc":
		// Go back to invoice list
		m.state = invoiceListView
		m.invoiceDetail = m.invoiceDetail.ClearMessage()
		return m, nil
	}

	// Update invoice detail
	var cmd tea.Cmd
	m.invoiceDetail, cmd = m.invoiceDetail.Update(msg)
	return m, cmd
}

// updateActiveView updates the currently active view
func (m mainModel) updateActiveView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case dashboardView:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case invoiceListView:
		m.invoiceList, cmd = m.invoiceList.Update(msg)
	case invoiceDetailView:
		m.invoiceDetail, cmd = m.invoiceDetail.Update(msg)
	}

	return m, cmd
}

// View renders the application
func (m mainModel) View() string {
	if m.quitting {
		return styles.SuccessStyle.Render("👋 Goodbye!\n")
	}

	if m.loading {
		return styles.InfoBoxStyle().Render("🔄 Loading invoice data...")
	}

	if m.err != nil {
		errorMsg := fmt.Sprintf("❌ Error: %v\n\nPress 'q' to quit", m.err)
		return styles.ErrorBoxStyle().Render(errorMsg)
	}

	// Show status message if any
	var content string
	if m.message != "" {
		msgBox := styles.SuccessBoxStyle().Render(m.message)
		content = msgBox + "\n\n"
	}

	// Render active view
	switch m.state {
	case dashboardView:
		content += m.dashboard.View()
	case invoiceListView:
		content += m.invoiceList.View()
	case invoiceDetailView:
		content += m.invoiceDetail.View()
	}
	return styles.DocStyle(m.width, m.height).Render(content)
}

// loadDataCmd loads invoice data from database
func loadDataCmd(db *models.Database) tea.Cmd {
	return func() tea.Msg {
		invoices, err := db.GetAllInvoices()
		if err != nil {
			return dataLoadedMsg{err: err}
		}
		summary, _ := db.GetSummaryStats()
		return dataLoadedMsg{invoices: invoices, summary: summary}
	}
}

// generatePDFCmd generates PDF for an invoice
func generatePDFCmd(invoiceNumber string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := os.Getwd()
		if err != nil {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("failed to get working directory: %v", err),
			}
		}

		parentDir := filepath.Dir(currentDir)
		scriptPath := filepath.Join(parentDir, "invoice_cli.py")

		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("invoice_cli.py not found at %s", scriptPath),
			}
		}

		shellCmd := fmt.Sprintf("cd %s && /usr/bin/python3 invoice_cli.py generate %s", parentDir, invoiceNumber)
		cmd := exec.Command("/bin/sh", "-c", shellCmd)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("Python exec failed: %v\nOutput: %s", err, string(output)),
			}
		}

		pdfPath := filepath.Join(parentDir, "invoices", fmt.Sprintf("invoice_%s.pdf", invoiceNumber))
		if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("PDF file not found at %s\nOutput: %s", pdfPath, string(output)),
			}
		}

		return pdfGeneratedMsg{
			invoiceNumber: invoiceNumber,
			success:       true,
			message:       fmt.Sprintf("✅ PDF generated: invoices/invoice_%s.pdf", invoiceNumber),
		}
	}
}

// markSubmittedCmd marks invoice as submitted
func markSubmittedCmd(db *models.Database, invoiceNumber string) tea.Cmd {
	return func() tea.Msg {
		err := db.MarkAsSubmitted(invoiceNumber)
		if err != nil {
			return invoiceUpdatedMsg{success: false, err: err}
		}
		return invoiceUpdatedMsg{
			success: true,
			message: fmt.Sprintf("✅ Invoice %s marked as submitted", invoiceNumber),
		}
	}
}

// markPaidCmd marks invoice as paid
func markPaidCmd(db *models.Database, invoiceNumber string) tea.Cmd {
	return func() tea.Msg {
		err := db.MarkAsPaid(invoiceNumber)
		if err != nil {
			return invoiceUpdatedMsg{success: false, err: err}
		}
		return invoiceUpdatedMsg{
			success: true,
			message: fmt.Sprintf("✅ Invoice %s marked as paid", invoiceNumber),
		}
	}
}

// sendInvoiceEmailCmdGo sends invoice via email and marks as submitted

// sendInvoiceEmailCmdGo sends invoice via Go SMTP
func sendInvoiceEmailCmdGo(db *models.Database, invoice *models.Invoice) tea.Cmd {
	return func() tea.Msg {
		pdfPath := filepath.Join("..", "invoices", fmt.Sprintf("invoice_%s.pdf", invoice.InvoiceNumber))
		
		err := sendInvoiceEmail(invoice, pdfPath)
		if err != nil {
			return invoiceUpdatedMsg{
				success: false,
				err:     fmt.Errorf("Email failed: %v", err),
			}
		}
		
		err = db.MarkAsSubmitted(invoice.InvoiceNumber)
		if err != nil {
			return invoiceUpdatedMsg{
				success: false,
				err:     fmt.Errorf("Email sent but DB update failed: %v", err),
			}
		}
		
		return invoiceUpdatedMsg{
			success: true,
			message: fmt.Sprintf("✅ Invoice %s emailed to InfoObjects (CC: you)", invoice.InvoiceNumber),
		}
	}
}
