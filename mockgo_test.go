package mockgo

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestGenerate(t *testing.T) {
	mockData, err := Generate(`
{
  "data|10": [
    {
      "id": "@increment",
      "title|1": [
        "Church in the Glass",
        "Destiny's Hero",
		],
      "pubdate": "@date",
	  "price": "@float(60, 100, 2, 2)",
	  "isbn": "@natural(9781782910284, 9981782910284)",
      "pages": "@integer(60, 100)",
	  "images|+1" : [
	 	"https://i.loli.net/2021/01/18/LSFhuWVlMGcAUOB.png",
        "https://i.loli.net/2021/01/18/9F3ND7KRdnuLxaX.png",
	  ]
	}
	]
}
	`)
	if err != nil {
		t.Error(err)
		return
	}
	if gjson.Get(mockData, "data.#").Int() != 10 {
		t.Error("Generate data.# fail")
		return
	}
	if val := gjson.Get(mockData, "data.1.id"); val.Type != gjson.Number && int(val.Num) != 1 {
		t.Error("Generate data.1.id fail")
		return
	}
	if val := gjson.Get(mockData, "data.1.images"); !val.IsArray() {
		t.Error("Generate data.images fail")
		return
	}
	if val := gjson.Get(mockData, "data.1.price"); val.Type != gjson.String || len(val.Str) != 5 {
		t.Error("Generate data.price fail")
		return
	}
	if val := gjson.Get(mockData, "data.1.isbn"); val.Type != gjson.String || len(val.Str) != 13 {
		t.Error("Generate data.price fail")
		return
	}
	if val := gjson.Get(mockData, "data.1.pages"); val.Type != gjson.Number || val.Num < 60 || val.Num > 100 {
		t.Error("Generate data.pages fail")
		return
	}
}
