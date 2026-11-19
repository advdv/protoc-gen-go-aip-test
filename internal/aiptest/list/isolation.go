package list

import (
	"github.com/advdv/protoc-gen-go-aip-test/internal/ident"
	"github.com/advdv/protoc-gen-go-aip-test/internal/onlyif"
	"github.com/advdv/protoc-gen-go-aip-test/internal/suite"
	"github.com/advdv/protoc-gen-go-aip-test/internal/util"
	"github.com/stoewer/go-strcase"
	"go.einride.tech/aip/reflect/aipreflect"
	"google.golang.org/protobuf/compiler/protogen"
)

//nolint:gochecknoglobals
var isolation = suite.Test{
	Name: "isolation",
	Doc: []string{
		"If parent is provided the method must only return resources",
		"under that parent.",
	},

	OnlyIf: suite.OnlyIfs(
		onlyif.HasMethod(aipreflect.MethodTypeList),
		onlyif.HasParent,
	),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		listMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeList)
		responseResources := strcase.UpperCamelCase(string(util.FindResourceField(
			listMethod.Output.Desc,
			scope.Resource,
		).Name()))
		util.MethodList{
			Resource: scope.Resource,
			Method:   listMethod,
			Parent:   "parent",
			PageSize: "999",
		}.Generate(f, "req", "response", "err", ":=", apiMode)
		f.P(ident.AssertNilError, "(t, err)")
		f.P(ident.AssertDeepEqual, "(")
		f.P("t,")
		f.P(ident.FixtureNormalizeAll, "(parentMsgs),")
		f.P(
			ident.FixtureNormalizeAll,
			"(",
			ident.FixtureListed,
			"(t, parent, parentMsgs, ",
			util.FieldGet("response", responseResources, apiMode),
			")),",
		)
		f.P(ident.CmpoptsSortSlices, "(func(a,b *", scope.Message.GoIdent, ") bool {")
		f.P("return ", util.FieldGet("a", "Name", apiMode), " < ", util.FieldGet("b", "Name", apiMode))
		f.P("}),")
		f.P(ident.ProtocmpTransform, "(),")
		f.P(")")
		return nil
	},
}
