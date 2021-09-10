package main

import (
	do "do/digitalocean"
	"fmt"
	"os"

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

	r, err := client.Servers()

	if err != nil {
		fmt.Printf("Error, %v\n", err)
		return
	}

	fmt.Println("r", r)
}
