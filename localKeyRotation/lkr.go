package localKeyRotation

import (
	"fmt"
	"local-key-rotation/cmd"
	"os"
)

func init() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
