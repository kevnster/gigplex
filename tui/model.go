package tui

import (
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
	g              *gigplex.Gigplex
	err            error
}

func New(backend gigplex.Backend, g *gigplex.Gigplex) Model {
	return Model{
		backend:    backend,
		g:          g,
		agentSteps: []string{},
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}
