package cr.fr.saucisseroyale.miko;


public interface FutureInputMessage {

  public void execute(Miko miko);
  // TODO Peut-être que Miko n'est pas le plus adapté, peut-être une classe faite pour serait mieux,
  // ou alors en plus rajouter un moyen d'envoyer des paquets de manière brute
}
