package theme

import "github.com/charmbracelet/lipgloss"

var Colors = struct {
	Accent lipgloss.Color
	Subtle lipgloss.Color
}{
	Accent: lipgloss.Color("205"),
	Subtle: lipgloss.Color("240"),
}

var Styles = struct {
	Panel           lipgloss.Style
	ChatHistory     lipgloss.Style
	Context         lipgloss.Style
	InputBorder     lipgloss.Style
	InputCursorLine lipgloss.Style
	Help            lipgloss.Style
}{
	Panel:           lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(Colors.Subtle).Padding(0, 1),
	ChatHistory:     lipgloss.NewStyle(),
	Context:         lipgloss.NewStyle(),
	InputBorder:     lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(Colors.Subtle),
	InputCursorLine: lipgloss.NewStyle().Background(lipgloss.Color("235")),
	Help:            lipgloss.NewStyle().Foreground(Colors.Subtle),
}
