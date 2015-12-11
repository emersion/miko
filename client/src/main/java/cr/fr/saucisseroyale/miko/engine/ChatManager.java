package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.util.Pair;

import java.util.Iterator;
import java.util.LinkedList;
import java.util.List;
import java.util.stream.Stream;

/**
 * Un gestionnaire des messages de chat reçus par le client.
 *
 */
class ChatManager {

  private static final long MESSAGES_DURATION = 20 * 1000000000L; // 20 seconds
  private List<Pair.Long<String>> messages = new LinkedList<>();

  /**
   * Ajoute un message à la liste des messages du chat.
   *
   * @param message Le message à ajouter.
   */
  public void addMessage(String message) {
    long time = System.nanoTime();
    messages.add(new Pair.Long<>(time, message));
  }

  /**
   * Renvoit un stream des messages du chat, dans l'ordre d'ajout (les messages trop anciens ne
   * seront pas renvoyés).
   *
   * @return Les messages récents.
   */
  public Stream<String> getMessages() {
    long oldMessageslimit = System.nanoTime() - MESSAGES_DURATION;
    Iterator<Pair.Long<String>> messagesIterator = messages.iterator();
    while (messagesIterator.hasNext()) {
      Pair.Long<String> message = messagesIterator.next();
      if (message.getFirst() >= oldMessageslimit) {
        break;
      }
      messagesIterator.remove();
    }
    return messages.stream().map(Pair.Long<String>::getSecond);
  }
}
