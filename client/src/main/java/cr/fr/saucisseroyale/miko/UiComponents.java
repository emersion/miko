package cr.fr.saucisseroyale.miko;

import edu.umd.cs.findbugs.annotations.SuppressFBWarnings;
import fr.delthas.uitest.*;

import java.util.function.BiConsumer;
import java.util.prefs.Preferences;

@SuppressWarnings("serial")
@SuppressFBWarnings("ICAST_IDIV_CAST_TO_DOUBLE")
final class UiComponents {
  private static final Preferences uiPrefsNode = Preferences.userRoot().node("miko.ui");

  private UiComponents() {}

  public static class Connect extends Layer {
    private Label statusLabel;

    public Connect(Runnable backCallback, String address, int port, BiConsumer<String, Integer> connectCallback, Runnable optionsCallback) {
      setOpaque(true);
      addComponent(new Component() {
        @Override
        protected boolean pushKeyButton(double x, double y, int key, boolean down) {
          if (key != Ui.KEY_ESCAPE || !down) {
            return false;
          }
          backCallback.run();
          return true;
        }
      });
      addComponent(Ui.getWidth() / 10, Ui.getHeight() / 2 + 10, Ui.getWidth() * 3 / 10, 40, new Label("Adresse"));
      addComponent(Ui.getWidth() / 10, Ui.getHeight() / 2 - 50, Ui.getWidth() * 3 / 10, 40, new Label("Port"));
      TextField addressField = new TextField(address);
      addComponent(Ui.getWidth() * 6 / 10, Ui.getHeight() / 2 + 10, Ui.getWidth() * 3 / 10, 40, addressField);
      TextField portField = new TextField(Integer.toString(port));
      portField.setPredicate(s -> {
        try {
          int i = Integer.parseInt(s);
          return i >= 0 && i < 32768;
        } catch (NumberFormatException ignore) {
          return false;
        }
      });
      addComponent(Ui.getWidth() * 6 / 10, Ui.getHeight() / 2 - 50, Ui.getWidth() * 3 / 10, 40, portField);
      Button connectButton = new Button("Se connecter");
      connectButton.setListener((x, y) -> {
        connectCallback.accept(addressField.getText(), Integer.parseInt(portField.getText()));
      });
      addComponent(Ui.getWidth() / 10, 200, Ui.getWidth() * 3 / 10, 30, connectButton);
      Button optionsButton = new Button("Options");
      optionsButton.setListener((x, y) -> {
        optionsCallback.run();
      });
      addComponent(Ui.getWidth() * 6 / 10, 200, Ui.getWidth() * 3 / 10, 30, optionsButton);
      statusLabel = new Label();
      addComponent(0, 10, Ui.getWidth(), 30, statusLabel);
    }

    public void setStatusText(String text) {
      statusLabel.setText(text);
    }
  }

  public static class Login extends Layer {
    private Label statusLabel;
    private Button loginButton;

    public Login(Runnable backCallback, BiConsumer<String, String> registerCallback, BiConsumer<String, String> loginCallback) {
      setOpaque(true);
      addComponent(new Component() {
        @Override
        protected boolean pushKeyButton(double x, double y, int key, boolean down) {
          if (key != Ui.KEY_ESCAPE || !down) {
            return false;
          }
          backCallback.run();
          return true;
        }
      });
      addComponent(Ui.getWidth() / 10, Ui.getHeight() / 2 + 10, Ui.getWidth() * 3 / 10, 40, new Label("Nom d'utilisateur"));
      addComponent(Ui.getWidth() / 10, Ui.getHeight() / 2 - 50, Ui.getWidth() * 3 / 10, 40, new Label("Mot de passe"));
      TextField userField = new TextField();
      userField.setHintText("user");
      addComponent(Ui.getWidth() * 6 / 10, Ui.getHeight() / 2 + 10, Ui.getWidth() * 3 / 10, 40, userField);
      TextField passField = new TextField();
      passField.setHidden(new String(new int[]{8226}, 0, 1));
      passField.setHintText("pwd");
      addComponent(Ui.getWidth() * 6 / 10, Ui.getHeight() / 2 - 50, Ui.getWidth() * 3 / 10, 40, passField);
      loginButton = new Button("Se connecter");
      loginButton.setListener((x, y) -> {
        loginCallback.accept(userField.getText(), passField.getText());
      });
      addComponent(Ui.getWidth() / 10, 200, Ui.getWidth() * 3 / 10, 30, loginButton);
      Button registerButton = new Button("S'inscrire");
      registerButton.setListener((x, y) -> {
        registerCallback.accept(userField.getText(), passField.getText());
      });
      addComponent(Ui.getWidth() * 6 / 10, 200, Ui.getWidth() * 3 / 10, 30, registerButton);
      statusLabel = new Label();
      addComponent(0, 10, Ui.getWidth(), 30, statusLabel);
    }

