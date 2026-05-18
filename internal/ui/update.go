package ui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raffaelramalhorosa/tui-management/internal/components"
	"github.com/raffaelramalhorosa/tui-management/internal/theme"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(collectCmd(), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		if !m.paused {
			return m, tea.Batch(collectCmd(), tickCmd())
		}
		return m, tickCmd()

	case metricsMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.metrics = msg.metrics
			m.cpuHistory = appendHistory(m.cpuHistory, msg.metrics.CPU)
		}

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// kill confirm intercepts all input
	if m.confirmKill {
		switch msg.String() {
		case "y", "Y":
			pid := components.GetSortedPID(m.metrics.Processes, m.selectedRow, m.sortBy, m.sortDesc)
			if pid > 0 {
				killProcess(pid)
			}
			m.confirmKill = false
		case "n", "N", "esc":
			m.confirmKill = false
		}
		return m, nil
	}

	// help intercepts all input
	if m.showHelp {
		m.showHelp = false
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab":
		m.activeTab = (m.activeTab + 1) % 3
		m.selectedRow = 0
		m.scrollOff = 0

	case "1":
		m.activeTab = 0
	case "2":
		m.activeTab = 1
		m.selectedRow = 0
		m.scrollOff = 0
	case "3":
		m.activeTab = 2

	case "up", "k":
		if m.activeTab == 1 {
			m.selectedRow--
			if m.selectedRow < 0 {
				m.selectedRow = 0
			}
			m.adjustScroll()
		}

	case "down", "j":
		if m.activeTab == 1 {
			maxRow := len(m.metrics.Processes) - 1
			m.selectedRow++
			if m.selectedRow > maxRow {
				m.selectedRow = maxRow
			}
			m.adjustScroll()
		}

	case "c":
		if m.sortBy == "cpu" {
			m.sortDesc = !m.sortDesc
		} else {
			m.sortBy = "cpu"
			m.sortDesc = true
		}
		m.selectedRow = 0
		m.scrollOff = 0

	case "m":
		if m.sortBy == "mem" {
			m.sortDesc = !m.sortDesc
		} else {
			m.sortBy = "mem"
			m.sortDesc = true
		}
		m.selectedRow = 0
		m.scrollOff = 0

	case "p":
		if m.sortBy == "pid" {
			m.sortDesc = !m.sortDesc
		} else {
			m.sortBy = "pid"
			m.sortDesc = true
		}
		m.selectedRow = 0
		m.scrollOff = 0

	case "K":
		if m.activeTab == 1 && len(m.metrics.Processes) > 0 {
			m.confirmKill = true
		}

	case " ":
		m.paused = !m.paused

	case "?":
		m.showHelp = true

	case "t":
		m.activeTheme = (m.activeTheme + 1) % len(theme.All)
	}

	return m, nil
}

func (m *Model) adjustScroll() {
	visible := m.height - 10
	if visible < 1 {
		visible = 1
	}
	if m.selectedRow < m.scrollOff {
		m.scrollOff = m.selectedRow
	}
	if m.selectedRow >= m.scrollOff+visible {
		m.scrollOff = m.selectedRow - visible + 1
	}
}

func appendHistory(history [][]float64, percents []float64) [][]float64 {
	if len(history) == 0 && len(percents) > 0 {
		history = make([][]float64, len(percents))
	}
	for i, pct := range percents {
		if i >= len(history) {
			break
		}
		history[i] = append(history[i], pct)
		if len(history[i]) > 60 {
			history[i] = history[i][1:]
		}
	}
	return history
}

func killProcess(pid int32) {
	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return
	}
	proc.Kill()
}

