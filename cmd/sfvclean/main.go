package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"sfvclean/internal/app"
	"sfvclean/internal/utils"
)

const (
	CmdAnalyze = "analyze"
	CmdClean   = "clean"
)

func main() {
	log.SetFlags(0)

	analyzeCmd := flag.NewFlagSet(CmdAnalyze, flag.ExitOnError)
	analyzeKeep := analyzeCmd.String("keep", "3M", "Duration for keeping history files")
	analyzeVerbose := analyzeCmd.Bool("verbose", false, "Verbose logging")

	cleanCmd := flag.NewFlagSet(CmdClean, flag.ExitOnError)
	cleanForce := cleanCmd.Bool("force", false, "Delete the data if set to true")

	if len(os.Args) < 2 {
		help()
	}

	var err error
	switch os.Args[1] {
	case CmdAnalyze:
		if err := analyzeCmd.Parse(os.Args[2:]); err != nil {
			help()
		}

		var dur utils.Duration
		dur, err = utils.ParseDuration(*analyzeKeep)
		if err != nil {
			log.Fatalf("invalid duration: %v", err)
		}

		path := analyzeCmd.Arg(0)
		if path == "" {
			path = "."
		}

		err = app.Analyze(path, dur, *analyzeVerbose)

	case CmdClean:
		if err := cleanCmd.Parse(os.Args[2:]); err != nil {
			help()
		}

		err = app.Clean(*cleanForce)

	default:
		help()
	}

	if err != nil {
		log.Fatalf("failed to execute %s: %v", os.Args[1], err)
	}
}

func help() {
	fmt.Println("sfvclean analyze [-verbose] [-keep=DURATION] [PATH]")
	fmt.Println("sfvclean clean [-force]")
	os.Exit(2)
}
