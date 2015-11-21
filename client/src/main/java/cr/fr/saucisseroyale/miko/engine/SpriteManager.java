package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.util.Pair;

import java.awt.Graphics2D;
import java.awt.GraphicsConfiguration;
import java.awt.Transparency;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.file.DirectoryStream;
import java.nio.file.FileSystem;
import java.nio.file.FileSystems;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;

import javax.imageio.ImageIO;

/**
 * Un gestionnaire de sprite gérant des sprites en mémoire et pouvant les afficher sur des
 * graphiques.
 *
 */
public class SpriteManager {

  private final GraphicsConfiguration graphicsConfiguration;
  private Map<SpriteType, Sprite> sprites;

  public SpriteManager(GraphicsConfiguration graphicsConfiguration) {
    this.graphicsConfiguration = graphicsConfiguration;
  }

  /**
   * Affiche le sprite spécifié sur les graphiques, à l'endroit spécifié.
   *
   * @param graphics Les graphiques sur lesquels peindre l'image.
   * @param spriteType Le sprite identifiant l'animation de laquelle peindre.
   * @param spriteTime Le temps en ticks depuis que le sprite a été sélectionné.
   * @param x La position en x du centre de l'image à afficher.
   * @param y La position en y du centre de l'image à afficher.
   */
  public void drawSpriteType(Graphics2D graphics, SpriteType spriteType, long spriteTime, int x, int y) {
    Sprite sprite = sprites.get(spriteType);
    if (sprite == null) {
      throw new IllegalArgumentException("SpriteType " + spriteType + " doesn't have a sprite.");
    }
    BufferedImage image = sprite.getImage(spriteTime);
    graphics.drawImage(image, x - image.getWidth() / 2, y - image.getHeight() / 2, null);
  }

  /**
   * Charge les sprites en mémoire pour pouvoir les utiliser après. Doit être appelé avant tout
   * appel à {@link #drawSpriteType(Graphics2D, SpriteType, long, int, int)}.
   *
   * @throws IOException S'il y a une erreur lors de l'accès ou de la lecture des sprites.
   */
  public void loadImages() throws IOException {
    try {
      URI uri = SpriteManager.class.getProtectionDomain().getCodeSource().getLocation().toURI();
      if (uri.getPath().endsWith(".jar")) {
        try (FileSystem fs = FileSystems.newFileSystem(URI.create("jar:file:" + uri.getPath()), Collections.emptyMap())) {
          loadImages(fs.getRootDirectories().iterator().next());
        }
      } else {
        loadImages(Paths.get(SpriteManager.class.getResource("/").toURI()));
      }
    } catch (URISyntaxException e) {
      // should never happen
      throw new RuntimeException(e);
    }
  }

  private void loadImages(Path root) throws IOException {
    try (DirectoryStream<Path> pathStream = Files.newDirectoryStream(root.resolve("sprite"))) {
      for (Path spritePath : pathStream) {
        if (Files.isDirectory(spritePath)) {
          String directoryName = spritePath.getFileName().toString();
          int id;
          try {
            id = Integer.parseUnsignedInt(directoryName);
          } catch (NumberFormatException e) {
            continue;
          }
          SpriteType spriteType = SpriteType.getType(id);
          if (spriteType == null) {
            continue;
          }
          try (DirectoryStream<Path> imagePathStream = Files.newDirectoryStream(spritePath)) {
            List<Pair<Integer, BufferedImage>> images = new LinkedList<>();
            for (Path imagePath : imagePathStream) {
              String imageName = imagePath.getFileName().toString();
              if (!imageName.endsWith(".png")) {
                continue;
              }
              int timecode;
              try {
                // 4 == ".png".length()
                String timecodeString = imageName.substring(0, imageName.length() - 4);
                timecode = Integer.parseUnsignedInt(timecodeString);
              } catch (NumberFormatException e) {
                continue;
              }
              try (InputStream is = Files.newInputStream(imagePath)) {
                BufferedImage image = readImage(is);
                images.add(new Pair<>(timecode, image));
              }
            }
            sprites.put(spriteType, new Sprite(images));
          }
        } else {
          String imageName = spritePath.getFileName().toString();
          if (!imageName.endsWith(".png")) {
            continue;
          }
          int id;
          try {
            // 4 == ".png".length()
            String idString = imageName.substring(0, imageName.length() - 4);
            id = Integer.parseUnsignedInt(idString);
          } catch (NumberFormatException e) {
            continue;
          }
          SpriteType spriteType = SpriteType.getType(id);
          if (spriteType == null) {
            continue;
          }
          if (sprites.containsKey(spriteType)) {
            continue;
          }
          try (InputStream is = Files.newInputStream(spritePath)) {
            BufferedImage image = readImage(is);
            sprites.put(spriteType, new Sprite(image));
          }
        }
      }
    }
  }

  private BufferedImage readImage(InputStream is) throws IOException {
    return readImage(ImageIO.read(is));
  }

  private BufferedImage readImage(BufferedImage image) {
    BufferedImage compatibleImage;
    if (image.getColorModel().equals(graphicsConfiguration.getColorModel())
        && (image.getTransparency() == Transparency.BITMASK || image.getTransparency() == Transparency.OPAQUE)) {
      compatibleImage = image;
    } else {
      compatibleImage = graphicsConfiguration.createCompatibleImage(image.getWidth(), image.getHeight(), Transparency.BITMASK);
      Graphics2D g2d = (Graphics2D) compatibleImage.getGraphics();
      g2d.drawImage(image, 0, 0, null);
      g2d.dispose();
      image.flush();
    }
    return compatibleImage;
  }

}
