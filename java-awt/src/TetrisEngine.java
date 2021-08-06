import java.util.List;

public class TetrisEngine {
  
  int counter = 0;

  int[][] grid;
  int[][] overlay = null;

  public TetrisEngine() {
    this.grid = new int[10][];
    for (int x = 0; x < 10; ++x) {
      this.grid[x] = new int[20];
      for (int y = 0; y < 20; ++y) {
        this.grid[x][y] = (x + y) % 4;
      }
    }
  }

  public void update(List<String> events) {

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
        if (colorId != 0) {
          color = colors[colorId];
          canvas.rectangle(boardLeft + tileSize * x, boardTop + tileSize * y, tileSize, tileSize, color[0], color[1], color[2]);
        }
      }
    }
  }
}