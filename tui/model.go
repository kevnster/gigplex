package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevnster/gigplex"
)

type tickMsg struct{}
type agentStepMsg string
type agentDoneMsg string

type Model struct {
	workers        []gigplex.WorkerInfo
	stats          gigplex.Stats
	recentJobs     []gigplex.Job
	selectedWorker int
	width          int
	height         int
	viewport       viewport.Model
	investigating  bool
	agentSteps     []string
	agentOutput    string
	backend        gigplex.Backend
	err            error
}

func New(backend gigplex.Backend) Model {
	return Model{
		backend:        backend,
		selectedWorker: 0,
		agentSteps:     []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}