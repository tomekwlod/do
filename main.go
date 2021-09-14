package main

import (
	do "do/digitalocean"
	"do/vpc"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/tomekwlod/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// type ClientInterface interface {
// 	Servers(interface{}) error
// }

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	sshpath := os.Getenv("SSH_PATH")
	if sshpath == "" {
		log.Fatal("No ssh path provived, expected SSH_PATH in env file (where are your ssh files?)")
	}

	client, err := do.NewDigitaloceanClient("https://api.digitalocean.com/", os.Getenv("DO_TOKEN"))

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	// initialize an empty model struct
	servers := do.RemoteEntities{}

	err = client.Servers(&servers)

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	key, err := ioutil.ReadFile(path.Join(sshpath, "id_rsa"))
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	hostKeyCallback, err := knownhosts.New(path.Join(sshpath, "known_hosts"))
	if err != nil {
		log.Fatal("could not create hostkeycallback function: ", err)
	}

	for _, server := range servers.Droplet {

		IP := server.PublicIP()

		if IP == "" {
			continue
		}

		if utils.SliceContains(server.Tags, "k8s") {

			continue
		}

		config := &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: hostKeyCallback,
		}

		v := vpc.NewVPC(IP, "22", config)

		cmd := "df /home | awk '{ print $5 }' | tail -n 1 | sed 's/%//'"

		res, err := v.ExecuteCommand(cmd)

		if err != nil {
			fmt.Printf("Problem with executing command '%s': %v\n", cmd, err)
			return
		}

		fmt.Printf("%s %s - %s", server.Name, IP, res)
	}

}
