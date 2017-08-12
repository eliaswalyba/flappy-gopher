package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/ttf"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)

	if err != nil {
		return fmt.Errorf("Could not intialize SDL: %v", err)
	}

	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Could not initialize TTF: %v", err)
	}

	defer ttf.Quit()

	win, ren, err := sdl.CreateWindowAndRenderer(1000, 600, sdl.WINDOW_SHOWN)

	if err != nil {
		return fmt.Errorf("Could not create window: %v", err)
	}
	defer win.Destroy()

	if err := drawTitle(ren); err != nil {
		return fmt.Errorf("Could not draw title: %v", err)
	}

	time.Sleep(5 * time.Second)

	s, err := newScene(ren)

	if err != nil {
		return fmt.Errorf("Could not create scene: %v", err)
	}

	defer s.destroy()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	select {
	case <-s.run(ctx, ren):
		return err
	case <-time.After(5 * time.Second):
		return nil
	}

}

func drawTitle(r *sdl.Renderer) error {

	r.Clear()

	f, err := ttf.OpenFont("res/fonts/flappy.ttf", 500)

	if err != nil {
		return fmt.Errorf("Could not load font: %v", err)
	}

	defer f.Close()

	s, err := f.RenderUTF8_Solid("Flappy Gopher by Elias Waly BA", sdl.Color{R: 255, G: 100, B: 0, A: 255})

	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}

	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)

	if err != nil {
		return fmt.Errorf("Could not create texture %v", err)
	}

	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	r.Present()

	return nil

}
