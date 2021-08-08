package main

import ("fmt"
				"math/rand"
				"time"
        "github.com/veandco/go-sdl2/sdl")

type Tetris struct {
	grid [][]int
	overlay [][]int
	overlayX int
	overlayY int
	overlayUsesTranspose bool
	fallCounter float64
	linesCleared int
	clearingLines map[int]bool
	fallPressed bool
}

func (tetris *Tetris) Init() {
	tetris.grid = make([][]int, 10)
	for x := 0; x < 10; x++ {
		col := make([]int, 20)
		for y := 0; y < 20; y++ {
			col[y] = 0
		}
		tetris.grid[x] = col
	}
	tetris.overlay = nil
	tetris.fallCounter = 0
	tetris.linesCleared = 0
	tetris.clearingLines = nil
}

func createOverlayHelper(pieceId int) []int {
	switch pieceId {
	case 0:
		return []int{
			0, 1, 0, 0,
			0, 1, 0, 0,
			0, 1, 0, 0,
			0, 1, 0, 0,
		}
	case 1:
		return []int{
			0, 0, 0, 0,
			0, 1, 1, 0,
			0, 1, 1, 0,
			0, 0, 0, 0,
		}
	case 2:
		return []int{
			0, 1, 0, 0,
			1, 1, 1, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		}
	case 3:
		return []int{
			0, 2, 0, 0,
			0, 2, 0, 0,
			0, 2, 2, 0,
			0, 0, 0, 0,
		}
	case 4:
		return []int{
			2, 0, 0, 0,
			2, 2, 0, 0,
			0, 2, 0, 0,
			0, 0, 0, 0,
		}
	case 5:
		return []int{
			0, 3, 0, 0,
			0, 3, 0, 0,
			3, 3, 0, 0,
			0, 0, 0, 0,
		}
	case 6:
		return []int{
			0, 3, 0, 0,
			3, 3, 0, 0,
			3, 0, 0, 0,
			0, 0, 0, 0,
		}
	}
		
	return nil // panic?
}
func createOverlay(pieceId int) [][]int {
	flat := createOverlayHelper(pieceId)
	cols := make([][]int, 4)
	for x := 0; x < 4; x++ {
		row := make([]int, 4)
		cols[x] = row;
		for y := 0; y < 4; y++ {
			cols[x][y] = flat[x + y * 4]
		}
	}
	return cols
}

func (tetris *Tetris) Update(events []string) {
	if tetris.clearingLines != nil {

	} else if tetris.overlay == nil {
		overlayId := rand.Intn(7)
		tetris.overlay = createOverlay(overlayId)
		tetris.fallCounter = 30.0
		tetris.overlayUsesTranspose = overlayId < 2
		tetris.overlayX = 3
		tetris.overlayY = 0
	} else {
		for _, e := range events {
			switch e {
			case "left:press":
				tetris.TryMoveOverlay(-1, 0)
			case "right:press":
				tetris.TryMoveOverlay(1, 0)
			case "space:press":
				tetris.RotateOverlay(false)
			case "up:press":
				tetris.RotateOverlay(true)
			case "down:press":
				tetris.fallPressed = true
			case "down:release":
				tetris.fallPressed = false
			}
		}

		if tetris.fallPressed {
			tetris.fallCounter -= 6.0
		}

		if tetris.fallCounter <= 0 {
			tetris.fallCounter = 30.0
			if !tetris.TryMoveOverlay(0, 1) {
				tetris.FlattenOverlay()
				tetris.overlay = nil
			}
		} else {
			tetris.fallCounter -= 1.0
		}
	}
}

func (tetris *Tetris) RotateOverlay(direction bool) {
	if tetris.overlay == nil {
		return
	}

	if tetris.overlayUsesTranspose {
		tetris.OverlayTranspose()
		if !tetris.IsOverlayValid() {
			tetris.OverlayTranspose()
		}
		return
	}

	if direction {
		tetris.OverlayTranspose()
		tetris.OverlayFlipColumns()
	} else {
		tetris.OverlayFlipColumns()
		tetris.OverlayTranspose()
	}

	if (!tetris.IsOverlayValid()) {
		if direction {
			tetris.OverlayFlipColumns()
			tetris.OverlayTranspose()
		} else {
			tetris.OverlayTranspose()
			tetris.OverlayFlipColumns()
		}
	}
}

func (tetris *Tetris) OverlayTranspose() {
	for y := 0; y < 4; y++ {
		for x := y + 1; x < 4; x++ {
			t := tetris.overlay[x][y]
			tetris.overlay[x][y] = tetris.overlay[y][x]
			tetris.overlay[y][x] = t
		}
	}
}

