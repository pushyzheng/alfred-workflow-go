package alfred

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

const None = "none"

type StringMap map[string]string

type AnyMap map[string]interface{}

type Workflow struct {
	Cmd   string
	Query string
	Items []Item
}

type Item struct {
	Title        string             `json:"title"`
	SubTitle     string             `json:"subtitle"`
	Arg          string             `json:"arg"`
	Autocomplete string             `json:"autocomplete"`
	Icon         string             `json:"icon"`
	QuickLookUrl string             `json:"quicklookurl"`
	Mods         map[string]ModItem `json:"mods"`
}

type ModItem struct {
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
	Arg      string `json:"arg"`
}

type Queries struct {
	Values []string
}

func (wf *Workflow) GetQuery() (string, bool) {
	if !wf.HasQuery() {
		return "", false
	}
	return wf.Query, true
}

func (wf *Workflow) GetQueries() (Queries, bool) {
	var queries Queries
	if !wf.HasQuery() {
		return queries, false
	}
	q := wf.Query
	if len(q) == 0 {
		return queries, false
	}
	values := strings.Split(q, " ")
	if len(values) == 0 {
		return queries, false
	}
	return Queries{Values: values}, true
}

func (wf *Workflow) HasQuery() bool {
	return len(wf.Query) > 0 && wf.Query != None
}

func (wf *Workflow) MustEnv(k string) string {
	if s, err := wf.GetEnv(k); err != nil {
		panic(err)
	} else {
		return s
	}
}

func (wf *Workflow) GetEnv(k string) (string, error) {
	if len(k) == 0 {
		return "", errors.New("the key cannot be empty")
	}
	v := os.Getenv(k)
	if len(v) == 0 {
		return "", fmt.Errorf("the '%s' key not in env", k)
	}
	return v, nil
}

func (wf *Workflow) Add(item Item) {
	wf.Items = append(wf.Items, item)
}

func (wf *Workflow) AddTitleItem(title string) {
	item := Item{
		Title:        title,
		Arg:          title,
		Autocomplete: title,
	}
	wf.Add(item)
}

func (wf *Workflow) AddItem(title string, subtitle string, arg string) {
	item := Item{
		Title:    title,
		SubTitle: subtitle,
		Arg:      arg,
	}
	wf.Add(item)
}

func (wf *Workflow) Render() string {
	if len(wf.Items) == 0 {
		errMsg := NoAnyResult
		if q, exists := wf.GetQuery(); exists {
			errMsg = fmt.Sprintf("%s: %s", NoAnyResult, q)
		}
		return wf.RenderError(errors.New(errMsg))
	}
	m := AnyMap{}
	m["items"] = wf.Items
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (wf *Workflow) RenderError(err error) string {
	m := AnyMap{}
	item := StringMap{}
	item["title"] = err.Error()
	m["items"] = []StringMap{item}
	b, err := json.Marshal(m)
	if err != nil {
		// join json string directly
		return fmt.Sprintf("{\"items\":[{\"title\":\"%s\"}]}", err.Error())
	}
	return string(b)
}

func (i *Item) SetCmd(item ModItem) {
	i.setMod("cmd", item)
}

func (i *Item) SetAlt(item ModItem) {
	i.setMod("alt", item)
}

func (i *Item) SetShift(item ModItem) {
	i.setMod("shift", item)
}

func (i *Item) setMod(k string, item ModItem) {
	if i.Mods == nil {
		i.Mods = make(map[string]ModItem)
	}
	if item.Arg == "" {
		item.Arg = item.SubTitle
	}
	i.Mods[k] = item
}

func (q *Queries) Len() int {
	return len(q.Values)
}

func (q *Queries) First() string {
	return q.Values[0]
}

func (q *Queries) Second() string {
	if q.Len() < 2 {
		panic("lack of second param")
	}
	return q.Values[1]
}

func (q *Queries) Third() string {
	if q.Len() < 3 {
		panic("lack of third param")
	}
	return q.Values[2]
}

func (q *Queries) Get(i int) string {
	if q.Len() != i+1 {
		panic(fmt.Sprintf("lack of param: %d", i))
	}
	return q.Values[i]
}
