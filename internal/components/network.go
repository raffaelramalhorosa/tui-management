package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/system"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func formatRate(bps float64) string {
	switch {
	case bps >= 1024*1024:
		return fmt.Sprintf("%.1f MB/s", bps/1024/1024)
	case bps >= 1024:
		return fmt.Sprintf("%.1f KB/s", bps/1024)
	default:
		return fmt.Sprintf("%.0f B/s", bps)
	}
}

func RenderNetwork(nets []system.NetInfo, width int, t theme.Theme) string {
	title := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render("Network")

	if len(nets) == 0 {
		return title + "\n  No interfaces"
	}

	header := lipgloss.NewStyle().Foreground(t.Subtext).Render(
		fmt.Sprintf("  %-20s %14s %14s", "Interface", "↑ Upload", "↓ Download"),
	)

	lines := []string{title, header}
	for _, n := range nets {
		name := n.Name
		if len(name) > 18 {
			name = name[:15] + "..."
		}

		sentStr := lipgloss.NewStyle().Foreground(t.Warning).Render(fmt.Sprintf("%14s", formatRate(n.SentRate)))
		recvStr := lipgloss.NewStyle().Foreground(t.Success).Render(fmt.Sprintf("%14s", formatRate(n.RecvRate)))

		row := fmt.Sprintf("  %-20s %s %s",
			lipgloss.NewStyle().Foreground(t.Text).Render(name),
			sentStr,
			recvStr,
		)
		lines = append(lines, row)
	}

	_ = width
	return strings.Join(lines, "\n")
}
