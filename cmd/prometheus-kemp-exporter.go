package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the main Cobra Command.
var RootCmd = &cobra.Command{
	Use:   "prometheus-kemp-exporter",
	Short: "prometheus-kemp-exporter exports Kemp statistics to Prometheus",
}
