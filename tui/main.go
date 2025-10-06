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

type sessionState int

const (
	dashboardView sessionState = iota
	invoiceListView
	invoiceDetailView
)

type mainModel struct {
	state           sessionState
	db              *models.Database
	invoices        []models.Invoice
	summary         *models.InvoiceSummary
	selectedInvoice *models.Invoice
	width           int
	height          int
	err             error

	dashboard     views.DashboardModel
	invoiceList   views.InvoiceListModel
	invoiceDetail views.InvoiceDetailModel

	message  string
	loading  bool
	quitting bool
}

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
	dbPath := "../invoices.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("Error: Database not found at %s\n", dbPath)
		fmt.Println("\nUsage: invoice-tui [path/to/invoices.db]")
		os.Exit(1)
	}

	db, err := models.OpenDatabase(dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	p := tea.NewProgram(
		initialModel(db),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func initialModel(db *models.Database) mainModel {
	return mainModel{
		state:         dashboardView,
		db:            db,
		loading:       true,
		dashboard:     views.NewDashboard(80, 24),
		invoiceList:   views.NewInvoiceList(80, 24),
		invoiceDetail: views.NewInvoiceDetail(80, 24),
	}
}

func (m mainModel) Init() tea.Cmd {
	return loadDataCmd(m.db)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
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
		m.invoices, m.summary = msg.invoices, msg.summary
		m.dashboard = m.dashboard.SetSummary(msg.summary)
		m.invoiceList = m.invoiceList.SetInvoices(msg.invoices)
		return m, nil

	case invoiceUpdatedMsg:
		if msg.success {
			m.message = msg.message
			return m, loadDataCmd(m.db)
		}
		m.err = msg.err
		return m, nil

	case pdfGeneratedMsg:
		if msg.success {
			m.message = msg.message
			m.invoiceDetail = m.invoiceDetail.SetMessage(msg.message)
		} else {
			m.err = msg.err
		}
		return m, nil
	}

	return m.updateActiveView(msg)
}

func (m mainModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.state == dashboardView {
			m.quitting = true
			return m, tea.Quit
		}
		m.state = dashboardView
		m.message, m.err = "", nil
		return m, nil

	case "esc":
		m.state = dashboardView
		m.message, m.err = "", nil
		return m, nil
	}

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

func (m mainModel) handleDashboardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "enter" {
		switch m.dashboard.SelectedAction() {
		case "📋 View All Invoices":
			m.state = invoiceListView
			m.invoiceList = m.invoiceList.SetFilter("all")
		case "📄 Generate Invoice":
			m.state = invoiceListView
			m.invoiceList = m.invoiceList.SetFilter("pending")
			m.message = "Select an invoice to generate PDF"
		case "✅ Approve Invoices":
			m.state = invoiceListView
			m.invoiceList = m.invoiceList.SetFilter("pending")
			m.message = "Select invoices to approve"
		case "🚪 Exit":
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.dashboard, cmd = m.dashboard.Update(msg)
	return m, cmd
}

func (m mainModel) handleInvoiceListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if selected := m.invoiceList.SelectedInvoice(); selected != nil {
			m.selectedInvoice = selected
			m.invoiceDetail = views.NewInvoiceDetail(m.width, m.height).SetInvoice(selected)
			m.state = invoiceDetailView
		}
		return m, nil

	case "f":
		filters := []string{"all", "pending", "submitted", "paid"}
		current := m.invoiceList.GetFilter()
		for i, f := range filters {
			if f == current {
				m.invoiceList = m.invoiceList.SetFilter(filters[(i+1)%len(filters)])
				break
			}
		}
		return m, nil

	case "g":
		if selected := m.invoiceList.SelectedInvoice(); selected != nil {
			return m, generatePDFCmd(selected.InvoiceNumber)
		}
		return m, nil

	case "s":
		if selected := m.invoiceList.SelectedInvoice(); selected != nil && !selected.Submitted {
			return m, sendInvoiceEmailCmdGo(m.db, selected)
		}
		return m, nil

	case "p":
		if selected := m.invoiceList.SelectedInvoice(); selected != nil && !selected.Paid {
			return m, markPaidCmd(m.db, selected.InvoiceNumber)
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.invoiceList, cmd = m.invoiceList.Update(msg)
	return m, cmd
}

func (m mainModel) handleInvoiceDetailKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "g":
		if m.selectedInvoice != nil {
			return m, generatePDFCmd(m.selectedInvoice.InvoiceNumber)
		}

	case "s":
		if m.selectedInvoice != nil && !m.selectedInvoice.Submitted {
			return m, sendInvoiceEmailCmdGo(m.db, m.selectedInvoice)
		}

	case "p":
		if m.selectedInvoice != nil && !m.selectedInvoice.Paid {
			return m, markPaidCmd(m.db, m.selectedInvoice.InvoiceNumber)
		}

	case "esc":
		m.state = invoiceListView
		m.invoiceDetail = m.invoiceDetail.ClearMessage()
		return m, nil
	}

	var cmd tea.Cmd
	m.invoiceDetail, cmd = m.invoiceDetail.Update(msg)
	return m, cmd
}

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

