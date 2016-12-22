package cr.fr.saucisseroyale.miko.engine;

import cr.fr.saucisseroyale.miko.protocol.SpriteType;
import cr.fr.saucisseroyale.miko.util.Pair.Int;
import fr.delthas.uitest.Drawer;
import fr.delthas.uitest.Image;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.file.*;
import java.util.*;

/**
 * Un gestionnaire de sprite gérant des sprites en mémoire et pouvant les afficher sur des
 * graphiques.
 */
public class SpriteManager {
  private Map<SpriteType, Sprite> sprites = new HashMap<>();

  /**
   * Affiche le sprite spécifié sur les graphiques, à l'endroit spécifié.
   *
   * @param drawer     Le drawer avec lequel peindre l'image.
   * @param spriteType Le sprite identifiant l'animation de laquelle peindre.
   * @param spriteTime Le temps en ticks depuis que le sprite a été sélectionné.
   * @param x          La position en x du centre de l'image à afficher.
   * @param y          La position en y du centre de l'image à afficher.
   */
  public void drawSpriteType(Drawer drawer, SpriteType spriteType, long spriteTime, int x, int y) {
    Sprite sprite = sprites.get(spriteType);
    if (sprite == null) {
      throw new IllegalArgumentException("SpriteType " + spriteType + " doesn't have a sprite.");
    }
    Image image = sprite.getImage(spriteTime);
    drawer.drawImage(x, y, image);
  }

  /**
   * Charge les sprites en mémoire pour pouvoir les utiliser après. Doit être appelé avant tout
   * appel à {@link #drawSpriteType(Drawer, SpriteType, long, int, int)}.
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
          } catch (NumberFormatException ignore) {
            continue;
          }
          SpriteType spriteType = SpriteType.getType(id);
          if (spriteType == null) {
            continue;
          }
          try (DirectoryStream<Path> imagePathStream = Files.newDirectoryStream(spritePath)) {
            List<Int<Image>> images = new LinkedList<>();
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
              } catch (NumberFormatException ignore) {
                continue;
              }
              images.add(new Int<>(timecode, Image.createImage(imagePath.toString())));
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
          } catch (NumberFormatException ignore) {
            continue;
          }
          SpriteType spriteType = SpriteType.getType(id);
          if (spriteType == null) {
            continue;
          }
          if (sprites.containsKey(spriteType)) {
            continue;
          }
          sprites.put(spriteType, new Sprite(Image.createImage(spritePath.toString())));
        }
      }
    }
  }
}
