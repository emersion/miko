package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.engine.Engine;
import cr.fr.saucisseroyale.miko.network.FutureInputMessage;
import cr.fr.saucisseroyale.miko.network.NetworkClient;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;
import cr.fr.saucisseroyale.miko.protocol.Action;
import cr.fr.saucisseroyale.miko.protocol.ChunkPoint;
import cr.fr.saucisseroyale.miko.protocol.Config;
import cr.fr.saucisseroyale.miko.protocol.EntityDataUpdate;
import cr.fr.saucisseroyale.miko.protocol.ExitType;
import cr.fr.saucisseroyale.miko.protocol.LoginResponseType;
import cr.fr.saucisseroyale.miko.protocol.RegisterResponseType;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.awt.Graphics2D;
import java.awt.Point;
import java.awt.Toolkit;
import java.io.IOException;
import java.util.List;

public class Miko implements MessageHandler {

  private enum MikoState {
    NETWORK, CONNECTION_REQUEST, CONNECTION, CONFIG, LOGIN_REQUEST, REGISTER, LOGIN, JOIN, EXIT;
  }

  public static final int PROTOCOL_VERSION = 7;
  private static final String DEFAULT_SERVER_ADDRESS = "127.0.0.1";
  private static final int DEFAULT_SERVER_PORT = 9999;
  public static final int TICK_TIME = 20; // milliseconds
  private static final int SERVER_TIMEOUT = 10 * 1000000000; // seconds
  private long lastMessageReceived = Long.MAX_VALUE;
  private boolean pingSent;
  private String username;
  private Config config;
  private Engine engine;
  private UiComponents.Connect uiConnect;
  private UiComponents.Login uiLogin;
  private MikoState state = MikoState.NETWORK;
  private UiWindow window;
  private boolean closeRequested = false;
  private NetworkClient networkClient;
  private KeyStateManager keyStateManager;
  private float alpha; // for drawing, updated each loop

  private void exit() {
    if (window != null) {
      window.close();
    }
    if (networkClient != null) {
      networkClient.disconnect();
    }
    System.exit(0);
  }

  private void initWindow() {
    // TODO uncomment this when stylesheet design is done
    // SynthLookAndFeel lookAndFeel = new SynthLookAndFeel();
    // try {
    // lookAndFeel.load(Miko.class.getResourceAsStream("/style.xml"), Miko.class);
    // } catch (ParseException e) {
    // // should never happen on normal releases
    // throw new RuntimeException(e);
    // }
    // try {
    // UIManager.setLookAndFeel(lookAndFeel);
    // } catch (UnsupportedLookAndFeelException e) {
    // // will never be thrown, SynthLookAndFeel#isSupportedLookAndFeel() never returns false
    // throw new RuntimeException(e);
    // }

    window = new UiWindow();
    window.setRenderable(this::render);
    uiConnect = new UiComponents.Connect(DEFAULT_SERVER_ADDRESS, DEFAULT_SERVER_PORT, this::connectRequested);
    uiLogin = new UiComponents.Login(this::registerRequested, this::loginRequested);
    window.addUi(uiConnect);
    window.addUi(uiLogin);

    keyStateManager = new KeyStateManager(window::getMousePosition);
    window.setKeyListener(keyStateManager);

    window.initAndShow();
  }

