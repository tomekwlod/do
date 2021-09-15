package main

import (
	do "do/digitalocean"
	"do/vpc"
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/tomekwlod/utils"
)

//
// TODO:
// - let's standarize the way we receive argument/envs!!
// - concurrency/goroutines
// - standarize error handling
//

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

	for _, server := range servers.Droplet {

		IP := server.PublicIP()

		if IP == "" {
			continue
		}

		if utils.SliceContains(server.Tags, "k8s") {
			continue
		}

		sshpath := os.Getenv("SSH_PATH")
		sshuser := os.Getenv("SSH_USER")

		vpc, err := vpc.NewVPC(IP, 22, sshuser, vpc.AuthWithKey(path.Join(sshpath, "id_rsa"), path.Join(sshpath, "known_hosts")))

		if err != nil {
			fmt.Printf("[%s] Error: %v\n", server.Name, err)
			continue
		}

		defer vpc.Close()

		cmd := "df /home | awk '{ print $5 }' | tail -n 1 | sed 's/%//'"

		res, err := vpc.ExecuteCommand(cmd)

		if err != nil {
			fmt.Printf("Problem with executing command '%s': %v\n", cmd, err)
			return
		}

		fmt.Printf("%s %s - %s", server.Name, IP, res)
	}

}
