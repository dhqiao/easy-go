package controllers

import (
	"encoding/json"
)

// 基类
type BaseController struct{}

// 返回值结构体
type JsonResponse struct {
	Status       int         `json:"status"`
	ErrorCode    int         `json:"errorCode"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

// 成功返回的数据
func (base *BaseController) SuccessData(data interface{}) string {
	b, _ := json.Marshal(JsonResponse{0, 0, "ok", data})
	return string(b)
}

// 失败返回的数据
func (base *BaseController) ErrorData(data interface{}) string {
	if data == nil {
		return base.ErrorNullData()
	}
	b, _ := json.Marshal(JsonResponse{0, -1, "error", data})
	return string(b)
}

// 空数据错误返回
func (base *BaseController) ErrorNullData() string {
	b, _ := json.Marshal(JsonResponse{0, -1, "error", ""})
	return string(b)
}

// 接口不支持穿数组参数
func (base *BaseController) FilterMap(from map[string][]string) map[string]string {
	to := make(map[string]string)
	if len(from) != 0 {
		for key, value := range from {
			to[key] = value[0]
		}
	}
	return to
}
