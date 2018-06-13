package controllers

import (
	"easy-go/library"
	"easy-go/models"
	"fmt"
	"net/http"
)

type MemberController struct {
	BaseController
}

// 日志实例
var logger = library.Logger

func (member *MemberController) getMemberList() {
	logger.Println("in controller member")
}

// 获取列表
func (member *MemberController) GetList(writer http.ResponseWriter, request *http.Request) {
	//params := request.URL.Query()

	members, err := models.GetListByUser()

	if err != nil {
		fmt.Fprintf(writer, member.ErrorData(err.Error()))
		return
	}

	fmt.Fprintf(writer, member.SuccessData(members))
}
