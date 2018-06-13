package models

import (
	"easy-go/data"
	"easy-go/library"
)

type Member struct {
	UserId   string
	ObjectId int
	Type     int
	Status   int
}

// 数据库实例
var DB = data.DBConnections["user_db"]

// 日志实例
var logger = library.Logger

// 获取列表
func GetListByUser() (memers []Member, err error) {
	rows, err := DB.Query("SELECT iUserID, iObjectID, iType, iStatus from t_user_follow")
	if err != nil {
		logger.Println("get List error")
		return
	}
	for rows.Next() {
		member := Member{}
		if err = rows.Scan(&member.UserId, &member.ObjectId, &member.Type, &member.Status); err != nil {
			return
		}
		memers = append(memers, member)
	}
	rows.Close()
	return
}
