package cr.fr.saucisseroyale.miko.network;

import java.io.IOException;
import java.net.Socket;
import java.util.Queue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.LinkedBlockingQueue;

import javax.net.ssl.SSLSocketFactory;

/**
 * Client de connexion à un serveur Miko fonctionnant sur la couche des messages.
 *
 * @see FutureInputMessage
 * @see FutureOutputMessage
 *
 */
public class NetworkClient {

  private Socket socket;
  private ReceiverThread receiverThread;
  private SenderThread senderThread;
  private Queue<FutureInputMessage> inputMessages = new ConcurrentLinkedQueue<>();
  private BlockingQueue<FutureOutputMessage> outputMessages = new LinkedBlockingQueue<>();

  /**
   * Se connecte au serveur spécifié et démarre les threads d'envoi et de réception des messages.
   *
   * @param address L'adresse du serveur auquel se connecter (IP ou nom d'hôte).
   * @param port Le port auquel se connecter.
   * @throws IOException S'il y a des erreurs quelconques d'IO lors de la connexion.
   */
  public void connect(String address, int port) throws IOException {
    socket = SSLSocketFactory.getDefault().createSocket(address, port);
    receiverThread = new ReceiverThread(socket.getInputStream(), inputMessages);
    senderThread = new SenderThread(socket.getOutputStream(), outputMessages);
  }


  /**
   * Termine définitivement la connection et l'envoi et la réception des messages.
   *
   * @throws IOException S'il y a des erreurs quelconques d'IO lors de l'extinction.
   */
  public void exit() throws IOException {
    receiverThread.interrupt();
    senderThread.interrupt();
    socket.close();
  }

  /**
   * Renvoie un message de la liste de réception des messages.
   *
   * @return Le message reçu le plus ancien, ou <code>null</code> si aucun message n'est en attente.
   *
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

}
