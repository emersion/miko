package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.network.FutureInputMessage;
import cr.fr.saucisseroyale.miko.network.NetworkClient;

import java.awt.Color;
import java.awt.Graphics2D;
import java.awt.GraphicsDevice;
import java.awt.GraphicsEnvironment;
import java.awt.image.BufferStrategy;
import java.io.IOException;

import javax.swing.JFrame;

public class Miko {

  public static final String GAME_NAME = "Miko";
  private static final String SERVER_ADDRESS = "localhost";
  private static final int SERVER_PORT = 9997;
  private static final int TICK_TIME = 10;
  private JFrame frame;
  private BufferStrategy strategy;
  private boolean closeRequested = false;
  private NetworkClient networkClient = new NetworkClient();
  private MessageHandler messageHandler = new MessageHandler(networkClient);

  private void connect() {
    // TODO
  }

  private void exit() {
    frame.setVisible(false);
    frame.dispose();
    System.exit(0);
  }

  private void initFrame() {
    frame = new JFrame(GAME_NAME);
    frame.setUndecorated(true);
    frame.setResizable(false);
    frame.setIgnoreRepaint(true);
    // TODO visible
    // frame.setVisible(true);
    GraphicsDevice device =
        GraphicsEnvironment.getLocalGraphicsEnvironment().getDefaultScreenDevice();
    device.setFullScreenWindow(frame);
    frame.createBufferStrategy(2);
    strategy = frame.getBufferStrategy();
    // TODO mettre ratio de rendu
    // TODO ajouter key listener
  }

  private void initNetwork() throws IOException {
    networkClient.connect(SERVER_ADDRESS, SERVER_PORT);
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
      float alpha = (float) accumulator / (TICK_TIME * 1000000);
      render(alpha);
      frame.getToolkit().sync();
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

  private void render(float alpha) {
    Graphics2D g = (Graphics2D) strategy.getDrawGraphics();
    g.setBackground(Color.BLACK);
    g.setColor(Color.WHITE);
    // TODO
    g.dispose();
    strategy.show();
  }

  private void start() throws IOException {
    initNetwork();
    initFrame();
    connect();
    loop();
  }

  public static void main(String... args) throws IOException {
    Miko miko = new Miko();
    miko.start();
  }
}
