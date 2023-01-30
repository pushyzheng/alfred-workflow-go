package alfred

import (
	"log"
	"strings"
)

const (
	SearchContainsMode = "contains"
	SearchStartMode    = "start"
)

type ViewFunc func(wf *Workflow)

type ViewHandler struct {
	Name string
	Func ViewFunc
}

var views = map[string]ViewHandler{}

func RegisterView(name string, fn ViewFunc) {
	if len(name) == 0 {
		log.Fatalln("the name cannot be empty")
	}
	if fn == nil {
		log.Fatalln("the view function cannot be null")
	}
	if _, exists := views[name]; exists {
		log.Fatalln("the view function is registered already")
	}
	views[name] = ViewHandler{
		Name: name,
		Func: fn,
	}
}

func GetView(name string) (ViewHandler, bool) {
	vh, ok := views[name]
	return vh, ok
}

func SearchView(q string, mode string) []ViewHandler {
	if mode == "" {
		mode = SearchContainsMode
	}
	var res []ViewHandler
	for name, h := range views {
		var ok bool
		if mode == SearchContainsMode {
			ok = strings.Contains(name, q)
		} else if mode == SearchStartMode {
			ok = strings.HasPrefix(name, q)
		}
		if ok {
			res = append(res, h)
		}
	}
	return res
}
