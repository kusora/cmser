package cmd

import (
	"encoding/json"
	"net/url"
	"github.com/kusora/cmser/util"
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


func  CalcuGroupMeanRelation(group []string) float64 {
	return 0.0
}

func GetRelations(key string, values []string) ([]float64, error) {
	data, _ := json.Marshal(values)

	resp := make([]byte, 0)
	status, resp, err := util.HttpPostUrlValuesRawResult(http.DefaultClient, "http://192.168.59.100:8080/api/similarity", url.Values{
		"key":  []string{key},
		"value": []string{string(data)},
	})
	if err != nil || status != http.StatusOK {
		dlog.Error("%+v, %+v", err, status)
		return nil, err
	}

	result := make([]float64, 0)
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}