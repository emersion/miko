package cr.fr.saucisseroyale.miko;

import cr.fr.saucisseroyale.miko.UiComponents.Connect;
import cr.fr.saucisseroyale.miko.UiComponents.Login;
import cr.fr.saucisseroyale.miko.UiComponents.MikoLayer;
import cr.fr.saucisseroyale.miko.UiComponents.Options;
import cr.fr.saucisseroyale.miko.engine.Chunk;
import cr.fr.saucisseroyale.miko.engine.Engine;
import cr.fr.saucisseroyale.miko.network.FutureInputMessage;
import cr.fr.saucisseroyale.miko.network.NetworkClient;
import cr.fr.saucisseroyale.miko.network.OutputMessageFactory;
import cr.fr.saucisseroyale.miko.network.TimeClient;
import cr.fr.saucisseroyale.miko.protocol.*;
import cr.fr.saucisseroyale.miko.util.Pair;
import cr.fr.saucisseroyale.miko.util.Pair.Int;
import fr.delthas.uitest.Drawer;
import fr.delthas.uitest.Icon;
import fr.delthas.uitest.InputState;
import fr.delthas.uitest.Ui;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.awt.*;
import java.io.IOException;
import java.util.List;
import java.util.prefs.Preferences;

public class Miko implements MessageHandler {
  public static final int PROTOCOL_VERSION = 10;
  public static final int TICK_TIME = 20 * 1000000; // milliseconds
  private static final long SERVER_TIMEOUT = 20 * 1000000000L; // seconds
  private static final String DEFAULT_SERVER_ADDRESS = "localhost";
  private static final int DEFAULT_SERVER_PORT = 9999;
  private static final Preferences uiPrefsNode = Preferences.userRoot().node("miko.ui");
  private static Logger logger = LogManager.getLogger("miko.main");
  private static Connect connect;
  private static Login login;
  private static Options options;
  private static MikoLayer mikoLayer;
  private long lastMessageReceived = Long.MAX_VALUE;
  private boolean pingSent;
  private String username;
  private Engine engine;
  private MikoState state = MikoState.NETWORK;
  private boolean closeRequested = false;
  private NetworkClient networkClient;
  private TimeClient timeClient;
  private InputStateManager inputStateManager = new InputStateManager();
  private Config config;
  private long accumulator;
  @SuppressWarnings("FieldCanBeLocal")
  private long lastFrame;
  private float alpha; // for #render(), updated each loop

  public static void main(String... args) throws Exception {
    if (args.length == 1) {
      if (args[0].equalsIgnoreCase("fs")) {
        uiPrefsNode.putBoolean("fullscreen", true);
      } else if (args[0].equalsIgnoreCase("nfs")) {
        uiPrefsNode.putBoolean("fullscreen", false);
      }
    }
    logger.info("Starting Miko version {}", PROTOCOL_VERSION);
    Miko miko = new Miko();
    miko.run();
  }

  private void exit() {
    logger.info("Starts exiting");
    Ui.getUi().destroy();
    if (networkClient != null) {
      networkClient.disconnect();
    }
    System.exit(0);
  }

  private void disconnect() {
    networkClient.disconnect();
    timeClient.disconnect();
    lastMessageReceived = Long.MAX_VALUE;
    changeStateTo(MikoState.CONNECTION_REQUEST);
  }

  private void initWindow() throws IOException {

    Ui.getUi().create("Miko", Icon.createIcon("icon.png"));
    connect = new Connect(this::backRequested, DEFAULT_SERVER_ADDRESS, DEFAULT_SERVER_PORT, this::connectRequested, this::optionsRequested);
    login = new Login(this::backRequested, this::registerRequested, this::loginRequested);
    options = new Options(this::backRequested, true);
    mikoLayer = new MikoLayer(this::backRequested, this::render, inputStateManager);

    // boolean fullscreen = uiPrefsNode.getBoolean("fullscreen", false);
  }

  private void logic() {
    switch (state) {
      case LOGIN_REQUEST:
        login.setLoginEnabled(timeClient.getClockDifference().isPresent());
        break;
      case EXIT:
        while (accumulator >= TICK_TIME) {
          engine.processNextTick(inputStateManager.getEventsAndFlush());
          accumulator -= TICK_TIME;
        }
        break;
      default:
    }
  }

