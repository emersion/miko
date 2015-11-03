package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.network.FutureInputMessage;
import cr.fr.saucisseroyale.miko.network.NetworkClient;

import java.awt.Graphics2D;
import java.awt.Toolkit;
import java.io.IOException;
import java.text.ParseException;

import javax.swing.UnsupportedLookAndFeelException;
import javax.swing.plaf.synth.SynthLookAndFeel;

public class Miko {

  public static final int PROTOCOL_VERSION = 2;
  private static final String SERVER_ADDRESS = "miko.emersion.fr";
  private static final int SERVER_PORT = 9999;
  private static final int TICK_TIME = 10; // milliseconds
  private UiWindow window;
  private boolean closeRequested = false;
  private NetworkClient networkClient;
  private MessageHandler messageHandler;
  private float alpha; // for drawing, updated each loop

  private void connect() {
    // TODO
  }

  private void exit() {
    window.close();
    System.exit(0);
  }

  private void initFrame() throws ParseException, UnsupportedLookAndFeelException {
    SynthLookAndFeel lookAndFeel = new SynthLookAndFeel(); // TODO handle exc better than throws
    lookAndFeel.load(Miko.class.getResourceAsStream("/style.xml"), Miko.class);
    // UIManager.setLookAndFeel(lookAndFeel);
    // TODO enable synth l&f once stylesheet created
    window = new UiWindow();
    window.setRenderable((graphics) -> render(graphics));
    // TODO add ui panels
  }

  private void initNetwork() throws IOException {
    networkClient = new NetworkClient();
    networkClient.connect(SERVER_ADDRESS, SERVER_PORT);
    messageHandler = new DebugMessageHandler(networkClient); // TODO MikoMessageHandler
  }

  private void logic() {
    // TODO
  }

  private void loop() {
    long lastFrame = System.nanoTime();
    long accumulator = 0;
    while (!closeRequested) {
      long newTime = System.nanoTime();
      long deltaTime = newTime - lastFrame;
      lastFrame = newTime;
      accumulator += deltaTime;
      network();
      while (accumulator >= TICK_TIME * 1000000) {
        logic();
        accumulator -= TICK_TIME * 1000000;
      }
      alpha = (float) accumulator / (TICK_TIME * 1000000);
      window.render(); // calls render(graphics)
      Toolkit.getDefaultToolkit().sync();
    }
    exit();
  }

  private void network() {
    FutureInputMessage fim = networkClient.getMessage();
    while (fim != null) {
      fim.execute(messageHandler);
      fim = networkClient.getMessage();
    }
  }

  private void render(Graphics2D graphics) {
    // TODO render w/ alpha
  }

  private void run() throws Exception {
    initNetwork(); // TODO handle better than throws Exc
    initFrame();
    connect();
    loop();
  }

  public static void main(String... args) throws Exception {
    Miko miko = new Miko();
    miko.run();
  }
}
