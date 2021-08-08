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
  linesToKeepForClear map[int]bool
  clearCounter int
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
  tetris.linesToKeepForClear = nil
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

  for _, e := range events {
    switch e {
    case "left:press":
      if tetris.overlay != nil {
        tetris.TryMoveOverlay(-1, 0)
      }
    case "right:press":
      if tetris.overlay != nil {
        tetris.TryMoveOverlay(1, 0)
      }
    case "space:press":
      if tetris.overlay != nil {
        tetris.RotateOverlay(false)
      }
    case "up:press":
      if tetris.overlay != nil {
        tetris.RotateOverlay(true)
      }
    case "down:press":
      tetris.fallPressed = true
    case "down:release":
      tetris.fallPressed = false
    }
  }

  if tetris.linesToKeepForClear != nil {
    const clearCounterMax = 30
    clearUpToX := 10 * tetris.clearCounter / clearCounterMax
    for y := 0; y < 20; y++ {
      if _, ok := tetris.linesToKeepForClear[y]; !ok {
        for x := 0; x < clearUpToX; x++ {
          tetris.grid[x][y] = 0
        }
      }
    }
    if tetris.clearCounter == clearCounterMax {
      tetris.DropClearedLines()
      tetris.linesToKeepForClear = nil
    } else {
      tetris.clearCounter++
    }
  } else if tetris.overlay == nil {
    overlayId := rand.Intn(7)
    tetris.overlay = createOverlay(overlayId)
    tetris.fallCounter = 30.0
    tetris.overlayUsesTranspose = overlayId < 2
    tetris.overlayX = 3
    tetris.overlayY = 0
  } else {

    if tetris.fallPressed {
      tetris.fallCounter -= 6.0
    }

    if tetris.fallCounter <= 0 {
      tetris.fallCounter = 30.0
      if !tetris.TryMoveOverlay(0, 1) {
        tetris.FlattenOverlay()
        tetris.overlay = nil
        tetris.DoClearLineCheck()
      }
    } else {
      tetris.fallCounter -= 1.0
    }
  }
}

func (tetris *Tetris) DropClearedLines() {
  actualLine := 19
  for y := 19; y >= 0; y-- {
    if _, ok := tetris.linesToKeepForClear[y]; ok {
      for x := 0; x < 10; x++ {
        tetris.grid[x][actualLine] = tetris.grid[x][y]
      }
      actualLine--
    }
  }

  for actualLine >= 0 {
    for x := 0; x < 10; x++ {
      tetris.grid[x][actualLine] = 0
    }
    actualLine--
    tetris.linesCleared++
  }
}

func (tetris *Tetris) DoClearLineCheck() {
  keepTheseLines := make(map[int]bool)
  for y := 0; y < 20; y++ {
    hasEmpty := false
    for x := 0; x < 10; x++ {
      if tetris.grid[x][y] == 0 {
        hasEmpty = true
        break
      }
    }
    if hasEmpty {
      keepTheseLines[y] = true
    }
  }

  if len(keepTheseLines) < 20 {
    tetris.linesToKeepForClear = keepTheseLines
    tetris.clearCounter = 0
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

var WHITE = []int{ 240, 240, 240 }
var CERULEAN = []int{ 0, 128, 255 };
var GREEN = []int{ 0, 128, 50 };
var ORANGE = []int{ 255, 128, 0 };
var YELLOW = []int{ 255, 240, 0 };
var RED = []int{ 255, 0, 30 };
var PURPLE = []int{ 128, 0, 140 };
var MAGENTA = []int{ 255, 40, 255 };
var BLUE = []int{ 0, 0, 235 };
var LIME = []int{ 50, 255, 0 };
var BROWN = []int{ 128, 64, 0 };
var TAN = []int{ 200, 150, 100 };
var PINK = []int{ 255, 180, 225 };
var CYAN = []int{ 0, 255, 255 };

func (tetris Tetris) GetCurrentColors() [][]int {

  var colorA []int
  var colorB []int
  
  switch (tetris.linesCleared / 10) % 10 {
  case 0:
    colorA = CERULEAN
    colorB = GREEN
  case 1:
    colorA = ORANGE
    colorB = YELLOW
  case 2:
    colorA = RED
    colorB = PURPLE
  case 3:
    colorA = BLUE
    colorB = MAGENTA
  case 4:
    colorA = GREEN
    colorB = LIME
  case 5:
    colorA = BROWN
    colorB = TAN
  case 6:
    colorA = PINK
    colorB = YELLOW
  case 7:
    colorA = GREEN
    colorB = PURPLE
  case 8:
    colorA = BLUE
    colorB = CYAN
  case 9:
    colorA = ORANGE
    colorB = RED
  }
  colors := make([][]int, 4)
  colors[0] = nil
  colors[1] = WHITE
  colors[2] = colorA
  colors[3] = colorB

  for i := 1; i <= 3; i++ {
    color := colors[i]
    newColor := make([]int, 9)
    r := color[0]
    g := color[1]
    b := color[2]
    newColor[0] = r
    newColor[1] = g
    newColor[2] = b
    newColor[3] = 255 - (255 - r) * 2 / 3
    newColor[4] = 255 - (255 - g) * 2 / 3
    newColor[5] = 255 - (255 - b) * 2 / 3
    newColor[6] = r * 2 / 3
    newColor[7] = g * 2 / 3
    newColor[8] = b * 2 / 3
    colors[i] = newColor
  }

  return colors
}

func (tetris *Tetris) Render(rend *GameRenderer) {
  rend.Fill(40, 40, 40)
  var tileSize int32 = 20
  var boardWidth int32 = tileSize * 10
  var boardHeight int32 = tileSize * 20
  var boardLeft int32 = (int32(rend.width) - boardWidth) / 2
  var boardTop int32 = (int32(rend.height) - boardHeight) / 2
  rend.Rectangle(boardLeft, boardTop, boardWidth, boardHeight, 0, 0, 0)

  colors := tetris.GetCurrentColors()
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
        color := colors[colorId]
        var left int32 = int32(x) * tileSize + boardLeft
        var top int32 = int32(y) * tileSize + boardTop
        rend.Rectangle(left, top, tileSize, tileSize, uint32(color[6]), uint32(color[7]), uint32(color[8]))
        rend.Rectangle(left, top, tileSize - 2, tileSize - 2, uint32(color[3]), uint32(color[4]), uint32(color[5]))
        rend.Rectangle(left + 2, top + 2, tileSize - 4, tileSize - 4, uint32(color[0]), uint32(color[1]), uint32(color[2]))
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

  rand.Seed(time.Now().UTC().UnixNano())
  rend := GameRenderer{surface, width, height}
  tetris := Tetris{nil, nil, 0, 0, false, 0, 0, nil, 0, false}
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
