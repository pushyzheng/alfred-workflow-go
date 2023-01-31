package alfred

import (
	"os"
	"sort"
	"strings"
)

func DisplayCommands(wf *Workflow) {
	if q, ok := wf.GetQueries(); ok {
		if q.Len() == 1 {
			views = FilterMap(views, func(name string, v *View) bool {
				return name != ListCommandsName && strings.Contains(name, q.First())
			})
			// complete match, invoke view directly
			if len(views) == 1 && views[q.First()] != nil && !views[q.First()].NeedsQuery {
				invoke(q)
			}
		} else {
			invoke(q)
		}
	}
	display(wf)
}

func display(wf *Workflow) {
	var names []string
	for k := range views {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, name := range names {
		if name == ListCommandsName {
			continue
		}
		view := views[name]
		item := Item{
			Title:        name,
			SubTitle:     view.Desc,
			Autocomplete: name,
			Arg:          name,
		}
		wf.Add(item)
	}
}

func invoke(q Queries) {
	name := q.First()
	var newWf Workflow
	if q.Len() > 1 {
		newWf = Workflow{
			Cmd:   name,
			Query: q.Second(),
		}
	} else {
		newWf = Workflow{
			Cmd:   name,
			Query: "",
		}
	}
	execute(&newWf)
	os.Exit(0)
}
