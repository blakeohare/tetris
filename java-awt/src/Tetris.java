import java.util.ArrayList;
import java.util.List;

public class Tetris {
  public static void main(String[] args) {
    TetrisEngine engine = new TetrisEngine();

    GameWindow window = new GameWindow("Tetris", 800, 600, 60) {
      @Override
      public void update(List<String> events) {
        ArrayList<String> gameEvents = new ArrayList<>();
        for (String e : events) {
          switch (e) {
            case "LEFT:PRESS":
              gameEvents.add("left");
              break;
            case "RIGHT:PRESS":
              gameEvents.add("right");
              break;
            case "UP:PRESS":
              gameEvents.add("rotate-right");
              break;
            case "SPACE:PRESS":
              gameEvents.add("rotate-left");
              break;
          }
        }
        if (isPressed("DOWN")) {
          gameEvents.add("drop");
        }
        engine.update(gameEvents);
      }

      @Override
      public void render(GameRenderSurface canvas) {
        engine.render(canvas);
      }
    };
    window.show();
  }
}
