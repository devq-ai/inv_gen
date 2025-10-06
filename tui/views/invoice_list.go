package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devqai/invoice-tui/models"
	"github.com/devqai/invoice-tui/styles"
)

// InvoiceListModel represents the invoice list view
type InvoiceListModel struct {
	table     table.Model
	invoices  []models.Invoice
	width     int
	height    int
	err       error
	loading   bool
	filter    string // "all", "pending", "submitted", "paid"
}

// NewInvoiceList creates a new invoice list view
func NewInvoiceList(width, height int) InvoiceListModel {
	columns := []table.Column{
		{Title: "Invoice", Width: 10},
		{Title: "Created", Width: 12},
		{Title: "Due Date", Width: 12},
		{Title: "Week Ending", Width: 16},
		{Title: "Amount", Width: 12},
		{Title: "Status", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.BorderColor).
		BorderBottom(true).
		Bold(true).
		Foreground(styles.AccentColor)
	s.Selected = s.Selected.
		Foreground(styles.NeonPink).
		Bold(true).
		Background(lipgloss.Color("#FCE4EC"))

	t.SetStyles(s)

	return InvoiceListModel{
		table:   t,
		width:   width,
		height:  height,
		loading: true,
		filter:  "all",
	}
}

// SetSize updates the list dimensions
func (m InvoiceListModel) SetSize(width, height int) InvoiceListModel {
	m.width = width
	m.height = height
	m.table.SetWidth(width - 8)
	m.table.SetHeight(height - 12)
	return m
}

// SetInvoices updates the invoice data
func (m InvoiceListModel) SetInvoices(invoices []models.Invoice) InvoiceListModel {
	m.invoices = invoices
	m.loading = false
	m.updateTableRows()
	return m
}

// SetError sets an error message
func (m InvoiceListModel) SetError(err error) InvoiceListModel {
	m.err = err
	m.loading = false
	return m
}

// SetFilter changes the active filter
func (m InvoiceListModel) SetFilter(filter string) InvoiceListModel {
	m.filter = filter
	m.updateTableRows()
	return m
}

// updateTableRows rebuilds the table rows based on current filter
func (m *InvoiceListModel) updateTableRows() {
	var rows []table.Row

	for _, inv := range m.invoices {
		// Apply filter
		if m.filter == "pending" && (inv.Submitted || inv.Paid) {
			continue
		}
		if m.filter == "submitted" && (!inv.Submitted || inv.Paid) {
			continue
		}
		if m.filter == "paid" && !inv.Paid {
			continue
		}

		// Format status with icon
		status := inv.StatusString()
		statusIcon := styles.GetStatusIcon(status)
		statusDisplay := fmt.Sprintf("%s %s", statusIcon, status)

		row := table.Row{
			inv.InvoiceNumber,
			inv.InvoiceCreateDate,
			inv.DueDate,
			inv.GetWeekEnding(),
			models.FormatAmountCompact(inv.LineTotal),
			statusDisplay,
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
}

// Update handles list events
func (m InvoiceListModel) Update(msg tea.Msg) (InvoiceListModel, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// SelectedInvoice returns the currently selected invoice
func (m InvoiceListModel) SelectedInvoice() *models.Invoice {
	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.invoices) {
		return nil
	}

	// Account for filtering
	visibleIndex := 0
	for i, inv := range m.invoices {
		// Apply same filter logic
		if m.filter == "pending" && (inv.Submitted || inv.Paid) {
			continue
		}
		if m.filter == "submitted" && (!inv.Submitted || inv.Paid) {
			continue
		}
		if m.filter == "paid" && !inv.Paid {
			continue
		}

		if visibleIndex == cursor {
			return &m.invoices[i]
		}
		visibleIndex++
	}

	return nil
}

// View renders the invoice list
func (m InvoiceListModel) View() string {
	if m.loading {
		return styles.InfoBoxStyle().Render("Loading invoices...")
	}

	if m.err != nil {
		return styles.ErrorBoxStyle().Render(fmt.Sprintf("Error: %v", m.err))
	}

	var sections []string

	// Title with filter indicator
	title := styles.TitleStyle.
		Width(m.width - 4).
		Align(lipgloss.Center).
		Render(m.getTitle())
	sections = append(sections, title)

	// Filter tabs
	filterTabs := m.renderFilterTabs()
	sections = append(sections, filterTabs)

	// Table
	tableView := styles.PanelStyle.Width(m.width - 8).Render(m.table.View())
	sections = append(sections, tableView)

	// Summary line
	summary := m.renderSummary()
	sections = append(sections, summary)

	// Help text
	help := styles.HelpStyle.Render(
		"[↑/↓] Navigate  [Enter] View Details  [f] Filter  [g] Generate PDF  [s] Submit  [p] Mark Paid  [Esc] Back  [q] Quit",
	)
	sections = append(sections, help)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return styles.DocStyle(m.width, m.height).Render(content)
}

// getTitle returns the title based on the current filter
func (m InvoiceListModel) getTitle() string {
	switch m.filter {
	case "pending":
		return "📋 Pending Invoices"
	case "submitted":
		return "📤 Submitted Invoices"
	case "paid":
		return "✅ Paid Invoices"
	default:
		return "📋 All Invoices"
	}
}

// renderFilterTabs creates the filter selection tabs
func (m InvoiceListModel) renderFilterTabs() string {
	tabs := []string{"All", "Pending", "Submitted", "Paid"}
	filterMap := map[string]string{
		"All":       "all",
		"Pending":   "pending",
		"Submitted": "submitted",
		"Paid":      "paid",
	}

	var tabViews []string
	for _, tab := range tabs {
		var style lipgloss.Style
		if filterMap[tab] == m.filter {
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.NeonPink).
				Background(lipgloss.Color("#FCE4EC")).
				Padding(0, 2)
		} else {
			style = lipgloss.NewStyle().
				Foreground(styles.DimTextColor).
				Padding(0, 2)
		}
		tabViews = append(tabViews, style.Render(tab))
	}

	return lipgloss.NewStyle().
		MarginBottom(1).
		Render(strings.Join(tabViews, " | "))
}

// renderSummary creates a summary line
func (m InvoiceListModel) renderSummary() string {
	// Calculate visible invoices statistics
	var count int
	var totalAmount float64

	for _, inv := range m.invoices {
		// Apply filter
		if m.filter == "pending" && (inv.Submitted || inv.Paid) {
			continue
		}
		if m.filter == "submitted" && (!inv.Submitted || inv.Paid) {
			continue
		}
		if m.filter == "paid" && !inv.Paid {
			continue
		}

		count++
		totalAmount += inv.LineTotal
	}

	summary := fmt.Sprintf(
		"Showing %d invoices | Total: %s",
		count,
		styles.CurrencyStyle.Render(models.FormatAmountCompact(totalAmount)),
	)

	return styles.DimStyle.Render(summary)
}

// GetTable returns the underlying table model (for testing)
func (m InvoiceListModel) GetTable() table.Model {
	return m.table
}

// GetFilter returns the current filter
func (m InvoiceListModel) GetFilter() string {
	return m.filter
}
