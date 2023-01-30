package alfred

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (wf *Workflow) GetQuery() (string, bool) {
	if !wf.HasQuery() {
		return "", false
	}
	return wf.Query, true
}

func (wf *Workflow) HasQuery() bool {
	return len(wf.Query) > 0 && wf.Query != None
}

func (wf *Workflow) Add(item Item) {
	wf.Items = append(wf.Items, item)
}

func (wf *Workflow) AddTitleItem(title string) {
	item := Item{
		Title: title,
		Arg:   title,
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
	i.Mods[k] = item
}
