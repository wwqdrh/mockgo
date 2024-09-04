package mockgo

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wwqdrh/gokit/logger"
)

type JsonMockHandler struct {
	Url     string
	Mock    string
	Handler http.HandlerFunc
}

// 的所在路径作为url，然后文件内容作为mock数据
// 使用http注册url以及返回对应的mock数据

func GetHandler(p string) []*JsonMockHandler {
	result := getJsonList(p)
	for _, item := range result {
		item.Handler = handlerMock(*item)
	}
	return result
}

func handlerMock(m JsonMockHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		mockData, err := Generate(m.Mock)
		if err != nil {
			logger.DefaultLogger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockData))
		}
	}
}

// 传入一个本地文件的路径, 将该文件夹下的json文件, 返回相对路径以及文件内容
func getJsonList(target string) []*JsonMockHandler {
	if target == "" {
		return []*JsonMockHandler{}
	}

	result := []*JsonMockHandler{}
	filepath.Walk(target, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(p) == ".json" {
			url, err := filepath.Rel(target, p)
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return nil
			}
			content, err := os.ReadFile(p)
			if err != nil {
				logger.DefaultLogger.Error(err.Error())
			}
			result = append(result, &JsonMockHandler{Url: strings.TrimSuffix(url, ".json"), Mock: string(content)})
			return nil
		}
		return nil
	})
	return result
}
