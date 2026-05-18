package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/system"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func RenderDisk(disks []system.DiskInfo, width int, t theme.Theme) string {
	title := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render("Disk")

	if len(disks) == 0 {
		return title + "\n  No data"
	}

	barWidth := width - 16
	if barWidth < 4 {
		barWidth = 4
	}

	lines := []string{title}
	for _, d := range disks {
		label := d.Path
		if len(label) > 6 {
			label = label[:5] + "…"
		}
		for len(label) < 6 {
			label += " "
		}

		usedGB := float64(d.Used) / 1024 / 1024 / 1024
		totalGB := float64(d.Total) / 1024 / 1024 / 1024

		bar := renderBar(label, d.Percent, barWidth, t)
		detail := lipgloss.NewStyle().Foreground(t.Subtext).Render(
			fmt.Sprintf("  %.1f GB / %.1f GB", usedGB, totalGB),
		)
		lines = append(lines, bar, detail)
	}

	return strings.Join(lines, "\n")
}
