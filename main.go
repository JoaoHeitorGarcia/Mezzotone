package main

import (
	"fmt"
	"os"

	"codeberg.org/JoaoGarcia/Mezzotone/internal/app"
	"codeberg.org/JoaoGarcia/Mezzotone/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	err := services.InitLogger("logs.log")
	if err != nil {
		return
	}

	p := tea.NewProgram(app.NewRouterModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
