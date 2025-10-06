package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Status colors
	PendingColor   = lipgloss.Color("#A4C2F4")
	SubmittedColor = lipgloss.Color("#F4A6C0")
	PaidColor      = lipgloss.Color("#B5A0E3")
	ErrorColor     = lipgloss.Color("#E69999")
	SuccessColor   = lipgloss.Color("#A1D9A0")
	WarningColor   = lipgloss.Color("#FFE599")

	// UI colors
	BorderColor  = lipgloss.Color("#E1E4E8")
	TextColor    = lipgloss.Color("#FFFFFF")
	DimTextColor = lipgloss.Color("#00FFFF")
	HeaderColor  = lipgloss.Color("#1A1A1A")

	// Neon accents
	NeonPink = lipgloss.Color("#FF10F0")
	NeonCyan = lipgloss.Color("#00FFFF")
	NeonGreen = lipgloss.Color("#39FF14")
)

// Base styles
var (
	TitleStyle = lipgloss.NewStyle().Bold(true).Foreground(NeonCyan).Padding(0, 1).MarginBottom(1)
	SubtitleStyle = lipgloss.NewStyle().Bold(true).Foreground(NeonCyan).Padding(0, 1)
	HelpStyle = lipgloss.NewStyle().Foreground(DimTextColor).Padding(1, 0)

	ErrorStyle = lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Padding(0, 1)
	SuccessStyle = lipgloss.NewStyle().Foreground(SuccessColor).Bold(true).Padding(0, 1)
	WarningStyle = lipgloss.NewStyle().Foreground(WarningColor).Bold(true).Padding(0, 1)

	StatusPendingStyle = lipgloss.NewStyle().Foreground(PendingColor).Bold(true)
	StatusSubmittedStyle = lipgloss.NewStyle().Foreground(SubmittedColor).Bold(true)
	StatusPaidStyle = lipgloss.NewStyle().Foreground(PaidColor).Bold(true)

	HighlightStyle = lipgloss.NewStyle().Foreground(NeonPink).Bold(true)
	DimStyle = lipgloss.NewStyle().Foreground(DimTextColor)
	BoldStyle = lipgloss.NewStyle().Bold(true).Foreground(TextColor)

	HeaderStyle = lipgloss.NewStyle().Background(HeaderColor).Foreground(NeonCyan).Bold(true).Padding(0, 1)
	CurrencyStyle = lipgloss.NewStyle().Foreground(NeonGreen).Bold(true)
	LabelStyle = lipgloss.NewStyle().Foreground(DimTextColor).Bold(true)
	ValueStyle = lipgloss.NewStyle().Foreground(TextColor)

	ActiveItemStyle = lipgloss.NewStyle().Foreground(NeonPink).Bold(true).PaddingLeft(2)
	InactiveItemStyle = lipgloss.NewStyle().Foreground(TextColor).PaddingLeft(2)

	PanelStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(BorderColor).Padding(1, 2).MarginTop(1).MarginBottom(1)
	CardStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(BorderColor).Padding(0, 1)
	FooterStyle = lipgloss.NewStyle().Foreground(DimTextColor).Background(HeaderColor).Padding(0, 1)

	TableHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(NeonCyan).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).BorderForeground(BorderColor)
	TableCellStyle = lipgloss.NewStyle().Foreground(TextColor).Padding(0, 1)
)

func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Paid", "paid":
		return StatusPaidStyle
	case "Submitted", "submitted":
		return StatusSubmittedStyle
	default:
		return StatusPendingStyle
	}
}

func GetStatusIcon(status string) string {
	switch status {
	case "Paid", "paid":
		return "✅"
	case "Submitted", "submitted":
		return "📤"
	default:
		return "⏳"
	}
}

func FormatStatusWithIcon(status string) string {
	return GetStatusIcon(status) + " " + GetStatusStyle(status).Render(status)
}

func DocStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2).Width(width).Height(height)
}

func BoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(BorderColor).Padding(1, 2)
}

func InfoBoxStyle() lipgloss.Style {
	return BoxStyle().BorderForeground(NeonCyan).Foreground(TextColor)
}

func SuccessBoxStyle() lipgloss.Style {
	return BoxStyle().BorderForeground(SuccessColor).Foreground(TextColor)
}

func ErrorBoxStyle() lipgloss.Style {
	return BoxStyle().BorderForeground(ErrorColor).Foreground(TextColor)
}

func WarningBoxStyle() lipgloss.Style {
	return BoxStyle().BorderForeground(WarningColor).Foreground(lipgloss.Color("#000000"))
}
