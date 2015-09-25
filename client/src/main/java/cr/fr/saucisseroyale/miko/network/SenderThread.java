package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.FutureOutputMessage;

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
        break;
      }
      try {
        fom.writeTo(dos);
      } catch (IOException e) {
        // TODO log cet évènement/notifier le networking
        e.printStackTrace();
        break;
      }
    }
  }

}
