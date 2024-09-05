package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/wwqdrh/gokit/logger"
	"github.com/wwqdrh/mockgo"
)

// 命令行参数
var (
	path string
	port int
)

func init() {
	flag.StringVar(&path, "path", "", "path")
	flag.IntVar(&port, "port", 8080, "port")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	for _, item := range mockgo.GetHandler(path, true) {
		logger.DefaultLogger.Infox("register url: /%s", nil, item.Url)

		mux.HandleFunc("/"+item.Url, item.Handler)
	}
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}
