package screens

import (
	"codeberg.org/JoaoGarcia/Mezzotone/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
)

type ImagePreview struct {
	loadingAnimation components.LoadingScreen
}

func NewImagePreview() ImagePreview {
	//experiment with animation - change later
	loadingAnimation := components.NewLoadingScreen(
		[]string{
			"      \n      \n      \n      \n      \n      ",
			"o     \n      \n      \n      \n      \n      ",
			"oo    \no     \n      \n      \n      \n      ",
			"oooo  \noo    \no     \n      \n      \n      ",
			"ooooo \nooo   \noo    \no     \n      \n      ",
			"oooooo\noooo  \nooo   \noo    \no     \n      ",
			" ooooo\nooooo \noooo  \nooo   \noo    \no     ",
			"  oooo\noooooo\nooooo \noooo  \nooo   \noo    ",
			"   ooo\n ooooo\noooooo\nooooo \noooo  \nooo   ",
			"    oo\n  oooo\n ooooo\noooooo\nooooo \noooo  ",
			"     o\n   ooo\n  oooo\n ooooo\noooooo\nooooo ",
			"      \n    oo\n   ooo\n  oooo\n ooooo\noooooo",
			"      \n     o\n    oo\n   ooo\n  oooo\n ooooo",
			"      \n      \n      \n    oo\n   ooo\n  oooo",
			"      \n      \n      \n      \n    oo\n   ooo",
			"      \n      \n      \n      \n      \n    oo",
			"      \n      \n      \n      \n      \n     o",
			"      \n      \n      \n      \n      \n      ",
		},
		10,
	)

	return ImagePreview{
		loadingAnimation,
	}
}

func (m ImagePreview) Init() tea.Cmd {
	return m.loadingAnimation.Spinner.Tick
}

func (m ImagePreview) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmd tea.Cmd
	m.loadingAnimation.Spinner, cmd = m.loadingAnimation.Spinner.Update(msg)
	return m, cmd
}

func (m ImagePreview) View() string {
	return m.loadingAnimation.Spinner.View()
}
