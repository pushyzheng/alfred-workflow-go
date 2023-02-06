package alfred

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pushyzheng/diskache"
	"log"
	"time"
)

const (
	defaultLoopInterval    = time.Second
	schedulerTimeTableName = "scheduler-time-table"
	schedulerRunCmd        = "scheduler_run"
	schedulerListCmd       = "scheduler_list"
)

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
	Name string
	Op   *TaskOption
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
		Name: name,
		Op:   op,
	}
	scheduler.tasks[name] = &task
	return nil
}

func (s *TaskScheduler) loop(times int) error {
	if len(s.tasks) == 0 {
		log.Println("[warn] no any tasks, program exits.")
		return nil
	}
	var i = 0
	for {
		if times != -1 && i >= times {
			break
		}
		s.runTasks()
		i++
		log.Printf("scheduler loop, next time: %d", GetTimestamp()+s.LoopInterval.Milliseconds())
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
				log.Printf("the task should be scheduled, last = %d, interval = %d", last, task.Op.Interval.Milliseconds())
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
	log.Printf("run task, name = %s", t.Name)
	if t.Op.Executor == nil {
		return
	}
	start := GetTimestamp()
	go func() {
		t.Op.Executor(t)
		log.Printf("task finished, name = %s, cost = %dms", t.Name, GetTimestamp()-start)
	}()
}

func displayTask(wf *Workflow) {
	table, err := scheduler.getTimeTable()
	if err != nil {
		panic(err)
	}
	i := 0
	for name, task := range scheduler.tasks {
		if last, ok := table[name]; ok {
			fmt.Printf("[%d] %s - next time: %d\n", i, name, last+task.Op.Interval.Milliseconds())
		} else {
			fmt.Printf("[%d] %s - next time: %d\n", i, name, GetTimestamp()+task.Op.Interval.Milliseconds())
		}
		i++
	}
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
