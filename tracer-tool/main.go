package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anurse/dotnet-diag-go/client"
)

func main() {
	processId := flag.Int("processId", 0, "The process to trace.")
	flag.Parse()

	if processId == nil || *processId == 0 {
		fmt.Fprintf(os.Stderr, "Missing required flag '-processId'\n")
		os.Exit(1)
	}

	port, err := client.GetPortFor(*processId)
	if noPortErr, ok := err.(client.NoPortForProcessError); ok {
		fmt.Fprintf(os.Stderr, "The process %d is not a .NET Core process, or is not running a new enough runtime.\n", int(noPortErr))
		os.Exit(1)
	}

	fmt.Printf("Connecting to port: %s\n", port.Port)
}
