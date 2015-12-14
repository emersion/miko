package cr.fr.saucisseroyale.miko;

import java.awt.AWTEvent;
import java.awt.AWTException;
import java.awt.BufferCapabilities;
import java.awt.BufferCapabilities.FlipContents;
import java.awt.Color;
import java.awt.Component;
import java.awt.EventQueue;
import java.awt.Graphics;
import java.awt.Graphics2D;
import java.awt.GraphicsConfiguration;
import java.awt.GraphicsDevice;
import java.awt.GraphicsEnvironment;
import java.awt.ImageCapabilities;
import java.awt.KeyboardFocusManager;
import java.awt.MouseInfo;
import java.awt.Point;
import java.awt.PointerInfo;
import java.awt.Rectangle;
import java.awt.Toolkit;
import java.awt.Window;
import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;
import java.awt.event.WindowAdapter;
import java.awt.event.WindowEvent;
import java.awt.image.BufferStrategy;
import java.util.function.Consumer;
import java.util.function.IntConsumer;

import javax.swing.JComponent;
import javax.swing.JFrame;
import javax.swing.RepaintManager;
import javax.swing.SwingUtilities;
import javax.swing.WindowConstants;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

/**
 * Fenêtre en mode plein écran ou fenêtré sans bordures avec des composants d'UI et un composant
 * principal, peint manuellement.
 * <p>
 * Render actif en mode exclusif en plein écran, ou en mode fenêtré sans bordures. Les composants
 * d'UI peuvent être cachés ou affichés, et seront affichés les uns au dessus des autres dans
 * l'ordre dans lequel ils sont affichés. Le composant principal sera toujours en dessous de tous
 * les composants d'UI visibles.
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

  // we won't serialize it
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

  private static class RepaintDisabler extends RepaintManager {
    // @noformatting
    public RepaintDisabler() {
      setDoubleBufferingEnabled(false);
    }
    @Override
    public void addDirtyRegion(JComponent c, int x, int y, int w, int h) {}
    @Override
    public void addDirtyRegion(Window window, int x, int y, int w, int h) {}
    @Override
    public synchronized void addInvalidComponent(JComponent invalidComponent) {}
    @Override
    public void markCompletelyClean(JComponent aComponent) {}
    @Override
    public void markCompletelyDirty(JComponent aComponent) {}
    @Override
    public void paintDirtyRegions() {}
    @Override
    public Rectangle getDirtyRegion(JComponent aComponent) {
      // pretend the component is completely dirty
      return new Rectangle(0, 0, Integer.MAX_VALUE, Integer.MAX_VALUE);
    }
    @Override
    public boolean isCompletelyDirty(JComponent aComponent) {
      return true;
    }
    @Override
    public synchronized void removeInvalidComponent(JComponent component) {}
    @Override
    public void validateInvalidComponents() {}
    // @formatting
  }

  private static final Integer UI_VISIBLE_LAYER = Integer.valueOf(1);
  private static final Integer UI_HIDDEN_LAYER = Integer.valueOf(-1);
  private static final Integer MAIN_LAYER = Integer.valueOf(0);
  private static Logger logger = LogManager.getLogger("miko.ui");
  private Runnable closeListener;
  private IntConsumer globalKeyDownListener;
  private JFrame frame; // will be null when we're disposed
  private GraphicsDevice device;
  private BufferStrategy strategy;
  private Consumer<Graphics2D> renderable;
  private EventCatcherComponent mainComponent;
  private Object uiLock = new Object();
  private boolean fullscreen;
  private int width;
  private int height;

  /**
   * Tente de créer et afficher la fenêtre sur l'écran par défaut.
   * <p>
   * La fenêtre ne sera pas forcément affichée en plein écran ; récupérer l'état final avec
   * {@link #isFullscreen()}.
   *
   * @param fullscreen Si la fenêtre doit être en plein écran exclusif ou en fenêtré sans bordures.
   *
   * @throws IllegalStateException Si {@code show} ou {@link #dispose()} a déjà été appelé.
   */
  public UiWindow(boolean fullscreen) {
    this(GraphicsEnvironment.getLocalGraphicsEnvironment().getDefaultScreenDevice(), fullscreen);
  }

  /**
   * Tente de créer et afficher la fenêtre sur l'écran spécifié.
   * <p>
   * La fenêtre ne sera pas forcément affichée en plein écran ; récupérer l'état final avec
   * {@link #isFullscreen()}.
   *
   * @param device L'écran sur lequel tenter d'afficher la fenêtre.
   * @param requestedFullscreen Si la fenêtre doit être en plein écran exclusif ou en fenêtré sans
   *        bordures.
   *
   * @throws IllegalStateException Si {@code show} ou {@link #dispose()} a déjà été appelé.
   */
  public UiWindow(GraphicsDevice device, boolean requestedFullscreen) {
    this.device = device;
    // make sure we have control over repainting before trying anything
    Toolkit.getDefaultToolkit().getSystemEventQueue().push(new SynchronizedEventQueue(uiLock));
    RepaintManager.setCurrentManager(new RepaintDisabler());
    Class<?> openglGraphicsConfigClass;
    try {
      openglGraphicsConfigClass = Class.forName("sun.java2d.opengl.OGLGraphicsConfig");
    } catch (ClassNotFoundException e) {
      String errorMessage = "Interface OGLGraphicsConfig not found.";
      logger.fatal(errorMessage);
      throw new InternalError(errorMessage);
    }
    GraphicsConfiguration bestConfiguration = null;
    for (GraphicsConfiguration gc : device.getConfigurations()) {
      // drop non opengl configurations
      if (!openglGraphicsConfigClass.isAssignableFrom(gc.getClass())) {
        continue;
      }
      // drop non page-flipping configurations
      if (!gc.getBufferCapabilities().isPageFlipping()) {
        continue;
      }
      // additionally drop fullscreen-only page-flipping configurations if we're requesting windowed
      // mode
      if (!requestedFullscreen && gc.getBufferCapabilities().isFullScreenRequired()) {
        continue;
      }
      // drop non accelerated configurations
      if (!gc.getImageCapabilities().isAccelerated()) {
        continue;
      }
      // drop non double accelereated configurations
      if (!gc.getBufferCapabilities().getFrontBufferCapabilities().isAccelerated()
          || !gc.getBufferCapabilities().getBackBufferCapabilities().isAccelerated()) {
        continue;
      }
      // drop null flip contents to make sure we do support vsync
      if (gc.getBufferCapabilities().getFlipContents() == null) {
        continue;
      }
      bestConfiguration = gc;
      break;
    }
    if (bestConfiguration == null) {
      String errorMessage = "No suitable GraphicsConfiguration found. Try updating your video card drivers.";
      logger.fatal(errorMessage);
      throw new RuntimeException(errorMessage);
    }
    width = bestConfiguration.getBounds().width;
    height = bestConfiguration.getBounds().height;
    mainComponent = new EventCatcherComponent(width, height);
    // create the frame
    frame = new JFrame(bestConfiguration);
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
    KeyboardFocusManager.getCurrentKeyboardFocusManager().addKeyEventDispatcher(e -> {
      if (e.getID() == KeyEvent.KEY_PRESSED && globalKeyDownListener != null) {
        globalKeyDownListener.accept(e.getKeyCode());
      }
      return false;
    });
    // show the frame
    if (requestedFullscreen && device.isFullScreenSupported()) {
      device.setFullScreenWindow(frame);
    }
    fullscreen = device.getFullScreenWindow() != null;
    if (!fullscreen) {
      frame.setSize(width, height);
      Rectangle deviceBounds = device.getDefaultConfiguration().getBounds();
      frame.setLocation(deviceBounds.x + deviceBounds.width - width, deviceBounds.y + deviceBounds.height - height);
      frame.setVisible(true);
    }
    // create double accelerated (volatile) page-flipping buffer with vsync
    // the only way to have vsync on with the opengl pipeline on windows and linux is to have
    // flipcontents on COPIED so that the VolatileSurfaceManager create WGLVSyncOffScreenSurfaceData
    // (on Windows) or GLXVSyncOffScreenSurfaceData (on Linux) as backbuffers
    // see initAcceleratedSurface() in WGLVolatileSurfaceManager and GLXVolatileSurfaceManager
    BufferCapabilities bc = new BufferCapabilities(new ImageCapabilities(true), new ImageCapabilities(true), FlipContents.COPIED);
    // we must use a sun.* class in order to tell the surfacemanager we want vsync
    @SuppressWarnings("restriction")
    BufferCapabilities vsyncBuffer =
    new sun.java2d.pipe.hw.ExtendedBufferCapabilities(bc, sun.java2d.pipe.hw.ExtendedBufferCapabilities.VSyncType.VSYNC_ON);
    try {
      frame.createBufferStrategy(2, vsyncBuffer);
    } catch (AWTException e) {
      e.printStackTrace();
    }
    strategy = frame.getBufferStrategy();
    frame.getLayeredPane().add(mainComponent, MAIN_LAYER);
    logger.debug("Initialized and created window");
  }

  /**
   * Affiche les composants sur la fenêtre.
   *
   * @throws IllegalStateException Si {@code show} n'a pas été appelé, ou {@link #dispose()} a déjà
   *         été appelé.
   */
  public void render() {
    checkFrameNotDisposed();
    do {
      do {
        Graphics graphics = strategy.getDrawGraphics();
        paintComponents((Graphics2D) graphics);
        graphics.dispose();
      } while (strategy.contentsRestored());
      long beforeShowTime = System.nanoTime();
      strategy.show(); // vsync (ie sleep) here
      System.out.println("Temps de sleep en microsecondes : " + (System.nanoTime() - beforeShowTime) / 1000);
    } while (strategy.contentsLost());
  }

  /**
   * Ferme la fenêtre et libère les ressources associées à celle-ci. Une fois fermée, elle ne pourra
   * plus être réaffichée.
   */
  public void dispose() {
    if (frame != null) {
      frame.dispose();
      frame = null;
      logger.debug("Closed and hid window");
    }
  }

  /**
   * Ajoute un composant d'UI à la fenêtre, en le redimensionnant à la taille de la fenêtre, caché.
   *
   * @param component Le composant à ajouter à la fenêtre.
   */
  public void addUi(JComponent component) {
    checkFrameNotDisposed();
    addUi(component, true);
  }

  /**
   * Ajoute un composant d'UI à la fenêtre, caché.
   *
   * @param component Le composant à ajouter à la fenêtre.
   * @param resize Si le composant doit être redimensionné à la taille de la fenêtre.
   */
  public void addUi(JComponent component, boolean resize) {
    checkFrameNotDisposed();
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
    checkFrameNotDisposed();
    component.setVisible(true);
    frame.getLayeredPane().setLayer(component, UI_VISIBLE_LAYER, 0);
  }

  /**
   * Cache un composant d'UI.
   *
   * @param component Le composant à cacher.
   */
  public void hideUi(JComponent component) {
    checkFrameNotDisposed();
    component.setVisible(false);
    frame.getLayeredPane().setLayer(component, UI_HIDDEN_LAYER, 0);
  }

  /**
   * Cache tous les composants d'UI.
   */
  public void hideAllUi() {
    checkFrameNotDisposed();
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
    checkFrameNotDisposed();
    this.renderable = renderable;
  }

  /**
   * Définit le listener de fermeture demandée de fenêtre. Appelé lors de la demande de fermeture de
   * la {@link UiWindow}.
   *
   * @param closeListener Le listener de fermeture.
   */
  public void setCloseRequestedListener(Runnable closeListener) {
    checkFrameNotDisposed();
    this.closeListener = closeListener;
  }

  /**
   * Définit le listener de clavier sur le composant principal.
   *
   * @param keyListener Le listener de clavier.
   */
  public void setKeyListener(KeyListener keyListener) {
    checkFrameNotDisposed();
    mainComponent.setKeyListener(keyListener);
  }

  /**
   * Définit le listener de touche pressée sur toute la fenêtre. Le listener sera appelé avec
   * l'entier renvoyé par {@link KeyEvent#getKeyCode()} pour chaque évènement de touche pressée.
   * <p>
   * Les évènements de touche pressée atteignant le composant principal seront à la fois envoyés à
   * ce listener et au listener principal défini par {@link #setKeyListener(KeyListener)}.
   *
   * @param listener Le listener de touche pressée.
   */
  public void setGlobalKeyDownListener(IntConsumer listener) {
    checkFrameNotDisposed();
    globalKeyDownListener = listener;
  }

  /**
   * @return La largeur de la fenêtre.
   */
  public int getWidth() {
    checkFrameNotDisposed();
    return width;
  }

  /**
   * @return La hauteur de la fenêtre.
   */
  public int getHeight() {
    checkFrameNotDisposed();
    return height;
  }

  /**
   * @return Le device sur lequel la frame est affichée.
   */
  public GraphicsDevice getDevice() {
    checkFrameNotDisposed();
    return device;
  }

  /**
   * @return true si la fenêtre est en plein écran exclusif.
   */
  public boolean isFullscreen() {
    checkFrameNotDisposed();
    return fullscreen;
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
    checkFrameNotDisposed();
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
    // render main component
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

  private void checkFrameNotDisposed() {
    if (frame == null) {
      throw new IllegalStateException("The window has already been disposed.");
    }
  }

}
