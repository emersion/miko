package cr.fr.saucisseroyale.miko.network;

import java.io.BufferedOutputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.util.concurrent.BlockingQueue;
import java.util.function.Consumer;

/**
 * Thread pour émettre des données sur un flux en boucle, de manière parallèle.
 * <p>
 * Cette classe supporte {@link #interrupt()}, pour terminer une instance de cette classe, appeler
 * cette méthode.
 */
class SenderThread extends Thread {
  private DataOutputStream dos;
  private BlockingQueue<FutureOutputMessage> outputMessages;
  private Consumer<Exception> errorCallback;

  public SenderThread(OutputStream os, BlockingQueue<FutureOutputMessage> outputMessages, Consumer<Exception> errorCallback) {
    dos = new DataOutputStream(new BufferedOutputStream(os));
    this.outputMessages = outputMessages;
    this.errorCallback = errorCallback;
    setName("Miko Network Sender");
  }

  @Override
  public void run() {
    while (true) {
      FutureOutputMessage fom;
      try {
        fom = outputMessages.take();
      } catch (InterruptedException e) {
        // close requested, quit
        return;
      }
      try {
        fom.writeTo(dos);
        dos.flush();
      } catch (IOException e) {
        errorCallback.accept(e);
        return;
      }
    }
  }
}
