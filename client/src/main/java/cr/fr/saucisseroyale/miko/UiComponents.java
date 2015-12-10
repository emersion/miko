package cr.fr.saucisseroyale.miko;

import java.awt.BorderLayout;
import java.awt.GridLayout;
import java.text.NumberFormat;
import java.text.ParseException;
import java.util.function.BiConsumer;
import java.util.prefs.Preferences;

import javax.swing.JButton;
import javax.swing.JCheckBox;
import javax.swing.JFormattedTextField;
import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.JPasswordField;
import javax.swing.JTextField;
import javax.swing.text.NumberFormatter;

@SuppressWarnings("serial")
class UiComponents {

  private static final Preferences uiPrefsNode = Preferences.userRoot().node("miko.ui");

  public static class Connect extends JPanel {

    private JLabel statusField;

    public Connect(String defaultAddress, int defaultPort, BiConsumer<String, Integer> connectCallback, Runnable optionsCallback) {
      setLayout(new GridLayout(0, 2, 10, 10));
      add(new JLabel("Adresse"));
      JTextField addressField = new JTextField(defaultAddress);
      add(addressField);
      add(new JLabel("Port"));
      NumberFormat longFormat = NumberFormat.getIntegerInstance();
      longFormat.setGroupingUsed(false);
      NumberFormatter numberFormatter = new NumberFormatter(longFormat);
      numberFormatter.setAllowsInvalid(false);
      numberFormatter.setValueClass(Integer.class);
      numberFormatter.setMinimum(0);
      JFormattedTextField portField = new JFormattedTextField(numberFormatter);
      portField.setValue(Integer.valueOf(defaultPort));
      add(portField);
      JButton connectButton = new JButton("Se connecter");
      connectButton.addActionListener(event -> {
        try {
          portField.commitEdit();
        } catch (ParseException e) {
          // will never happen since we disallowed invalids
          throw new RuntimeException(e);
        }
        connectCallback.accept(addressField.getText(), (Integer) portField.getValue());
      });
      add(connectButton);
      JButton optionsButton = new JButton("Options");
      optionsButton.addActionListener(event -> {
        optionsCallback.run();
      });
      add(optionsButton);
      statusField = new JLabel();
      add(statusField);
    }

    public void setStatusText(String text) {
      statusField.setText(text);
    }
  }

  public static class Login extends JPanel {

    private JLabel statusField;

    public Login(BiConsumer<String, String> registerCallback, BiConsumer<String, String> loginCallback) {
      setLayout(new GridLayout(4, 2, 10, 10));
      add(new JLabel("Nom d'utilisateur"));
      JTextField usernameField = new JTextField();
      add(usernameField);
      add(new JLabel("Mot de passe"));
      JPasswordField passwordField = new JPasswordField();
      add(passwordField);
      JButton registerButton = new JButton("S'inscrire");
      registerButton.addActionListener(e -> registerCallback.accept(usernameField.getText(), new String(passwordField.getPassword())));
      add(registerButton);
      JButton loginButton = new JButton("Se connecter");
      loginButton.addActionListener(e -> loginCallback.accept(usernameField.getText(), new String(passwordField.getPassword())));
      add(loginButton);
      statusField = new JLabel();
      add(statusField);
    }

    public void setStatusText(String text) {
      statusField.setText(text);
    }
  }

  public static class Options extends JPanel {

    public Options(boolean fullscreen) {
      setLayout(new BorderLayout());
      add(new JLabel("Les changements prendront effet au prochain lancement de l'application."), BorderLayout.SOUTH);
      JPanel panel = new JPanel();
      panel.setLayout(new GridLayout());
      add(panel, BorderLayout.CENTER);
      JCheckBox fullscreenCheckBox = new JCheckBox(fullscreen ? "Plein écran" : "Fenêtré sans bordures", fullscreen);
      fullscreenCheckBox.addActionListener(l -> {
        fullscreenCheckBox.setText(fullscreenCheckBox.isSelected() ? "Plein écran" : "Fenêtré sans bordures");
        uiPrefsNode.putBoolean("fullscreen", fullscreenCheckBox.isSelected());
      });
      panel.add(fullscreenCheckBox);
    }
  }

}