  private void loop() {
    lastFrame = System.nanoTime();
    accumulator = 0;
    while (!closeRequested) {
      long newTime = System.nanoTime();
      long deltaTime = newTime - lastFrame;
      lastFrame = newTime;
      accumulator += deltaTime;
      network();
      Ui.getUi().input();
      logic();
      alpha = (float) accumulator / TICK_TIME; // update alpha for #render()
      Ui.getUi().render(); // calls render(graphics)
      postLoop();
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

  private void render(InputState inputState, Drawer drawer) {
    if (state != MikoState.EXIT) {
      return;
    }
    engine.render(drawer, alpha, new Point.Double(inputState.getMouseX(), inputState.getMouseY()));
  }

  private void postLoop() {
    if (state != MikoState.NETWORK && state != MikoState.CONNECTION_REQUEST && state != MikoState.CONNECTION) {
      long currentTime = System.nanoTime();
      if (currentTime - lastMessageReceived > SERVER_TIMEOUT) {
        logger.error("Ping timeout");
        exit(ExitType.PING_TIMEOUT);
      }
      if (currentTime - lastMessageReceived > SERVER_TIMEOUT / 2 && !pingSent) {
        logger.warn("Sent ping to avoid timeout (last message received {} millis ago)", (currentTime - lastMessageReceived) / 1000000);
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
      logger.info("Connecting to {} at port {}", address, port);
      networkClient.connect(address, port);
    } catch (IOException e) {
      logger.warn("Connection error, disconnecting");
      connect.setStatusText("Erreur d'établissement de connexion: " + e.getClass().getCanonicalName() + ": " + e.getLocalizedMessage());
      disconnect();
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
    logger.info("Logging in as {}", username);
    networkClient.putMessage(OutputMessageFactory.login(username, password));
  }

  private void registerRequested(String username, String password) {
    if (state != MikoState.LOGIN_REQUEST) {
      return;
    }
    changeStateTo(MikoState.REGISTER);
    logger.info("Registering {}", username);
    networkClient.putMessage(OutputMessageFactory.register(username, password));
  }

  private void optionsRequested() {
    if (state != MikoState.CONNECTION_REQUEST) {
      return;
    }
    changeStateTo(MikoState.OPTIONS);
  }

  private void backRequested() {
    switch (state) {
      case NETWORK:
      case CONNECTION_REQUEST:
        closeRequested = true;
        break;
      case OPTIONS:
        changeStateTo(MikoState.CONNECTION_REQUEST);
        break;
      case EXIT:
        // send exit message...
        networkClient.putMessage(OutputMessageFactory.exit(ExitType.CLIENT_QUIT));
        //$FALL-THROUGH$
      case CONNECTION:
      case JOIN:
      case LOGIN:
      case LOGIN_REQUEST:
      case REGISTER:
      case CONFIG:
        // ...disconnect and return to connection panel
        logger.warn("Exited due to user request");
        connect.setStatusText("Déconnecté par l'utilisateur");
        disconnect();
        break;
      default:
        // ignore
    }
  }

  private void changeStateTo(MikoState newState) {
    logger.trace("Changing state from {} to {}", state, newState);
    state = newState;
    Ui.getUi().pop();
    switch (state) {
      case CONNECTION_REQUEST:
        Ui.getUi().push(connect);
        break;
      case OPTIONS:
        Ui.getUi().push(options);
        break;
      case CONNECTION:
        break;
      case CONFIG:
        break;
      case LOGIN_REQUEST:
        Ui.getUi().push(login);
        break;
      case LOGIN:
        break;
      case REGISTER:
        break;
      case JOIN:
        break;
      case EXIT:
        Ui.getUi().push(mikoLayer);
        break;
      default:
        break;
    }
  }

  private void run() throws IOException {
    logger.debug("Initializing window");
    initWindow();
    networkClient = new NetworkClient();
    timeClient = new TimeClient();
    changeStateTo(MikoState.CONNECTION_REQUEST);
    logger.debug("Starting game loop");
    loop();
  }

  @Override
  public void actions(int tickRemainder, List<Int<Action>> actions) {
    messageReceived();
    if (state != MikoState.EXIT) {
      logger.warn("Ignored actions message received in state {}", state);
      return;
    }
    engine.actions(tickRemainder, actions);
  }

  @Override
  public void chatReceived(int tickRemainder, int entityIdChat, String chatMessage) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      logger.warn("Ignored chatreceived message received in state {}", state);
      return;
    }
    engine.chatReceived(tickRemainder, entityIdChat, chatMessage);
  }

  @Override
  public void chunksUpdate(int tickRemainder, List<Pair<ChunkPoint, Chunk>> chunks) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      logger.warn("Ignored chunksupdate message received in state {}", state);
      return;
    }
    engine.chunksUpdate(tickRemainder, chunks);
  }

