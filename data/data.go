package data

import (
	"database/sql"
	"strconv"
	"easy-go/conf"
	"easy-go/library"

	_ "github.com/go-sql-driver/mysql"
)

// 数据库链接池
var DBConnections = make(map[string]*sql.DB)

var logger = library.Logger

// 初始化数据库链接
func init() {
	dbList := conf.DBConfigMap
	for dbName, dbConfig := range dbList {
		DBConnections[dbName] = connectDB(dbConfig)
	}
}

// 数据库链接
func connectDB(dbConfig conf.DBConfig) *sql.DB {
	connectStr := dbConfig.User + ":" + dbConfig.Password + "@" + "tcp(" + dbConfig.Host + ":" + strconv.Itoa(dbConfig.Port) + ")/" + dbConfig.DbName + "?" + "parseTime=true&charset=utf8,utf8mb4"
	//fmt.Println(connectStr)
	oneConn, err := sql.Open("mysql", connectStr)
	if err != nil {
		logger.Fatalln("mysql connect error")
	}
	return oneConn
}
