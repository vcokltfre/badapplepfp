package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"runtime"
	"sync"

	"github.com/nfnt/resize"
	"golang.org/x/sync/semaphore"
)

const (
	FRAMERATE     = 30
	SIZE_MODIFIER = 48
	PIX_WIDTH     = 1444 / SIZE_MODIFIER
	PIX_HEIGHT    = 1080 / SIZE_MODIFIER
	WIDTH         = PIX_WIDTH * SIZE_MODIFIER
	HEIGHT        = PIX_HEIGHT * SIZE_MODIFIER
	LENGTH        = 6962
)

var CPUs = runtime.NumCPU()

func getImage(file string, width, height int) image.Image {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	src, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return resize.Resize(uint(width), uint(height), src, resize.Lanczos3)
}

func renderFrame(num int, frame image.Image, profile image.Image) {
	output := image.NewNRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	for y := 0; y < PIX_HEIGHT; y++ {
		for x := 0; x < PIX_WIDTH; x++ {
			pStart := image.Point{x * SIZE_MODIFIER, y * SIZE_MODIFIER}
			pEnd := image.Point{pStart.X + SIZE_MODIFIER, pStart.Y + SIZE_MODIFIER}
			rect := image.Rectangle{pStart, pEnd}

			r, g, b, _ := frame.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				draw.Draw(output, rect, profile, image.Point{}, draw.Src)
			} else {
				draw.Draw(output, rect, frame, image.Point{x * SIZE_MODIFIER, y * SIZE_MODIFIER}, draw.Src)
			}
		}
	}

	file, err := os.Create(fmt.Sprintf("_output/%04d.png", num))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, output)
	if err != nil {
		panic(err)
	}
}

func main() {
	profile := getImage("profile.png", SIZE_MODIFIER, SIZE_MODIFIER)

	wg := sync.WaitGroup{}
	wg.Add(LENGTH)

	sm := semaphore.NewWeighted(int64(CPUs))

	for i := 1; i <= LENGTH; i++ {
		go func(n int) {
			sm.Acquire(context.Background(), 1)
			defer sm.Release(1)
			defer wg.Done()

			frame := getImage(fmt.Sprintf("_assets/frames/%04d.png", n), PIX_WIDTH, PIX_HEIGHT)
			renderFrame(n, frame, profile)

			fmt.Printf("Rendered frame %04d\n", n)
		}(i)
	}

	wg.Wait()
}
