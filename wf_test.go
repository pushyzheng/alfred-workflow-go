package alfred

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueries_Get(t *testing.T) {
	wf := Workflow{
		Cmd:   "foo",
		Query: "1 2 3",
	}
	q, exists := wf.GetQueries()
	assert.True(t, exists)
	assert.Equal(t, 3, q.Len())
	assert.Equal(t, "1", q.First())
	assert.Equal(t, "2", q.Second())
	assert.Equal(t, "3", q.Third())
}

func TestWorkflow_Render(t *testing.T) {
	wf := Workflow{Cmd: "foo"}
	wf.AddTitleItem("Hello World")
	wf.AddItem("title", "body", "arg")
	result := wf.Render()
	m := make(map[string][]map[string]string)
	err := json.Unmarshal([]byte(result), &m)
	if err != nil {
		t.Error(err)
	}
	if items, ok := m["items"]; !ok {
		t.Error("the result don't has items key")
	} else {
		assert.True(t, ok)
		assert.Equal(t, 2, len(items))
		item := items[0]
		assert.Equal(t, "Hello World", item["title"])
		assert.Equal(t, "", item["subtitle"])
		assert.Equal(t, "Hello World", item["arg"])

		item2 := items[1]
		assert.Equal(t, "title", item2["title"])
		assert.Equal(t, "body", item2["subtitle"])
		assert.Equal(t, "arg", item2["arg"])
	}
}
