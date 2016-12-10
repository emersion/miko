package cr.fr.saucisseroyale.miko.network;

import java.io.BufferedInputStream;
import java.io.DataInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Queue;
import java.util.function.Consumer;

/**
 * Thread pour recevoir les données d'un flux en boucle, de manière parallèle.
 * <p>
 * Cette classe supporte {@link #interrupt()}, pour terminer une instance de cette classe, appeler
 * cette méthode.
 */
class ReceiverThread extends Thread {
  private DataInputStream dis;
  private Queue<FutureInputMessage> inputMessages;
  private Consumer<Exception> errorCallback;

  public ReceiverThread(InputStream is, Queue<FutureInputMessage> inputMessages, Consumer<Exception> errorCallback) {
    dis = new DataInputStream(new BufferedInputStream(is));
    this.inputMessages = inputMessages;
    this.errorCallback = errorCallback;
    setName("Miko Network Receiver");
  }

  @Override
  public void run() {
    while (true) {
      FutureInputMessage fim;
      try {
        fim = InputMessageFactory.parseMessage(dis);
        inputMessages.add(fim);
      } catch (IOException e) {
        errorCallback.accept(e);
        return;
      }
    }
  }
}
