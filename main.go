package main

import (
	"fmt"
	"os"

	"github.com/EmreKb/qssh/pkg/ui"
)

func main() {
	err := ui.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
