package cr.fr.saucisseroyale.miko.network;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.IOException;
import java.net.*;
import java.util.LinkedList;
import java.util.List;
import java.util.Optional;

/**
 * Client de serveur de temps Miko.
 */
public class TimeClient {
  private static Logger logger = LogManager.getLogger("miko.network");
  private volatile DatagramSocket socket;
  private Runnable receiverRunnable;
  private volatile boolean clockDifferenceReady;
  private volatile long clockDifference;

  public TimeClient() {
    Thread timeSenderThread = new Thread(() -> {
      while (true) {
        SocketAddress address;
        if (socket != null && !socket.isClosed() && (address = socket.getRemoteSocketAddress()) != null) {
          byte[] data = longToByteArray(System.nanoTime() / 1000);
          try {
            socket.send(new DatagramPacket(data, data.length, address));
          } catch (IOException e) {
            logger.warn("Failed sending time datagram to time server", e);
          }
        }
        try {
          Thread.sleep(20);
        } catch (InterruptedException e) {
          return;
        }
      }
    });
    timeSenderThread.setName("Miko Time Sender");
    timeSenderThread.setDaemon(true);
    timeSenderThread.start();
    receiverRunnable = () -> {
      List<TimeMessage> messages = new LinkedList<>();
      DatagramPacket packet = new DatagramPacket(new byte[16], 16);
      clockDifferenceReady = false;
      while (!socket.isClosed()) {
        try {
          socket.receive(packet);
        } catch (IOException e) {
          // socket closed to interrupt thread
          break;
        }
        long localTimestamp = byteArrayToLong(packet.getData(), 0);
        long remoteTimestamp = byteArrayToLong(packet.getData(), 8);
        long newTimestamp = System.nanoTime() / 1000;
        long latency = (newTimestamp - localTimestamp) / 2;
        messages.add(new TimeMessage(latency, localTimestamp + latency - remoteTimestamp));
        clockDifference = (long) messages.stream().mapToLong(t -> t.clockDifference).average().getAsDouble();
        if (messages.size() > 10) {
          clockDifferenceReady = true;
        }
      }
      clockDifferenceReady = false;
    };
  }

  private static long byteArrayToLong(byte[] array, int offset) {
    return ((long) array[offset] & 0xFFL) << 56 | ((long) array[offset + 1] & 0xFFL) << 48 | ((long) array[offset + 2] & 0xFFL) << 40 | ((long) array[offset + 3] & 0xFFL) << 32 | ((long) array[offset + 4] & 0xFFL) << 24 | ((long) array[offset + 5] & 0xFFL) << 16 | ((long) array[offset + 6] & 0xFFL) << 8 | (long) array[offset + 7] & 0xFFL;
  }

  private static byte[] longToByteArray(long l) {
    return new byte[]{(byte) (l >> 56 & 0xFFL), (byte) (l >> 48 & 0xFFL), (byte) (l >> 40 & 0xFFL), (byte) (l >> 32 & 0xFFL), (byte) (l >> 24 & 0xFFL), (byte) (l >> 16 & 0xFFL), (byte) (l >> 8 & 0xFFL), (byte) (l & 0xFFL)};
  }

  /**
   * Se connecte au serveur de temps spécifié et démarre les échanges de messages de temps.
   *
   * @param address L'adresse du serveur auquel se connecter (IP ou nom d'hôte).
   * @param port    Le port auquel se connecter.
   */
  public void connect(String address, int port) {
    disconnect();
    logger.debug("Starting connection to time server {} at port {}", address, port);
    try {
      socket = new DatagramSocket();
      socket.connect(InetAddress.getByName(address), port);
    } catch (Exception e) {
      logger.fatal("Failed creating datagram socket while initating time client", e);
      throw new RuntimeException(e);
    }
    try {
      socket.setTrafficClass(0x10); // LOWDELAY
    } catch (SocketException e) {
      // ignore if it fails
    }
    Thread timeReceiverThread = new Thread(receiverRunnable);
    timeReceiverThread.setName("Miko Time Receiver");
    timeReceiverThread.setDaemon(true);
    timeReceiverThread.start();
    logger.info("Connected to time server {} at port {}", address, port);
  }

  /**
   * Déconnecte le client du serveur de temps (arrête l'échange de message de temps).
   */
  public synchronized void disconnect() {
    if (socket != null) {
      socket.close();
    }
    logger.info("Disconnected from time server");
  }

  /**
   * @return la différence de temps entre l'horloge locale et l'horloge du serveur : locale - serveur
   */
  public Optional<Long> getClockDifference() {
    if (clockDifferenceReady) {
      return Optional.of(clockDifference);
    }
    return Optional.empty();
  }

  private static class TimeMessage {
    public final long latency;
    public final long clockDifference;

    public TimeMessage(long latency, long clockDifference) {
      this.latency = latency;
      this.clockDifference = clockDifference;
    }
  }
}
