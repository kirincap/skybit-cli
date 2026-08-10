package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/kirincap/skybit-cli/internal/theme"
)

type ChatModel struct {
	width  int
	height int

	history viewport.Model
	input   textarea.Model

	messages []string
}

func NewChatModel() ChatModel {
	ta := textarea.New()
	ta.Placeholder = "Type a message and press Enter"
	ta.FocusedStyle.CursorLine = theme.Styles.InputCursorLine
	ta.Prompt = "› "
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.SetHeight(3)
	ta.SetWidth(40)

	vp := viewport.New(40, 10)
	vp.Style = theme.Styles.ChatHistory

	return ChatModel{
		history:  vp,
		input:    ta,
		messages: []string{"Welcome to Skybit CLI!", "Press Tab to switch panes.", "Press Ctrl+C to quit."},
	}
}

func (m ChatModel) WithSize(width, height int) ChatModel {
	m.width, m.height = width, height
	// Outer border consumes space
	outer := theme.Styles.Panel.GetHorizontalFrameSize()
	innerWidth := max(20, width-outer)
	// Reserve input height + border spacing
	inputHeight := 5
	historyHeight := max(3, height-inputHeight-outer)

	m.input.SetWidth(innerWidth - 2)
	m.history.Width = innerWidth - 2
	m.history.Height = historyHeight - 2
	m.refreshHistory()
	return m
}

func (m ChatModel) Init() tea.Cmd { return textarea.Blink }

func (m ChatModel) Update(msg tea.Msg) (ChatModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			text := strings.TrimSpace(m.input.Value())
			if text != "" {
				m.messages = append(m.messages, "> "+text)
				// echo a placeholder reply
				m.messages = append(m.messages, "…ack: "+text)
				m.refreshHistory()
				m.input.Reset()
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m ChatModel) View() string {
	header := lipgloss.NewStyle().Bold(true).Foreground(theme.Colors.Accent).Render("Chat")
	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.history.View(),
		theme.Styles.InputBorder.Render(m.input.View()),
		theme.Styles.Help.Render("Enter to send · Tab to switch panes · Ctrl+C to quit"),
	)
	return theme.Styles.Panel.Render(content)
}

func (m *ChatModel) refreshHistory() {
	m.history.SetContent(strings.Join(m.messages, "\n"))
	m.history.GotoBottom()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
