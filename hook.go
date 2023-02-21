package alfred

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

var (
	workerLog = NewFileLogger("worker")
	hooks     = map[string]HookFn{}
)

type HookCtx struct {
	Args   string
	Logger *logrus.Entry
}

type HookFn func(hookCtx *HookCtx)

// RegisterWorker this is beta feature
func RegisterWorker(name string, fn HookFn) {
	if len(name) == 0 {
		panic("the hook name cannot be empty")
	}
	if fn == nil {
		panic("the hook function cannot be null")
	}
	hooks[name] = fn
}

// Trigger this is beta feature
func Trigger(name string, args string) error {
	mainName := os.Args[0]
	err := startSubprocess(mainName, "-cmd", "hook", "-query", fmt.Sprintf("%s %s", name, args))
	if err != nil {
		return err
	}
	return nil
}

func hookCli(wf *Workflow) {
	if qs, ok := wf.GetQueries(); !ok {
		workerLog.Warn("no specific hook")
	} else {
		name := qs.First()
		var args string
		if qs.Len() >= 2 {
			args = qs.Second()
		}
		if fn, ok := hooks[name]; !ok {
			workerLog.Warnf("the hook cannot be found: %s", name)
		} else {
			logEntry := workerLog.WithField("name", name)
			logEntry.Info("exec hook, args: ", args)
			hookCtx := HookCtx{
				Args:   args,
				Logger: logEntry,
			}
			fn(&hookCtx)
		}
	}
}

func startSubprocess(name string, arg ...string) error {
	workerLog.WithFields(logrus.Fields{
		"name": name,
		"arg":  arg,
	}).Infof("Start worker")

	cmd := exec.Command(name, arg...)
	cmd.Stdout = workerLog.Out
	cmd.Stderr = workerLog.Out
	cmd.Stdin = os.Stdin
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
