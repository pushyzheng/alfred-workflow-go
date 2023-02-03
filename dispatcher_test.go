package alfred

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecute(t *testing.T) {
	wf := Workflow{Cmd: "foo"}
	RegisterView("foo", func(wf *Workflow) {
		wf.AddTitleItem("Hello World")
	})
	execute(&wf)

	assert.Equal(t, 1, len(wf.Items))
	assert.Equal(t, "Hello World", wf.Items[0].Title)
}
