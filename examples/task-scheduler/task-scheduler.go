package main

import (
	"github.com/pushyzheng/alfred-workflow-go"
	"log"
	"time"
)

func Foo(wf *alfred.Workflow) {
	data, exists := alfred.GetCacheData("cache-key")
	if !exists {
		wf.AddTitleItem("no data")
	} else {
		wf.AddTitleItem(string(data))
	}
}

func refreshCache(_ *alfred.Task) {
	time.Sleep(time.Second)
	body := "Hello World"
	alfred.SetCacheData("cache-key", -1, []byte(body))
}

func main() {
	alfred.Run()
}

func init() {
	alfred.RegisterView("foo", Foo)

	err := alfred.RegisterTask("foo", &alfred.TaskOption{
		Interval: time.Second * 3,
		Executor: refreshCache,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
