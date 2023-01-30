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
	defer handleErr(&wf)
	if *cmd == "" {
		panic(errors.New("command cannot be empty"))
	}
	handler, ok := GetView(*cmd)
	if !ok {
		panic(fmt.Errorf("unknown cmd: %s", *cmd))
	}
	handler.Func(&wf)
	fmt.Println(wf.Render())
}

func handleErr(wf *Workflow) {
	if err := recover(); err != nil {
		log.Println("main exec error:", err)
		switch err.(type) {
		case error:
			fmt.Println(wf.RenderError(err.(error)))
		case string:
			fmt.Println(wf.RenderError(errors.New(err.(string))))
		default:
			fmt.Println(wf.RenderError(errors.New("unexpected type")))
		}
	}
}
