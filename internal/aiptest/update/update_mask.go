package update

import (
	"strconv"

	"github.com/advdv/protoc-gen-go-aip-test/internal/ident"
	"github.com/advdv/protoc-gen-go-aip-test/internal/onlyif"
	"github.com/advdv/protoc-gen-go-aip-test/internal/suite"
	"github.com/advdv/protoc-gen-go-aip-test/internal/util"
	"go.einride.tech/aip/reflect/aipreflect"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// updateMaskOnlyIf are the conditions of every test of update masks that updates a resource.
func updateMaskOnlyIf(extra ...suite.OnlyIf) suite.OnlyIf {
	return suite.OnlyIfs(append([]suite.OnlyIf{
		onlyif.HasMethod(aipreflect.MethodTypeUpdate),
		onlyif.MethodNotLRO(aipreflect.MethodTypeUpdate),
		onlyif.HasUpdateMask,
		onlyif.HasMutableFields,
	}, extra...)...)
}

//nolint:gochecknoglobals
var updateMaskPath = suite.Test{
	Name: "update mask path",
	Doc: []string{
		"For each field that Update can change, an update_mask that names only the field",
		"should give it the value of the request, and keep the others.",
	},
	OnlyIf: updateMaskOnlyIf(),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		rangeMutableFields(f, scope, nil, func(field protoreflect.FieldDescriptor) {
			generateMaskedUpdate(f, scope, apiMode, maskedUpdate{paths: []string{string(field.Name())}})
			generateRequireFields(f, scope, func(other protoreflect.FieldDescriptor) string {
				if other == field {
					return "msg"
				}
				return "created"
			})
		})
		return nil
	},
}

//nolint:gochecknoglobals
var updateMaskPathClears = suite.Test{
	Name: "update mask path clears",
	Doc: []string{
		"For each optional field that Update can change, an update_mask that names the field",
		"should clear it when the request leaves it empty.",
	},
	OnlyIf: updateMaskOnlyIf(onlyif.HasMutableFieldsWith(false)),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		required := false
		rangeMutableFields(f, scope, &required, func(field protoreflect.FieldDescriptor) {
			generateMaskedUpdate(f, scope, apiMode, maskedUpdate{
				clear: field.Name(),
				paths: []string{string(field.Name())},
			})
			generateRequireFields(f, scope, func(other protoreflect.FieldDescriptor) string {
				if other == field {
					return "&" + f.QualifiedGoIdent(scope.Message.GoIdent) + "{}"
				}
				return "created"
			})
		})
		return nil
	},
}

//nolint:gochecknoglobals
var updateMaskUnset = suite.Test{
	Name: "update mask unset",
	Doc: []string{
		"For each field that Update can change, an update without update_mask should keep the",
		"field when the request leaves it empty, as it updates only the fields the request sets.",
	},
	OnlyIf: updateMaskOnlyIf(),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		rangeMutableFields(f, scope, nil, func(field protoreflect.FieldDescriptor) {
			generateMaskedUpdate(f, scope, apiMode, maskedUpdate{clear: field.Name()})
			generateRequireFields(f, scope, func(other protoreflect.FieldDescriptor) string {
				if other == field {
					return "created"
				}
				return ""
			})
		})
		return nil
	},
}

//nolint:gochecknoglobals
var updateMaskEmpty = suite.Test{
	Name: "update mask empty",
	Doc: []string{
		"An update_mask without paths should change no field.",
	},
	OnlyIf: updateMaskOnlyIf(),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		generateMaskedUpdate(f, scope, apiMode, maskedUpdate{emptyMask: true})
		generateRequireFields(f, scope, func(protoreflect.FieldDescriptor) string { return "created" })
		return nil
	},
}

//nolint:gochecknoglobals
var updateMaskPathRequired = suite.Test{
	Name: "update mask path required",
	Doc: []string{
		"For each required field that Update can change, the method should fail with",
		"InvalidArgument if an update_mask names the field and the request leaves it empty.",
	},
	OnlyIf: suite.OnlyIfs(
		onlyif.HasMethod(aipreflect.MethodTypeUpdate),
		onlyif.HasUpdateMask,
		onlyif.HasMutableFieldsWith(true),
	),
	Generate: func(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode) error {
		updateMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeUpdate)
		required := true
		rangeMutableFields(f, scope, &required, func(field protoreflect.FieldDescriptor) {
			f.P("msg := ", ident.FixtureWithoutField, "(created, ", strconv.Quote(string(field.Name())), ")")
			util.MethodUpdate{
				Resource:   scope.Resource,
				Method:     updateMethod,
				Msg:        "msg",
				UpdateMask: []string{strconv.Quote(string(field.Name()))},
			}.Generate(f, "req", "_", "err", ":=", apiMode)
			ident.AssertCode(f, scope.Runtime, codes.InvalidArgument)
		})
		return nil
	},
}

