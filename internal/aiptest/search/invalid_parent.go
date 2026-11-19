package search

import (
	"strconv"

	"github.com/advdv/protoc-gen-go-aip-test/internal/ident"
	"github.com/advdv/protoc-gen-go-aip-test/internal/onlyif"
	"github.com/advdv/protoc-gen-go-aip-test/internal/suite"
	"github.com/advdv/protoc-gen-go-aip-test/internal/util"
	"go.einride.tech/aip/reflect/aipreflect"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/compiler/protogen"
)

//nolint:gochecknoglobals
var invalidParent = suite.Test{
	Name: "invalid parent",
	Doc: []string{
		"Method should fail with InvalidArgument if provided parent is invalid.",
	},

	OnlyIf: suite.OnlyIfs(
		onlyif.HasMethod(aipreflect.MethodTypeSearch),
		onlyif.HasParent,
	),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		searchMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeSearch)
		util.MethodSearch{
			Resource: scope.Resource,
			Method:   searchMethod,
			Parent:   strconv.Quote("invalid resource name"),
		}.Generate(f, "req", "_", "err", ":=", apiMode)
		ident.AssertCode(f, scope.Runtime, codes.InvalidArgument)
		return nil
	},
}
