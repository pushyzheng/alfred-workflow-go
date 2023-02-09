package alfred

import (
	"fmt"
	"log"
)

func deleteCache(wf *Workflow) {
	if k, ok := wf.GetQuery(); !ok {
		log.Fatalln("error: no cache key")
	} else {
		ok = cache.Delete(k)
		if !ok {
			log.Fatalln("fail to delete cache, key =", k)
		} else {
			log.Println("delete cache succeed, key =", k)
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
