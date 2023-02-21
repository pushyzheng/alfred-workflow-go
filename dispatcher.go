package alfred

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	cmd     = flag.String("cmd", "", "Input the cmd type")
	query   = flag.String("query", None, "Input the query string")
	debug   = flag.Bool("debug", false, "Input the debug mode")
	mainLog *logrus.Logger
)

func Run() {
	flag.Parse()
	if isDebug() {
		mainLog.SetLevel(logrus.DebugLevel)
		mainLog.Debug("Debug mode is open")
	}
	mainLog.WithFields(logrus.Fields{"cmd": *cmd, "query": *query}).Info("main exec")

	var names []string
	for name := range views {
		names = append(names, name)
	}
	mainLog.Debug("register views:", names)
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
		mainLog.Error("main exec error:", err)
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
	mainLog = logrus.New()
	mainLog.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	mainLog.Out = os.Stderr
}
