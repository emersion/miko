package cr.fr.saucisseroyale.miko.network;

import java.io.BufferedOutputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.util.concurrent.BlockingQueue;


/**
 * Thread pour émettre des données sur un flux en boucle, de manière parallèle.
 * <p>
 * Cette classe supporte {@link #interrupt()}, pour terminer une instance de cette classe, appeler
 * cette méthode.
 *
 */
class SenderThread extends Thread {

  private DataOutputStream dos;
  private BlockingQueue<FutureOutputMessage> outputMessages;

  public SenderThread(OutputStream os, BlockingQueue<FutureOutputMessage> outputMessages) {
    dos = new DataOutputStream(new BufferedOutputStream(os));
    this.outputMessages = outputMessages;
  }

  @Override
  public void run() {
    while (true) {
      FutureOutputMessage fom;
      try {
        fom = outputMessages.take();
      } catch (InterruptedException e) {
        // On a demandé notre interruption, quitter
        // TODO enlever le stack trace
        e.printStackTrace();
        break;
      }
      try {
        fom.writeTo(dos);
        dos.flush();
      } catch (IOException e) {
        // TODO log cet évènement/notifier le networking
        e.printStackTrace();
        break;
      }
    }
  }

}
