package cr.fr.saucisseroyale.miko.network;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLSocketFactory;
import javax.net.ssl.TrustManagerFactory;
import java.io.IOException;
import java.io.InputStream;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.security.GeneralSecurityException;
import java.security.KeyStore;
import java.util.Queue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.LinkedBlockingQueue;

/**
 * Client de connexion à un serveur Miko fonctionnant sur la couche des messages.
 *
 * @see FutureInputMessage
 * @see FutureOutputMessage
 */
public class NetworkClient {
  private static Logger logger = LogManager.getLogger("miko.network");
  private SSLSocketFactory socketFactory;
  private Socket socket;
  private ReceiverThread receiverThread;
  private SenderThread senderThread;
  private String lastAddressString;
  private InetSocketAddress lastAddress;
  private Queue<FutureInputMessage> inputMessages = new ConcurrentLinkedQueue<>();
  private BlockingQueue<FutureOutputMessage> outputMessages = new LinkedBlockingQueue<>();

  public NetworkClient() {
    try {
      socketFactory = createSocketFactory();
      logger.debug("Created socket factory");
    } catch (Exception e) {
      logger.fatal("Failed creating socket factory while initating network client", e);
      throw new RuntimeException(e);
    }
  }

  private static SSLSocketFactory createSocketFactory() throws GeneralSecurityException, IOException {
    KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
    try (InputStream keyStoreStream = NetworkClient.class.getResourceAsStream("/keystore")) {
      keyStore.load(keyStoreStream, "keypass".toCharArray());
    }
    TrustManagerFactory tmf = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
    tmf.init(keyStore);
    SSLContext context = SSLContext.getInstance("TLS");
    context.init(null, tmf.getTrustManagers(), null);
    return context.getSocketFactory();
  }

  /**
   * Se connecte au serveur spécifié et démarre les threads d'envoi et de réception des messages.
   *
   * @param address L'adresse du serveur auquel se connecter (IP ou nom d'hôte).
   * @param port    Le port auquel se connecter.
   * @throws IOException S'il y a des erreurs quelconques d'IO lors de la connexion.
   */
  public void connect(String address, int port) throws IOException {
    logger.debug("Starting connection to server {} at port {}", address, port);
    socket = socketFactory.createSocket(address, port);
    socket.setTcpNoDelay(true);
    socket.setTrafficClass(0x10); // LOWDELAY
    receiverThread = new ReceiverThread(socket.getInputStream(), inputMessages, this::networkError);
    senderThread = new SenderThread(socket.getOutputStream(), outputMessages, this::networkError);
    receiverThread.start();
    senderThread.start();
    lastAddress = (InetSocketAddress) socket.getRemoteSocketAddress();
    lastAddressString = address;
    logger.info("Connected to server {} at port {}", address, port);
  }

  /**
   * Déconnecte le client du serveur.
   * <p>
   * Le client peut se reconnecter ensuite à un serveur avec {@link #connect(String, int)}. Si des
   * messages sont encores reçus ou ajoutés à la liste d'envoi, ils seront ignorés.
   */
  public void disconnect() {
    if (receiverThread != null) {
      receiverThread.interrupt();
      receiverThread = null;
    }
    if (senderThread != null) {
      senderThread.interrupt();
      senderThread = null;
    }
    if (socket != null) {
      logger.info("Disconnected from server");
      try {
        socket.close();
      } catch (IOException e) {
        // connection closing failed, ignore.
      }
      socket = null;
    }
  }

  /**
   * @return La dernière adresse de serveur auquel le client était connecté, ou null s'il n'a jamais été connecté.
   */
  public InetSocketAddress getLastAddress() {
    return lastAddress;
  }

  /**
   * @return La dernière adresse de serveur auquel le client était connecté, tel qu'il était entré, ou null s'il n'a jamais été connecté.
   */
  public String getLastAddressString() {
    return lastAddressString;
  }

  /**
   * Renvoie un message de la liste de réception des messages.
   *
   * @return Le message reçu le plus ancien, ou <code>null</code> si aucun message n'est en attente.
   */
  public FutureInputMessage getMessage() {
    return inputMessages.poll();
  }

  /**
   * Ajoute un message à la liste d'envoi des messages.
   *
   * @param fom Le message à envoyer.
   */
  public void putMessage(FutureOutputMessage fom) {
    outputMessages.add(fom);
  }

  private void networkError(Exception e) {
    inputMessages.add(InputMessageFactory.networkError(e));
    disconnect();
  }
}
