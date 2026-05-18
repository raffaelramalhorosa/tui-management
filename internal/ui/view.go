package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/components"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	t := theme.All[m.activeTheme]

	if m.showHelp {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			components.RenderHelp(t),
		)
	}

	if m.confirmKill {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			renderKillConfirm(m, t),
		)
	}

	header := renderHeader(m, t)
	tabBar := renderTabBar(m, t)
	footer := renderFooter(m, t)

	headerH := lipgloss.Height(header)
	tabH := lipgloss.Height(tabBar)
	footerH := lipgloss.Height(footer)
	contentH := m.height - headerH - tabH - footerH - 2

	var content string
	switch m.activeTab {
	case 0:
		content = renderOverview(m, t, contentH)
	case 1:
		content = renderProcessesTab(m, t, contentH)
	case 2:
		content = renderNetworkTab(m, t, contentH)
	}

	if m.err != nil {
		content = lipgloss.NewStyle().Foreground(t.Danger).Render(
			fmt.Sprintf("  Error collecting metrics: %v", m.err),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabBar,
		content,
		footer,
	)
}

func renderHeader(m Model, t theme.Theme) string {
	clock := lipgloss.NewStyle().Foreground(t.Subtext).Render(time.Now().Format("15:04:05"))
	pause := ""
	if m.paused {
		pause = lipgloss.NewStyle().Foreground(t.Warning).Bold(true).Render(" ⏸ PAUSED")
	}
	themeName := lipgloss.NewStyle().Foreground(t.Primary).Render("[t] " + t.Name)

	title := headerStyle(t).Render("  GoMonitor — System Monitor")

	right := lipgloss.JoinHorizontal(lipgloss.Center, pause, "  ", clock, "  ", themeName, "  ")
	rightW := lipgloss.Width(right)

	titleW := m.width - rightW
	if titleW < 0 {
		titleW = 0
	}

	titlePadded := lipgloss.NewStyle().Width(titleW).Render(title)
	return lipgloss.JoinHorizontal(lipgloss.Center, titlePadded, right)
}

func renderTabBar(m Model, t theme.Theme) string {
	tabs := []string{"Overview", "Processes", "Network"}
	rendered := make([]string, len(tabs))
	for i, name := range tabs {
		if i == m.activeTab {
			rendered[i] = activeTabStyle(t).Render(name)
		} else {
			rendered[i] = inactiveTabStyle(t).Render(name)
		}
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Left, rendered...)
	return "  " + bar
}

func renderFooter(m Model, t theme.Theme) string {
	base := "[q] quit  [?] help  [space] pause  [t] theme  [tab] switch"
	extra := ""
	switch m.activeTab {
	case 1:
		extra = "  [↑↓] navigate  [c/m/p] sort  [K] kill"
	}
	return labelStyle(t).Render("  " + base + extra)
}

func renderOverview(m Model, t theme.Theme, height int) string {
	halfW := m.width/2 - 2
	innerW := halfW - 4

	cpu := components.RenderCPU(m.metrics.CPU, m.cpuHistory, innerW, t)
	mem := components.RenderMemory(m.metrics.Memory, innerW, t)
	dsk := components.RenderDisk(m.metrics.Disks, innerW, t)

	right := lipgloss.JoinVertical(lipgloss.Left, mem, "", dsk)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		panelStyle(t).Width(halfW).Render(cpu),
		panelStyle(t).Width(halfW).Render(right),
	)

	_ = height
	return topRow
}

func renderProcessesTab(m Model, t theme.Theme, height int) string {
	innerW := m.width - 6
	maxVisible := height - 5
	if maxVisible < 1 {
		maxVisible = 1
	}

	content := components.RenderProcesses(
		m.metrics.Processes,
		m.selectedRow,
		m.scrollOff,
		maxVisible,
		m.sortBy,
		m.sortDesc,
		innerW,
		t,
	)

	return panelStyle(t).Width(m.width - 4).Render(content)
}

func renderNetworkTab(m Model, t theme.Theme, height int) string {
	innerW := m.width - 6
	content := components.RenderNetwork(m.metrics.Network, innerW, t)
	_ = height
	return panelStyle(t).Width(m.width - 4).Render(content)
}

func renderKillConfirm(m Model, t theme.Theme) string {
	name := components.GetSortedName(m.metrics.Processes, m.selectedRow, m.sortBy, m.sortDesc)
	pid := components.GetSortedPID(m.metrics.Processes, m.selectedRow, m.sortBy, m.sortDesc)

	warning := lipgloss.NewStyle().Foreground(t.Danger).Bold(true).Render("⚠  Kill Process")
	proc := lipgloss.NewStyle().Foreground(t.Text).Render(fmt.Sprintf("%s (PID %d)", name, pid))
	confirm := lipgloss.NewStyle().Foreground(t.Subtext).Render("Confirm? ") +
		lipgloss.NewStyle().Foreground(t.Success).Bold(true).Render("[y] yes") +
		lipgloss.NewStyle().Foreground(t.Subtext).Render("  ") +
		lipgloss.NewStyle().Foreground(t.Danger).Bold(true).Render("[n] cancel")

	body := lipgloss.JoinVertical(lipgloss.Center, warning, "", proc, "", confirm)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Danger).
		Padding(1, 4).
		Render(body)
}
