package services

import (
	"errors"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	xdraw "golang.org/x/image/draw"
)

//const asciiRamp := "$@B%8&WM#*oahkbdpqwmZO0QLCJUYXzcvunxrjft/\|()1{}[]?-_+~<>i!lI;:,"^`. "
//const unicodeRamp := "█▉▊▋▌▍▎▏▓▒░■□@&%$#*+=-~:;!,\".^`' "

type ConvertDoneMsg struct {
	Err error
}

func ConvertImageToString(filePath string) error {
	isDebugEnvAny, _ := Shared().Get("debug")
	isDebugEnv := isDebugEnvAny.(bool)

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_ = Logger().Info("Successfully Loaded: " + filePath)

	inputImg, format, err := image.Decode(f)
	if err != nil {
		return err
	}
	_ = Logger().Info("format: " + format)

	textSizeAny, ok := Shared().Get("textSize")
	if !ok || textSizeAny == nil {
		return errors.New("textSize is nil")
	}
	textSize, ok := textSizeAny.(int)
	if !ok || textSize <= 0 {
		textSize = 8
	}

	fontAspectAny, ok := Shared().Get("fontAspect")
	if !ok || fontAspectAny == nil {
		return errors.New("fontAspect is nil")
	}
	fontAspect, ok := fontAspectAny.(float64)
	if !ok || fontAspect <= 0 {
		fontAspect = 2
	}

	grayImg := convertRgbaToGray(inputImg)
	if isDebugEnv {
		err := saveImageToDebugDir(grayImg, "grayscale_img", "")
		if err != nil {
			return err
		}
	}

	cols, rows := gridFromTextSize(grayImg, textSize, fontAspect)
	gridImage := downScaleToGrid(grayImg, cols, rows)
	if isDebugEnv {
		_ = saveImageToDebugDir(gridImage, "grid_gray", "")
	}

	//TODO if directionalRender is true Sobel filter on `small` (size cols×rows)

	splitImages, err := splitImage(textSize, fontAspect, gridImage, isDebugEnv)
	if err != nil {
		return err
	}

	splitImages[1].Bounds()

	return nil
}

func convertRgbaToGray(img image.Image) *image.Gray {
	var (
		bounds = img.Bounds()
		gray   = image.NewGray(bounds)
	)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	return gray
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

	cols = imgW / charW
	rows = imgH / charH
	if cols < 1 {
		cols = 1
	}
	if rows < 1 {
		rows = 1
	}

	return cols, rows
}

func downScaleToGrid(img image.Image, cols, rows int) *image.Gray {
	dst := image.NewGray(image.Rect(0, 0, cols, rows))
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), xdraw.Over, nil)
	return dst
}

func splitImage(textSize int, fontAspect float64, inputImg image.Image, isDebugEnv bool) ([]image.Image, error) {

	imgBounds := inputImg.Bounds()
	imgWidth, imgHeight := imgBounds.Dx(), imgBounds.Dy()
	_ = Logger().Info("imgWidth: " + strconv.Itoa(imgWidth) + " imgHeight: " + strconv.Itoa(imgHeight))

	charWidth := textSize
	charHeight := int(float64(textSize) * fontAspect)

	var rects []image.Rectangle
	for y := 0; y < imgHeight; y += charHeight {
		y1 := y + charHeight
		if y1 > imgHeight {
			y1 = imgHeight
		}
		for x := 0; x < imgWidth; x += charWidth {
			x1 := x + charWidth
			if x1 > imgWidth {
				x1 = imgWidth
			}

			rects = append(rects, image.Rect(
				imgBounds.Min.X+x, imgBounds.Min.Y+y,
				imgBounds.Min.X+x1, imgBounds.Min.Y+y1,
			))
		}
	}

	//save images to debug folder if flag true
	var tiles []image.Image
	for i, r := range rects {
		var tile image.Image

		if si, ok := inputImg.(interface {
			SubImage(r image.Rectangle) image.Image
		}); ok {
			tile = si.SubImage(r)
		} else {
			// Fallback: copy crop if subimage fails
			dst := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
			draw.Draw(dst, dst.Bounds(), inputImg, r.Min, draw.Src)
			tile = dst
		}
		tiles = append(tiles, tile)

		//If Debug true - save images
		if isDebugEnv {
			err := saveImageToDebugDir(tile, "image_"+strconv.Itoa(i), "cropped_img")
			if err != nil {
				return nil, err
			}
		}
	}

	return tiles, nil
}

func saveImageToDebugDir(img image.Image, outputName string, subFolderName string) error {
	if filepath.Ext(outputName) == "" {
		outputName += ".png"
	}
	if err := os.MkdirAll("debugFolder/"+subFolderName, 0o755); err != nil {
		return err
	}

	outPath := filepath.Join("debugFolder/"+subFolderName, outputName)
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	return png.Encode(out, img)
}
