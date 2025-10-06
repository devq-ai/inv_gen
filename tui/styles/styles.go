package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color definitions matching the specification
var (
	// Status Colors
	PendingColor   = lipgloss.Color("#A4C2F4") // pastel_blue
	SubmittedColor = lipgloss.Color("#F4A6C0") // pastel_pink
	PaidColor      = lipgloss.Color("#B5A0E3") // pastel_purple
	ErrorColor     = lipgloss.Color("#E69999") // pastel_red
	SuccessColor   = lipgloss.Color("#A1D9A0") // pastel_green
	WarningColor   = lipgloss.Color("#FFE599") // pastel_yellow
	NoteColor      = lipgloss.Color("#FFE599") // pastel_yellow

	// UI Element Colors
	BorderColor     = lipgloss.Color("#E1E4E8") // gray
	AccentColor     = lipgloss.Color("#00FFFF") // neon cyan
	TextColor       = lipgloss.Color("#FFFFFF") // white
	DimTextColor    = lipgloss.Color("#00FFFF") // neon cyan
	BackgroundColor = lipgloss.Color("#000000") // black
	HeaderColor     = lipgloss.Color("#1A1A1A") // dark gray

	// Neon Colors (for highlights)
	NeonPink   = lipgloss.Color("#FF10F0")
	NeonPurple = lipgloss.Color("#9D00FF")
	NeonGreen  = lipgloss.Color("#39FF14")
	NeonBlue   = lipgloss.Color("#1B03A3")
	NeonCyan   = lipgloss.Color("#00FFFF")
)

// Base Styles
var (
	// TitleStyle for main headings
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(NeonCyan).
			Padding(0, 1).
			MarginBottom(1)

	// SubtitleStyle for section headings
	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(NeonCyan).
			Padding(0, 1)

	// PanelStyle for bordered containers
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	// CardStyle for smaller containers
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	// HelpStyle for help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(DimTextColor).
			Padding(1, 0)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true).
			Padding(0, 1)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true).
			Padding(0, 1)

	// WarningStyle for warning messages
	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true).
			Padding(0, 1)

	// StatusPendingStyle for pending status
	StatusPendingStyle = lipgloss.NewStyle().
				Foreground(PendingColor).
				Bold(true)

	// StatusSubmittedStyle for submitted status
	StatusSubmittedStyle = lipgloss.NewStyle().
				Foreground(SubmittedColor).
				Bold(true)

	// StatusPaidStyle for paid status
	StatusPaidStyle = lipgloss.NewStyle().
			Foreground(PaidColor).
			Bold(true)

	// HighlightStyle for highlighted text
	HighlightStyle = lipgloss.NewStyle().
			Foreground(NeonPink).
			Bold(true)

	// DimStyle for dimmed text
	DimStyle = lipgloss.NewStyle().
			Foreground(DimTextColor)

	// BoldStyle for bold text
	BoldStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	// HeaderStyle for table/list headers
	HeaderStyle = lipgloss.NewStyle().
			Background(HeaderColor).
			Foreground(NeonCyan).
			Bold(true).
			Padding(0, 1)

	// CurrencyStyle for monetary amounts
	CurrencyStyle = lipgloss.NewStyle().
			Foreground(NeonGreen).
			Bold(true)

	// LabelStyle for field labels
	LabelStyle = lipgloss.NewStyle().
			Foreground(DimTextColor).
			Bold(true)

	// ValueStyle for field values
	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	// ActiveItemStyle for selected list items
	ActiveItemStyle = lipgloss.NewStyle().
			Foreground(NeonPink).
			Bold(true).
			PaddingLeft(2)

	// InactiveItemStyle for unselected list items
	InactiveItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				PaddingLeft(2)

	// SeparatorStyle for visual separators
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(BorderColor)

	// FooterStyle for footer text
	FooterStyle = lipgloss.NewStyle().
			Foreground(DimTextColor).
			Background(HeaderColor).
			Padding(0, 1)
)

// GetStatusStyle returns the appropriate style for a status string
func GetStatusStyle(status string) lipgloss.Style {
	switch status {
	case "Paid", "paid":
		return StatusPaidStyle
	case "Submitted", "submitted":
		return StatusSubmittedStyle
	case "Pending", "pending":
		return StatusPendingStyle
	default:
		return DimStyle
	}
}

// GetStatusColor returns the color for a status string
func GetStatusColor(status string) lipgloss.Color {
	switch status {
	case "Paid", "paid":
		return PaidColor
	case "Submitted", "submitted":
		return SubmittedColor
	case "Pending", "pending":
		return PendingColor
	default:
		return DimTextColor
	}
}

// GetStatusIcon returns an emoji icon for a status
func GetStatusIcon(status string) string {
	switch status {
	case "Paid", "paid":
		return "✅"
	case "Submitted", "submitted":
		return "📤"
	case "Pending", "pending":
		return "⏳"
	default:
		return "❓"
	}
}

// FormatStatusWithIcon formats a status with icon and color
func FormatStatusWithIcon(status string) string {
	icon := GetStatusIcon(status)
	style := GetStatusStyle(status)
	return icon + " " + style.Render(status)
}

// DocStyle is the main document style
func DocStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Width(width).
		Height(height)
}

// AppBorderStyle creates a styled border for the entire application
func AppBorderStyle(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(AccentColor).
		Width(width).
		Padding(0, 1)
}

// SpinnerStyle for loading indicators
var SpinnerStyle = lipgloss.NewStyle().
	Foreground(NeonPink)

// ProgressBarStyle for progress indicators
var ProgressBarStyle = lipgloss.NewStyle().
	Foreground(SuccessColor)

// TableHeaderStyle for table column headers
var TableHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(AccentColor).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(BorderColor)

// TableCellStyle for table cells
var TableCellStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFFFF")).
	Padding(0, 1)

// SelectedRowStyle for selected table rows
var SelectedRowStyle = lipgloss.NewStyle().
	Foreground(NeonPink).
	Bold(true).
	Background(lipgloss.Color("#FCE4EC"))

// AlignRight creates a right-aligned style
func AlignRight(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Right)
}

// AlignCenter creates a center-aligned style
func AlignCenter(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center)
}

// BoxStyle creates a simple box
func BoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor).
		Padding(1, 2)
}

// InfoBoxStyle creates an info box
func InfoBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonCyan).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(1, 2)
}

// SuccessBoxStyle creates a success box
func SuccessBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SuccessColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(1, 2)
}

// ErrorBoxStyle creates an error box
func ErrorBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ErrorColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(1, 2)
}

// WarningBoxStyle creates a warning box
func WarningBoxStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(WarningColor).
		Foreground(lipgloss.Color("#000000")).
		Padding(1, 2)
}
