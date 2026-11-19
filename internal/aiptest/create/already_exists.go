package create

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
var alreadyExists = suite.Test{
	Name: "already exists",
	Doc: []string{
		"If method support user settable IDs and the same ID is reused",
		"the method should return AlreadyExists.",
	},

	OnlyIf: suite.OnlyIfs(
		onlyif.HasMethod(aipreflect.MethodTypeCreate),
		onlyif.MethodNotLRO(aipreflect.MethodTypeCreate),
		onlyif.HasUserSettableID,
	),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		createMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeCreate)
		if util.HasParent(scope.Resource) {
			f.P("parent := ", ident.FixtureNextParent, "(t, false)")
		}
		f.P("userSetID := ", ident.FixtureUserSettableID, "(", strconv.Quote("alreadyexists"), ")")
		m := util.MethodCreate{
			Resource:       scope.Resource,
			Method:         createMethod,
			Parent:         "parent",
			UserSettableID: "userSetID",
		}
		m.Generate(f, "req1", "_", "err", ":=", apiMode)
		f.P(ident.AssertNilError, "(t, err)")
		m.Generate(f, "req2", "_", "err", "=", apiMode)
		ident.AssertCode(f, scope.Runtime, codes.AlreadyExists)
		return nil
	},
}
