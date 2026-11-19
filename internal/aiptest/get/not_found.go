package get

import (
	"github.com/advdv/protoc-gen-go-aip-test/internal/ident"
	"github.com/advdv/protoc-gen-go-aip-test/internal/onlyif"
	"github.com/advdv/protoc-gen-go-aip-test/internal/suite"
	"github.com/advdv/protoc-gen-go-aip-test/internal/util"
	"go.einride.tech/aip/reflect/aipreflect"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/compiler/protogen"
)

//nolint:gochecknoglobals
var notFound = suite.Test{
	Name: "not found",
	Doc: []string{
		"Method should fail with NotFound if the resource does not exist.",
	},

	OnlyIf: suite.OnlyIfs(
		onlyif.HasMethod(aipreflect.MethodTypeGet),
		onlyif.IsNotSingletonResource,
	),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		getMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeGet)

		if util.HasParent(scope.Resource) {
			f.P("parent := ", ident.FixtureNextParent, "(t, false)")
			f.P("created := fx.create(t, parent)")
		} else {
			f.P("created := fx.create(t)")
		}
		util.MethodGet{
			Resource: scope.Resource,
			Method:   getMethod,
			Name:     ident.FixtureNotFoundName + "(" + util.FieldGet("created", "Name", apiMode) + ")",
		}.Generate(f, "req", "_", "err", ":=", apiMode)
		ident.AssertCode(f, scope.Runtime, codes.NotFound)
		return nil
	},
}
