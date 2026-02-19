package tui

import (
	"fmt"
	"image"
	"net/http"
	"sync"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	iconRenderCache = map[string]string{}
	iconCacheMu     sync.RWMutex
)

func renderIconUnicodeFromURL(iconURL string, cols int, rows int) (string, error) {
	if iconURL == "" {
		return "", fmt.Errorf("empty icon url")
	}

	iconCacheMu.RLock()
	if cached, ok := iconRenderCache[iconURL]; ok {
		iconCacheMu.RUnlock()
		return cached, nil
	}
	iconCacheMu.RUnlock()

	client := &http.Client{Timeout: 2500 * time.Millisecond}
	resp, err := client.Get(iconURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("http status %d", resp.StatusCode)
	}

	src, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", err
	}

	ansi := imageToUnicodeANSI(src, cols, rows)
	iconCacheMu.Lock()
	iconRenderCache[iconURL] = ansi
	iconCacheMu.Unlock()
	return ansi, nil
}

// imageToUnicodeANSI converts image to truecolor unicode blocks using "▀" char.
func imageToUnicodeANSI(src image.Image, cols int, rows int) string {
	if cols < 2 {
		cols = 2
	}
	if rows < 2 {
		rows = 2
	}

	targetH := rows * 2
	var out string
	for y := 0; y < rows; y++ {
		topY := y * 2
		botY := topY + 1
		for x := 0; x < cols; x++ {
			tr, tg, tb := sampleNearest(src, x, topY, cols, targetH)
			br, bg, bb := sampleNearest(src, x, botY, cols, targetH)
			out += fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", tr, tg, tb, br, bg, bb)
		}
		out += "\x1b[0m\n"
	}
	return out
}

func sampleNearest(src image.Image, x int, y int, targetW int, targetH int) (int, int, int) {
	b := src.Bounds()
	srcW := b.Dx()
	srcH := b.Dy()
	if srcW <= 0 || srcH <= 0 {
		return 255, 255, 255
	}

	sx := b.Min.X + (x*srcW)/targetW
	sy := b.Min.Y + (y*srcH)/targetH

	r, g, bl, _ := src.At(sx, sy).RGBA()
	return int(r >> 8), int(g >> 8), int(bl >> 8)
}
