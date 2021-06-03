package main

import (
	"github.com/hirvitek/aws-tools/localKeyRotation/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("[LKR] There was an error: %v", err.Error())
	}
}
