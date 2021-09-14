package vpc

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type VPC struct {
	Config *ssh.ClientConfig
	IP     string
	Port   string
}

func NewVPC(ip, port string, config *ssh.ClientConfig) *VPC {

	return &VPC{IP: ip, Port: port, Config: config}
}

func (v *VPC) dial() (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", v.IP, v.Port), v.Config)

	if err != nil {
		return nil, err
	}

	return client, nil

	// defer client.Close()
}

func (v *VPC) createSession(client *ssh.Client) (*ssh.Session, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()

	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %v", err)
	}

	return session, nil
}

func (v *VPC) ExecuteCommand(cmd string) (string, error) {
	client, err := v.dial()

	if err != nil {
		return "", fmt.Errorf("Failed to dial: %v", err)
	}

	defer client.Close()

	session, err := v.createSession(client)

	if err != nil {
		return "", fmt.Errorf("Failed to create session: %v", err)
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("Failed to run: %+v", err)
	}

	return b.String(), nil
}
