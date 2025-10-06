package views

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devqai/invoice-tui/models"
	"github.com/devqai/invoice-tui/styles"
)

// InvoiceDetailModel represents the invoice detail view
type InvoiceDetailModel struct {
	invoice *models.Invoice
	width   int
	height  int
	err     error
	loading bool
	message string
}

// NewInvoiceDetail creates a new invoice detail view
func NewInvoiceDetail(width, height int) InvoiceDetailModel {
	return InvoiceDetailModel{
		width:   width,
		height:  height,
		loading: true,
	}
}

// SetSize updates the detail view dimensions
func (m InvoiceDetailModel) SetSize(width, height int) InvoiceDetailModel {
	m.width = width
	m.height = height
	return m
}

// SetInvoice updates the invoice data
func (m InvoiceDetailModel) SetInvoice(invoice *models.Invoice) InvoiceDetailModel {
	m.invoice = invoice
	m.loading = false
	return m
}

// SetError sets an error message
func (m InvoiceDetailModel) SetError(err error) InvoiceDetailModel {
	m.err = err
	m.loading = false
	return m
}

// SetMessage sets a status message
func (m InvoiceDetailModel) SetMessage(msg string) InvoiceDetailModel {
	m.message = msg
	return m
}

// Update handles detail view events
func (m InvoiceDetailModel) Update(msg tea.Msg) (InvoiceDetailModel, tea.Cmd) {
	// No internal updates needed for static detail view
	return m, nil
}

// View renders the invoice detail
func (m InvoiceDetailModel) View() string {
	if m.loading {
		return styles.InfoBoxStyle().Render("Loading invoice details...")
	}

	if m.err != nil {
		return styles.ErrorBoxStyle().Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.invoice == nil {
		return styles.WarningBoxStyle().Render("No invoice selected")
	}

	var sections []string

	// Title
	title := styles.TitleStyle.
		Width(m.width - 4).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("📄 Invoice Details: %s", m.invoice.InvoiceNumber))
	sections = append(sections, title)

	// Status message if any
	if m.message != "" {
		msgBox := styles.SuccessBoxStyle().Width(m.width - 8).Render(m.message)
		sections = append(sections, msgBox)
	}

	// PDF Status Indicator
	pdfStatus := m.renderPDFStatus()
	sections = append(sections, pdfStatus)

	// Basic information
	basicInfo := m.renderBasicInfo()
	sections = append(sections, basicInfo)

	// Parties
	parties := m.renderParties()
	sections = append(sections, parties)

	// Work details
	workDetails := m.renderWorkDetails()
	sections = append(sections, workDetails)

	// Totals
	totals := m.renderTotals()
	sections = append(sections, totals)

	// Help text
	help := styles.HelpStyle.Render(
		"[g] Generate PDF  [s] Submit  [p] Mark Paid  [Esc] Back  [q] Quit",
	)
	sections = append(sections, help)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return styles.DocStyle(m.width, m.height).Render(content)
}

// renderBasicInfo renders the basic invoice information
func (m InvoiceDetailModel) renderBasicInfo() string {
	inv := m.invoice

	var lines []string
	lines = append(lines, styles.SubtitleStyle.Render("Invoice Information"))
	lines = append(lines, "")

	// Status with icon and color
	statusDisplay := styles.FormatStatusWithIcon(inv.StatusString())
	lines = append(lines, fmt.Sprintf("%-20s %s", styles.LabelStyle.Render("Status:"), statusDisplay))

	// Dates
	lines = append(lines, fmt.Sprintf("%-20s %s",
		styles.LabelStyle.Render("Created:"),
		styles.ValueStyle.Render(inv.InvoiceCreateDate)))
	lines = append(lines, fmt.Sprintf("%-20s %s",
		styles.LabelStyle.Render("Due Date:"),
		styles.ValueStyle.Render(inv.DueDate)))
	lines = append(lines, fmt.Sprintf("%-20s %s",
		styles.LabelStyle.Render("Payment Terms:"),
		styles.ValueStyle.Render(fmt.Sprintf("Net %d days", inv.PaymentTerms))))
	lines = append(lines, fmt.Sprintf("%-20s %s",
		styles.LabelStyle.Render("Week Ending:"),
		styles.ValueStyle.Render(inv.GetWeekEnding())))

	content := strings.Join(lines, "\n")
	return styles.PanelStyle.Width(m.width - 8).Render(content)
}

// renderParties renders the payee and payor information
func (m InvoiceDetailModel) renderParties() string {
	inv := m.invoice

	// Payee section
	payeeLines := []string{
		styles.SubtitleStyle.Render("FROM (Payee):"),
		"",
		styles.BoldStyle.Render(inv.Payee),
		styles.DimStyle.Render(inv.PayeeAddress),
	}

	// Payor section
	payorLines := []string{
		styles.SubtitleStyle.Render("TO (Payor):"),
		"",
		styles.BoldStyle.Render(inv.Payor),
		styles.DimStyle.Render(inv.PayorAddress),
		styles.DimStyle.Render(inv.PayorPhone),
	}

	payeeBox := lipgloss.NewStyle().
		Width((m.width - 12) / 2).
		Render(strings.Join(payeeLines, "\n"))

	payorBox := lipgloss.NewStyle().
		Width((m.width - 12) / 2).
		Render(strings.Join(payorLines, "\n"))

	parties := lipgloss.JoinHorizontal(lipgloss.Top, payeeBox, "  ", payorBox)

	return styles.PanelStyle.Width(m.width - 8).Render(parties)
}

