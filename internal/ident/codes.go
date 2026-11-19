package ident

import (
	"github.com/advdv/protoc-gen-go-aip-test/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/compiler/protogen"
)

const connectImportPath = "connectrpc.com/connect/v2"

// AssertCode writes an assertion that the error err has code, in the runtime's codes.
func AssertCode(f *protogen.GeneratedFile, runtime util.Runtime, code codes.Code) {
	f.P(AssertEqual, "(t, ", codeIdent(runtime, code), ", ", statusCodeIdent(runtime), "(err), err)")
}

// codeIdent returns the identifier of code in the runtime's codes.
func codeIdent(runtime util.Runtime, code codes.Code) protogen.GoIdent {
	if runtime == util.RuntimeConnect {
		// Connect names its codes like gRPC, prefixed with Code, and uses the same numbers.
		return protogen.GoIdent{
			GoName:       "Code" + code.String(),
			GoImportPath: connectImportPath,
		}
	}
	return protogen.GoIdent{
		GoName:       code.String(),
		GoImportPath: "google.golang.org/grpc/codes",
	}
}

// statusCodeIdent returns the identifier of the runtime's function that reads the code of an error.
func statusCodeIdent(runtime util.Runtime) protogen.GoIdent {
	if runtime == util.RuntimeConnect {
		return protogen.GoIdent{GoName: "CodeOf", GoImportPath: connectImportPath}
	}
	return protogen.GoIdent{GoName: "Code", GoImportPath: "google.golang.org/grpc/status"}
}
