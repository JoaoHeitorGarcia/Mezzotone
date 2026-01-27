package components

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type LoadingScreen struct {
	Spinner  spinner.Model
	Quitting bool
	Err      error
}

func NewLoadingScreen(animation []string, fps int) LoadingScreen {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: animation,
		FPS:    time.Second / time.Duration(fps),
	}
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return LoadingScreen{
		Spinner: s,
	}
}
