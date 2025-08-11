package pages

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kirincap/skybit-cli/internal/theme"
)

type ContextModel struct {
	width  int
	height int

	vp viewport.Model
}

func NewContextModel() ContextModel {
	vp := viewport.New(40, 10)
	vp.Style = theme.Styles.Context
	m := ContextModel{vp: vp}
	m.refresh()
	return m
}

func (m ContextModel) WithSize(width, height int) ContextModel {
	m.width, m.height = width, height
	frame := theme.Styles.Panel.GetHorizontalFrameSize()
	m.vp.Width = max(20, width-frame-2)
	m.vp.Height = max(3, height-frame-3)
	return m
}

func (m ContextModel) Init() tea.Cmd { return nil }

func (m ContextModel) Update(msg tea.Msg) (ContextModel, tea.Cmd) {
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m ContextModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(theme.Colors.Accent).Render("Context")
	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.vp.View(),
	)
	return theme.Styles.Panel.Render(body)
}

func (m *ContextModel) refresh() {
	// Placeholder quotes/orders/positions
	now := time.Now().Format(time.Kitchen)
	m.vp.SetContent(fmt.Sprintf("Quotes @ %s\n\nBTC-USD  67890.12\nETH-USD  3456.78\nSOL-USD   145.67\n\nPositions: none\nOrders: none", now))
}
