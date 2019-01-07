package main

import (
	"fmt"
	"os"

	"github.com/majandres/prometheus-kemp-exporter/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
