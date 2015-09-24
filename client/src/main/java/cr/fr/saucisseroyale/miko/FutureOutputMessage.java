package cr.fr.saucisseroyale.miko;

import java.io.DataOutputStream;
import java.io.IOException;

public interface FutureOutputMessage {

  public void writeTo(DataOutputStream dos) throws IOException;
}
