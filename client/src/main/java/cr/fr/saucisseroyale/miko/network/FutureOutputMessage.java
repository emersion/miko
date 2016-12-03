package cr.fr.saucisseroyale.miko.network;

import java.io.DataOutputStream;
import java.io.IOException;

@FunctionalInterface
public interface FutureOutputMessage {
  void writeTo(DataOutputStream dos) throws IOException;
}
