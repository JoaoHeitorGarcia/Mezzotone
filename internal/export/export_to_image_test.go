package export

import (
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestASCIIToPNGCreatesValidPNG(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.png")

	err := ASCIIToPNG("hello\nworld", outPath, ASCIIExportOptions{
		FontSize:     14,
		DPI:          300,
		BG:           color.Black,
		FG:           color.White,
		TargetAspect: 1.0 / 2.3,
	})
	if err != nil {
		t.Fatalf("ASCIIToPNG failed: %v", err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open png output: %v", err)
	}
	defer f.Close()

	cfg, err := png.DecodeConfig(f)
	if err != nil {
		t.Fatalf("failed to decode png config: %v", err)
	}
	if cfg.Width < 1 || cfg.Height < 1 {
		t.Fatalf("invalid png dimensions: %dx%d", cfg.Width, cfg.Height)
	}
}

func TestASCIIFramesToGIFCreatesAnimatedGIF(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.gif")

	frames := []ASCIIGIFFrame{
		{ASCII: "frame one", Duration: 40 * time.Millisecond},
		{ASCII: "frame two", Duration: 90 * time.Millisecond},
	}

	err := ASCIIFramesToGIF(frames, outPath, ASCIIExportOptions{
		FontSize:     14,
		DPI:          300,
		BG:           color.Black,
		FG:           color.White,
		TargetAspect: 1.0 / 2.3,
	})
	if err != nil {
		t.Fatalf("ASCIIFramesToGIF failed: %v", err)
	}

	f, err := os.Open(outPath)
	if err != nil {
		t.Fatalf("failed to open gif output: %v", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		t.Fatalf("failed to decode animated gif: %v", err)
	}
	if len(g.Image) != len(frames) {
		t.Fatalf("expected %d gif frames, got %d", len(frames), len(g.Image))
	}

	expectedDelays := []int{4, 9}
	for i, want := range expectedDelays {
		if g.Delay[i] != want {
			t.Fatalf("frame %d delay mismatch: want %d got %d", i, want, g.Delay[i])
		}
	}
}

func TestASCIIFramesToGIFNoFramesReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "out.gif")

	err := ASCIIFramesToGIF(nil, outPath, ASCIIExportOptions{
		FontSize: 14,
		DPI:      300,
		BG:       color.Black,
		FG:       color.White,
	})
	if err == nil {
		t.Fatalf("expected error when exporting gif with no frames")
	}
}
