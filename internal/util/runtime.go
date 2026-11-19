package util

// Runtime is the RPC runtime the generated test suites call services through.
type Runtime int

const (
	// RuntimeGRPC tests a gRPC <Service>Server from the proto's Go package (default).
	RuntimeGRPC Runtime = iota
	// RuntimeConnect tests any value with the service's standard methods in Connect v2's
	// simple form, such as a Connect client or handler, and generates the suites into a
	// separate <package>aiptest package.
	RuntimeConnect
)
