/*
Copyright © 2024 Floom,  dev@floom.ai
*/
package main

import (
	"FloomCLI/cmd"
	"FloomCLI/config"
	"log"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to load floom configuration file: %v", err)
	}
	cmd.Execute()
}
