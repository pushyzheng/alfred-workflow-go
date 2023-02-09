package alfred

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var httpLog = NewFileLogger("https")

func HttpGet(url string, headers map[string]string) string {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	resp := doRequest(req, headers)
	httpLog.WithFields(logrus.Fields{
		"method":   "GET",
		"url":      url,
		"response": resp,
	}).Infof("HttpGet")
	return resp
}

func HttpGetJsonMap(url string, headers map[string]string) map[string]any {
	data := make(map[string]any)
	HttpGetJson(url, headers, &data)
	return data
}

func HttpGetJsonMapArray(url string, headers map[string]string) []map[string]any {
	var data []map[string]any
	HttpGetJson(url, headers, &data)
	return data
}

func HttpGetJson[T any](url string, headers map[string]string, data *T) {
	body := HttpGet(url, headers)
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}
}

func HttpPost(url string, headers map[string]string, raw string) string {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(raw))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	resp := doRequest(req, headers)
	httpLog.WithFields(logrus.Fields{
		"method":   "POST",
		"url":      url,
		"raw":      raw,
		"response": resp,
	}).Infof("HttpPost")
	return resp
}

func HttpPostJson[T any](url string, headers map[string]string, raw string, data *T) {
	resp := HttpPost(url, headers, raw)
	err := json.Unmarshal([]byte(resp), &data)
	if err != nil {
		panic(err)
	}
}

func HttpPostJsonBody(url string, headers map[string]string, jsonBody any) string {
	b, err := json.Marshal(jsonBody)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return doRequest(req, headers)
}

func doRequest(req *http.Request, headers map[string]string) string {
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(b)
}
