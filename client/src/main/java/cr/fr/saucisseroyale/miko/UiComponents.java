package cr.fr.saucisseroyale.miko;

import java.awt.GridLayout;
import java.util.function.BiConsumer;

import javax.swing.JButton;
import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.JPasswordField;
import javax.swing.JTextField;

@SuppressWarnings("serial")
public class UiComponents {

  public static class Connect extends JPanel {

    private JLabel statusField;

    public Connect(String defaultAddress, String defaultPort,
        BiConsumer<String, String> connectCallback) {
      setLayout(new GridLayout(4, 2, 10, 10));
      add(new JLabel("Adresse"));
      JTextField addressField = new JTextField(defaultAddress, 20);
      add(addressField);
      add(new JLabel("Port"));
      JTextField portField = new JTextField(defaultPort, 20);
      add(portField);
      JButton connectButton = new JButton("Se connecter");
      connectButton.addActionListener((e) -> connectCallback.accept(addressField.getText(),
          portField.getText()));
      add(connectButton);
      add(new JLabel());
      statusField = new JLabel();
      add(statusField);
    }

    public void setStatusText(String text) {
      statusField.setText(text);
    }
  }

  public static class Login extends JPanel {

    private JLabel statusField;

    public Login(BiConsumer<String, String> registerCallback,
        BiConsumer<String, String> loginCallback) {
      setLayout(new GridLayout(4, 2, 10, 10));
      add(new JLabel("Nom d'utilisateur"));
      JTextField usernameField = new JTextField(20);
      add(usernameField);
      add(new JLabel("Mot de passe"));
      JPasswordField passwordField = new JPasswordField(20);
      add(passwordField);
      JButton registerButton = new JButton("S'inscrire");
      registerButton.addActionListener((e) -> registerCallback.accept(usernameField.getText(),
          passwordField.getPassword().toString()));
      add(registerButton);
      JButton loginButton = new JButton("Se connecter");
      loginButton.addActionListener((e) -> loginCallback.accept(usernameField.getText(),
          passwordField.getPassword().toString()));
      add(loginButton);
      statusField = new JLabel();
      add(statusField);
    }

    public void setStatusText(String text) {
      statusField.setText(text);
    }
  }

}
