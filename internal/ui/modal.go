package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/raffaelramalhorosa/tui-management/internal/system"
)

type tickMsg time.Time

type metricsMsg struct {
	metrics system.Metrics
	err     error
}

type Model struct {
	metrics    system.Metrics
	cpuHistory [][]float64 // [core][time], ring buffer de 60 leituras
	err        error
	width      int
	height     int

	activeTab   int    // 0=Overview 1=Processes 2=Network
	selectedRow int    // linha selecionada na tabela de processos
	scrollOff   int    // scroll offset da tabela
	sortBy      string // "cpu" | "mem" | "pid"
	sortDesc    bool

	paused      bool
	showHelp    bool
	confirmKill bool
	activeTheme int
}

func NewModel() Model {
	return Model{
		sortBy:   "cpu",
		sortDesc: true,
	}
}

func collectCmd() tea.Cmd {
	return func() tea.Msg {
		m, err := system.Collect()
		return metricsMsg{metrics: m, err: err}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