// rangeMutableFields generates a subtest, named after the field, for each field that Update can
// change, or only for those that are or are not REQUIRED when required is set, with body
// generating its body.
func rangeMutableFields(
	f *protogen.GeneratedFile,
	scope suite.Scope,
	required *bool,
	body func(protoreflect.FieldDescriptor),
) {
	util.RangeMutableFields(scope.Message.Desc, func(field protoreflect.FieldDescriptor) {
		if required != nil && util.HasFieldAnnotation(field, annotations.FieldBehavior_REQUIRED) != *required {
			return
		}
		f.P("t.Run(", strconv.Quote(string(field.Name())), ", func(t *", ident.TestingT, ") {")
		f.P(ident.FixtureMaybeSkip, "(t)")
		body(field)
		f.P("})")
	})
}

// maskedUpdate describes the update a test of update masks sends: the Update value of the
// config, with the field called clear cleared, and update_mask with paths, or without paths
// when emptyMask is set, or unset otherwise.
type maskedUpdate struct {
	clear     protoreflect.Name
	paths     []string
	emptyMask bool
}

// generateMaskedUpdate generates a test that creates a resource called created, sends the
// update u with a resource called msg, and reads the result, called updated, and, when the
// service has a Get method, the stored resource, called persisted.
func generateMaskedUpdate(f *protogen.GeneratedFile, scope suite.Scope, apiMode util.APIMode, u maskedUpdate) {
	updateMethod, _ := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeUpdate)
	if util.HasParent(scope.Resource) {
		f.P("parent := ", ident.FixtureNextParent, "(t, false)")
		f.P("created := fx.create(t, parent)")
		f.P("msg := fx.Update(parent)")
	} else {
		f.P("created := fx.create(t)")
		f.P("msg := fx.Update()")
	}
	if u.clear != "" {
		f.P("msg = ", ident.FixtureWithoutField, "(msg, ", strconv.Quote(string(u.clear)), ")")
	}
	f.P(util.FieldSet("msg", "name", util.FieldGet("created", "name", apiMode), apiMode))
	if util.HasEtagField(scope.Message.Desc) {
		f.P(util.FieldSet("msg", "etag", util.FieldGet("created", "etag", apiMode), apiMode))
	}
	paths := make([]string, 0, len(u.paths))
	for _, path := range u.paths {
		paths = append(paths, strconv.Quote(path))
	}
	util.MethodUpdate{
		Resource:   scope.Resource,
		Method:     updateMethod,
		Msg:        "msg",
		UpdateMask: paths,
		EmptyMask:  u.emptyMask,
	}.Generate(f, "req", "updated", "err", ":=", apiMode)
	f.P(ident.AssertNilError, "(t, err)")
	if getMethod, ok := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeGet); ok {
		util.MethodGet{
			Resource: scope.Resource,
			Method:   getMethod,
			Name:     util.FieldGet("updated", "name", apiMode),
		}.Generate(f, "getReq", "persisted", "err", ":=", apiMode)
		f.P(ident.AssertNilError, "(t, err)")
	}
}

// generateRequireFields generates assertions that each field that Update can change has, in
// updated and, when the service has a Get method, in persisted, the value it has in the
// resource that want returns for it. It generates none for a field for which want returns "".
func generateRequireFields(
	f *protogen.GeneratedFile,
	scope suite.Scope,
	want func(protoreflect.FieldDescriptor) string,
) {
	results := "updated"
	if _, ok := util.StandardMethod(scope.Service, scope.Resource, aipreflect.MethodTypeGet); ok {
		results += ", persisted"
	}
	f.P("for _, got := range []*", scope.Message.GoIdent, "{", results, "} {")
	util.RangeMutableFields(scope.Message.Desc, func(field protoreflect.FieldDescriptor) {
		w := want(field)
		if w == "" {
			return
		}
		name := strconv.Quote(string(field.Name()))
		f.P(
			ident.AssertDeepEqual, "(t, ",
			ident.FixtureField, "(", w, ", ", name, "), ",
			ident.FixtureField, "(got, ", name, "), ",
			ident.ProtocmpTransform, "())",
		)
	})
	f.P("}")
}
