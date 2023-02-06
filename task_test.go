package alfred

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTask_Run(t *testing.T) {
	runFlag := false
	fn := func(task *Task) {
		runFlag = true
		assert.Equal(t, "foo", task.Name)
	}
	err := RegisterTask("foo", &TaskOption{
		Interval: 2 * time.Second,
		Executor: fn,
	})
	if err != nil {
		t.Error(err)
	}
	err = scheduler.loop(3)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, runFlag)
}
