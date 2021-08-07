package main

// This is not Tetris.
// This is just some sample boilerplate code I found for SDL2 in Go

import ("time"
        "github.com/veandco/go-sdl2/sdl")

func main() {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Tetris", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	running := true

	const fps = 60;
	const nspf int64 = 1000000000 / fps

	var x int32 = 0
	var y int32 = 0

	for running {

		frameStart := time.Now().UnixNano()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			}
		}

		surface.FillRect(nil, 0)
		rect := sdl.Rect{x, y, 200, 200}
		surface.FillRect(&rect, 0xffff0000)
		window.UpdateSurface()

		x += 2
		y += 1
	
		if x > 600 {
			x = 0
		}
		if y > 400 {
			y = 0
		}

		frameEnd := time.Now().UnixNano()
		diff := frameEnd - frameStart
		delay := nspf - diff
		if delay > 0 {
			time.Sleep(time.Duration(delay))
		}
	}
}
