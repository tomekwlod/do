package vpc

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// AuthMethod is an interface and all types that want to satisfy it need to
// provide method() function
type AuthMethod interface {
	method() (ssh.AuthMethod, ssh.HostKeyCallback, error)
}

// Password Authentication
func AuthWithPassword(password string) AuthMethod {
	return passwordCallback(func() (string, error) {
		return password, nil
	})
}

type passwordCallback func() (password string, err error)

func (p passwordCallback) method() (ssh.AuthMethod, ssh.HostKeyCallback, error) {

	password, _ := p()

	auth := ssh.Password(password)

	var hostKey ssh.PublicKey

	hostKeyCallback := ssh.FixedHostKey(hostKey)

	return auth, hostKeyCallback, nil
}

// SSHKey authentication
func AuthWithKey(publicKeyPath, knownHostsPath string) AuthMethod {
	return keyCallback(func() (string, string, error) {
		return publicKeyPath, knownHostsPath, nil
	})
}

type keyCallback func() (publicKeyPath, knownHostsPath string, err error)

func (k keyCallback) method() (ssh.AuthMethod, ssh.HostKeyCallback, error) {

	publicKeyPath, knownHostsPath, _ := k()

	if publicKeyPath == "" {
		return nil, nil, fmt.Errorf("Please provide a proper location to your pubic ssh key")
	}

	key, err := ioutil.ReadFile(publicKeyPath)

	if err != nil {
		return nil, nil, fmt.Errorf("Unable to read a private key: %v", err)
	}

	// Create the Signer for the private key.
	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		return nil, nil, fmt.Errorf("Unable to parse private key: %v", err)
	}

	auth := ssh.PublicKeys(signer)

	hostKeyCallback, err := knownhosts.New(knownHostsPath)

	if err != nil {
		return nil, nil, fmt.Errorf("Could not create hostkeycallback function: %v", err)
	}

	return auth, hostKeyCallback, nil
}