func (m mainModel) View() string {
	if m.quitting {
		return styles.SuccessStyle.Render("👋 Goodbye!\n")
	}
	if m.loading {
		return styles.InfoBoxStyle().Render("🔄 Loading invoice data...")
	}
	if m.err != nil {
		return styles.ErrorBoxStyle().Render(fmt.Sprintf("❌ Error: %v\n\nPress 'q' to quit", m.err))
	}

	content := ""
	if m.message != "" {
		content = styles.SuccessBoxStyle().Render(m.message) + "\n\n"
	}

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

func generatePDFCmd(invoiceNumber string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := os.Getwd()
		if err != nil {
			return pdfGeneratedMsg{invoiceNumber: invoiceNumber, success: false, err: err}
		}

		parentDir := filepath.Dir(currentDir)
		scriptPath := filepath.Join(parentDir, "invoice_cli.py")

		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("invoice_cli.py not found"),
			}
		}

		shellCmd := fmt.Sprintf("cd %s && /usr/bin/python3 invoice_cli.py generate %s", parentDir, invoiceNumber)
		cmd := exec.Command("/bin/sh", "-c", shellCmd)
		output, err := cmd.CombinedOutput()

		if err != nil {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("generation failed: %v", err),
			}
		}

		pdfPath := filepath.Join(parentDir, "invoices", fmt.Sprintf("invoice_%s.pdf", invoiceNumber))
		if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
			return pdfGeneratedMsg{
				invoiceNumber: invoiceNumber,
				success:       false,
				err:           fmt.Errorf("PDF not found: %s", string(output)),
			}
		}

		return pdfGeneratedMsg{
			invoiceNumber: invoiceNumber,
			success:       true,
			message:       fmt.Sprintf("✅ PDF generated: invoice_%s.pdf", invoiceNumber),
		}
	}
}

func markPaidCmd(db *models.Database, invoiceNumber string) tea.Cmd {
	return func() tea.Msg {
		if err := db.MarkAsPaid(invoiceNumber); err != nil {
			return invoiceUpdatedMsg{success: false, err: err}
		}
		return invoiceUpdatedMsg{
			success: true,
			message: fmt.Sprintf("✅ Invoice %s marked as paid", invoiceNumber),
		}
	}
}

func sendInvoiceEmailCmdGo(db *models.Database, invoice *models.Invoice) tea.Cmd {
	return func() tea.Msg {
		pdfPath := filepath.Join("..", "invoices", fmt.Sprintf("invoice_%s.pdf", invoice.InvoiceNumber))

		if err := sendInvoiceEmail(invoice, pdfPath); err != nil {
			return invoiceUpdatedMsg{success: false, err: fmt.Errorf("email failed: %v", err)}
		}

		if err := db.MarkAsSubmitted(invoice.InvoiceNumber); err != nil {
			return invoiceUpdatedMsg{success: false, err: fmt.Errorf("email sent but DB update failed: %v", err)}
		}

		return invoiceUpdatedMsg{
			success: true,
			message: fmt.Sprintf("✅ Invoice %s emailed and marked submitted", invoice.InvoiceNumber),
		}
	}
}
