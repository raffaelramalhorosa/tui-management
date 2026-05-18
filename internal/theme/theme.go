package theme

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name    string
	Primary lipgloss.Color
	Success lipgloss.Color
	Warning lipgloss.Color
	Danger  lipgloss.Color
	Text    lipgloss.Color
	Subtext lipgloss.Color
}

var All = []Theme{
	{
		Name:    "Default",
		Primary: lipgloss.Color("#5f87ff"),
		Success: lipgloss.Color("#00ff87"),
		Warning: lipgloss.Color("#ffff00"),
		Danger:  lipgloss.Color("#ff5f5f"),
		Text:    lipgloss.Color("#ffffff"),
		Subtext: lipgloss.Color("#626262"),
	},
	{
		Name:    "Catppuccin",
		Primary: lipgloss.Color("#cba6f7"),
		Success: lipgloss.Color("#a6e3a1"),
		Warning: lipgloss.Color("#f9e2af"),
		Danger:  lipgloss.Color("#f38ba8"),
		Text:    lipgloss.Color("#cdd6f4"),
		Subtext: lipgloss.Color("#6c7086"),
	},
	{
		Name:    "Dracula",
		Primary: lipgloss.Color("#bd93f9"),
		Success: lipgloss.Color("#50fa7b"),
		Warning: lipgloss.Color("#f1fa8c"),
		Danger:  lipgloss.Color("#ff5555"),
		Text:    lipgloss.Color("#f8f8f2"),
		Subtext: lipgloss.Color("#6272a4"),
	},
}

func (t Theme) BarColor(pct float64) lipgloss.Color {
	switch {
	case pct >= 80:
		return t.Danger
	case pct >= 60:
		return t.Warning
	default:
		return t.Success
	}
}
