package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/system"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func RenderMemory(mem system.MemInfo, width int, t theme.Theme) string {
	title := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render("Memory")

	barWidth := width - 12
	if barWidth < 4 {
		barWidth = 4
	}

	usedGB := float64(mem.Used) / 1024 / 1024 / 1024
	totalGB := float64(mem.Total) / 1024 / 1024 / 1024

	bar := renderBar("RAM   ", mem.UsedPercent, barWidth, t)
	detail := lipgloss.NewStyle().Foreground(t.Subtext).Render(
		fmt.Sprintf("  %.1f GB / %.1f GB", usedGB, totalGB),
	)

	return strings.Join([]string{title, bar, detail}, "\n")
}
