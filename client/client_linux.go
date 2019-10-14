// +build linux darwin

package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
)

var diagPortFormat = regexp.MustCompile(`dotnet-diagnostic-(\d+)-(\d+)-socket`)

type NoPortForProcessError int

func (n NoPortForProcessError) Error() string {
	return fmt.Sprintf("No port for process: %d", int(n))
}

func GetPortFor(processId int) (DotNetDiagnosticsPort, error) {
	ports, err := GetDiagnosticsPorts()
	if err != nil {
		return DotNetDiagnosticsPort{}, err
	}

	for _, port := range ports {
		if port.ProcessId == processId {
			return port, nil
		}
	}

	return DotNetDiagnosticsPort{}, NoPortForProcessError(processId)
}

func GetDiagnosticsPorts() ([]DotNetDiagnosticsPort, error) {
	tempBase := os.TempDir()
	candidates, err := ioutil.ReadDir(tempBase)
	if err != nil {
		return nil, err
	}

	ports := make([]DotNetDiagnosticsPort, 0)
	for _, candidate := range candidates {
		if candidate.IsDir() {
			continue
		}

		matches := diagPortFormat.FindStringSubmatch(candidate.Name())
		if len(matches) == 0 {
			continue
		}

		// It's a valid candidate!
		processId, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			continue
		}

		// TODO: Log a warning or something if the process ID can't be parsed?
		// Valid process ID!
		ports = append(ports, DotNetDiagnosticsPort{
			ProcessId: int(processId),
			Port:      path.Join(tempBase, candidate.Name()),
		})
	}

	return ports, nil
}
