package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

var sparkChars = []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

func sparkline(values []float64, width int) string {
	if len(values) == 0 || width <= 0 {
		return strings.Repeat(" ", width)
	}
	if len(values) > width {
		values = values[len(values)-width:]
	}
	result := make([]rune, len(values))
	for i, v := range values {
		idx := int(v / 100 * float64(len(sparkChars)-1))
		if idx >= len(sparkChars) {
			idx = len(sparkChars) - 1
		}
		if idx < 0 {
			idx = 0
		}
		result[i] = sparkChars[idx]
	}
	pad := width - len(result)
	return strings.Repeat(" ", pad) + string(result)
}

func renderBar(label string, pct float64, barWidth int, t theme.Theme) string {
	if barWidth < 1 {
		barWidth = 1
	}
	filled := int(float64(barWidth) * pct / 100)
	if filled > barWidth {
		filled = barWidth
	}
	bar := lipgloss.NewStyle().Foreground(t.BarColor(pct)).Render(strings.Repeat("█", filled))
	empty := lipgloss.NewStyle().Foreground(t.Subtext).Render(strings.Repeat("░", barWidth-filled))
	pctStr := lipgloss.NewStyle().Foreground(t.Text).Render(fmt.Sprintf("%5.1f%%", pct))
	lbl := lipgloss.NewStyle().Foreground(t.Subtext).Render(label)
	return fmt.Sprintf("%s %s%s %s", lbl, bar, empty, pctStr)
}

func RenderCPU(percents []float64, history [][]float64, width int, t theme.Theme) string {
	title := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render("CPU")

	if len(percents) == 0 {
		return title + "\n  No data"
	}

	sparkWidth := 20
	pctWidth := 8
	labelWidth := 8
	barWidth := width - labelWidth - pctWidth - sparkWidth - 4
	if barWidth < 4 {
		barWidth = 4
	}

	lines := []string{title}
	for i, pct := range percents {
		label := fmt.Sprintf("Core%-2d", i)
		bar := renderBar(label, pct, barWidth, t)

		var spark string
		if i < len(history) {
			spark = lipgloss.NewStyle().Foreground(t.BarColor(pct)).Render(sparkline(history[i], sparkWidth))
		} else {
			spark = strings.Repeat(" ", sparkWidth)
		}

		lines = append(lines, bar+" "+spark)
	}
	return strings.Join(lines, "\n")
}
