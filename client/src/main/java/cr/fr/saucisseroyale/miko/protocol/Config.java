package cr.fr.saucisseroyale.miko.protocol;

/**
 * Une configuration immutable envoyée par le serveur pendant la connexion.
 */
public class Config {
  private final int maxRollbackTicks;
  private final int timeServerPort;
  private final float defaultPlayerSpeed;
  private final int playerBallCooldown;
  private final float defaultBallSpeed;
  private final int defaultBallLifespan;

  public Config(int maxRollbackTicks, int timeServerPort, float defaultPlayerSpeed, int playerBallCooldown, float defaultBallSpeed, int defaultBallLifespan) {
    this.maxRollbackTicks = maxRollbackTicks;
    this.timeServerPort = timeServerPort;
    this.defaultPlayerSpeed = defaultPlayerSpeed;
    this.playerBallCooldown = playerBallCooldown;
    this.defaultBallSpeed = defaultBallSpeed;
    this.defaultBallLifespan = defaultBallLifespan;
  }

  /**
   * @return Le nombre maximum de ticks pour le rollback.
   */
  public int getMaxRollbackTicks() {
    return maxRollbackTicks;
  }

  /**
   * @return Le port du serveur de temps.
   */
  public int getTimeServerPort() {
    return timeServerPort;
  }

  /**
   * @return La vitesse par défault du joueur.
   */
  public float getDefaultPlayerSpeed() {
    return defaultPlayerSpeed;
  }

  /**
   * @return Le cooldown de lancer de boules du joueur.
   */
  public int getPlayerBallCooldown() {
    return playerBallCooldown;
  }

  /**
   * @return La vitesse par défaut d'une boule.
   */
  public float getDefaultBallSpeed() {
    return defaultBallSpeed;
  }

  /**
   * @return La durée de vie par défaut d'une boule.
   */
  public int getDefaultBallLifespan() {
    return defaultBallLifespan;
  }
}
