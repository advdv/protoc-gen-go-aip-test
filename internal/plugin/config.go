package plugin

import (
	"flag"
	"fmt"

	"github.com/einride/protoc-gen-go-aip-test/internal/util"
)

// APIMode represents the protobuf API mode for generated code.
type APIMode int

const (
	// APIModeOpen generates code for the Open Struct API (default).
	APIModeOpen APIMode = iota
	// APIModeOpaque generates code for the Opaque API.
	APIModeOpaque
)

func (m APIMode) String() string {
	switch m {
	case APIModeOpen:
		return "API_OPEN"
	case APIModeOpaque:
		return "API_OPAQUE"
	default:
		return ""
	}
}

func (m *APIMode) Set(s string) error {
	switch s {
	case "API_OPAQUE":
		*m = APIModeOpaque
	case "API_OPEN":
		*m = APIModeOpen
	default:
		return fmt.Errorf("invalid api_mode '%s'. Must be either 'API_OPEN' or 'API_OPAQUE'", s)
	}
	return nil
}

// Runtime represents the RPC runtime the generated test suites call services through.
type Runtime util.Runtime

func (r Runtime) String() string {
	switch util.Runtime(r) {
	case util.RuntimeGRPC:
		return "grpc"
	case util.RuntimeConnect:
		return "connect"
	default:
		return ""
	}
}

func (r *Runtime) Set(s string) error {
	switch s {
	case "grpc":
		*r = Runtime(util.RuntimeGRPC)
	case "connect":
		*r = Runtime(util.RuntimeConnect)
	default:
		return fmt.Errorf("invalid runtime '%s'. Must be either 'grpc' or 'connect'", s)
	}
	return nil
}

// Config for the protoc-gen-go-aip-cli plugin.
type Config struct {
	// APIMode specifies which protobuf API style to generate code for.
	APIMode APIMode
	// Runtime specifies which RPC runtime the generated test suites call services through.
	Runtime Runtime
}

// AddToFlagSet adds the config to a pflag.FlagSet.
func (c *Config) AddToFlagSet(flags *flag.FlagSet) {
	flags.Var(
		&c.APIMode,
		"api_mode",
		"the protobuf api mode. Allowed values are 'API_OPEN' and 'API_OPAQUE'",
	)
	flags.Var(
		&c.Runtime,
		"runtime",
		"the rpc runtime the tests call services through. Allowed values are 'grpc' and 'connect'",
	)
}

func (c Config) toUtilAPIMode() util.APIMode {
	switch c.APIMode {
	case APIModeOpen:
		return util.APIModeOpen
	case APIModeOpaque:
		return util.APIModeOpaque
	}
	return util.APIModeOpen
}
