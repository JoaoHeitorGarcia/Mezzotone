package services

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

const asciiRampDarkToBrightStr = "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\\|()1{}[]?-_+~<>i!lI;:,^`. "
const unicodeRampDarkToBrightStr = "█▉▊▋▌▍▎▏▓▒░■□@&%$#*+=-~:;!,\".^`' "

const asciiRampBrightToDarkStr = " .`^,:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
const unicodeRampBrightToDarkStr = " '`^.\",!;:~-=+*#$%&@□■░▒▓▏▎▍▌▋▊▉█"

/*
TODO if unicode is true i could make a bunch of different ramps for the user to choose from.
	Example:
	█▇▆▅▄▃▂▁ ▁▂▃▄▅▆▇█
	█▓▒░ ░▒▓█
	⣿⣷⣧⣇⣆⣄⣀ ⣀⣄⣆⣇⣧⣷⣿
	●∙•·  ·•∙●
*/

func ConvertImageToString(filePath string) ([][]rune, error) {
	var outputChars [][]rune

	f, err := os.Open(filePath)
	if err != nil {
		return outputChars, err
	}
	defer func() { _ = f.Close() }()

	_ = Logger().Info(fmt.Sprintf("Successfully Loaded: %s", filePath))

	inputImg, format, err := image.Decode(f)
	if err != nil {
		return outputChars, err
	}
	_ = Logger().Info(fmt.Sprintf("format: %s", format))

	textSizeAny, ok := Shared().Get("textSize")
	if !ok || textSizeAny == nil {
		return outputChars, errors.New("textSize is nil")
	}
	textSize, ok := textSizeAny.(int)
	if !ok || textSize <= 0 {
		textSize = 8
	}

	fontAspectAny, ok := Shared().Get("fontAspect")
	if !ok || fontAspectAny == nil {
		return outputChars, errors.New("fontAspect is nil")
	}
	fontAspect, ok := fontAspectAny.(float64)
	if !ok || fontAspect <= 0 {
		fontAspect = 2
	}

	highContrastAny, ok := Shared().Get("highContrast")
	if !ok || highContrastAny == nil {
		return outputChars, errors.New("highContrast is nil")
	}
	highContrast := highContrastAny.(bool)

	cols, rows := gridFromTextSize(inputImg, textSize, fontAspect)
	outputChars = make([][]rune, rows)
	for r := 0; r < rows; r++ {
		outputChars[r] = make([]rune, cols)
	}

	lumaGrid, err := buildLumianceGrid(inputImg, cols, rows, highContrast)
	if err != nil {
		return outputChars, err
	}
	_ = Logger().Info(fmt.Sprintf("Successfully Build LumaGrid for %s", filePath))

	directionalRenderAny, ok := Shared().Get("directionalRender")
	if !ok || directionalRenderAny == nil {
		return outputChars, errors.New("directionalRender is nil")
	}
	directionalRender := directionalRenderAny.(bool)
	if directionalRender {
		//TODO apply Sobel filter
	}

	_ = Logger().Info(fmt.Sprintf("Beginning image conversion"))
	useUnicodeAny, ok := Shared().Get("useUnicode")
	if !ok || useUnicodeAny == nil {
		return outputChars, errors.New("useUnicode is nil")
	}
	useUnicode := useUnicodeAny.(bool)

	reverseCharsAny, ok := Shared().Get("reverseChars")
	if !ok || reverseCharsAny == nil {
		return outputChars, errors.New("reverseChars is nil")
	}
	reverseChars := reverseCharsAny.(bool)

	for i := 0; i < len(lumaGrid); i++ {
		for j := 0; j < len(lumaGrid[i]); j++ {
			outputChars[i][j] = getCharForLuminanceValue(lumaGrid[i][j], useUnicode, reverseChars)
		}
	}
	_ = Logger().Info(fmt.Sprintf("Finished image conversion"))

	return outputChars, nil
}

