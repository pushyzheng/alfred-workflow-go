package alfred

import (
	"log"
	"strings"
)

const (
	ListCommandsName   = "list"
	SearchContainsMode = "contains"
	SearchStartMode    = "start"
)

type ViewFunc func(wf *Workflow)

type View struct {
	Name       string
	Func       ViewFunc
	Desc       string
	NeedsQuery bool
	IsCli      bool
}

// views: the registry implemented with memory map.
var views = map[string]*View{}

// RegisterView Registers a view function to registry with name
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
	v := View{
		Name: name,
		Func: fn,
	}
	Register(&v)
}

// Register Registers a View to registry
func Register(view *View) {
	views[view.Name] = view
}

func getView(name string) (*View, bool) {
	vh, ok := views[name]
	return vh, ok
}

func searchView(q string, mode string) []*View {
	if mode == "" {
		mode = SearchContainsMode
	}
	var res []*View
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

// init Registers some built-in views.
func init() {
	RegisterView(ListCommandsName, DisplayCommands)

	// cli
	Register(&View{
		Name: schedulerRunCmd,
		Func: func(wf *Workflow) {
			err := scheduler.loop(-1)
			if err != nil {
				log.Fatalln(err)
			}
		},
		IsCli: true,
	})
	Register(&View{
		Name:  "hook",
		Func:  hookCli,
		IsCli: true,
	})
	Register(&View{
		Name:  schedulerListCmd,
		Func:  displayTask,
		IsCli: true,
	})
	Register(&View{
		Name:  "cache_delete",
		Func:  deleteCache,
		IsCli: true,
	})
	Register(&View{
		Name:  "cache_get",
		Func:  getCache,
		IsCli: true,
	})
}
