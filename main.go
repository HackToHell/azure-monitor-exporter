package main

import (
	"flag"
	"os"

	"github.com/jkroepke/azure-monitor-exporter/pkg/cmd/config"
	"github.com/jkroepke/azure-monitor-exporter/pkg/cmd/exporter"
)

func main() {
	generate := flag.Bool("generate", false, "Generate something")

	// Parse the flags
	flag.Parse()

	// Check if the generate flag is passed
	if *generate {
		os.Exit(config.Run())
	}
	os.Exit(exporter.Run())
}