// renderWorkDetails renders the daily work breakdown
func (m InvoiceDetailModel) renderWorkDetails() string {
	inv := m.invoice

	var lines []string
	lines = append(lines, styles.SubtitleStyle.Render("Work Details"))
	lines = append(lines, "")

	// Table header
	header := fmt.Sprintf("%-12s %-16s %-8s %-8s %-8s %-10s %10s",
		"Day", "Date", "In", "Out", "Hours", "Rate", "Total")
	lines = append(lines, styles.HeaderStyle.Render(header))
	lines = append(lines, styles.SeparatorStyle.Render(strings.Repeat("─", 80)))

	// Days of the week
	days := []struct {
		name  string
		date  string
		in    string
		out   string
		hours float64
		rate  float64
		total float64
	}{
		{"Monday", inv.MondayDate, inv.MondayIn, inv.MondayOut, inv.MondayHours, inv.MondayRate, inv.MondayTotal},
		{"Tuesday", inv.TuesdayDate, inv.TuesdayIn, inv.TuesdayOut, inv.TuesdayHours, inv.TuesdayRate, inv.TuesdayTotal},
		{"Wednesday", inv.WednesdayDate, inv.WednesdayIn, inv.WednesdayOut, inv.WednesdayHours, inv.WednesdayRate, inv.WednesdayTotal},
		{"Thursday", inv.ThursdayDate, inv.ThursdayIn, inv.ThursdayOut, inv.ThursdayHours, inv.ThursdayRate, inv.ThursdayTotal},
		{"Friday", inv.FridayDate, inv.FridayIn, inv.FridayOut, inv.FridayHours, inv.FridayRate, inv.FridayTotal},
	}

	for _, day := range days {
		row := fmt.Sprintf("%-12s %-16s %-8s %-8s %8.1f  $%8.2f %10s",
			day.name,
			day.date,
			day.in,
			day.out,
			day.hours,
			day.rate,
			models.FormatAmountCompact(day.total))
		lines = append(lines, row)
	}

	content := strings.Join(lines, "\n")
	return styles.PanelStyle.Width(m.width - 8).Render(content)
}

// renderTotals renders the invoice totals
func (m InvoiceDetailModel) renderTotals() string {
	inv := m.invoice

	var lines []string
	lines = append(lines, "")
	lines = append(lines, styles.SeparatorStyle.Render(strings.Repeat("═", 50)))

	// Total hours
	totalHoursLine := fmt.Sprintf("%-30s %8.1f hours",
		styles.BoldStyle.Render("TOTAL HOURS:"),
		inv.TotalHours)
	lines = append(lines, totalHoursLine)

	// Total amount
	totalAmountLine := fmt.Sprintf("%-30s %s",
		styles.BoldStyle.Render("TOTAL AMOUNT:"),
		styles.CurrencyStyle.Render(models.FormatAmountCompact(inv.LineTotal)))
	lines = append(lines, totalAmountLine)

	lines = append(lines, styles.SeparatorStyle.Render(strings.Repeat("═", 50)))

	content := strings.Join(lines, "\n")
	return styles.PanelStyle.Width(m.width - 8).Render(content)
}

// GetInvoice returns the current invoice
func (m InvoiceDetailModel) GetInvoice() *models.Invoice {
	return m.invoice
}

// renderPDFStatus renders PDF file status indicator
func (m InvoiceDetailModel) renderPDFStatus() string {
	inv := m.invoice

	// Check if PDF exists
	pdfPath := fmt.Sprintf("../invoices/invoice_%s.pdf", inv.InvoiceNumber)
	_, err := os.Stat(pdfPath)
	pdfExists := err == nil

	var statusLine string
	if pdfExists {
		statusLine = fmt.Sprintf("PDF Status: %s %s",
			styles.BoldStyle.Foreground(lipgloss.Color("#39FF14")).Render("✓ EXISTS"),
			styles.DimStyle.Render(fmt.Sprintf("(invoices/invoice_%s.pdf)", inv.InvoiceNumber)))
	} else {
		statusLine = fmt.Sprintf("PDF Status: %s %s",
			styles.BoldStyle.Foreground(lipgloss.Color("#FF3131")).Render("✗ NOT GENERATED"),
			styles.DimStyle.Render("(press 'g' to generate)"))
	}

	return styles.PanelStyle.Width(m.width - 8).Render(statusLine)
}

// ClearMessage clears the status message
func (m InvoiceDetailModel) ClearMessage() InvoiceDetailModel {
	m.message = ""
	return m
}
