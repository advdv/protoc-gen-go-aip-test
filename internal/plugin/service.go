package plugin

import (
	"strconv"

	"github.com/einride/protoc-gen-go-aip-test/internal/aip/method"
	"github.com/einride/protoc-gen-go-aip-test/internal/util"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type serviceGenerator struct {
	service   *protogen.Service
	resources []*annotations.ResourceDescriptor
	messages  []*protogen.Message
	config    Config
}

func (s *serviceGenerator) Generate(f *protogen.GeneratedFile) error {
	s.generateServiceInterface(f)
	s.generateConfigProvider(f)
	s.generateMainTestFunction(f)
	s.generateTestFunctions(f)
	s.generateFixture(f)
	s.generateTestMethods(f)
	for i, resource := range s.resources {
		message := s.messages[i]
		generator := resourceGenerator{
			service:  s.service,
			resource: resource,
			message:  message,
			config:   s.config,
		}
		if err := generator.Generate(f); err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceGenerator) generateConfigProvider(f *protogen.GeneratedFile) {
	t := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "T",
		GoImportPath: "testing",
	})
	name := serviceTestConfigProviderName(s.service.Desc)
	f.P("// ", name, " is the interface to implement to decide which resources")
	f.P("// that should be tested and how it's configured.")
	f.P("type ", name, " interface {")
	for _, resource := range s.resources {
		resourceFx := serviceResourceName(s.service.Desc, resource)
		f.P("// ", resourceFx, " should return a config, or nil, which means that the tests will be skipped.")
		f.P(resourceFx, "(t *", t, ") *", resourceTestSuiteConfigName(s.service.Desc, resource), "")
	}
	f.P("}")
	f.P()
}

func (s *serviceGenerator) generateMainTestFunction(f *protogen.GeneratedFile) {
	t := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "T",
		GoImportPath: "testing",
	})
	funcName := "test" + string(s.service.Desc.Name())
	f.P("// ", funcName, " is the main entrypoint for starting the AIP tests.")
	f.P("func ", funcName, "(t *", t, ",s ", serviceTestConfigProviderName(s.service.Desc), ") {")
	for _, resource := range s.resources {
		name := serviceResourceName(s.service.Desc, resource)
		f.P("test", name, "(t, s)")
	}
	f.P("}")
	f.P()
}

func (s *serviceGenerator) generateTestFunctions(f *protogen.GeneratedFile) {
	t := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "T",
		GoImportPath: "testing",
	})
	context := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "Context",
		GoImportPath: "context",
	})
	background := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "Background",
		GoImportPath: "context",
	})
	for _, resource := range s.resources {
		name := serviceResourceName(s.service.Desc, resource)
		f.P("func test", name, "(t *", t, ",s ", serviceTestConfigProviderName(s.service.Desc), ") {")
		f.P("t.Run(", strconv.Quote(resourceType(resource)), ", func(t *", t, ") {")
		f.P("config := s.", serviceResourceName(s.service.Desc, resource), "(t)")
		f.P("if (config == nil) {")
		f.P("t.Skip(\"Method ", name, " not implemented\")")
		f.P("}")
		f.P("if (config.Service == nil) {")
		f.P("t.Skip(\"Method ", name, ".Service() not implemented\")")
		f.P("}")
		f.P("if (config.Context == nil) {")
		f.P("config.Context = func() ", context, " { return ", background, "() }")
		f.P("}")
		f.P("config.test(t)")
		f.P("})")
		f.P("}")
		f.P()
	}
}

func (s *serviceGenerator) generateFixture(f *protogen.GeneratedFile) {
	testingT := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "T",
		GoImportPath: "testing",
	})

	service := serviceType(f, s.service, s.config)

	f.P("type ", serviceTestSuiteName(s.service.Desc), " struct {")
	f.P("T *", testingT)
	f.P("// Server to test.")
	f.P("Server  ", service)
	f.P()

	f.P("}")
	f.P()
}

func (s *serviceGenerator) generateTestMethods(f *protogen.GeneratedFile) {
	context := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "Context",
		GoImportPath: "context",
	})
	testingT := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "T",
		GoImportPath: "testing",
	})
	service := serviceType(f, s.service, s.config)
	serviceFx := serviceTestSuiteName(s.service.Desc)
	for _, resource := range s.resources {
		resourceFx := resourceTestSuiteConfigName(s.service.Desc, resource)
		f.P("func (fx ", serviceFx, ") Test", resourceType(resource), "(ctx ", context, ", options ", resourceFx, ") {")
		f.P("fx.T.Run(", strconv.Quote(resourceType(resource)), ", func(t *", testingT, ") {")
		f.P("options.Context = func() ", context, " { return ctx }")
		f.P("options.Service = func() ", service, " { return fx.Server", "}")
		f.P("options.test(t)")
		f.P("})")
		f.P("}")
		f.P()
	}
}

// serviceType returns the type the test suites call the service through: the gRPC server
// interface in the proto's Go package, or for Connect the interface generateServiceInterface writes.
func serviceType(f *protogen.GeneratedFile, service *protogen.Service, config Config) string {
	if util.Runtime(config.Runtime) == util.RuntimeConnect {
		return service.GoName
	}
	return f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       service.GoName + "Server",
		GoImportPath: service.Methods[0].Input.GoIdent.GoImportPath,
	})
}

// generateServiceInterface writes, for Connect, an interface with the service's standard
// methods in Connect v2's simple form, which both its generated client and handler implement.
func (s *serviceGenerator) generateServiceInterface(f *protogen.GeneratedFile) {
	if util.Runtime(s.config.Runtime) != util.RuntimeConnect {
		return
	}
	context := f.QualifiedGoIdent(protogen.GoIdent{
		GoName:       "Context",
		GoImportPath: "context",
	})
	standard := make(map[protoreflect.FullName]bool)
	for _, m := range method.NewMethods(s.service.Desc) {
		standard[m.Descriptor.FullName()] = true
	}
	f.P(
		"// ",
		s.service.GoName,
		" is the part of ",
		s.service.Desc.FullName(),
		" the tests call: its standard methods.",
	)
	f.P("// Both its Connect client and its Connect handler implement it.")
	f.P("type ", s.service.GoName, " interface {")
	for _, m := range s.service.Methods {
		if !standard[m.Desc.FullName()] {
			continue
		}
		f.P(m.GoName, "(", context, ", *", m.Input.GoIdent, ") (*", m.Output.GoIdent, ", error)")
	}
	f.P("}")
	f.P()
}
