package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kevnster/gigplex"
)

var (
	colorGreen  = lipgloss.Color("#10b981")
	colorRed    = lipgloss.Color("#ef4444")
	colorYellow = lipgloss.Color("#f59e0b")
	colorGray   = lipgloss.Color("#6b7280")
	colorWhite  = lipgloss.Color("#f9fafb")

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)

	styleMuted = lipgloss.NewStyle().
			Foreground(colorGray)

	styleBorderGreen = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorGreen).
				Padding(0, 1)

	styleBorderRed = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorRed).
			Padding(0, 1)

	styleBorderYellow = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorYellow).
				Padding(0, 1)

	styleBorderGray = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGray).
			Padding(0, 1)
)

func (m Model) View() string {
	if m.width == 0 {
		return "loading gigplex.ai...\n"
	}

	sections := []string{
		renderHeader(m),
		renderWorkers(m),
		renderStats(m),
		renderJobLog(m),
		renderHelp(m),
	}

	return strings.Join(sections, "\n")
}

func renderHeader(m Model) string {
	left := styleHeader.Render("👺 gigplex.ai") +
    	styleMuted.Render(fmt.Sprintf("  v%s  ·  memory  ·  %d workers", "0.1.0", len(m.workers)))

	return left + "\n" + strings.Repeat("━", m.width) + "\n"
}

func renderWorkers(m Model) string {
	if len(m.workers) == 0 {
		return styleMuted.Render("  no workers connected\n")
	}

	workerWidth := (m.width / len(m.workers)) - 2
	cards := make([]string, len(m.workers))

	for i, w := range m.workers {
		cards[i] = renderWorkerCard(w, i == m.selectedWorker, workerWidth)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, cards...) + "\n"
}

func renderWorkerCard(w gigplex.WorkerInfo, selected bool, width int) string {
	var status string
	var style lipgloss.Style

	secsSince := time.Since(w.LastBeat).Seconds()

	switch {
	case secsSince > 10:
		status = "✕ DEAD"
		style = styleBorderRed
	case w.InFlight == 0:
		status = "● IDLE"
		style = styleBorderYellow
	default:
		status = "● ACTIVE"
		style = styleBorderGreen
	}

	if selected {
		style = style.BorderForeground(lipgloss.Color("#818cf8"))
	}

	leader := ""
	if w.IsLeader {
		leader = " LEADER"
	}

	content := fmt.Sprintf(
		"%s%s\n%s\ndone:    %d\nactive:  %d\nbeat:    %s ago",
		status,
		styleMuted.Render(leader),
		styleMuted.Render(w.ID),
		w.JobsDone,
		w.InFlight,
		formatDuration(time.Since(w.LastBeat)),
	)

	return style.Width(width).Render(content)
}

func renderStats(m Model) string {
	bar := fmt.Sprintf(
		"  pending: %s   processing: %s   failed: %s   done: %s",
		lipgloss.NewStyle().Foreground(colorYellow).Render(fmt.Sprintf("%d", m.stats.Pending)),
		lipgloss.NewStyle().Foreground(colorGreen).Render(fmt.Sprintf("%d", m.stats.Processing)),
		lipgloss.NewStyle().Foreground(colorRed).Render(fmt.Sprintf("%d", m.stats.Failed)),
		lipgloss.NewStyle().Foreground(colorWhite).Render(fmt.Sprintf("%d", m.stats.Done)),
	)

	return "\n" + bar + "\n" + strings.Repeat("─", m.width) + "\n"
}

func renderJobLog(m Model) string {
	header := styleHeader.Render("  RECENT JOBS\n")

	if len(m.recentJobs) == 0 {
		return header + styleMuted.Render("  no jobs yet\n")
	}

	lines := []string{header}
	for _, job := range m.recentJobs {
		lines = append(lines, renderJobLine(job))
	}

	return strings.Join(lines, "\n") + "\n"
}

func renderJobLine(job gigplex.Job) string {
	var icon string
	var iconStyle lipgloss.Style

	switch job.Status {
	case gigplex.StatusDone:
		icon = "✓"
		iconStyle = lipgloss.NewStyle().Foreground(colorGreen)
	case gigplex.StatusFailed:
		icon = "✕"
		iconStyle = lipgloss.NewStyle().Foreground(colorRed)
	case gigplex.StatusProcessing:
		icon = "→"
		iconStyle = lipgloss.NewStyle().Foreground(colorYellow)
	default:
		icon = "·"
		iconStyle = lipgloss.NewStyle().Foreground(colorGray)
	}

	timestamp := styleMuted.Render(job.CreatedAt.Format("15:04:05"))
	jobType := fmt.Sprintf("%-24s", job.Type)

	return fmt.Sprintf("  %s  %s  %s",
		timestamp,
		iconStyle.Render(icon),
		jobType,
	)
}

func renderHelp(m Model) string {
	help := styleMuted.Render(
		"\n  [p] pause  [k] kill  [r] retry  [tab] select  [a] investigate  [q] quit",
	)
	return strings.Repeat("─", m.width) + help
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%.0fs", d.Seconds())
}