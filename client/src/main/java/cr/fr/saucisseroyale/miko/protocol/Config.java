package cr.fr.saucisseroyale.miko.protocol;

/**
 * Une configuration immutable envoyée par le serveur pendant la connexion.
 *
 */
public class Config {

  private final int maxRollbackTicks;

  public Config(int maxRollbackTicks) {
    this.maxRollbackTicks = maxRollbackTicks;
  }

  /**
   * @return Le nombre maximum de ticks pour le rollback.
   */
  public int getMaxRollbackTicks() {
    return maxRollbackTicks;
  }

}
