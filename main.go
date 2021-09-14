package main

import (
	do "do/digitalocean"
	"do/vpc"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/tomekwlod/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	dotoken := os.Getenv("DO_TOKEN")

	if dotoken == "" {
		fmt.Println("No DO_TOKEN env provived")
		return
	}

	client, err := do.NewDigitaloceanClient("https://api.digitalocean.com/", dotoken)

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	servers := do.RemoteEntities{}

	err = client.Servers(&servers)

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	// all good, now start with the VPCs

	sshpath := os.Getenv("SSH_PATH")

	if sshpath == "" {
		fmt.Println("No ssh path provived, expected SSH_PATH in env file (where are your ssh files?)")
		return
	}

	key, err := ioutil.ReadFile(path.Join(sshpath, "id_rsa"))

	if err != nil {
		fmt.Printf("Unable to read a private key: %v", err)
		return
	}

	// Create the Signer for the private key.
	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		fmt.Printf("Unable to parse private key: %v", err)
		return
	}

	hostKeyCallback, err := knownhosts.New(path.Join(sshpath, "known_hosts"))

	if err != nil {
		fmt.Printf("Could not create hostkeycallback function: %v", err)
		return
	}

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
	}

	for _, server := range servers.Droplet {

		IP := server.PublicIP()

		if IP == "" {
			continue
		}

		if utils.SliceContains(server.Tags, "k8s") {
			continue
		}

		v := vpc.NewVPC(IP, 22, config)

		cmd := "df /home | awk '{ print $5 }' | tail -n 1 | sed 's/%//'"

		res, err := v.ExecuteCommand(cmd)

		if err != nil {
			fmt.Printf("Problem with executing command '%s': %v\n", cmd, err)
			return
		}

		fmt.Printf("%s %s - %s", server.Name, IP, res)
	}

}
