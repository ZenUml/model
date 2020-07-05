package dsl

import (
	"net/url"
	"strings"

	"goa.design/goa/v3/eval"
	"goa.design/structurizr/expr"

	// Register code generators for the structurizr plugin
	_ "goa.design/structurizr/plugin"
)

// Workspace defines the workspace containing the models and views. Workspace
// must appear exactly once in a given design. A name must be provided if a
// description is.
//
// Workspace is a top-level DSL function.
//
// Workspace takes one to three arguments. The first argument is either a string
// or a function. If the first argument is a string then an optional description
// may be passed as second argument. The last argument must be a function that
// defines the models and views.
//
// The valid syntax for Workspace is thus:
//
//    Workspace(func())
//
//    Workspace("name", func())
//
//    Workspace("name", "description", func())
//
// Examples:
//
//    // Default workspace, no description
//    var _ = Workspace(func() {
//        SoftwareSystem("My Software System")
//    })
//
//    // Workspace with given name, no description
//    var _ = Workspace("name", func() {
//        SoftwareSystem("My Software System")
//    })
//
//    // Workspace with given name and description
//    var _ = Workspace("My Workspace", "A great architecture.", func() {
//        SoftwareSystem("My Software System")
//    })
//
func Workspace(args ...interface{}) {
	_, ok := eval.Current().(eval.TopExpr)
	if !ok {
		eval.IncompatibleDSL()
		return
	}
	nargs := len(args)
	if nargs == 0 {
		eval.ReportError("missing child DSL")
		return
	}
	dsl, ok := args[nargs-1].(func())
	if !ok {
		eval.ReportError("missing child DSL (last argument must be a func)")
		return
	}
	var name, desc string
	if nargs > 1 {
		name, ok = args[0].(string)
		if !ok {
			eval.InvalidArgError("string", args[0])
		}
	}
	if nargs > 2 {
		desc, ok = args[1].(string)
		if !ok {
			eval.InvalidArgError("string", args[1])
		}
	}
	if nargs > 3 {
		eval.ReportError("too many arguments")
		return
	}
	w := &expr.Workspace{Name: name, Description: desc, Model: &expr.Model{}}
	if !eval.Execute(dsl, w) {
		return
	}
	expr.Root = w
}

// Version specifies a version number for the workspace.
//
// Version must appear in a Workspace expression.
//
// Version takes exactly one argument: the version number.
//
// Example:
//
//    var _ = Workspace(func() {
//        Version("1.0")
//    })
//
func Version(v string) {
	w, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
	} else {
		w.Version = v
	}
}

// Enterprise defines a named "enterprise" (e.g. an organisation). On System
// Landscape and System Context diagrams, an enterprise is represented as a
// dashed box. Only a single enterprise can be defined within a model.
//
// Enterprise must appear in a Workspace expression.
//
// Enterprise takes exactly one argument: the enterprise name.
//
// Example:
//
//    var _ = Workspace(func() {
//        Enterprise("Goa Design")
//    })
//
func Enterprise(e string) {
	w, ok := eval.Current().(*expr.Workspace)
	if !ok {
		eval.IncompatibleDSL()
	} else {
		w.Model.Enterprise = &expr.Enterprise{Name: e}
	}
}

// Tag defines a set of tags on the given element. Tags are used in views to
// identify group of elements that should be rendered together for example.
//
// Tag may appear in Person, SoftwareSystem, Container, Component,
// DeploymentNode, InfrastructureNode, ContainerInstance.
//
// Tag accepts the set of tag values as argument. Tag may appear multiple times
// in the same expression in which case the tags accumulate.
//
// Example:
//
//    var _ = Workspace(func() {
//        System("My system", func() {
//            Tag("sharded", "critical")
//            Tag("blue team")
//         })
//    })
//
func Tag(first string, t ...string) {
	tags := first
	if len(t) > 0 {
		tags = tags + "," + strings.Join(t, ",")
	}
	setOrAppend := func(exist, new string) string {
		if exist == "" {
			return new
		}
		return exist + "," + new
	}
	switch e := eval.Current().(type) {
	case *expr.Person:
		e.Tags = setOrAppend(e.Tags, tags)
	case *expr.SoftwareSystem:
		e.Tags = setOrAppend(e.Tags, tags)
	case *expr.Container:
		e.Tags = setOrAppend(e.Tags, tags)
	case *expr.Component:
		e.Tags = setOrAppend(e.Tags, tags)
	case *expr.DeploymentNode:
		e.Tags = setOrAppend(e.Tags, tags)
	case *expr.InfrastructureNode:
		e.Tags = setOrAppend(e.Tags, tags)
	case *expr.ContainerInstance:
		e.Tags = setOrAppend(e.Tags, tags)
	default:
		eval.IncompatibleDSL()
	}
}

