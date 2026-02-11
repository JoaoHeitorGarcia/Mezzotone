package services

import (
	"strconv"
	"strings"

	"codeberg.org/JoaoGarcia/Mezzotone/internal/ui"
	"github.com/charmbracelet/x/ansi"
)

func TruncateLinesANSI(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = ansi.Truncate(lines[i], maxWidth, "…")
	}
	return strings.Join(lines, "\n")
}

func NormalizeRenderOptionsForService(settingsValues []ui.SettingItem) RenderOptions {
	var textSize int
	var fontAspect, edgeThreshold float64
	var directionalRender, reverseChars, highContrast bool
	var runeMode string

	for _, item := range settingsValues {
		switch item.Key {
		case "textSize":
			textSize, _ = strconv.Atoi(item.Value)

		case "fontAspect":
			fontAspect, _ = strconv.ParseFloat(item.Value, 2)

		case "edgeThreshold":
			edgeThreshold, _ = strconv.ParseFloat(item.Value, 2)

		case "directionalRender":
			directionalRender, _ = strconv.ParseBool(item.Value)

		case "reverseChars":
			reverseChars, _ = strconv.ParseBool(item.Value)

		case "highContrast":
			highContrast, _ = strconv.ParseBool(item.Value)

		case "runeMode":
			runeMode = item.Value
		}
	}
	options, err := NewRenderOptions(textSize, fontAspect, directionalRender, edgeThreshold, reverseChars, highContrast, runeMode)
	if err != nil {
		//TODO render Error and go back to renderOptionsMenu
	}
	return options
}
