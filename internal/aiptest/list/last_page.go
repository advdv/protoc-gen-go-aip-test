package list

import (
	"github.com/advdv/protoc-gen-go-aip-test/internal/ident"
	"github.com/advdv/protoc-gen-go-aip-test/internal/onlyif"
	"github.com/advdv/protoc-gen-go-aip-test/internal/suite"
	"github.com/advdv/protoc-gen-go-aip-test/internal/util"
	"go.einride.tech/aip/reflect/aipreflect"
	"google.golang.org/protobuf/compiler/protogen"
)

//nolint:gochecknoglobals
var lastPage = suite.Test{
	Name: "last page",
	Doc: []string{
		"If there are no more resources, next_page_token should not be set.",
	},

	OnlyIf: suite.OnlyIfs(
		onlyif.HasMethod(aipreflect.MethodTypeList),
		onlyif.HasParent,
	),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		listMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeList)
		util.MethodList{
			Resource: scope.Resource,
			Method:   listMethod,
			Parent:   "parent",
			PageSize: "int32(resourcesCount + fx.Preexisting)",
		}.Generate(f, "req", "response", "err", ":=", apiMode)
		f.P(ident.AssertNilError, "(t, err)")
		f.P(ident.AssertEqual, "(t, \"\", ", util.FieldGet("response", "NextPageToken", apiMode), ")")
		return nil
	},
}
