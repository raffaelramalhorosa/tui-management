package components

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/system"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func sortProcs(procs []system.ProcessInfo, sortBy string, sortDesc bool) []system.ProcessInfo {
	result := make([]system.ProcessInfo, len(procs))
	copy(result, procs)
	sort.Slice(result, func(i, j int) bool {
		switch sortBy {
		case "mem":
			if sortDesc {
				return result[i].Memory > result[j].Memory
			}
			return result[i].Memory < result[j].Memory
		case "pid":
			if sortDesc {
				return result[i].PID > result[j].PID
			}
			return result[i].PID < result[j].PID
		default: // "cpu"
			if sortDesc {
				return result[i].CPU > result[j].CPU
			}
			return result[i].CPU < result[j].CPU
		}
	})
	return result
}

func colorByValue(val float64, t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.BarColor(val))
}

func sortIndicator(col, sortBy string, sortDesc bool, t theme.Theme) string {
	if col != sortBy {
		return " "
	}
	if sortDesc {
		return lipgloss.NewStyle().Foreground(t.Primary).Render("▼")
	}
	return lipgloss.NewStyle().Foreground(t.Primary).Render("▲")
}

func RenderProcesses(procs []system.ProcessInfo, selected, scrollOff, maxVisible int, sortBy string, sortDesc bool, width int, t theme.Theme) string {
	title := lipgloss.NewStyle().Foreground(t.Primary).Bold(true).Render("Processes")

	cpuInd := sortIndicator("cpu", sortBy, sortDesc, t)
	memInd := sortIndicator("mem", sortBy, sortDesc, t)
	pidInd := sortIndicator("pid", sortBy, sortDesc, t)

	colHeader := lipgloss.NewStyle().Foreground(t.Subtext).Render(
		fmt.Sprintf("  %-7s %-30s %7s%s  %6s%s", "PID"+pidInd, "Name", "CPU%", cpuInd, "MEM%", memInd),
	)

	sorted := sortProcs(procs, sortBy, sortDesc)

	end := scrollOff + maxVisible
	if end > len(sorted) {
		end = len(sorted)
	}
	if scrollOff > len(sorted) {
		scrollOff = 0
	}

	lines := []string{title, colHeader}

	for i, p := range sorted[scrollOff:end] {
		absIdx := i + scrollOff
		name := p.Name
		if len(name) > 28 {
			name = name[:25] + "..."
		}

		cpuStr := colorByValue(float64(p.CPU), t).Render(fmt.Sprintf("%7.1f%%", p.CPU))
		memStr := colorByValue(float64(p.Memory), t).Render(fmt.Sprintf("%6.1f%%", p.Memory))
		row := fmt.Sprintf("  %-7d %-30s %s  %s", p.PID, name, cpuStr, memStr)

		if absIdx == selected {
			row = lipgloss.NewStyle().
				Background(t.Primary).
				Foreground(t.Text).
				Bold(true).
				Render(fmt.Sprintf("  %-7d %-30s %7.1f%%  %6.1f%%", p.PID, name, p.CPU, p.Memory))
		}

		lines = append(lines, row)
	}

	if len(sorted) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(t.Subtext).Render("  No processes"))
	}

	scrollInfo := ""
	if len(sorted) > maxVisible {
		scrollInfo = lipgloss.NewStyle().Foreground(t.Subtext).Render(
			fmt.Sprintf("  ↑↓ scroll  %d-%d of %d", scrollOff+1, end, len(sorted)),
		)
		lines = append(lines, scrollInfo)
	}

	_ = width
	return strings.Join(lines, "\n")
}

func GetSortedPID(procs []system.ProcessInfo, selected int, sortBy string, sortDesc bool) int32 {
	sorted := sortProcs(procs, sortBy, sortDesc)
	if selected < 0 || selected >= len(sorted) {
		return -1
	}
	return sorted[selected].PID
}

func GetSortedName(procs []system.ProcessInfo, selected int, sortBy string, sortDesc bool) string {
	sorted := sortProcs(procs, sortBy, sortDesc)
	if selected < 0 || selected >= len(sorted) {
		return ""
	}
	return sorted[selected].Name
}
