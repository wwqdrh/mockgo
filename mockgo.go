package mockgo

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/brianvoe/gofakeit/v7/source"
	"github.com/tidwall/gjson"
	"github.com/wwqdrh/gokit/logger"
)

var increment = gofakeit.NewFaker(source.NewDumb(1), true)

// 解析json格式数据，并且根据value值对应的生成规则生成mock数据
func Generate(in string) (string, error) {
	result := map[string]interface{}{}

	inJson := gjson.Parse(in)
	inJson.ForEach(func(key, value gjson.Result) bool {
		mocklabel, mockdata := mockData(key.String(), value)
		result[mocklabel] = mockdata
		return true
	})

	data, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func mockData(key string, value gjson.Result) (string, interface{}) {
	idx := strings.Index(key, "|")
	if idx == -1 {
		return key, mockSingle(value)
	}

	num, err := strconv.ParseInt(key[idx+1:], 10, 64)
	if err != nil {
		logger.DefaultLogger.Error(err.Error())
		num = 1
	}
	label := key[0:idx]
	if num > 1 || key[idx+1] == '+' {
		return label, mockArray(value, int(num))
	} else {
		return label, mockSingle(value)
	}
}

func mockArray(value gjson.Result, num int) []interface{} {
	result := []interface{}{}
	if value.IsArray() {
		res := value.Array()
		for i := 0; i < num; i++ {
			_, mockdata := mockData("", res[rand.Intn(len(res))])
			result = append(result, mockdata)
		}
	}
	return result
}

func mockSingle(value gjson.Result) interface{} {
	if value.IsArray() {
		res := value.Array()
		_, mockData := mockData("", res[rand.Intn(len(res))])
		return mockData
	} else if value.IsObject() {
		res := value.Map()
		curresult := map[string]interface{}{}
		for k, v := range res {
			mockLabel, mockData := mockData(k, v)
			curresult[mockLabel] = mockData
		}
		return curresult
	} else if value.Type == gjson.String {
		valueStr := value.String()
		if valueStr == "@increment" {
			return increment.Int64() - 1
		} else if valueStr == "@date" {
			return gofakeit.Date().String()
		} else if valueStr == "@cname" {
			return gofakeit.FirstName()
		} else if strings.HasPrefix(valueStr, "@natural") {
			// 从@natural(9781782910284, 9981782910284)格式解析最小和最大值, 并且返回为字符串格式
			valueStr = strings.TrimPrefix(valueStr, "@natural")
			valueStr = valueStr[1 : len(valueStr)-1]
			split := strings.Split(valueStr, ",")
			if len(split) != 2 {
				return fmt.Sprint(gofakeit.Int64())
			}
			min, err := strconv.Atoi(strings.TrimSpace(split[0]))
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return fmt.Sprint(gofakeit.Int64())
			}
			max, err := strconv.Atoi(strings.TrimSpace(split[1]))
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return fmt.Sprint(gofakeit.Int64())
			}
			return fmt.Sprint(gofakeit.IntRange(min, max))
		} else if strings.HasPrefix(valueStr, "@integer") {
			// 从@integer(60, 100)格式解析最小和最大值
			valueStr = strings.TrimPrefix(valueStr, "@integer")
			valueStr = valueStr[1 : len(valueStr)-1]
			split := strings.Split(valueStr, ",")
			if len(split) != 2 {
				return gofakeit.Int64()
			}
			min, err := strconv.Atoi(strings.TrimSpace(split[0]))
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return gofakeit.Int64()
			}
			max, err := strconv.Atoi(strings.TrimSpace(split[1]))
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return gofakeit.Int64()
			}
			return gofakeit.IntRange(min, max)
		} else if strings.HasPrefix(valueStr, "@float") {
			// 从@float(60, 100, 2, 2)格式解析最小和最大值, 小数点后位数最小和最大值
			valueStr = strings.TrimPrefix(valueStr, "@float")
			valueStr = valueStr[1 : len(valueStr)-1]
			split := strings.Split(valueStr, ",")
			if len(split) != 4 {
				return gofakeit.Float64()
			}
			min, err := strconv.ParseFloat(strings.TrimSpace(split[0]), 64)
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return gofakeit.Float64()
			}
			max, err := strconv.ParseFloat(strings.TrimSpace(split[1]), 64)
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return gofakeit.Float64()
			}
			precisionMin, err := strconv.Atoi(strings.TrimSpace(split[2]))
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return gofakeit.Float64()
			}
			precisionMax, err := strconv.Atoi(strings.TrimSpace(split[3]))
			if err != nil {
				logger.DefaultLogger.Warn(err.Error())
				return gofakeit.Float64()
			}
			result := gofakeit.Float64Range(min, max)
			// 计算result保留位数
			precision := rand.Intn(precisionMax-precisionMin+1) + precisionMin
			return strconv.FormatFloat(result, 'f', precision, 64)
		} else {
			return valueStr
		}
	}
	return nil
}
