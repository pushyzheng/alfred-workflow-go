package alfred

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

func Run() {
	var names []string
	for name := range views {
		names = append(names, name)
	}
	log.Println("register views:", names)

	cmd := flag.String("cmd", "", "Input the cmd type")
	query := flag.String("query", None, "Input the query string")
	flag.Parse()
	log.Printf("main exec, cmd = %s, query = %s", *cmd, *query)

	wf := Workflow{
		Cmd:   *cmd,
		Query: *query,
	}
	execute(&wf)
}

func execute(wf *Workflow) {
	defer handleErr(wf)
	if wf.Cmd == "" {
		panic(errors.New("command cannot be empty"))
	}
	handler, ok := GetView(wf.Cmd)
	if !ok {
		panic(fmt.Errorf("unknown cmd: %s", wf.Cmd))
	}
	handler.Func(wf)
	fmt.Println(wf.Render())
}

func handleErr(wf *Workflow) {
	if err := recover(); err != nil {
		log.Println("main exec error:", err)
		var msg string
		switch err.(type) {
		case error:
			msg = err.(error).Error()
		case string:
			msg = err.(string)
		default:
			msg = "unexpected type"
		}
		fmt.Println(wf.RenderError(errors.New(fmt.Sprintf("error: %s", msg))))
	}
}
