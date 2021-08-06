import java.util.List;

public class TetrisEngine {
  
  int counter = 0;
  
  public TetrisEngine() {

  }

  public void update(List<String> events) {
    
  }

  public void render(GameRenderSurface canvas) {
    this.counter++;
    if (this.counter > 255) this.counter = 0;

    canvas
      .fill(this.counter, 0, 255 - this.counter)
      .rectangle(10, 20, 50, 60, 255, 255, 0);
  }
}