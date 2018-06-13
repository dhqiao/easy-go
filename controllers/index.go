package controllers

import (
	"easy-go/data"
	"easy-go/library"
	"easy-go/request"
	"easy-go/request/form"
	"fmt"
	"net/http"
	"net/url"
)

type IndexController struct {
	BaseController
}

var redis = data.RedisClient

// 测试用接口
func (index *IndexController) Test(writer http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	filterMap := index.FilterMap(params)
	testData := &form.TestValidator{}
	library.FmtPrint("-----a-----", "incontroller")
	err := request.ParamsValidate(filterMap, testData)
	if err != nil {
		fmt.Fprintf(writer, index.ErrorData(err.Error()))
	} else {
		val, err := redis.Get("key").Result()
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(writer, index.SuccessData(val))
	}
}

// redis验证
func (index *IndexController) Redis(writer http.ResponseWriter, req *http.Request) {
	val, err := redis.Get("key").Result()
	if err != nil {
		fmt.Fprintf(writer, index.ErrorData(err.Error()))
		return
	}
	fmt.Fprintf(writer, index.SuccessData(val))
}

// validator 测试接口
func (index *IndexController) ValidatorTest(writer http.ResponseWriter, req *http.Request) {
	params, _ := url.ParseQuery(req.URL.RawQuery)
	filterMap := index.FilterMap(params)
	fmt.Fprintln(writer, index.ErrorData(filterMap))
}
