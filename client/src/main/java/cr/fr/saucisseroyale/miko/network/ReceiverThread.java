package cr.fr.saucisseroyale.miko.network;

import cr.fr.saucisseroyale.miko.FutureInputMessage;

import java.io.BufferedInputStream;
import java.io.DataInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.util.Queue;


/**
 * Thread pour recevoir les données d'un flux en boucle, de manière parallèle.
 * <p>
 * Cette classe supporte {@link #interrupt()}, pour terminer une instance de cette classe, appeler
 * cette méthode.
 *
 */
class ReceiverThread extends Thread {

  private DataInputStream dis;
  private Queue<FutureInputMessage> inputMessages;

  public ReceiverThread(InputStream is, Queue<FutureInputMessage> inputMessages) {
    dis = new DataInputStream(new BufferedInputStream(is));
    this.inputMessages = inputMessages;
  }

  @Override
  public void run() {
    while (true) {
      FutureInputMessage fim;
      try {
        fim = InputMessageFactory.parseMessage(dis);
      } catch (MessageParsingException e) {
        // Si l'on arrive ici, c'est sans doute qu'un parse a échoué et qu'on est dans un état
        // corrompu, ou qu'on un client maicieux envoit n'importe quoi
        e.printStackTrace();
        // TODO log cet évènement/notifier le networking
        break;
      } catch (IOException e) {
        e.printStackTrace();
        // TODO log cet évènement/notifier le networking
        break;
      }
      inputMessages.add(fim);
    }
  }

}
