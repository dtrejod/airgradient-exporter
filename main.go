package main

import (
	"os"

	"github.com/dtrejod/airgradient-exporter/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
