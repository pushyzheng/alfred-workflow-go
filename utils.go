package alfred

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	StandardDatetimeLayout = "2006-01-02 15:04:05"
	StandardDateLayout     = "2006-01-02"
	LogBaseDir             = "./logs"
)

func GetTimestamp() int64 {
	return time.Now().UnixMilli()
}

func FromTimestamp(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func FormatTimestamp(ts int64) string {
	t := FromTimestamp(ts)
	return t.Format(StandardDatetimeLayout)
}

func FormatTime(t time.Time) string {
	return t.Format(StandardDateLayout)
}

func ParseTime(s string) time.Time {
	date, err := time.Parse(StandardDatetimeLayout, s)
	if err != nil {
		panic(err)
	}
	return date
}

func FilterSlice[T any](data []T, fn func(v T) bool) []T {
	var result []T
	for _, v := range data {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func FilterMap[T any](data map[string]T, fn func(k string, v T) bool) map[string]T {
	result := make(map[string]T)
	for k, v := range data {
		if fn(k, v) {
			result[k] = v
		}
	}
	return result
}

func NewFileLogger(name string) *logrus.Logger {
	if _, err := os.Stat(LogBaseDir); os.IsNotExist(err) {
		err := os.MkdirAll(LogBaseDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
	logPath := path.Join(LogBaseDir, fmt.Sprintf("%s.log", name))
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	nl := logrus.New()
	if err == nil {
		nl.Out = file
	} else {
		fmt.Printf("error: Failed to log to file, name = %s", name)
	}
	nl.WithField("name", name).Infof("init log succeed")
	return nl
}
