package tui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		ctx := context.Background()

		workers, err := m.backend.Workers(ctx)
		if err == nil {
			m.workers = workers
		}

		stats, err := m.backend.Stats(ctx)
		if err == nil {
			m.stats = stats
		}

		jobs, err := m.backend.RecentJobs(ctx, 10)
		if err == nil {
			m.recentJobs = jobs
		}

		return m, tick()
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if len(m.workers) > 0 {
				m.selectedWorker = (m.selectedWorker + 1) % len(m.workers)
			}
		case "k":
			if len(m.workers) > 0 {
				_ = m.backend.KillWorker(context.Background(), m.workers[m.selectedWorker].ID)
			}
		case "r":
			_ = m.backend.RetryFailed(context.Background())
		}
	}

	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
