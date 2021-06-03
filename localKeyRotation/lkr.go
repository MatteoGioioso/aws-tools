package main

import (
	"fmt"
	"github.com/hirvitek/aws-tools/localKeyRotation/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
