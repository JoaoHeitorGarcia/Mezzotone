package main

import (
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
		_ = services.Logger().Error(err.Error())
		os.Exit(1)
	}
}
