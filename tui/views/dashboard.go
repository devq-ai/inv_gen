package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/devqai/invoice-tui/models"
	"github.com/devqai/invoice-tui/styles"
)

// DashboardModel represents the main dashboard view
type DashboardModel struct {
	summary    *models.InvoiceSummary
	list       list.Model
	width      int
	height     int
	err        error
	loading    bool
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// NewDashboard creates a new dashboard view
func NewDashboard(width, height int) DashboardModel {
	items := []list.Item{
		item{title: "📋 View All Invoices", desc: "Browse and manage all invoices"},
		item{title: "📄 Generate Invoice", desc: "Create PDF for an invoice"},
		item{title: "✅ Approve Invoices", desc: "Mark invoices as submitted"},
		item{title: "🚪 Exit", desc: "Quit the application"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(styles.NeonPink).
		BorderLeftForeground(styles.NeonPink)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(styles.DimTextColor)

	l := list.New(items, delegate, 0, 0)
	l.Title = "Quick Actions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = styles.TitleStyle
	l.Styles.TitleBar = lipgloss.NewStyle().Padding(0, 0, 1, 0)

	return DashboardModel{
		list:    l,
		width:   width,
		height:  height,
		loading: true,
	}
}

// SetSize updates the dashboard dimensions
func (m DashboardModel) SetSize(width, height int) DashboardModel {
	m.width = width
	m.height = height
	m.list.SetSize(width-8, 10)
	return m
}

// SetSummary updates the summary data
func (m DashboardModel) SetSummary(summary *models.InvoiceSummary) DashboardModel {
	m.summary = summary
	m.loading = false
	return m
}

// SetError sets an error message
func (m DashboardModel) SetError(err error) DashboardModel {
	m.err = err
	m.loading = false
	return m
}

// Update handles dashboard events
func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Return selected action
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// SelectedAction returns the currently selected action
func (m DashboardModel) SelectedAction() string {
	if i, ok := m.list.SelectedItem().(item); ok {
		return i.title
	}
	return ""
}

// View renders the dashboard
func (m DashboardModel) View() string {
	if m.loading {
		return styles.InfoBoxStyle().Render("Loading dashboard data...")
	}

	if m.err != nil {
		return styles.ErrorBoxStyle().Render(fmt.Sprintf("Error: %v", m.err))
	}

	var sections []string

	// Title
	title := styles.TitleStyle.
		Width(m.width - 4).
		Align(lipgloss.Center).
		Render("📋 INVOICE MANAGEMENT SYSTEM")
	sections = append(sections, title)

	// Financial Overview
	if m.summary != nil {
		overview := m.renderFinancialOverview()
		sections = append(sections, overview)
	}

	// Quick Actions List
	actions := m.list.View()
	sections = append(sections, actions)

	// Help text
	help := styles.HelpStyle.Render("[↑/↓] Navigate  [Enter] Select  [q] Quit")
	sections = append(sections, help)

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return styles.DocStyle(m.width, m.height).Render(content)
}

// renderFinancialOverview creates the financial statistics panel
func (m DashboardModel) renderFinancialOverview() string {
	s := m.summary

	// Calculate percentages
	submittedPct := 0.0
	paidPct := 0.0
	if s.TotalCount > 0 {
		submittedPct = float64(s.SubmittedCount) / float64(s.TotalCount) * 100
		paidPct = float64(s.PaidCount) / float64(s.TotalCount) * 100
	}

	// Build overview content
	var lines []string

	// Section title
	lines = append(lines, styles.SubtitleStyle.Render("Financial Overview"))
	lines = append(lines, "")

	// Invoice counts
	countSection := []string{
		fmt.Sprintf("%-20s %d", "Total Invoices:", s.TotalCount),
		fmt.Sprintf("%-20s %s",
			"Total Value:",
			styles.CurrencyStyle.Render(models.FormatAmountCompact(s.TotalAmount))),
		"",
	}
	lines = append(lines, countSection...)

	// Status breakdown with icons
	statusLines := []string{
		fmt.Sprintf("%s %-15s %3d    %s",
			styles.GetStatusIcon("submitted"),
			"Submitted:",
			s.SubmittedCount,
			styles.CurrencyStyle.Render(models.FormatAmountCompact(s.SubmittedAmount))),
		fmt.Sprintf("%s %-15s %3d    %s",
			styles.GetStatusIcon("paid"),
			"Paid:",
			s.PaidCount,
			styles.CurrencyStyle.Render(models.FormatAmountCompact(s.PaidAmount))),
		fmt.Sprintf("%s %-15s %3d    %s",
			styles.GetStatusIcon("pending"),
			"Pending:",
			s.PendingCount,
			styles.CurrencyStyle.Render(models.FormatAmountCompact(s.PendingAmount))),
	}
	lines = append(lines, statusLines...)
	lines = append(lines, "")

	// Progress indicators
	if s.TotalCount > 0 {
		lines = append(lines, fmt.Sprintf("Submitted: %.1f%% | Paid: %.1f%%", submittedPct, paidPct))
		lines = append(lines, m.renderProgressBar(submittedPct, 30))
	}

	content := strings.Join(lines, "\n")
	return styles.PanelStyle.Width(m.width - 8).Render(content)
}

// renderProgressBar creates a simple text-based progress bar
func (m DashboardModel) renderProgressBar(percentage float64, width int) string {
	filled := int(percentage / 100.0 * float64(width))
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return styles.ProgressBarStyle.Render(bar)
}

// GetListModel returns the underlying list model (for selection handling)
func (m DashboardModel) GetListModel() list.Model {
	return m.list
}
