package main

import (
	"fmt"
	"io"
	"os"

	"github.com/MihaiLupoiu/interview-exasol/config"
	"github.com/MihaiLupoiu/interview-exasol/miner"
)

const (
	// exitFail is the exit code if the program
	// fails.
	exitFail = 1
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run(args []string, stdout io.Writer) error {
	configuration := config.Get()

	minerCtx, err := miner.Init(configuration)
	if err != nil {
		return err
	}

	return minerCtx.Run()
}
