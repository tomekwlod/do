package vpc

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type VPC struct {
	client     *ssh.Client
	IP         string
	Port       int
	User       string
	AuthMethod AuthMethod
}

// NewVPC initializes and dials to the server
func NewVPC(ip string, port int, user string, authMethod AuthMethod) (*VPC, error) {

	auth, hostKeyCallback, err := authMethod.method() // call it prepare

	if err != nil {
		return nil, fmt.Errorf("Prepare method error: %v", err)
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)

	if err != nil {
		return nil, fmt.Errorf("Error while dialing: %v", err)
	}

	return &VPC{IP: ip, Port: port, User: user, client: client}, nil
}

// Close closes the client connection
func (vpc *VPC) Close() error {
	return vpc.client.Close()
}

// ExecuteCommand executes provided command and returns a response in string format
func (v *VPC) ExecuteCommand(cmd string) (string, error) {

	session, err := v.client.NewSession()

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
