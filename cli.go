package alfred

import (
	"fmt"
)

func getCache(wf *Workflow) {
	if k, ok := wf.GetQuery(); !ok {
		wf.Logger.Error("error: the cache key cannot be empty")
	} else {
		v, ok := cache.GetStr(k)
		if !ok {
			wf.Logger.Warnf("The cache of %s don't exists", k)
		} else {
			wf.Logger.Infof("The value of cache(%s) is: \n%s", k, v)
		}
	}
}

func deleteCache(wf *Workflow) {
	if k, ok := wf.GetQuery(); !ok {
		wf.Logger.Error("error: the cache key cannot be empty")
	} else {
		ok = cache.Delete(k)
		if !ok {
			wf.Logger.Error("fail to delete cache, key =", k)
		} else {
			wf.Logger.Info("delete cache succeed, key =", k)
		}
	}
}

func displayTask(_ *Workflow) {
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
