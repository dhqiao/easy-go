/*
this is the entrance of the filter params from query
*/
package request

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"easy-go/library"
)

type ValidateParam struct{}

var validateParam = &ValidateParam{}

var invoke = reflect.ValueOf(validateParam)

const (
	REQUIRED = "requried"
	VALID    = "valid"
)

var logger = library.Logger

func ParamsValidate(params map[string]string, validatePtr interface{}) error {

	// 判断validate必须是struct类型
	kind := reflect.ValueOf(validatePtr).Kind()
	if kind != reflect.Ptr {
		return errors.New("err type of validate struct")
	}

	// 获取validate参数数量
	fieldNum := reflect.TypeOf(validatePtr).Elem().NumField()
	if fieldNum <= 0 {
		return errors.New("validate struct is not nil")
	}

	// 循环匹配参数
	for fieldIndex := 0; fieldIndex < fieldNum; fieldIndex++ {

		fieldName := reflect.TypeOf(validatePtr).Elem().Field(fieldIndex).Name

		if len([]rune(fieldName)) == 0 {
			return errors.New("validate struct filed name is not nil")
		}
		fieldValue, ok := params[fieldName]
		if !ok {
			fieldValue = ""
		}
		err := Compare(validatePtr, fieldIndex, fieldName, fieldValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// 对比参数
func Compare(validatePtr interface{}, fieldIndex int, fieldName string, fieldValue string) error {
	fieldTag := reflect.TypeOf(validatePtr).Elem().Field(fieldIndex).Tag
	required, _ := strconv.ParseBool(fieldTag.Get(REQUIRED))

	// 判断不能为空
	if required && fieldValue == "" {
		return errors.New(fieldName + "不能为空")
	}
	fieldValidFuc := reflect.TypeOf(validatePtr).Elem().Field(fieldIndex).Tag.Get("valid")
	fieldValidFuc = "SetAndCheck" + fieldValidFuc
	outValue := library.InvokeObjectMethod(validateParam, fieldValidFuc, validatePtr, fieldIndex, fieldName, fieldValue, fieldTag)
	outValLen := len(outValue)
	if outValLen >= 1 && !outValue[outValLen-1].IsNil() {
		return outValue[outValLen-1].Elem().Interface().(error)
	}
	return nil

}

// 验证是否是正整数
func (validate *ValidateParam) SetAndCheckPosNO(validatePtr interface{}, fieldIndex int, fieldName string, fieldValue string, tag reflect.StructTag) error {

	name := GetErrorParamName(fieldName, tag)
	intValue, err := strconv.Atoi(fieldValue)
	if err != nil {
		return errors.New(name + "参数类型错误")
	}
	stringValue := strconv.Itoa(intValue)
	if fieldValue != stringValue {
		return errors.New(name + "参数类型错误")
	}
	reflect.ValueOf(validatePtr).Elem().Field(fieldIndex).SetInt(int64(intValue))
	return nil
}

// 字符串校验
func (validate *ValidateParam) SetAndCheckString(validatePtr interface{}, fieldIndex int, fieldName string, fieldValue string, tag reflect.StructTag) error {
	name := GetErrorParamName(fieldName, tag)
	fieldValue = strings.Replace(fieldValue, "'", "", -1)
	fieldValue = strings.Replace(fieldValue, " ", "", -1)
	fieldValue = strings.Replace(fieldValue, " ", "", -1)
	fieldValue = strings.Replace(fieldValue, "\\", "", -1)
	fieldValue = strings.Replace(fieldValue, "\"", "", -1)

	var minLen, maxLen int = 0, 0
	var err error = nil
	lenStr := strings.Split(tag.Get("len"), ",")
	if lenStr[0] != "" {
		minLen, err = strconv.Atoi(lenStr[0])
		if err != nil {
			return errors.New(name + ":验证字符串的最小长度参数输入有误!")
		}
	}
	if len(lenStr) == 2 && lenStr[1] != "" {
		maxLen, err = strconv.Atoi(lenStr[1])
		if err != nil {
			return errors.New(name + ":验证字符串的最大长度参数输入有误!")
		}
	}
	if minLen > maxLen {
		return errors.New(name + ":最小长度和最大长度冲突!")
	}
	fieldLen := len(fieldValue)
	if fieldLen < minLen {
		return errors.New(name + ":字符串长度过短!")
	}
	if maxLen != 0 && fieldLen > maxLen {
		return errors.New(name + ":字符串过长!")
	}
	reflect.ValueOf(validatePtr).Elem().Field(fieldIndex).SetString(fieldValue)
	return nil
}

// 如果验证失败需要返回错误，返回错误当前参数的名称
func GetErrorParamName(fieldName string, tag reflect.StructTag) string {
	name := tag.Get("name")
	if name == "" {
		name = fieldName
	}
	return name
}