// URL where more information about this element can be found.
// Or URL of health check when used within a HealthCheck expression.
//
// URL may appear in Person, SoftwareSystem, Container, Component,
// DeploymentNode, InfrastructureNode or HealthCheck.
//
// URL takes exactly one argument: a valid URL.
//
// Example:
//
//    var _ = Workspace(func() {
//        System("My system", func() {
//            URL("https://goa.design/docs/mysystem")
//        })
//    })
//
func URL(u string) {
	if _, err := url.Parse(u); err != nil {
		eval.ReportError("invalid URL %q: %s", u, err.Error())
	}
	switch e := eval.Current().(type) {
	case *expr.Person:
		e.URL = u
	case *expr.SoftwareSystem:
		e.URL = u
	case *expr.Container:
		e.URL = u
	case *expr.Component:
		e.URL = u
	case *expr.DeploymentNode:
		e.URL = u
	case *expr.InfrastructureNode:
		e.URL = u
	case *expr.HealthCheck:
		e.URL = u
	default:
		eval.IncompatibleDSL()
	}
}

// External indicates the element is external to the enterprise.
//
// External may appear in Person or SoftwareSystem.
//
// Example:
//
//    var _ = Workspace(func() {
//        System("My system", func() {
//            External()
//        })
//    })
//
func External() {
	switch e := eval.Current().(type) {
	case *expr.Person:
		e.Location = expr.LocationExternal
	case *expr.SoftwareSystem:
		e.Location = expr.LocationExternal
	default:
		eval.IncompatibleDSL()
	}
}

// Properties defines arbitrary key-value pairs. They are shown in the diagram
// tooltip and can be used to store metadata (e.g. team name).
//
// Properties must appear in Person, SoftwareSystem, Container, Component,
// DeploymentNode, InfrastructureNode or ContainerInstance.
//
// Properties accepts a single argument: a function that lists each property
// using Prop.
//
// Example:
//
//    var _ = Workspace(func() {
//        SoftwareSystem("MySystem", func() {
//            Properties(func() {
//                Prop("name", "value")
//            })
//        })
//    })
//
func Properties(dsl func()) {
	switch e := eval.Current().(type) {
	case *expr.Person:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	case *expr.SoftwareSystem:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	case *expr.Container:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	case *expr.Component:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	case *expr.DeploymentNode:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	case *expr.InfrastructureNode:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	case *expr.ContainerInstance:
		if e.Properties == nil {
			e.Properties = make(map[string]string)
		}
	default:
		eval.IncompatibleDSL()
		return
	}
	eval.Execute(dsl, eval.Current())
}

// Prop defines a single key-value pair.
//
// Prop must appear in a Properties expression.
//
// Prop accepts two arguments: the name and value.
//
// Example:
//
//    var _ = Workspace(func() {
//        SoftwareSystem("MySystem", func() {
//            Properties(func() {
//                Prop("name", "value")
//            })
//        })
//    })
//
func Prop(name, value string) {
	switch e := eval.Current().(type) {
	case *expr.Person:
		e.Properties[name] = value
	case *expr.SoftwareSystem:
		e.Properties[name] = value
	case *expr.Container:
		e.Properties[name] = value
	case *expr.Component:
		e.Properties[name] = value
	case *expr.DeploymentNode:
		e.Properties[name] = value
	case *expr.InfrastructureNode:
		e.Properties[name] = value
	case *expr.ContainerInstance:
		e.Properties[name] = value
	default:
		eval.IncompatibleDSL()
		return
	}
}