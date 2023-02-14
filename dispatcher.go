package alfred

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	cmd    = flag.String("cmd", "", "Input the cmd type")
	query  = flag.String("query", None, "Input the query string")
	debug  = flag.Bool("debug", false, "Input the debug mode")
	logger *logrus.Logger
)

func Run() {
	flag.Parse()
	if isDebug() {
		logger.SetLevel(logrus.DebugLevel)
		logger.Debug("Debug mode is open")
	}
	logger.WithFields(logrus.Fields{"cmd": *cmd, "query": *query}).Info("main exec")

	var names []string
	for name := range views {
		names = append(names, name)
	}
	logger.Info("register views:", names)
	execute(newWorkflow(*cmd, *query))
}

func execute(wf *Workflow) {
	defer handleErr(wf)
	if wf.Cmd == "" {
		panic(errors.New("command cannot be empty"))
	}
	view, ok := GetView(wf.Cmd)
	if !ok {
		panic(fmt.Errorf("unknown cmd: %s", wf.Cmd))
	}
	view.Func(wf)
	if !view.IsCli {
		if isDebug() {
			fmt.Println(wf.RenderDebug())
		} else {
			fmt.Println(wf.Render())
		}
	}
}

func handleErr(wf *Workflow) {
	if err := recover(); err != nil {
		logger.Error("main exec error:", err)
		var msg string
		switch err.(type) {
		case error:
			msg = err.(error).Error()
		case string:
			msg = err.(string)
		default:
			msg = "unexpected type"
		}
		e := errors.New(fmt.Sprintf("error: %s", msg))
		if isDebug() {
			fmt.Println(wf.RenderDebugError(e))
		} else {
			fmt.Println(wf.RenderError(e))
		}
	}
}

func isDebug() bool {
	return *debug
}

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
}