func (tetris *Tetris) OverlayFlipColumns() {
	for y := 0; y < 4; y++ {
		t := tetris.overlay[0][y]
		tetris.overlay[0][y] = tetris.overlay[2][y]
		tetris.overlay[2][y] = t
	}
}

func (tetris *Tetris) TryMoveOverlay(dx int, dy int) bool {
	tetris.overlayX += dx
	tetris.overlayY += dy
	if !tetris.IsOverlayValid() {
		tetris.overlayX -= dx
		tetris.overlayY -= dy
		return false
	}
	return true
}

func (tetris *Tetris) FlattenOverlay() {
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if tetris.overlay[x][y] != 0 {
				tetris.grid[x + tetris.overlayX][y + tetris.overlayY] = tetris.overlay[x][y]
			}
		}
	}
}

func (tetris Tetris) IsOverlayValid() bool {
	if tetris.overlay == nil {
		return true
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if tetris.overlay[x][y] != 0 {
				gridX := x + tetris.overlayX
				gridY := y + tetris.overlayY
				if gridX < 0 || gridX >= 10 || gridY < 0 || gridY >= 20 {
					return false
				}
				if tetris.grid[gridX][gridY] != 0 {
					return false
				}
			}
		}
	}
	return true
}

func (tetris *Tetris) Render(rend *GameRenderer) {
	rend.Fill(40, 40, 40)
	var tileSize int32 = 20
	var boardWidth int32 = tileSize * 10
	var boardHeight int32 = tileSize * 20
	var boardLeft int32 = (int32(rend.width) - boardWidth) / 2
	var boardTop int32 = (int32(rend.height) - boardHeight) / 2
	rend.Rectangle(boardLeft, boardTop, boardWidth, boardHeight, 0, 0, 0)

	color1 := []uint32 { 255, 255, 255 }
	color2 := []uint32 { 0, 128, 255 }
	color3 := []uint32 { 0, 128, 40 }
	var color []uint32
	for y := 0; y < 20; y++ {
		for x := 0; x < 10; x++ {
			colorId := tetris.grid[x][y]
			if colorId == 0 && tetris.overlay != nil {
				ox := x - tetris.overlayX
				oy := y - tetris.overlayY
				if ox >= 0 && ox < 4 && oy >= 0 && oy < 4 {
					colorId = tetris.overlay[ox][oy]
				}
			}
			if colorId > 0 {
				if colorId == 1 {
					color = color1
				} else if colorId == 2 {
					color = color2
				} else {
					color = color3
				}
				var left int32 = int32(x) * tileSize + boardLeft
				var top int32 = int32(y) * tileSize + boardTop
				rend.Rectangle(left, top, tileSize, tileSize, color[0], color[1], color[2])
			}
		}
	}
}

type GameRenderer struct {
	surface *sdl.Surface
	width int
	height int
}

func (rend *GameRenderer) Rectangle(x int32, y int32, width int32, height int32, red uint32, green uint32, blue uint32) {
	rect := sdl.Rect{x, y, width, height}
	var color uint32 = 0xff000000 | (red << 16) | (green << 8) | (blue)
	rend.surface.FillRect(&rect, color)
}

func (rend *GameRenderer) Fill(red uint32, green uint32, blue uint32) {
	var color uint32 = 0xff000000 | (red << 16) | (green << 8) | (blue)
	rend.surface.FillRect(nil, color)
}

func main() {

	const width = 800
	const height = 600

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Tetris", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	rend := GameRenderer{surface, width, height}

	tetris := Tetris{nil, nil, 0, 0, false, 0, 0, nil, false}
	tetris.Init()
	fmt.Println(tetris.fallCounter)

	running := true

	const fps = 60;
	const nspf int64 = 1000000000 / fps

	var x int32 = 0
	var y int32 = 0

	pressedKeys := make(map[string]bool)

	for running {

		pressedEvents := make([]string, 0)
		frameStart := time.Now().UnixNano()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				isPress := event.GetType() == sdl.KEYDOWN
				keyCode := t.Keysym.Sym
				var keyName string = ""
				switch keyCode {
				case sdl.K_LEFT:
					keyName = "left"
				case sdl.K_RIGHT:
					keyName = "right"
				case sdl.K_DOWN:
					keyName = "down"
				case sdl.K_UP:
					keyName = "up"
				case sdl.K_SPACE:
					keyName = "space"
				}
				var wasPressed = false
				if value, ok := pressedKeys[keyName]; ok {
					wasPressed = value
				}
				if wasPressed != isPress {
					pressedKeys[keyName] = isPress
					var newName string
					if isPress {
						newName = keyName + ":press"
					} else {
						newName = keyName + ":release"
					}
					pressedEvents = append(pressedEvents, newName)
				}
				break
			}
		}

		tetris.Update(pressedEvents)
		tetris.Render(&rend)
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
