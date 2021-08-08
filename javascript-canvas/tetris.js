const Tetris = () => {
  const WIDTH = 800;
  const HEIGHT = 600;

  const MAX_CLEAR_COUNTER = 30;

  let getGrid = (width, height, defaultValue) => {
    let grid = [];
    while (width --> 0) {
      grid.push(arrayOfSize(height).map(_ => defaultValue));
    }
    return grid;
  };

  let arrayOfSize = n => {
    let output = [];
    while (n --> 0) output.push(null);
    return output;
  };

  let grid = getGrid(10, 20, 0);
  let overlay = null;
  let overlayX = 3;
  let overlayY = 0;
  let fallCounter = 0;
  let clearCounter = 0;
  let linesToKeepForClearing = null;
  let lineCount = 0;
  let dropPressed = false;

  const COLORS = {
    WHITE: [240, 240, 240],
    CERULEAN: [0, 128, 255],
    GREEN: [0, 128, 50],
    ORANGE: [255, 128, 0],
    YELLOW: [255, 240, 0],
    RED: [255, 0, 30],
    PURPLE: [128, 0, 140],
    MAGENTA: [255, 40, 255],
    BLUE: [0, 0, 235],
    LIME: [50, 255, 0],
    BROWN: [128, 64, 0],
    TAN: [200, 150, 100],
    PINK: [255, 180, 225],
    CYAN: [0, 255, 255],
  };

  Object.values(COLORS).forEach(color => {
    let [r, g, b] = color;
    color.push(
      Math.floor(255 - (255 - r) * 2 / 3),
      Math.floor(255 - (255 - g) * 2 / 3),
      Math.floor(255 - (255 - b) * 2 / 3));
    color.push(
      Math.floor(r * 2 / 3),
      Math.floor(g * 2 / 3),
      Math.floor(b * 2 / 3));
  });

  let getCurrentColors = () => {
    let c;
    switch (Math.floor(lineCount / 10) % 10) {
      case 0: c = [COLORS.CERULEAN, COLORS.GREEN]; break;
      case 1: c = [COLORS.ORANGE, COLORS.YELLOW]; break;
      case 2: c = [COLORS.RED, COLORS.PURPLE]; break;
      case 3: c = [COLORS.BLUE, COLORS.MAGENTA]; break;
      case 4: c = [COLORS.GREEN, COLORS.LIME]; break;
      case 5: c = [COLORS.BROWN, COLORS.TAN]; break;
      case 6: c = [COLORS.PINK, COLORS.YELLOW]; break;
      case 7: c = [COLORS.GREEN, COLORS.PURPLE]; break;
      case 8: c = [COLORS.BLUE, COLORS.CYAN]; break;
      case 9: c = [COLORS.ORANGE, COLORS.RED]; break;
    }
    return [null, COLORS.WHITE, c[0], c[1]];
  };

  const render = (gfx, rc) => {
    gfx.fill(40, 40, 40);
    const TILE_SIZE = 20;
    const BOARD_WIDTH = TILE_SIZE * 10;
    const BOARD_HEIGHT = TILE_SIZE * 20;
    const BOARD_LEFT = Math.floor((WIDTH - BOARD_WIDTH) / 2);
    const BOARD_TOP = Math.floor((HEIGHT - BOARD_HEIGHT) / 2);

    gfx.rectangle(BOARD_LEFT, BOARD_TOP, BOARD_WIDTH, BOARD_HEIGHT, 0, 0, 0);

    let colors = getCurrentColors();
    for (let y = 0; y < 20; ++y) {
      for (let x = 0; x < 10; ++x) {
        let colorId = grid[x][y];
        if (colorId === 0 && overlay) {
          let ox = x - overlayX;
          let oy = y - overlayY;
          if (ox >= 0 && oy >= 0 && ox < 4 && oy < 4) {
            colorId = overlay[ox][oy];
          }
        }
        if (colorId > 0) {
          let color = colors[colorId];
          let left = BOARD_LEFT + TILE_SIZE * x;
          let top = BOARD_TOP + TILE_SIZE * y;
          gfx.rectangle(left, top, TILE_SIZE, TILE_SIZE, color[6], color[7], color[8]);
          gfx.rectangle(left, top, TILE_SIZE - 2, TILE_SIZE - 2, color[3], color[4], color[5]);
          gfx.rectangle(left + 2, top + 2, TILE_SIZE - 4, TILE_SIZE - 4, color[0], color[1], color[2]);
        }
      }
    }
  };

  const getOverlayImpl = id => {
    switch (id) {
      case 0: return [
        0, 1, 0, 0,
        0, 1, 0, 0,
        0, 1, 0, 0,
        0, 1, 0, 0,
      ];
      case 1: return [
        0, 0, 0, 0,
        0, 1, 1, 0,
        0, 1, 1, 0,
        0, 0, 0, 0,
      ];
      case 2: return [
        0, 1, 0, 0,
        1, 1, 1, 0,
        0, 0, 0, 0,
        0, 0, 0, 0,
      ];
      case 3: return [
        0, 2, 0, 0,
        0, 2, 0, 0,
        0, 2, 2, 0,
        0, 0, 0, 0,
      ];
      case 4: return [
        2, 0, 0, 0,
        2, 2, 0, 0,
        0, 2, 0, 0,
        0, 0, 0, 0,
      ];
      case 5: return [
        0, 3, 0, 0,
        0, 3, 0, 0,
        3, 3, 0, 0,
        0, 0, 0, 0,
      ];
      case 6: return [
        0, 3, 0, 0,
        3, 3, 0, 0,
        3, 0, 0, 0,
        0, 0, 0, 0,
      ];
      default: throw new Error();
    }
  };
  const getOverlay = id => {
    let overlayFlat = getOverlayImpl(id);
    let overlay = getGrid(4, 4, 0);
    for (let y = 0; y < 4; ++y) {
      for (let x = 0; x < 4; ++x) {
        overlay[x][y] = overlayFlat[x + y * 4];
      }
    }
    return overlay;
  };

  const update = (events) => {
    for (let e of events) {
      switch (e.key) {
        case "left":
          if (e.down) tryMoveOverlay(-1, 0);
          break;
        case "right":
          if (e.down) tryMoveOverlay(1, 0);
          break;
        case "up":
          if (e.down) tryRotate(true);
          break;
        case "space":
          if (e.down) tryRotate(false);
          break;
        case "down":
          dropPressed = e.down;
          break;
      }
    }

    if (linesToKeepForClearing) {
      let clearUpToX = Math.floor(10 * clearCounter / MAX_CLEAR_COUNTER);
      for (let y = 0; y < 20; ++y) {
        if (!linesToKeepForClearing[y])
        for (let x = 0; x < clearUpToX; ++x) {
          grid[x][y] = 0;
        }
      }

      if (clearCounter === MAX_CLEAR_COUNTER) {
        removeClearedLines();
        lineCount += 20 - Object.keys(linesToKeepForClearing).length;
        linesToKeepForClearing = null;
      }
      clearCounter++;

    } else if (!overlay) {
      let pieceId = Math.floor(Math.random() * 7);
      overlayUsesTransform = pieceId < 2;
      overlay = getOverlay(pieceId);
      overlayX = 3;
      overlayY = 0;
      fallCounter = 30;
    } else {

      fallCounter -= dropPressed ? 6.0 : 1.0;
      
      if (fallCounter <= 0) {
        fallCounter = 30;
        if (!tryMoveOverlay(0, 1)) {
          flattenOverlay();
          overlay = null;
          doClearRoutine();
        }
      }
    }
  };

  let tryRotate = dir => {
    if (overlayUsesTransform) {
      overlayTransform();
      if (!isOverlayValid()) {
        overlayTransform();
      }
      return;
    }
    if (dir) {
      overlayTransform();
      overlaySwapColumns();
    } else {
      overlaySwapColumns();
      overlayTransform();
    }

    if (!isOverlayValid()) {
      if (dir) {
        overlaySwapColumns();
        overlayTransform();
      } else {
        overlayTransform();
        overlaySwapColumns();
      }
    }
  };

  let overlayTransform = () => {
    for (let y = 0; y < 4; ++y) {
      for (let x = y + 1; x < 4; ++x) {
        let t = overlay[x][y];
        overlay[x][y] = overlay[y][x];
        overlay[y][x] = t;
      }
    }
  };

  let overlaySwapColumns = () => {
    for (let y = 0; y < 4; ++y) {
      let t = overlay[0][y];
      overlay[0][y] = overlay[2][y];
      overlay[2][y] = t;
    }
  };

  let tryMoveOverlay = (dx, dy) => {
    overlayX += dx;
    overlayY += dy;
    if (!isOverlayValid()) {
      overlayX -= dx;
      overlayY -= dy;
      return false;
    }
    return true;
  };

  let isOverlayValid = () => {
    if (!overlay) return true;
    for (let y = 0; y < 4; ++y) {
      for (let x = 0; x < 4; ++x) {
        if (overlay[x][y]) {
          let gx = x + overlayX;
          let gy = y + overlayY;
          if (gx < 0 || gy < 0 || gx >= 10 || gy >= 20) return false;
          if (grid[gx][gy]) return false;
        }
      }
    }
    return true;
  };

  let removeClearedLines = () => {
    let actualLine = 19;
    for (let y = 19; y >= 0; --y) {
      if (linesToKeepForClearing[y]) {
        for (let x = 0; x < 10; ++x) {
          grid[x][actualLine] = grid[x][y];
        }
        actualLine--;
      }
    }

    while (actualLine >= 0) {
      for (let x = 0; x < 10; ++x) {
        grid[x][actualLine] = 0;
      }
      actualLine--;
    }
  };

  let flattenOverlay = () => {
    for (let y = 0; y < 4; ++y) {
      for (let x = 0; x < 4; ++x) {
        if (overlay[x][y] !== 0) {
          let gx = overlayX + x;
          let gy = overlayY + y;
          grid[gx][gy] = overlay[x][y];
        }
      }
    }
  };

  let doClearRoutine = () => {
    linesToKeep = {};
    for (let y = 0; y < 20; ++y) {
      let hasEmpty = false;
      for (let x = 0; x < 10; ++x) {
        if (grid[x][y] === 0) {
          hasEmpty = true;
          break;
        }
      }
      if (hasEmpty) {
        linesToKeep[y] = true;
      }
    }

    let lines = Object.keys(linesToKeep).length;
    if (lines < 20) {
      linesToKeepForClearing = linesToKeep;
      clearCounter = 0;
    }
  };

  return {
    width: WIDTH,
    height: HEIGHT,
    render,
    update,
  };
};

function main() {
  const canvasHost = document.getElementById('canvas_host');
  let tetris = Tetris();
  Renderer.start(tetris.width, tetris.height, tetris.update, tetris.render, canvasHost);
}
