package library

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

var Utils struct{}

// 读文件
func ReadFile(filePath string, dataStruct interface{}) (interface{}, error) {
	fmt.Println()
	file, err := os.Open(filePath)
	if err != nil {
		Logger.Fatalln("connot open config file", err)
		return dataStruct, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dataStruct)

	if err != nil {
		log.Fatalln("connot get info from file", err)
	}
	return dataStruct, err

}

// 判断数据的类型
func GetType(data interface{}) reflect.Kind {
	return reflect.ValueOf(data).Kind()
}

// 通过反射调用类的函数
func InvokeObjectMethod(object interface{}, methodName string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	return reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
}
