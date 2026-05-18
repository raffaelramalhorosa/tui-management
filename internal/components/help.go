package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func RenderHelp(t theme.Theme) string {
	keyStyle := lipgloss.NewStyle().Foreground(t.Primary).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(t.Text)
	dimStyle := lipgloss.NewStyle().Foreground(t.Subtext)

	title := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render("Keyboard Shortcuts")
	sep := dimStyle.Render(strings.Repeat("─", 36))

	type binding struct{ key, desc string }

	groups := []struct {
		name     string
		bindings []binding
	}{
		{"Navigation", []binding{
			{"Tab / 1 / 2 / 3", "Switch tab"},
			{"↑ / ↓  or  k / j", "Navigate processes"},
		}},
		{"Sorting (Processes tab)", []binding{
			{"c", "Sort by CPU"},
			{"m", "Sort by Memory"},
			{"p", "Sort by PID"},
		}},
		{"Actions", []binding{
			{"K", "Kill selected process"},
			{"Space", "Pause / Resume updates"},
			{"t", "Cycle theme"},
			{"?", "Toggle this help"},
			{"q / Ctrl+C", "Quit"},
		}},
	}

	lines := []string{title, sep}
	for _, g := range groups {
		lines = append(lines, dimStyle.Render(g.name))
		for _, b := range g.bindings {
			lines = append(lines,
				"  "+keyStyle.Render(fmt.Sprintf("%-20s", b.key))+" "+descStyle.Render(b.desc),
			)
		}
		lines = append(lines, "")
	}

	lines = append(lines, dimStyle.Render("Press ? to close"))

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(1, 3).
		Render(strings.Join(lines, "\n"))

	return box
}