  private void logic() {
    if (state != MikoState.EXIT) {
      return;
    }
    engine.processNextTick(keyStateManager.getEventsAndFlush(), window.getMousePosition());
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
      postLoop();
      Toolkit.getDefaultToolkit().sync(); // vsync
    }
    exit();
  }

  private void network() {
    if (state == MikoState.NETWORK || state == MikoState.CONNECTION_REQUEST || state == MikoState.CONNECTION) {
      return;
    }
    FutureInputMessage fim = networkClient.getMessage();
    while (fim != null) {
      fim.execute(this);
      fim = networkClient.getMessage();
    }
  }

  private void render(Graphics2D graphics) {
    if (state != MikoState.EXIT) {
      return;
    }
    Point mousePosition = window.getMousePosition();
    engine.render(graphics, alpha, mousePosition);
  }

  private void postLoop() {
    if (state != MikoState.NETWORK && state != MikoState.CONNECTION_REQUEST && state != MikoState.CONNECTION) {
      long currentTime = System.nanoTime();
      if (currentTime - lastMessageReceived > SERVER_TIMEOUT) {
        exit(ExitType.PING_TIMEOUT);
      }
      if (currentTime - lastMessageReceived > SERVER_TIMEOUT / 2 && !pingSent) {
        pingSent = true;
        networkClient.putMessage(OutputMessageFactory.ping());
      }
    }
    if (state == MikoState.EXIT) {
      engine.freeTime();
    }
  }

  private void connect(String address, int port) {
    try {
      networkClient.connect(address, port);
    } catch (IOException e) {
      networkClient.disconnect();
      uiConnect.setStatusText("Erreur de connexion : erreur d'établissement de connexion inconnue.");
      changeStateTo(MikoState.CONNECTION_REQUEST);
      return;
    }
    changeStateTo(MikoState.CONFIG);
    networkClient.putMessage(OutputMessageFactory.version());
  }

  private void connectRequested(String address, int port) {
    if (state != MikoState.CONNECTION_REQUEST) {
      return;
    }
    changeStateTo(MikoState.CONNECTION);
    // use thread to avoid blocking swing event dispatching for too long
    new Thread(() -> connect(address, port)).start();
  }

  private void loginRequested(String username, String password) {
    if (state != MikoState.LOGIN_REQUEST) {
      return;
    }
    changeStateTo(MikoState.LOGIN);
    this.username = username;
    networkClient.putMessage(OutputMessageFactory.login(username, password));
  }

  private void registerRequested(String username, String password) {
    if (state != MikoState.LOGIN_REQUEST) {
      return;
    }
    changeStateTo(MikoState.REGISTER);
    networkClient.putMessage(OutputMessageFactory.register(username, password));
  }

  private void changeStateTo(MikoState newState) {
    state = newState;
    window.hideAllUi();
    switch (state) {
      case CONNECTION_REQUEST:
        window.showUi(uiConnect);
        break;
      case CONNECTION:
        break;
      case CONFIG:
        break;
      case LOGIN_REQUEST:
        window.showUi(uiLogin);
        break;
      case LOGIN:
        break;
      case REGISTER:
        break;
      case JOIN:
        // TODO montrer chargement
        break;
      case EXIT:
        break;
      default:
        break;
    }
  }

  private void run() throws Exception {
    initWindow();
    networkClient = new NetworkClient();
    changeStateTo(MikoState.CONNECTION_REQUEST);
    loop();
  }

  public static void main(String... args) throws Exception {
    Miko miko = new Miko();
    miko.run();
  }

  @Override
  public void actions(int tickRemainder, List<Pair<Integer, Action>> actions) {
    messageReceived();
    if (state != MikoState.EXIT) {
      return;
    }
    engine.actions(tickRemainder, actions);
  }

  @Override
  public void chatReceived(int tickRemainder, int entityIdChat, String chatMessage) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    engine.chatReceived(tickRemainder, entityIdChat, chatMessage);
  }

  @Override
  public void chunkUpdate(int tickRemainder, ChunkPoint chunkPoint, Chunk chunk) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    engine.chunkUpdate(tickRemainder, chunkPoint, chunk);
  }

  @Override
  public void entityIdChange(int oldEntityId, int newEntityId) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    engine.entityIdChange(oldEntityId, newEntityId);
  }

  @Override
  public void entityCreate(int tickRemainder, EntityDataUpdate entityDataUpdate) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    engine.entityCreate(tickRemainder, entityDataUpdate);
  }

  @Override
  public void entityDestroy(int tickRemainder, int entityId) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    engine.entityDestroy(tickRemainder, entityId);
  }

  @Override
  public void entitiesUpdate(int tickRemainder, List<EntityDataUpdate> entitiesUpdateList) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    engine.entitiesUpdate(tickRemainder, entitiesUpdateList);
  }

  @Override
  public void exit(ExitType exitType) {
    networkClient.disconnect();
    String statusMessage;
    switch (exitType) {
      case CLIENT_BANNED:
        statusMessage = "Déconnexion : client banni.";
        break;
      case CLIENT_KICKED:
        statusMessage = "Déconnexion : client kické.";
        break;
      case CLIENT_QUIT:
        statusMessage = "Déconnexion : client fermé.";
        break;
      case NETWORK_ERROR:
        statusMessage = "Déconnexion : erreur de réseau.";
        break;
      case PING_TIMEOUT:
        statusMessage = "Déconnexion : timeout.";
        break;
      case SERVER_CLOSED:
        statusMessage = "Déconnexion : serveur fermé.";
        break;
      case CLIENT_OUTDATED:
        statusMessage = "Déconnexion : client obsolète.";
        break;
      case SERVER_OUTDATED:
        statusMessage = "Déconnexion : serveur obsolète.";
        break;
      default:
        statusMessage = "Déconnexion : cause de déconnexion inconnue.";
        break;
    }
    uiConnect.setStatusText(statusMessage);
    changeStateTo(MikoState.CONNECTION_REQUEST);
  }

  @Override
  public void loginFail(LoginResponseType loginResponseType) {
    messageReceived();
    if (state != MikoState.LOGIN) {
      return;
    }
    String statusMessage;
    switch (loginResponseType) {
      case ALREADY_CONNECTED:
        statusMessage = "Erreur de connexion : utilisateur déjà connecté.";
        break;
      case PLAYER_LIMIT_REACHED:
        statusMessage = "Erreur de connexion : serveur plein.";
        break;
      case TOO_MANY_TRIES:
        statusMessage = "Erreur de connexion : trop d'essais de connexions.";
        break;
      case UNKNOWN_PSEUDO:
        statusMessage = "Erreur de connexion : nom d'utilisateur inconnu.";
        break;
      case WRONG_PASSWORD:
        statusMessage = "Erreur de connexion : mot de passe erronné.";
        break;
      default:
        statusMessage = "Erreur de connexion : erreur de connexion inconnue.";
        break;
    }
    uiLogin.setStatusText(statusMessage);
    changeStateTo(MikoState.LOGIN_REQUEST);
  }

  @Override
  public void loginSuccess(int tickRemainder) {
    messageReceived();
    if (state != MikoState.LOGIN) {
      return;
    }
    changeStateTo(MikoState.JOIN);
    try {
      engine =
          new Engine(config, window.getDevice().getDefaultConfiguration(), networkClient::putMessage, window.getWidth(), window.getHeight(),
              tickRemainder);
    } catch (IOException e) {
      uiConnect.setStatusText("Erreur lors de la création du jeu : erreur de lecture de données.");
      changeStateTo(MikoState.CONNECTION_REQUEST);
    }
  }

  @Override
  public void networkError(Exception e) {
    networkClient.disconnect();
    uiConnect.setStatusText("Déconnexion forcée : erreur de réseau : " + e.getLocalizedMessage());
    changeStateTo(MikoState.CONNECTION_REQUEST);
  }

  @Override
  public void ping() {
    messageReceived();
    networkClient.putMessage(OutputMessageFactory.pong());
  }

  @Override
  public void playerJoined(int tickRemainder, int entityId, String pseudo) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      return;
    }
    if (state == MikoState.JOIN && pseudo.equals(username)) {
      engine.setPlayerEntityId(entityId);
      engine.playerJoined(tickRemainder, entityId, pseudo);
      changeStateTo(MikoState.EXIT);
    } else {
      engine.playerJoined(tickRemainder, entityId, pseudo);
    }
  }

  @Override
  public void playerLeft(int tickRemainder, int entityId) {
    messageReceived();
    if (state != MikoState.EXIT) {
      return;
    }
    engine.playerLeft(tickRemainder, entityId);
  }

  @Override
  public void pong() {
    messageReceived();
    pingSent = false;
  }

  @Override
  public void registerResponse(RegisterResponseType registerResponseType) {
    messageReceived();
    if (state != MikoState.REGISTER) {
      return;
    }
    String statusMessage;
    switch (registerResponseType) {
      case OK:
        statusMessage = "Enregistré avec succès.";
        break;
      case INVALID_PASSWORD:
        statusMessage = "Erreur d'enregistrement : mot de passe invalide.";
        break;
      case INVALID_PSEUDO:
        statusMessage = "Erreur d'enregistrement : nom d'utilisateur invalide.";
        break;
      case REGISTER_DISABLED:
        statusMessage = "Erreur d'enregistrement : enregistrement désactivé.";
        break;
      case TOO_MANY_TRIES:
        statusMessage = "Erreur d'enregistrement : trop d'essais d'enregistrements.";
        break;
      case USED_PSEUDO:
        statusMessage = "Erreur d'enregistrement : nom d'utilisateur déjà utilisé.";
        break;
      default:
        statusMessage = "Erreur d'enregistrement : erreur d'enregistrement inconnue.";
        break;
    }
    uiLogin.setStatusText(statusMessage);
    changeStateTo(MikoState.LOGIN_REQUEST);
  }

  @Override
  public void config(Config config) {
    messageReceived();
    if (state != MikoState.CONFIG) {
      return;
    }
    this.config = config;
    uiLogin.setStatusText("Connexion réussie, connectez ou enregistrez-vous.");
    changeStateTo(MikoState.LOGIN_REQUEST);
  }

  private void messageReceived() {
    // called when any message is received
    lastMessageReceived = System.nanoTime();
  }

}
