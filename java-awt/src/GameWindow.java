import java.awt.Canvas;
import java.awt.Dimension;
import java.awt.Graphics2D;
import java.awt.RenderingHints;
import java.awt.event.KeyAdapter;
import java.awt.event.KeyEvent;
import java.awt.event.MouseEvent;
import java.awt.event.MouseListener;
import java.awt.event.MouseMotionListener;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.awt.image.BufferStrategy;
import java.awt.image.BufferedImage;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Timer;
import java.util.TimerTask;
import javax.swing.JFrame;
import javax.swing.JPanel;

public abstract class GameWindow {
  
  private String title;
  private int width;
  private int height;
  private int fps;
  private GameRenderSurface canvas;
  private BufferedImage canvasBackingImage;
  private BufferStrategy bufferStrategy;
  private ArrayList<String> events = new ArrayList<>();

  public GameWindow(String title, int width, int height, int fps) {
    this.title = title;
    this.width = width;
    this.height = height;
    this.fps = fps;
  }

  public void show() {

		int dpi = java.awt.Toolkit.getDefaultToolkit().getScreenResolution();
		int pixelWidth = this.width * dpi / 96;
		int pixelHeight = this.height * dpi / 96;

    this.canvasBackingImage = new BufferedImage(
        this.width, this.height, BufferedImage.TYPE_INT_ARGB);

    JFrame frame = new JFrame(this.title);

    frame.addWindowListener(new WindowAdapter() { 
      @Override 
      public void windowClosing(WindowEvent e) {
        closeWindow();
      }
    });    
    frame.setSize(pixelWidth, pixelHeight);
    frame.setVisible(true);
    // frame.setIconImage(ResourceReader.loadImageFromLocalFile("icon.png"));

    JPanel canvasHost = (JPanel) frame.getContentPane();
    canvasHost.setPreferredSize(new Dimension(pixelWidth, pixelHeight));
    canvasHost.setLayout(null);

    canvasHost.removeAll();
    Canvas canvas = new Canvas();
    canvas.setBounds(0, 0, pixelWidth, pixelHeight);
    canvas.setSize(pixelWidth, pixelHeight);
    canvasHost.add(canvas);
    canvas.setIgnoreRepaint(true);

    canvas.createBufferStrategy(2);
    this.bufferStrategy = canvas.getBufferStrategy();

    canvasHost.addKeyListener(new KeyAdapter() {
      @Override
      public void keyPressed(KeyEvent e) {
        handleKeyPress(e.getKeyCode(), true);
      }

      @Override
      public void keyReleased(KeyEvent e) { 
        handleKeyPress(e.getKeyCode(), false);
      }
    });
    canvasHost.setFocusTraversalKeysEnabled(false);
    canvasHost.requestFocus();

    this.canvas = new GameRenderSurface(this.canvasBackingImage);

    Timer timer = new Timer();
    timer.scheduleAtFixedRate(new TimerTask() {
        @Override
        public void run() {
            timerTick();
        }
    }, 0, (int)(1000 / fps));
  }

  private HashMap<String, Boolean> pressedKeys = new HashMap<>();

  public boolean isPressed(String keyId) {
    Boolean state = pressedKeys.get(keyId);
    if (state == null) return false;
    return (boolean)state;
  }

  private String keyCodeToId(int keyCode) {
    switch (keyCode) {
      case 17: return "CTRL";
      case 32: return "SPACE";
      case 37: return "LEFT";
      case 38: return "UP";
      case 39: return "RIGHT";
      case 40: return "DOWN";
      case 157: return "CLOVER";
      default: 
        if (keyCode >= 65 && keyCode <= 90) {
          return Character.toString((char)keyCode);
        }
        return "UNKNOWN";
    }
  }

  private void closeWindow() {
    System.exit(0);
  }

  private void handleKeyPress(int keyCode, boolean isPress) {

    String keyId = keyCodeToId(keyCode);

    if (isPress && isPressed(keyId)) {
      // don't want key-repeat events coming in as multiple press events
      return;
    }

    pressedKeys.put(keyId, isPress);

    // Ctrl + W should close the window, I suppose.
    if (isPress && keyId.equals("W") && (isPressed("CLOVER") || isPressed("CTRL"))) {
      closeWindow();
    } else {
      this.events.add(keyId + ":" + (isPress ? "PRESS" : "RELEASE"));
    }

  }

  private void timerTick() {

    this.update(this.events);
    this.events.clear();

    this.canvas.startSession();
    this.render(this.canvas);
    this.canvas.endSession();

		Graphics2D g = (Graphics2D) this.bufferStrategy.getDrawGraphics();
		g.setRenderingHint(
			RenderingHints.KEY_INTERPOLATION, 
			RenderingHints.VALUE_INTERPOLATION_NEAREST_NEIGHBOR);
    
		g.drawImage(
			this.canvasBackingImage, 0, 0, this.width, this.height,
			0, 0, this.canvasBackingImage.getWidth(), this.canvasBackingImage.getHeight(),
			null);
    g.dispose();
    this.bufferStrategy.show();
  }

  public abstract void update(List<String> events);

  public abstract void render(GameRenderSurface canvas);
}