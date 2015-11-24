package cr.fr.saucisseroyale.miko;

import java.awt.AWTEvent;
import java.awt.Color;
import java.awt.Component;
import java.awt.DisplayMode;
import java.awt.EventQueue;
import java.awt.Graphics;
import java.awt.Graphics2D;
import java.awt.GraphicsDevice;
import java.awt.GraphicsEnvironment;
import java.awt.MouseInfo;
import java.awt.Point;
import java.awt.PointerInfo;
import java.awt.Toolkit;
import java.awt.Window;
import java.awt.event.KeyListener;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.awt.image.BufferStrategy;
import java.util.function.Consumer;

import javax.swing.JComponent;
import javax.swing.JFrame;
import javax.swing.RepaintManager;
import javax.swing.SwingUtilities;
import javax.swing.WindowConstants;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

/**
 * Fenêtre en plein écran avec des composants d'UI et un composant principal, peint manuellement.
 * <p>
 * Render actif en mode exclusif en plein écran. Les composants d'UI peuvent être cachés ou
 * affichés, et seront affichés les uns au dessus des autres dans l'ordre dans lequel ils sont
 * affichés. Le composant principal sera toujours en dessous de tous les composants d'UI visibles.
 * <p>
 * <ul>
 * <li>Pour ajouter un composant d'UI à la fenêtre, utiliser {@link #addUi(JComponent)}.
 * <li>Pour le rendre visible, utiliser {@link #showUi(JComponent)}.
 * <li>Pour le cacher, utiliser {@link #hideUi(JComponent)}.
 * <li>Pour définir le composant principal, utiliser {@link #setRenderable(Consumer)}.
 * <li>Pour afficher les composants sur l'écran, utiliser {@link #render()}.
 * </ul>
 */
class UiWindow {

  @SuppressWarnings("serial")
  private static class EventCatcherComponent extends Component {

    private KeyListener keyListener;

    public EventCatcherComponent(int width, int height) {
      setSize(width, height);
      setFocusable(true);
    }

    public void setKeyListener(KeyListener keyListener) {
      if (this.keyListener != null) {
        removeKeyListener(this.keyListener);
      }
      this.keyListener = keyListener;
      addKeyListener(keyListener);
    }
  }

  /**
   * Une EventQueue qui se synchronise avec un lock pour éviter des problèmes de concurrence.
   */
  private static class SynchronizedEventQueue extends EventQueue {

    private Object lock;

    public SynchronizedEventQueue(Object lock) {
      this.lock = lock;
    }

    @Override
    protected void dispatchEvent(AWTEvent event) {
      synchronized (lock) {
        super.dispatchEvent(event);
      }
    }
  }

  private static final Integer UI_VISIBLE_LAYER = Integer.valueOf(1);
  private static final Integer UI_HIDDEN_LAYER = Integer.valueOf(-1);
  private static final Integer MAIN_LAYER = Integer.valueOf(0);
  private static Logger logger = LogManager.getLogger("miko.ui");
  private SynchronizedEventQueue eventQueue;
  private Runnable closeListener;
  private JFrame frame;
  private GraphicsDevice device;
  private DisplayMode displayMode;
  private BufferStrategy strategy;
  private Consumer<Graphics2D> renderable;
  private EventCatcherComponent mainComponent;
  private Object uiLock = new Object();
  private int width;
  private int height;

  /**
   * Construit un {@link UiWindow} avec les écran et mode d'affichage par défaut.
   */
  public UiWindow() {
    this(GraphicsEnvironment.getLocalGraphicsEnvironment().getDefaultScreenDevice(), GraphicsEnvironment.getLocalGraphicsEnvironment()
        .getDefaultScreenDevice().getDisplayMode());
  }