    public void setStatusText(String text) {
      statusLabel.setText(text);
    }

    public void setLoginEnabled(boolean enabled) {
      loginButton.setEnabled(enabled);
    }
  }

  public static class Options extends Layer {
    public Options(Runnable backCallback, boolean fullscreen) {
      setOpaque(true);
      addComponent(new Component() {
        @Override
        protected boolean pushKeyButton(double x, double y, int key, boolean down) {
          if (key != Ui.KEY_ESCAPE || !down) {
            return false;
          }
          backCallback.run();
          return true;
        }
      });
      addComponent(0, 10, Ui.getWidth(), 30, new Label("Les changements prendront effet au prochain lancement de l'application."));
      CheckBox fullscreenCheckBox = new CheckBox(fullscreen ? "Plein écran" : "Fenêtré sans bordures");
      fullscreenCheckBox.setChecked(fullscreen);
      fullscreenCheckBox.setListener((x, y) -> {
        fullscreenCheckBox.setText(fullscreenCheckBox.isChecked() ? "Plein écran" : "Fenêtré sans bordures");
        uiPrefsNode.putBoolean("fullscreen", fullscreenCheckBox.isChecked());
      });
      addComponent(Ui.getWidth() / 2 - 250, Ui.getHeight() / 2 - 20, 500, 40, fullscreenCheckBox);
    }
  }

  public static class MikoLayer extends Layer {
    public MikoLayer(Runnable backCallback, BiConsumer<InputState, Drawer> render, InputStateManager inputStateManager) {
      setOpaque(true);
      addComponent(new Component() {
        @Override
        protected boolean pushKeyButton(double x, double y, int key, boolean down) {
          if (key != Ui.KEY_ESCAPE || !down) {
            return false;
          }
          backCallback.run();
          return true;
        }
      });
      addComponent(new Component() {
        @Override
        protected boolean pushKeyButton(double x, double y, int key, boolean down) {
          inputStateManager.pushKeyButton(key, down);
          return true;
        }

        @Override
        protected boolean pushMouseButton(double x, double y, int button, boolean down) {
          inputStateManager.pushMouseButton(button, down);
          return true;
        }

        @Override
        protected boolean pushMouseMove(double x, double y) {
          inputStateManager.pushMouseMove(x, y);
          return true;
        }

        @Override
        protected void render(InputState inputState, Drawer drawer) {
          render.accept(inputState, drawer);
        }
      });
    }
  }

  public static class Loading extends Layer {
    private long loadingStart;

    public Loading() {
      setOpaque(true);
      addComponent(new Component() {
        @Override
        protected void render(InputState inputState, Drawer drawer) {
          double maxScreen = Double.min(Ui.getWidth(), Ui.getHeight()) / 2;
          double max = -maxScreen * Math.expm1(-(System.nanoTime() - loadingStart) / 2000000000.0);
          double temp = (System.nanoTime() - loadingStart) / (5000000.0 * 50);
          double offset = (temp - Math.floor(temp)) * 50.0;
          int n = (int) ((max + offset) / 50);
          for (int i = 1; i <= n; i++) {
            drawer.drawLine(Ui.getWidth() / 2 + i * 50 - offset, Ui.getHeight() / 2, Ui.getWidth() / 2, Ui.getHeight() / 2 + i * 50 - offset);
            drawer.drawLine(Ui.getWidth() / 2 + i * 50 - offset, Ui.getHeight() / 2, Ui.getWidth() / 2, Ui.getHeight() / 2 - i * 50 + offset);
            drawer.drawLine(Ui.getWidth() / 2 - i * 50 + offset, Ui.getHeight() / 2, Ui.getWidth() / 2, Ui.getHeight() / 2 + i * 50 - offset);
            drawer.drawLine(Ui.getWidth() / 2 - i * 50 + offset, Ui.getHeight() / 2, Ui.getWidth() / 2, Ui.getHeight() / 2 - i * 50 + offset);
            if (i == 0) {
              break;
            }
            if (i == n) {
              i = -1;
              offset = max;
            }
          }
        }
      });
    }

    public void restartLoading() {
      loadingStart = System.nanoTime();
    }
  }
}
