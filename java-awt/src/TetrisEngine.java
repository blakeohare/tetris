import java.util.List;
import java.util.Random;

public class TetrisEngine {
  
  private Random random = new Random();
  private int counter = 0;

  private int[][] grid;
  private int[][] overlay = null;
  private boolean overlayUsesTranspose = false;
  private int overlayX = 3;
  private int overlayY = 0;
  private double fallCounter = 30;

  public TetrisEngine() {
    this.grid = new int[10][];
    for (int x = 0; x < 10; ++x) {
      this.grid[x] = new int[20];
    }
  }

  public void update(List<String> events) {
    if (overlay == null) {
      this.createNewOverlay();
      this.overlayX = 3;
      this.overlayY = 0;
      this.fallCounter = 30;
    }

    this.fallCounter -= 1.0;

    if (this.fallCounter <= 0) {
      boolean valid = this.tryMoveOverlay(0, 1);
      if (!valid) {
        this.flattenOverlay();
        this.overlay = null;
      }
      this.fallCounter = 30;
    }

    for (String event : events) {
      switch (event) {
        case "drop":
          this.fallCounter -= 4;
          break;
        case "left":
          this.tryMoveOverlay(-1, 0);
          break;
        case "right":
          this.tryMoveOverlay(1, 0);
          break;
        case "rotate-right":
          this.tryRotate(true);
          break;
        case "rotate-left":
          this.tryRotate(false);
          break;
      }
    }
  }

  private void tryRotate(boolean clockwise) {
    if (this.overlay == null) return;

    if (this.overlayUsesTranspose) {
      this.overlayTranspose();
      if (!this.isOverlayValid()) {
        this.overlayTranspose();
      }
    } else {
      if (clockwise) {
        this.overlayTranspose();
        this.overlayFlip();
      } else {
        this.overlayFlip();
        this.overlayTranspose();
      }
      if (!this.isOverlayValid()) {
        if (clockwise) {
          this.overlayFlip();
          this.overlayTranspose();
        } else {
          this.overlayTranspose();
          this.overlayFlip();
        }
      }
    }
  }

  private void overlayTranspose() {
    for (int y = 0; y < 4; ++y) {
      for (int x = y + 1; x < 4; ++x) {
        int t = this.overlay[x][y];
        this.overlay[x][y] = this.overlay[y][x];
        this.overlay[y][x] = t;
      }
    }
  }

  private void overlayFlip() {
    for (int y = 0; y < 4; ++y) {
      int t = this.overlay[0][y];
      this.overlay[0][y] = this.overlay[2][y];
      this.overlay[2][y] = t;
    }
  }

  private void flattenOverlay() {
    for (int y = 0; y < 4; ++y) {
      for (int x = 0; x < 4; ++x) {
        int color = this.overlay[x][y];
        if (color > 0) {
          int gridX = this.overlayX + x;
          int gridY = this.overlayY + y;
          this.grid[gridX][gridY] = color;
        }
      }
    }
  }

  private boolean tryMoveOverlay(int dx, int dy) {
    this.overlayX += dx;
    this.overlayY += dy;
    if (isOverlayValid()) {
      return true;
    }
    this.overlayX -= dx;
    this.overlayY -= dy;
    return false;
  }

  private boolean isOverlayValid() {

    for (int y = 0; y < 4; ++y) {
      for (int x = 0; x < 4; ++x) {
        int color = this.overlay[x][y];
        if (color > 0) {
          int gridX = x + this.overlayX;
          int gridY = y + this.overlayY;
          if (gridX < 0 || gridX >= 10 ||
              gridY < 0 || gridY >= 20 ||
              this.grid[gridX][gridY] != 0) {
            
            return false;
          }
        }
      }
    }
    return true;
  }

  private int[] getPiece(int id) {
    switch (id) {
      case 0: return new int[] {
        0, 1, 0, 0, 
        0, 1, 0, 0, 
        0, 1, 0, 0, 
        0, 1, 0, 0, 
      };
      case 1: return new int[] {
        0, 0, 0, 0, 
        0, 1, 1, 0, 
        0, 1, 1, 0, 
        0, 0, 0, 0, 
      };
      case 2: return new int[] {
        0, 1, 0, 0, 
        1, 1, 1, 0, 
        0, 0, 0, 0, 
        0, 0, 0, 0, 
      };
      case 3: return new int[] {
        0, 2, 0, 0, 
        0, 2, 0, 0, 
        0, 2, 2, 0, 
        0, 0, 0, 0, 
      };
      case 4: return new int[] {
        2, 0, 0, 0, 
        2, 2, 0, 0, 
        0, 2, 0, 0, 
        0, 0, 0, 0, 
      };
      case 5: return new int[] {
        0, 3, 0, 0, 
        0, 3, 0, 0, 
        3, 3, 0, 0, 
        0, 0, 0, 0, 
      };
      case 6: return new int[] {
        0, 3, 0, 0, 
        3, 3, 0, 0, 
        3, 0, 0, 0, 
        0, 0, 0, 0, 
      };
      default: throw new RuntimeException("Bad piece ID");
    }
  }

  private void createNewOverlay() {
    int pieceId = this.random.nextInt(7);
    int[] piece = getPiece(pieceId);
    this.overlayUsesTranspose = pieceId <= 1 || pieceId == 4 || pieceId == 6;
    this.overlay = new int[4][];
    for (int x = 0; x < 4; ++x) {
      this.overlay[x] = new int[4];
      for (int y = 0; y < 4; ++y) {
        this.overlay[x][y] = piece[x + y * 4];
      }
    }
  }

  public void render(GameRenderSurface canvas) {
    this.counter++;
    if (this.counter > 255) this.counter = 0;

    int tileSize = 20;
    int boardWidth = tileSize * 10;
    int screenWidth = 800;
    int boardLeft = (screenWidth - boardWidth) / 2;
    int boardHeight = tileSize * 20;
    int screenHeight = 600;
    int boardTop = (screenHeight - boardHeight) / 2;

    int[] color1 = new int[] { 255, 255, 255 };
    int[] color2 = new int[] { 0, 128, 255 };
    int[] color3 = new int[] { 0, 128, 50 };
    int[][] colors = new int[][] { 
      null,
      color1,
      color2, 
      color3,
    };

    canvas
      .fill(40, 40, 40)
      .rectangle(boardLeft, boardTop, boardWidth, boardHeight, 0, 0, 0);
    
    int colorId;
    int[] color;
    for (int y = 0; y < 20; ++y) {
      for (int x = 0; x < 10; ++x) {
        colorId = this.grid[x][y];
        if (colorId == 0 && 
            this.overlay != null &&
            y >= this.overlayY && y < this.overlayY + 4 &&
            x >= this.overlayX && x < this.overlayX + 4) {
          colorId = this.overlay[x - this.overlayX][y - this.overlayY];
        }

        if (colorId != 0) {
          color = colors[colorId];
          canvas.rectangle(boardLeft + tileSize * x, boardTop + tileSize * y, tileSize, tileSize, color[0], color[1], color[2]);
        }
      }
    }
  }
}