  /**
   * Construit un {@link UiWindow} avec les écran et mode d'affichage spécifié, sans afficher la
   * fenêtre.
   *
   * @param device L'écran sur lequel afficher la fenêtre.
   * @param displayMode La mode d'affichage à utiliser pour afficher la fenêtre.
   *
   * @throws IllegalArgumentException Si l'écran ne supporte pas le mode d'affichage spécifié.
   */
  public UiWindow(GraphicsDevice device, DisplayMode displayMode) {
    if (!isSupportedDisplayMode(device, displayMode)) {
      throw new IllegalArgumentException("DisplayMode unsupported on this GraphicsDevice");
    }
    this.device = device;
    this.displayMode = displayMode;
    width = displayMode.getWidth();
    height = displayMode.getHeight();
    eventQueue = new SynchronizedEventQueue(uiLock);
    Toolkit.getDefaultToolkit().getSystemEventQueue().push(eventQueue);
    frame = new JFrame();
    frame.setUndecorated(true);
    frame.setResizable(false);
    frame.setFocusTraversalKeysEnabled(false);
    frame.setIgnoreRepaint(true);
    frame.setDefaultCloseOperation(WindowConstants.DO_NOTHING_ON_CLOSE);
    frame.addWindowListener(new WindowAdapter() {
      @Override
      public void windowClosing(WindowEvent e) {
        if (closeListener != null) {
          closeListener.run();
        }
      }
    });
    mainComponent = new EventCatcherComponent(width, height);
  }

  /**
   * Initialise et affiche la fenêtre.
   */
  public void initAndShow() {
    RepaintManager repaintDisabler = new RepaintManager() {
      @Override
      public void addDirtyRegion(JComponent c, int x, int y, int w, int h) {}

      @Override
      public void addDirtyRegion(Window window, int x, int y, int w, int h) {}
    };
    repaintDisabler.setDoubleBufferingEnabled(false);
    RepaintManager.setCurrentManager(repaintDisabler);
    device.setFullScreenWindow(frame);
    device.setDisplayMode(displayMode);
    frame.createBufferStrategy(2);
    strategy = frame.getBufferStrategy();
    frame.getLayeredPane().add(mainComponent, MAIN_LAYER);
    logger.info("Initialized and created window");
  }

  /**
   * Affiche les composants sur la fenêtre.
   */
  public void render() {
    do {
      do {
        Graphics graphics = strategy.getDrawGraphics();
        paintComponents((Graphics2D) graphics);
        graphics.dispose();
      } while (strategy.contentsRestored());
      strategy.show();
    } while (strategy.contentsLost());
  }

  /**
   * Ferme la fenêtre et libère les ressources associées à celle-ci.
   */
  public void close() {
    device.setFullScreenWindow(null);
    frame.setVisible(false);
    frame.dispose();
    logger.info("Closed and hid window");
  }

  /**
   * Ajoute un composant d'UI à la fenêtre, en le redimensionnant à la taille de la fenêtre, caché.
   *
   * @param component Le composant à ajouter à la fenêtre.
   */
  public void addUi(JComponent component) {
    addUi(component, true);
  }

  /**
   * Ajoute un composant d'UI à la fenêtre, caché.
   *
   * @param component Le composant à ajouter à la fenêtre.
   * @param resize Si le composant doit être redimensionné à la taille de la fenêtre.
   */
  public void addUi(JComponent component, boolean resize) {
    if (resize) {
      component.setSize(getWidth(), getHeight());
    }
    component.setVisible(false);
    frame.getLayeredPane().add(component, UI_HIDDEN_LAYER);
  }

  /**
   * Rend un composant d'UI visible et le met au premier plan.
   *
   * @param component Le composant à rendre visible.
   */
  public void showUi(JComponent component) {
    component.setVisible(true);
    frame.getLayeredPane().setLayer(component, UI_VISIBLE_LAYER, 0);
  }

  /**
   * Cache un copposant d'UI.
   *
   * @param component Le composant à cacher.
   */
  public void hideUi(JComponent component) {
    component.setVisible(false);
    frame.getLayeredPane().setLayer(component, UI_HIDDEN_LAYER, 0);
  }

  /**
   * Cache tous les composants d'UI.
   */
  public void hideAllUi() {
    for (Component c : frame.getLayeredPane().getComponentsInLayer(UI_VISIBLE_LAYER)) {
      hideUi((JComponent) c);
    }
  }

