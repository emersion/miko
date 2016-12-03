package cr.fr.saucisseroyale.miko.engine;

import javax.imageio.ImageIO;
import java.awt.*;
import java.awt.image.BufferedImage;
import java.io.IOException;
import java.io.InputStream;

public class ImageReader {
  private final GraphicsConfiguration graphicsConfiguration;

  public ImageReader(GraphicsDevice device) {
    graphicsConfiguration = device.getDefaultConfiguration();
  }

  public BufferedImage read(InputStream is) throws IOException {
    return read(ImageIO.read(is));
  }

  public BufferedImage read(BufferedImage image) {
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
