package cmd

import (
	"encoding/json"
	"net/url"
	"github.com/kusora/cmser/util"
	qu "qxf-backend/util"
	"sync"
	"net/http"
	"github.com/kusora/dlog"
	"io/ioutil"
)

/*
	文本聚类, 先放出来看看
 */

func Groups(input []string) [][]string {
	// 类似于找团算法
	//ioutil.WriteFile(name, data, 0666)
	relation := make([][]float64, len(input), len(input))
	for id, str := range input {
		relations, err := GetRelations(str, input)
		if err != nil {
			dlog.Error("failed to get relation for %s, %+v", str, err)
			continue
		}

		relation[id] = relations
	}

	relationData, _ := json.Marshal(relation)
	err := ioutil.WriteFile("relations", relationData, 0666)
	if err != nil {
		dlog.Error("failed to save relation data")
		return nil
	}

	// todo 直接修改relationData， 将》0.8的点修改为1, 这样进行聚类

	return nil
}

func CalcuGroupMeanRelation(group []string) float64 {
	return 0.0
}

var BATCH_SIZE = 400
var MAX_LENGTH = 65088

func GetRelations(key string, values []string) ([]float64, error) {
	result := make([]float64, len(values))
	end := 0
	lock := &sync.Mutex{}
	executor := qu.NewExecutor(10, 10000)
	executor.Start()
	roundRobin := 0
	for start := 0; start < len(values); start = end {
		end = start + BATCH_SIZE
		if end > len(values) {
			end = len(values)
		}
		data, _ := json.Marshal(values[start:end])
		if len(key) + len(string(data)) > MAX_LENGTH {
			end = start + BATCH_SIZE - 200
			data, _ = json.Marshal(values[start:end])
		}

		newStart := start
		executor.AddTask(func() {
			resp := make([]byte, 0)
			//status, resp, err := util.HttpPostUrlValuesRawResult(http.DefaultClient, "http://192.168.59.100:8080/api/similarity", url.Values{
			server := "http://10.143.248.75:666/nlnop/api/similarity"
			if roundRobin%2 == 0 {
				server = "http://localhost:8080/api/similarity"
			}
			roundRobin++
			status, resp, err := util.HttpPostUrlValuesRawResult(http.DefaultClient, server, url.Values{
				"key":  []string{key},
				"value": []string{string(data)},
			})
			if err != nil || status != http.StatusOK {
				dlog.Error("%+v, %+v", err, status)
				return
			}

			part := make([]float64, end - start)
			err = json.Unmarshal(resp, &part)
			if err != nil {
				dlog.Error("marshal error %+v", err)
				return
			}
			lock.Lock()
			defer lock.Unlock()
			for i, _ := range part {
				result[newStart + i] = part[i]
			}
		})
	}
	executor.Close()
	return result, nil
}