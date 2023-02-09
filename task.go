package alfred

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pushyzheng/diskache"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const (
	defaultLoopInterval    = time.Second
	schedulerTimeTableName = "scheduler-time-table"
	schedulerRunCmd        = "scheduler_run"
	schedulerListCmd       = "scheduler_list"
)

var taskLog = NewFileLogger("task")

type TaskScheduler struct {
	LoopInterval time.Duration
	tasks        map[string]*Task
	cache        *diskache.Diskache
}

type TaskOption struct {
	Interval time.Duration
	Executor func(task *Task)
}

type Task struct {
	Name   string
	Op     *TaskOption
	Logger *logrus.Logger
}

var scheduler *TaskScheduler

func RegisterTask(name string, op *TaskOption) error {
	if scheduler == nil {
		return nil
	}
	if op.Executor == nil {
		return errors.New("the executor func cannot be null")
	}
	if op.Interval < 0 {
		return fmt.Errorf("illegal interval: %d", op.Interval)
	}
	task := Task{
		Name:   name,
		Op:     op,
		Logger: taskLog,
	}
	scheduler.tasks[name] = &task
	return nil
}

func (s *TaskScheduler) loop(times int) error {
	if len(s.tasks) == 0 {
		taskLog.Warn("no any tasks, program exits.")
		return nil
	}
	var i = 0
	for {
		if times != -1 && i >= times {
			break
		}
		s.runTasks()
		i++
		taskLog.Infof("scheduler loop, next time: %d", GetTimestamp()+s.LoopInterval.Milliseconds())
		time.Sleep(s.LoopInterval)
	}
	return nil
}

func (s *TaskScheduler) runTasks() {
	timeTable, err := s.getTimeTable()
	if err != nil {
		log.Fatalln(err)
	}
	for name, task := range s.tasks {
		now := GetTimestamp()
		if last, ok := timeTable[name]; ok {
			if (now - last) >= task.Op.Interval.Milliseconds() {
				taskLog.WithFields(logrus.Fields{
					"name":     name,
					"last":     last,
					"interval": task.Op.Interval.Milliseconds(),
				}).Info("the task should be scheduled")
				task.run()
				timeTable[name] = now
			}
		} else {
			timeTable[name] = now
		}
	}
	_ = s.updateTimeTable(timeTable)
}

func (s *TaskScheduler) getTimeTable() (map[string]int64, error) {
	data, exists := s.cache.Get(schedulerTimeTableName)
	if !exists {
		return map[string]int64{}, nil
	}
	var res map[string]int64
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TaskScheduler) updateTimeTable(data map[string]int64) error {
	return s.cache.SetJson(schedulerTimeTableName, data)
}

func (s *TaskScheduler) getTask(name string) *Task {
	if t, ok := s.tasks[name]; ok {
		return t
	}
	return nil
}

func (t *Task) run() {
	taskLog.WithField("name", t.Name).Info("run task")
	if t.Op.Executor == nil {
		return
	}
	start := GetTimestamp()
	go func() {
		t.Op.Executor(t)
		taskLog.WithFields(logrus.Fields{
			"name": t.Name,
			"cost": GetTimestamp() - start,
		}).Info("task finished")
	}()
}

func newTaskScheduler() *TaskScheduler {
	opts := diskache.Opts{
		Directory: "cache",
	}
	var err error
	c, err := diskache.New(&opts)
	if err != nil {
		log.Fatalln(err)
	}
	return &TaskScheduler{
		LoopInterval: defaultLoopInterval,
		tasks:        make(map[string]*Task),
		cache:        c,
	}
}

func init() {
	scheduler = newTaskScheduler()
}