  /**
   * Définit le composant principal.
   * <p>
   * Il devra s'afficher sur les graphiques lorsque sa méthode {@link Consumer#accept(Object)
   * accept(Graphics2D)} sera appelée.
   *
   * @param renderable Le composant principal.
   */
  public void setRenderable(Consumer<Graphics2D> renderable) {
    this.renderable = renderable;
  }

  /**
   * Définit le listener de fermeture demandée de fenêtre. Appelé lors de la demande de fermeture de
   * la {@link UiWindow}.
   *
   * @param closeListener Le listener de fermeture.
   */
  public void setCloseRequestedListener(Runnable closeListener) {
    this.closeListener = closeListener;
  }

  /**
   * Définit le listener de clavier sur le composant principal.
   *
   * @param keyListener Le listener de clavier.
   */
  public void setKeyListener(KeyListener keyListener) {
    mainComponent.setKeyListener(keyListener);
  }

  /**
   * @return La largeur de la fenêtre.
   */
  public int getWidth() {
    return width;
  }

  /**
   * @return La hauteur de la fenêtre.
   */
  public int getHeight() {
    return height;
  }

  /**
   * @return Le device sur lequel la frame est affichée.
   */
  public GraphicsDevice getDevice() {
    return device;
  }

  /**
   * Retourne la position de la souris au-dessus du composant principal.
   * <p>
   * La position de la souris est renvoyée avec un système de coordonnées <b>différent</b> de celui
   * de Swing : les coordonnées sont dans le sens X vers la droite, et Y vers le haut, avec (0;0) au
   * <b>centre</b> du composant principal.
   *
   * @return La position de la souris par rapport au composant principal, ou null si la souris n'est
   *         pas au-dessus.
   */
  public Point getMousePosition() {
    PointerInfo pointerInfo = MouseInfo.getPointerInfo();
    if (pointerInfo.getDevice() != device) {
      return null;
    }
    Point mousePoint = pointerInfo.getLocation();
    for (Component component : frame.getLayeredPane().getComponentsInLayer(UI_VISIBLE_LAYER)) {
      Point componentPoint = new Point(mousePoint);
      SwingUtilities.convertPointFromScreen(mousePoint, component);
      if (component.contains(componentPoint)) {
        return null;
      }
    }
    // change coordinate system
    Point offsetMousePoint = new Point(mousePoint.x - width / 2, height / 2 - mousePoint.y);
    return offsetMousePoint;
  }

  private void paintComponents(Graphics2D graphics) {
    Component[] uiComponents = frame.getLayeredPane().getComponentsInLayer(UI_VISIBLE_LAYER);
    for (int i = 0; i < uiComponents.length; i++) { // top to bottom
      Component uiComponent = uiComponents[i];
      if (uiComponent.isOpaque() && uiComponent.getWidth() == width && uiComponent.getHeight() == height) {
        // fully opaque component found
        // only render ui in front of it + itself
        for (int j = i; j >= 0; j--) { // bottom to top
          synchronized (uiLock) {
            uiComponents[j].paint(graphics);
          }
        }
        return;
      }
    }
    // no fully opaque component found
    // draw opaque background
    graphics.setBackground(Color.BLACK);
    graphics.clearRect(0, 0, width, height);
    // render game
    if (renderable != null) {
      renderable.accept(graphics);
    }
    // render ui
    for (int i = uiComponents.length - 1; i >= 0; i--) { // bottom to top
      synchronized (uiLock) {
        uiComponents[i].paint(graphics);
      }
    }
  }

  private static final boolean isSupportedDisplayMode(GraphicsDevice device, DisplayMode displayMode) {
    for (DisplayMode supportedDisplayMode : device.getDisplayModes()) {
      if (displayMode.getWidth() == supportedDisplayMode.getWidth()
          && displayMode.getHeight() == supportedDisplayMode.getHeight()
          && displayMode.getBitDepth() == supportedDisplayMode.getBitDepth()
          && (displayMode.getRefreshRate() == DisplayMode.REFRESH_RATE_UNKNOWN || displayMode.getRefreshRate() == supportedDisplayMode
          .getRefreshRate())) {
        return true;
      }
    }
    return false;
  }

}
