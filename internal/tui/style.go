package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	borderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63"))

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("24")).
			Bold(true)

	tabsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")).
			Bold(true)

	metaOnlineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120"))

	metaOfflineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("203"))

	panelTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("230")).
			Background(lipgloss.Color("60")).
			Bold(true)

	statusOKStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Bold(true)

	statusWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220")).
			Bold(true)

	statusErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Bold(true)

	statusInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")).
			Bold(true)

	selectedLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("230")).
				Background(lipgloss.Color("57")).
				Bold(true)

	backLineStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("225")).
			Background(lipgloss.Color("25")).
			Bold(true)

	keysStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	verifiedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("120")).
			Bold(true)
)

func paintLayout(layout string) string {
	if layout == "" {
		return layout
	}

	lines := strings.Split(layout, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "┌"), strings.HasPrefix(line, "├"), strings.HasPrefix(line, "└"):
			lines[i] = borderStyle.Render(line)
		case strings.Contains(line, "ST-8508 Reader TUI"):
			lines[i] = headerStyle.Render(line)
		case strings.Contains(line, "▣ ") || strings.Contains(line, "□ "):
			lines[i] = tabsStyle.Render(line)
		case strings.Contains(line, "Reader ONLINE"):
			lines[i] = metaOnlineStyle.Render(line)
		case strings.Contains(line, "Reader OFFLINE"):
			lines[i] = metaOfflineStyle.Render(line)
		case strings.Contains(line, "[ OK  ]"):
			lines[i] = statusOKStyle.Render(line)
		case strings.Contains(line, "[WARN ]"):
			lines[i] = statusWarnStyle.Render(line)
		case strings.Contains(line, "[ERROR]"):
			lines[i] = statusErrStyle.Render(line)
		case strings.Contains(line, "[INFO ]"):
			lines[i] = statusInfoStyle.Render(line)
		case strings.Contains(line, "│ ▶ "):
			lines[i] = selectedLineStyle.Render(line)
		case strings.Contains(line, "Back to Home"):
			lines[i] = backLineStyle.Render(line)
		case strings.Contains(line, "[VERIFIED]"):
			lines[i] = verifiedStyle.Render(line)
		case isPanelTitleLine(line):
			lines[i] = panelTitleStyle.Render(line)
		case strings.Contains(line, "Keys:"):
			lines[i] = keysStyle.Render(line)
		}
	}

	return strings.Join(lines, "\n")
}

func isPanelTitleLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "│ [") || !strings.HasSuffix(trimmed, "] │") {
		return false
	}
	return !strings.Contains(trimmed, "[ OK  ]") &&
		!strings.Contains(trimmed, "[WARN ]") &&
		!strings.Contains(trimmed, "[ERROR]") &&
		!strings.Contains(trimmed, "[INFO ]")
}
