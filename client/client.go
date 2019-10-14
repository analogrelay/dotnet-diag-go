package client

// Represents a diagnostics "port" which generally
// refers to a single process.
type DotNetDiagnosticsPort struct {
	// The ID of the process this port can access
	ProcessId int

	// The platform-specific port name (named pipe on Windows, Unix socket on Unix, etc.)
	Port string
}
