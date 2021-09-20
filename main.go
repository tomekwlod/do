package main

import (
	do "do/digitalocean"
	"do/vpc"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tomekwlod/utils"
)

//
// TODO:
// - concurrency/goroutines
// - tests
// - maybe sending emails/teams when >80% detected
// - save history to a database
//

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("A problem reading .env file:\n%v", err)
		os.Exit(2)
	}

	dotoken := os.Getenv("DO_TOKEN")

	if dotoken == "" {
		fmt.Println("No DO_TOKEN env provived")
		os.Exit(2)
	}

	client, err := do.NewDigitaloceanClient(http.DefaultClient, "https://api.digitalocean.com/", dotoken)

	if err != nil {
		fmt.Printf("An error occured while connecting to DO:\n%v", err)
		os.Exit(2)
	}

	servers := do.RemoteEntities{}

	err = client.Servers(&servers)

	if err != nil {
		fmt.Printf("Error when fetching the servers:\n%v", err)
		os.Exit(2)
	}

	for _, server := range servers.Droplet {

		IP := server.PublicIP()

		if IP == "" {
			continue
		}

		if utils.SliceContains(server.Tags, "k8s") {
			continue
		}

		publickey := os.Getenv("PUBLIC_KEY_PATH")
		knownhosts := os.Getenv("KNOWN_HOSTS_PATH")
		sshuser := os.Getenv("SSH_USER")

		vpc, err := vpc.NewVPC(IP, 22, sshuser, vpc.AuthWithKey(publickey, knownhosts))
		defer vpc.Close()

		if err != nil {
			fmt.Printf("[%s] An error occured when creating an instance of VPC: %v\n", server.Name, err)
			continue
		}

		cmd := "df /home | awk '{ print $5 }' | tail -n 1 | sed 's/%//'"

		res, err := vpc.ExecuteCommand(cmd)

		if err != nil {
			fmt.Printf("[%s] Problem with executing '%s' command: %v\n", server.Name, cmd, err)
			continue
		}

		fmt.Printf("%s %s - %s", server.Name, IP, res)
	}

}
