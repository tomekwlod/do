package main

import (
	do "do/digitalocean"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	client, err := do.NewDigitaloceanClient("https://api.digitalocean.com/", os.Getenv("DO_TOKEN"))

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

	for _, server := range servers.Droplets {
		fmt.Printf("-> %s - %s (tags: %v)\n", server.Name, server.Networks.Version[0].IP, strings.Join(server.Tags, "; "))
	}

}
