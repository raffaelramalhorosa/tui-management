package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func panelStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(0, 1)
}

func titleStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true)
}

func labelStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Subtext)
}

func valueStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Text)
}

func headerStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Success).
		Bold(true).
		Padding(0, 1)
}

func activeTabStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Text).
		Background(t.Primary).
		Padding(0, 2).
		Bold(true)
}

func inactiveTabStyle(t theme.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.Subtext).
		Padding(0, 2)
}