func gridFromTextSize(img image.Image, textSize int, fontAspect float64) (cols, rows int) {
	b := img.Bounds()
	imgW, imgH := b.Dx(), b.Dy()

	charW := textSize
	charH := int(float64(textSize) * fontAspect)
	if charW <= 0 {
		charW = 8
	}
	if charH <= 0 {
		charH = 16
	}

	cols = (imgW + charW - 1) / charW
	rows = (imgH + charH - 1) / charH

	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	return cols, rows
}

func buildLumianceGrid(inputImg image.Image, cols, rows int, highContrast bool) ([][]float64, error) {

	imgBounds := inputImg.Bounds()
	imgWidth, imgHeight := imgBounds.Dx(), imgBounds.Dy()

	cellWidth := imgWidth / cols
	cellHeight := imgHeight / rows
	if cellWidth <= 0 {
		cellWidth = 8
	}
	if cellHeight <= 0 {
		cellHeight = 16
	}

	grid := make([][]float64, rows)
	for gridRow := 0; gridRow < rows; gridRow++ {
		grid[gridRow] = make([]float64, cols)
	}

	for gridRow := 0; gridRow < rows; gridRow++ {
		cellRowPixelStartY := gridRow * cellHeight
		cellRowPixelEndY := cellRowPixelStartY + cellHeight
		if cellRowPixelStartY >= imgHeight {
			cellRowPixelStartY = imgHeight
		}
		if cellRowPixelEndY > imgHeight {
			cellRowPixelEndY = imgHeight
		}

		for gridCol := 0; gridCol < cols; gridCol++ {
			cellColPixelStartX := gridCol * cellWidth
			cellColPixelEndX := cellColPixelStartX + cellWidth
			if cellColPixelStartX >= imgWidth {
				cellColPixelStartX = imgWidth
			}
			if cellColPixelEndX > imgWidth {
				cellColPixelEndX = imgWidth
			}

			//fallback (should not happen)
			if cellColPixelEndX <= cellColPixelStartX || cellRowPixelEndY <= cellRowPixelStartY {
				grid[gridRow][gridCol] = 0
				continue
			}

			var lumaSum float64
			var sampleCount float64

			for cellCol := cellRowPixelStartY; cellCol < cellRowPixelEndY; cellCol++ {
				for CellRow := cellColPixelStartX; CellRow < cellColPixelEndX; CellRow++ {
					c := color.NRGBAModel.Convert(
						inputImg.At(imgBounds.Min.X+CellRow, imgBounds.Min.Y+cellCol),
					).(color.NRGBA)

					if c.A < 10 {
						continue
					}

					pixelLuminance := calculateLuminance(c.R, c.G, c.B) // expects 0..1
					lumaSum += pixelLuminance
					sampleCount++
				}
			}

			var cellLuma float64
			if sampleCount == 0 {
				cellLuma = 0
			} else {
				cellLuma = lumaSum / sampleCount
			}

			if highContrast {
				cellLuma = applyContrast(cellLuma, 1.7)
			}

			grid[gridRow][gridCol] = clamp01(cellLuma)
		}
	}

	return grid, nil
}

func calculateLuminance(red uint8, green uint8, blue uint8) float64 {
	luminance := 0.2126*float64(red) + 0.7152*float64(green) + 0.0722*float64(blue)
	//normalize to 0...1
	return luminance / 255.0
}

func getCharForLuminanceValue(luminance float64, useUnicode bool, reverseChars bool) rune {
	var ramp []rune
	if useUnicode {
		if reverseChars {
			ramp = []rune(unicodeRampBrightToDarkStr)
		} else {
			ramp = []rune(unicodeRampDarkToBrightStr)
		}
	} else {
		if reverseChars {
			ramp = []rune(asciiRampBrightToDarkStr)
		} else {
			ramp = []rune(asciiRampDarkToBrightStr)
		}
	}

	luminance = clamp01(luminance)
	index := int(luminance * float64(len(ramp)-1))

	_ = Logger().Info(
		fmt.Sprintf(
			"brightness: %.2f | character: %s | character index: %d",
			luminance, string(ramp[index]), index,
		),
	)

	return ramp[index]
}

func clamp01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func applyContrast(l float64, factor float64) float64 {
	return clamp01((l-0.5)*factor + 0.5)
}
