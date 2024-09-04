# mockgo

this is a toolkit for generate fake data like mockjs

> base on https://github.com/brianvoe/gofakeit

## usage

> go get github.com/wwqdrh/mockgo

```go
import "github.com/wwqdrh/mockgo"

func main() {
    mockData, err := mockgo.Generate(`
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
      "pages": "@integer(60, 100)",
	  "images|+1" : [
	 	"https://i.loli.net/2021/01/18/LSFhuWVlMGcAUOB.png",
        "https://i.loli.net/2021/01/18/9F3ND7KRdnuLxaX.png",
	  ]
	}
	]
}
	`)
    // ...
}
```