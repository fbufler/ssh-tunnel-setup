package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

const keySize = 4096
const newSuffix = ".new"

// MakeKeyPair generates a new RSA key pair and writes the public key to pubKeyPath and the private key to privateKeyPath
func MakeKeyPair(keyPath, keyName, userReference string) error {
	slog.Debug("Generating key pair")
	pubKeyPath := fmt.Sprintf("%s/%s.pub", keyPath, keyName)
	privateKeyPath := fmt.Sprintf("%s/%s", keyPath, keyName)
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}
	slog.Debug("Generating and writing private key as PEM")
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		privateKeyFile.Close()
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}
	privateKeyFile.Close()
	err = os.Chmod(privateKeyPath, 0600)
	if err != nil {
		return err
	}

	slog.Debug("Generating and writing public key")
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	pubKeyStr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pub)))
	pubKeyWithComment := fmt.Sprintf("%s %s", pubKeyStr, userReference)
	return os.WriteFile(pubKeyPath, []byte(pubKeyWithComment), 0644)
}

func PrepareKeyDirectory(keyPath string) error {
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		slog.Debug(fmt.Sprintf("Creating directory: %s", keyPath))
		err := os.MkdirAll(keyPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

type RemoteAuth struct {
	User           string
	Password       string
	KeyPath        string
	TrustedHostKey string
}

func NewRemoteAuth(user, password, keyPath, trustedHostKey string) RemoteAuth {
	return RemoteAuth{
		User:           user,
		Password:       password,
		KeyPath:        keyPath,
		TrustedHostKey: trustedHostKey,
	}
}

// authMethods returns a slice of ssh.AuthMethod based on the provided RemoteAuth
// KeyPath is preferred over Password
func (ra RemoteAuth) authMethods() ([]ssh.AuthMethod, error) {
	if ra.KeyPath != "" {
		key, err := os.ReadFile(ra.KeyPath)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
	}
	if ra.Password != "" {
		return []ssh.AuthMethod{ssh.Password(ra.Password)}, nil
	}
	return nil, fmt.Errorf("no password or key provided")
}

func (ra RemoteAuth) test(remote string) error {
	slog.Debug("Testing remote authentication")
	authMethods, err := ra.authMethods()
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", remote, &ssh.ClientConfig{
		User:            ra.User,
		Auth:            authMethods,
		HostKeyCallback: trustedHostKeyCallback(ra.TrustedHostKey),
	})
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

// AuthorizePublicKeyOnRemote appends the public key at pubKeyPath to the remote's authorized_keys file
func AuthorizePublicKeyOnRemote(privateKeyPath, remote, remoteUser string, ra RemoteAuth) error {
	slog.Debug("Authorizing public key on remote")
	publicKey, err := os.ReadFile(fmt.Sprintf("%s.pub", privateKeyPath))
	if err != nil {
		return err
	}
	authMethods, err := ra.authMethods()
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", remote, &ssh.ClientConfig{
		User:            ra.User,
		Auth:            authMethods,
		HostKeyCallback: trustedHostKeyCallback(ra.TrustedHostKey),
	})
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	command := fmt.Sprintf(`echo "%s" >> /home/%s/.ssh/authorized_keys`, string(publicKey), remoteUser)
	err = session.Run(command)
	if err != nil {
		slog.Error("Error authorizing public key on remote")
		return err
	}

	slog.Debug("Testing new key on remote")
	newRa := RemoteAuth{
		User:           remoteUser,
		KeyPath:        privateKeyPath,
		TrustedHostKey: ra.TrustedHostKey,
	}
	err = newRa.test(remote)
	if err != nil {
		slog.Error("Error testing public key on remote")
		return err
	}
	return nil
}

func UnauthorizedPublicKeyOnRemote(privateKeyPath, remote, remoteUser string, auth RemoteAuth) error {
	slog.Debug("Unauthorizing public key on remote")
	publicKey, err := os.ReadFile(fmt.Sprintf("%s.pub", privateKeyPath))
	if err != nil {
		return err
	}
	authMethods, err := auth.authMethods()
	if err != nil {
		return err
	}
	client, err := ssh.Dial("tcp", remote, &ssh.ClientConfig{
		User:            auth.User,
		Auth:            authMethods,
		HostKeyCallback: trustedHostKeyCallback(auth.TrustedHostKey),
	})
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	escapedPublicKey := strings.ReplaceAll(string(publicKey), "/", "\\/")
	command := fmt.Sprintf(`sed -i '/%s/d' /home/%s/.ssh/authorized_keys`, escapedPublicKey, remoteUser)
	err = session.Run(command)
	if err != nil {
		slog.Error("Error unauthorizing public key on remote")
		return err
	}

	slog.Debug("Testing new key on remote")
	err = auth.test(remote)
	if err != nil {
		slog.Error("Error testing public key on remote")
		return err
	}
	return nil
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}

func trustedHostKeyCallback(trustedHostKey string) ssh.HostKeyCallback {
	if trustedHostKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			slog.Warn(fmt.Sprintf("SSH-key verification is *NOT* in effect: to fix, add this trusted-host-key: %q", keyString(k)))
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedHostKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedHostKey, ks)
		}
		return nil
	}
}

// RotateKeyPair generates a new key pair, authorizes the public key on the remote, and removes the old key pair
func RotateKeyPair(user, remote, keyPath, keyName, keyUser, trustedHostKey string) error {
	newPubKeyPath := fmt.Sprintf("%s/%s%s.pub", keyPath, keyName, newSuffix)
	newPrivateKeyPath := fmt.Sprintf("%s/%s%s", keyPath, keyName, newSuffix)
	pubKeyPath := fmt.Sprintf("%s/%s.pub", keyPath, keyName)
	privateKeyPath := fmt.Sprintf("%s/%s", keyPath, keyName)

	slog.Debug("Generating new key pair")
	if err := MakeKeyPair(keyPath, keyName+newSuffix, keyUser); err != nil {
		slog.Error("Error making new key pair")
		return err
	}

	slog.Debug("Authorizing new public key on remote")
	ra := NewRemoteAuth(user, "", privateKeyPath, trustedHostKey)
	if err := AuthorizePublicKeyOnRemote(newPrivateKeyPath, remote, user, ra); err != nil {
		slog.Error("Error authorizing new public key on remote")
		removeKeyPair(newPubKeyPath, newPrivateKeyPath)
		return err
	}

	slog.Debug("Unauthorizing old public key on remote")
	ra = NewRemoteAuth(user, "", newPrivateKeyPath, trustedHostKey)
	if err := UnauthorizedPublicKeyOnRemote(privateKeyPath, remote, user, ra); err != nil {
		slog.Error("Error unauthorizing old public key on remote")
		removeKeyPair(newPubKeyPath, newPrivateKeyPath)
		return err
	}

	slog.Debug("Replacing old key pair with new key pair")
	if err := os.Rename(newPubKeyPath, pubKeyPath); err != nil {
		slog.Error("Error replacing old public key with new public key")
		return err
	}
	if err := os.Rename(newPrivateKeyPath, privateKeyPath); err != nil {
		slog.Error("Error replacing old private key with new private key")
		return err
	}

	return nil
}

func removeKeyPair(pubKeyPath, privateKeyPath string) error {
	slog.Debug("Removing key pair")
	if err := os.Remove(pubKeyPath); err != nil {
		slog.Error(fmt.Sprintf("Error removing public key %s, please remove manually", pubKeyPath))
	}
	if err := os.Remove(privateKeyPath); err != nil {
		slog.Error(fmt.Sprintf("Error removing private key %s, please remove manually", privateKeyPath))
	}
	return nil
}
