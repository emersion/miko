package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.Miko;
import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.awt.image.BufferedImage;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

/**
 * Une animation d'entité. Chaque entité ne possède qu'un {@link SpriteType} et cet objet est donc
 * partagé et utilisé par {@link SpriteManager}.
 *
 */
public class Sprite {

  private final int spriteTimeDivider;
  private final List<Pair.Int<BufferedImage>> frames;

  /**
   * Crée une animation avec une seule image.
   *
   * @param image L'unique image de l'animation.
   */
  public Sprite(BufferedImage image) {
    frames = new ArrayList<>(1);
    frames.add(new Pair.Int<>(1, image));
    spriteTimeDivider = 1;
  }

  /**
   * Crée une animation avec plusieurs images, représentées par (timecode;image) avec timecode le
   * temps en millisecondes auquel on passera à l'image suivante.
   *
   * @param frames La liste des images avec leur timecode associé.
   */
  public Sprite(List<Pair.Int<BufferedImage>> frames) {
    this.frames = frames.stream().sorted((p1, p2) -> {
      if (p1.getFirst() < p2.getFirst()) {
        return -1;
      }
      if (p1.getFirst() > p2.getFirst()) {
        return +1;
      }
      return 0;
    }).collect(Collectors.toCollection(ArrayList::new));
    if (this.frames.get(0).getFirst() < 0) {
      throw new IllegalArgumentException("Timecodes must be positive");
    }
    spriteTimeDivider = this.frames.get(frames.size() - 1).getFirst();
  }

  /**
   * Retourne l'image de l'animation correspondant au temps depuis le choix de l'animation.
   *
   * @param spriteTime Le temps en ticks depuis que l'animation a été choisie.
   * @return L'image correspondant
   */
  public BufferedImage getImage(long spriteTime) {
    // convert ticks to milliseconds!
    long spriteTimeMillis = spriteTime * Miko.TICK_TIME / 1000000;
    int remainder = (int) (spriteTimeMillis % spriteTimeDivider);
    for (Pair.Int<BufferedImage> frame : frames) {
      if (remainder < frame.getFirst()) {
        return frame.getSecond();
      }
    }
    // will never get here because last getFirst() returns spriteTimeDivider
    return null;
  }

}
