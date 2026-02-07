package app

import (
	"fmt"
	"os"
	"strings"

	"codeberg.org/JoaoGarcia/Mezzotone/internal/services"
	"codeberg.org/JoaoGarcia/Mezzotone/internal/ui"
	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type MezzotoneModel struct {
	filePicker   filepicker.Model
	selectedFile string

	renderView    viewport.Model
	leftColumn    viewport.Model
	renderOptions ui.SettingsPanel

	style styleVariables

	currentActiveMenu int

	width  int
	height int

	err error
}

type styleVariables struct {
	windowMargin    int
	leftColumnWidth int
}

var optionsTableRowSize int

const (
	filePickerMenu = iota
	renderOptionsMenu
	renderViewText
)

func NewMezzotoneModel() MezzotoneModel {
	windowStyles := styleVariables{
		windowMargin: 2,
	}

	runeMode := []string{"ASCII", "UNICODE", "DOTS", "RECTANGLES", "BARS", "LOADING"}
	renderOptionItems := []ui.SettingItem{
		{Label: "Text Size", Key: "textSize", Type: ui.TypeInt, Value: "10"},
		{Label: "Font Aspect", Key: "fontAspect", Type: ui.TypeFloat, Value: "2.3"},
		{Label: "Directional Render", Key: "directionalRender", Type: ui.TypeBool, Value: "false"},
		{Label: "Edge Threshold Percentile", Key: "edgeThresholdPercentile", Type: ui.TypeFloat, Value: "0.6"},
		{Label: "Reverse Chars", Key: "reverseChars", Type: ui.TypeBool, Value: "true"},
		{Label: "High Contrast", Key: "highContrast", Type: ui.TypeBool, Value: "true"},
		{Label: "Rune Mode", Key: "runeMode", Type: ui.TypeEnum, Value: "ASCII", Enum: runeMode},
	}
	optionsTableRowSize = len(renderOptionItems)
	renderOptionsTable := ui.NewSettingsPanel("Render Options", renderOptionItems)
	renderOptionsTable.ClearActive()

	fp := filepicker.New()
	fp.AllowedTypes = []string{".png", ".jpg", ".jpeg", ".bmp", ".webp", ".tiff"}
	fp.CurrentDirectory, _ = os.UserHomeDir()
	fp.ShowPermissions = false
	fp.ShowSize = true
	fp.KeyMap = filepicker.KeyMap{
		Down:     key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		Up:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		PageUp:   key.NewBinding(key.WithKeys("K", "pgup"), key.WithHelp("pgup", "page up")),
		PageDown: key.NewBinding(key.WithKeys("J", "pgdown"), key.WithHelp("pgdown", "page down")),
		Back:     key.NewBinding(key.WithKeys("left", "backspace"), key.WithHelp("h", "back")),
		Open:     key.NewBinding(key.WithKeys("right", "enter"), key.WithHelp("l", "open")),
		Select:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
	}

	renderView := viewport.New(0, 0)
	renderView.SetContent("Placeholder Text")
	leftColumn := viewport.New(0, 0)

	return MezzotoneModel{
		filePicker:        fp,
		renderView:        renderView,
		style:             windowStyles,
		leftColumn:        leftColumn,
		renderOptions:     renderOptionsTable,
		currentActiveMenu: filePickerMenu,
	}
}

func (m MezzotoneModel) Init() tea.Cmd {
	return m.filePicker.Init()
}

func (m MezzotoneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		m.renderView.Height = m.height - m.style.windowMargin
		m.renderView.Width = m.width / 7 * 5

		m.style.leftColumnWidth = m.width / 7 * 2

		m.renderOptions.SetWidth(m.style.leftColumnWidth)
		m.renderOptions.SetHeight(optionsTableRowSize)

		computedFilePickerHeight := m.renderView.Height - (optionsTableRowSize + 4 /*renderOptions header and end*/) - m.style.windowMargin*2 - 2 //inputFile Title
		m.filePicker.SetHeight(computedFilePickerHeight)

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.currentActiveMenu == filePickerMenu {
				//TODO ask for confimation
				return m, tea.Quit
			} else if m.currentActiveMenu == renderOptionsMenu {
				if !m.renderOptions.Editing {
					m.currentActiveMenu--
					m.renderOptions.ClearActive()
				}
			}

		case "enter":
			if m.currentActiveMenu == filePickerMenu {
				m.renderOptions.SetActive(0)
			}
		}
	}

	if m.currentActiveMenu == filePickerMenu {
		m.filePicker, cmd = m.filePicker.Update(msg)
		cmds = append(cmds, cmd)
		if didSelect, path := m.filePicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path
			_ = services.Logger().Info(fmt.Sprintf("Selected File: %s", m.selectedFile))

			m.currentActiveMenu++
		}

		if didSelect, path := m.filePicker.DidSelectDisabledFile(msg); didSelect {
			//TODO maybe make a modal here with error ? or no modal but better error info
			m.renderView.SetContent("Selected file need to be an image.\nAllowed types: .png, .jpg, .jpeg, .bmp, .webp, .tiff")
			m.selectedFile = ""
			_ = services.Logger().Info(fmt.Sprintf("Tried Selecting File: %s", path))
		}
	} else if m.currentActiveMenu == renderOptionsMenu {
		m.renderOptions, cmd = m.renderOptions.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.renderView, cmd = m.renderView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MezzotoneModel) View() string {
	innerW := m.style.leftColumnWidth - 2

	//filePickerTitleStyle := lipgloss.NewStyle().SetString("Pick an image, gif or video to convert:")
	//filePickerTitleRender := truncateLinesANSI(filePickerTitleStyle.Render(), innerW)
	filePickerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		Width(m.style.leftColumnWidth)
	fpView := truncateLinesANSI(m.filePicker.View(), innerW)
	filePickerRender := filePickerStyle.Render( /*filePickerTitleRender + "\n\n" +*/ fpView)

	lefColumnRender := lipgloss.JoinVertical(lipgloss.Left, filePickerRender, m.renderOptions.View())

	renderViewStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder())
	renderViewRender := renderViewStyle.Render(m.renderView.View())

	return lipgloss.JoinHorizontal(lipgloss.Left, lefColumnRender, renderViewRender)
}

func truncateLinesANSI(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], maxWidth, "…")
	}
	return strings.Join(lines, "\n")
}
