package cache_data

import (
	"encoding/json"
	"github.com/pushyzheng/alfred-workflow-go"
	"io"
	"log"
	"net/http"
)

type HttpBinResp struct {
	Origin  string `json:"origin"`
	Url     string `json:"url"`
	Headers struct {
		Accept         string `json:"accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		UserAgent      string `json:"User-Agent"`
	}
}

func Foo(wf *alfred.Workflow) {
	data, exists := alfred.CacheData("http-bin", GetHttpBin)
	if exists {
		log.Println("hit cache")
	}
	wf.AddTitleItem(data.Origin)
	wf.AddTitleItem(data.Url)
	wf.AddTitleItem(data.Headers.AcceptEncoding)
	wf.AddTitleItem(data.Headers.Accept)
	wf.AddTitleItem(data.Headers.UserAgent)
}

func GetHttpBin() HttpBinResp {
	resp, err := http.Get("http://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res HttpBinResp
	err = json.Unmarshal(b, &res)
	if err != nil {
		panic(err)
	}
	return res
}

func main() {
	alfred.Run()
}

func init() {
	alfred.RegisterView("foo", Foo)
}
