package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kirincap/skybit-cli/internal/pages"
	"github.com/kirincap/skybit-cli/internal/theme"
)

type AppModel struct {
	width  int
	height int

	chat    pages.ChatModel
	context pages.ContextModel

	focusedLeft bool
}

func New() tea.Model {
	chat := pages.NewChatModel()
	ctx := pages.NewContextModel()
	return &AppModel{
		chat:        chat,
		context:     ctx,
		focusedLeft: true,
	}
}

func (m *AppModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focusedLeft = !m.focusedLeft
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resizePanes()
	}

	var cmds []tea.Cmd
	if m.focusedLeft {
		var cmd tea.Cmd
		m.chat, cmd = m.chat.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.context, cmd = m.context.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *AppModel) View() string {
	// Layout: two columns with a thin divider
	left := m.chat.View()
	right := m.context.View()

	divider := lipgloss.NewStyle().Foreground(theme.Colors.Subtle).Render("│")

	row := lipgloss.JoinHorizontal(lipgloss.Top, left, divider, right)
	return lipgloss.NewStyle().Width(m.width).Height(m.height).Render(row)
}

func (m *AppModel) resizePanes() {
	if m.width == 0 || m.height == 0 {
		return
	}
	// Borders and divider take space. Split ~60/40.
	dividerWidth := 1
	leftWidth := int(float32(m.width-dividerWidth) * 0.58)
	rightWidth := (m.width - dividerWidth) - leftWidth

	// Leave a row for status/help line in chat
	m.chat = m.chat.WithSize(leftWidth, m.height)
	m.context = m.context.WithSize(rightWidth, m.height)
}

// Debug helper optional to show size (not used)
func sizeString(w, h int) string { return fmt.Sprintf("%dx%d", w, h) }