  @Override
  public void entityIdChange(int oldEntityId, int newEntityId) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      logger.warn("Ignored entityidchange message received in state {}", state);
      return;
    }
    engine.entityIdChange(oldEntityId, newEntityId);
  }

  @Override
  public void entityCreate(int tickRemainder, EntityDataUpdate entityDataUpdate) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      logger.warn("Ignored entitycreate message received in state {}", state);
      return;
    }
    engine.entityCreate(tickRemainder, entityDataUpdate);
  }

  @Override
  public void entityDestroy(int tickRemainder, int entityId) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      logger.warn("Ignored entitydestroy message received in state {}", state);
      return;
    }
    engine.entityDestroy(tickRemainder, entityId);
  }

  @Override
  public void entitiesUpdate(int tickRemainder, List<EntityDataUpdate> entitiesUpdateList) {
    messageReceived();
    if (state != MikoState.JOIN && state != MikoState.EXIT) {
      logger.warn("Ignored entitiesupdate message received in state {}", state);
      return;
    }
    engine.entitiesUpdate(tickRemainder, entitiesUpdateList);
  }

  @Override
  public void exit(ExitType exitType) {
    logger.warn("Exited because of {} exit message", exitType);
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
    connect.setStatusText(statusMessage);
    disconnect();
  }

  @Override
  public void loginFail(LoginResponseType loginResponseType) {
    messageReceived();
    if (state != MikoState.LOGIN) {
      logger.warn("Ignored loginfail message received in state {}", state);
      return;
    }
    logger.warn("Login failed with type {}", loginResponseType);
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
    login.setStatusText(statusMessage);
    changeStateTo(MikoState.LOGIN_REQUEST);
  }

  @Override
  public void loginSuccess(int tickRemainder, long timestamp) {
    messageReceived();
    if (state != MikoState.LOGIN) {
      logger.warn("Ignored loginsuccess message received in state {}", state);
      return;
    }
    long clockDifference = timeClient.getClockDifference().get();
    timeClient.disconnect();
    long tickLocalStartTimestamp = timestamp + clockDifference;
    logger.info("Login success, starting engine at tick {}, timestamp {}", tickRemainder, tickLocalStartTimestamp);
    long currentTime = System.nanoTime();
    accumulator = currentTime - tickLocalStartTimestamp * 1000;
    lastFrame = currentTime;
    changeStateTo(MikoState.JOIN);
    try {
      engine = new Engine(config, networkClient::putMessage, tickRemainder);
    } catch (IOException e) {
      connect.setStatusText("Erreur lors de la création du jeu : erreur de lecture de données. " + e.getMessage());
      disconnect();
    }
  }

  @Override
  public void networkError(Exception e) {
    if (state == MikoState.NETWORK || state == MikoState.CONNECTION_REQUEST || state == MikoState.OPTIONS) {
      logger.warn("Ignored network error received in state {}", state);
      return;
    }
    logger.error("Network error, disconnecting", e);
    connect.setStatusText("Déconnexion forcée : erreur de réseau : " + e.getClass().getCanonicalName() + (e.getLocalizedMessage() != null ? ": " + e.getLocalizedMessage() : ""));
    disconnect();
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
      logger.warn("Ignored playerjoined message received in state {}", state);
      return;
    }
    if (state == MikoState.JOIN && pseudo.equals(username)) {
      logger.info("Received self player join, finished initializing engine, starting game");
      engine.setPlayerEntityId(entityId);
      engine.playerJoined(tickRemainder, entityId, pseudo);
      long ticksSinceStartup = accumulator / TICK_TIME;
      accumulator %= TICK_TIME;
      engine.endStartup(ticksSinceStartup);
      changeStateTo(MikoState.EXIT);
    } else {
      engine.playerJoined(tickRemainder, entityId, pseudo);
    }
  }

  @Override
  public void playerLeft(int tickRemainder, int entityId) {
    messageReceived();
    if (state != MikoState.EXIT) {
      logger.warn("Ignored playerleft message received in state {}", state);
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
      logger.warn("Ignored register message received in state {}", state);
      return;
    }
    if (registerResponseType == RegisterResponseType.OK) {
      logger.info("Register succeeded");
    } else {
      logger.warn("Register failed, reason: {}", registerResponseType);
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
    login.setStatusText(statusMessage);
    changeStateTo(MikoState.LOGIN_REQUEST);
  }

  @Override
  public void config(Config config) {
    messageReceived();
    if (state != MikoState.CONFIG) {
      logger.warn("Ignored config message received in state {}", state);
      return;
    }
    logger.debug("Received server config");
    this.config = config;
    login.setStatusText("Connexion réussie, connectez ou enregistrez-vous.");
    changeStateTo(MikoState.LOGIN_REQUEST);
    timeClient.connect(networkClient.getLastAddress().getAddress().getHostAddress(), config.getTimeServerPort());
  }

  private void messageReceived() {
    // called when any message is received
    lastMessageReceived = System.nanoTime();
  }

  private enum MikoState {
    NETWORK, CONNECTION_REQUEST, OPTIONS, CONNECTION, CONFIG, LOGIN_REQUEST, REGISTER, LOGIN, JOIN, EXIT
  }
}
