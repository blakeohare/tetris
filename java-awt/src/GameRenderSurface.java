import java.awt.image.BufferedImage;
import java.awt.Color;
import java.awt.Graphics2D;

public class GameRenderSurface {
  
  private int width;
  private int height;
  private Graphics2D g2d;
  private BufferedImage image;

  public GameRenderSurface(BufferedImage image) {
    this.image = image;
    this.width = image.getWidth();
    this.height = image.getHeight();
  }

  public void startSession() {
    this.g2d = this.image.createGraphics();
  }

  public void endSession() {
    this.g2d.dispose();
  }

  public GameRenderSurface fill(int red, int green, int blue) {
    this.g2d.setColor(new Color(red, green, blue));
    this.g2d.fillRect(0, 0, this.width, this.height);
    return this;
  }

  public GameRenderSurface rectangle(int x, int y, int width, int height, int red, int green, int blue) {
    this.g2d.setColor(new Color(red, green, blue));
    this.g2d.fillRect(x, y, width, height);
    return this;
  }
}
