/*
this file validate one struct, not friendly for query params
*/
package request

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"easy-go/library"
)

type Validator struct{}

var Validate = newValidator()

// 生成实例
func newValidator() *Validator {
	return &Validator{}
}

// 验证方法入口
func (validator *Validator) ValidateData(data interface{}) error {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	dataValueKind := dataValue.Kind()

	// 普通类型不解析
	if dataValueKind != reflect.Slice && dataValueKind != reflect.Ptr && dataValueKind != reflect.Struct {
		return nil
	}
	fmt.Println(dataType)

	err := validator.parseParam(dataType, dataValue, reflect.StructTag(""))
	return err
}

// 解析参数
func (validator *Validator) parseParam(t reflect.Type, v reflect.Value, tag reflect.StructTag) error {
	vKind := v.Kind()
	// 数组
	if vKind == reflect.Slice {
		for i, n := 0, v.Len(); i < n; i++ {
			t1 := reflect.TypeOf(v.Index(i).Interface())
			err := validator.parseParam(t1, v.Index(i), tag)
			if err != nil {
				return err
			}
		}
		// 指针
	} else if vKind == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
		err := validator.parseParam(t, v, tag)
		if err != nil {
			return err
		}
		// 结构体
	} else if vKind == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			t.FieldByName(f.Name)
			err := validator.parseParam(f.Type, v.FieldByName(f.Name), f.Tag)
			if err != nil {
				return err

			}
		}
		// 普通类型
	} else {
		invoker := reflect.ValueOf(validator)
		// 获取要进行验证的函数名
		methodsName := strings.Split(tag.Get("valid"), ",")
		for i, n := 0, len(methodsName); i < n; i++ {
			methodName := methodsName[i]
			if len(methodName) != 0 {
				// 调用需要验证的函数
				method := invoker.MethodByName("Check" + methodName)
				inVal := []reflect.Value{v, reflect.ValueOf(tag)}
				outVal := method.Call(inVal)
				outValLen := len(outVal)
				if outValLen != 1 {
					v.Set(outVal[0])
				}
				if !outVal[outValLen-1].IsNil() {
					// 返回错误信息
					return outVal[outValLen-1].Elem().Interface().(error)
				}
			}
		}
	}
	return nil
}

func (validator *Validator) CheckPosNO(data int, tag reflect.StructTag) error {
	if data <= 0 {
		return errors.New(tag.Get("name") + ":不能为负数!")
	}
	return nil
}

// 检查字符串
func (validator *Validator) CheckStr(data string, tag reflect.StructTag) (string, error) {

	data = strings.Replace(data, "'", "", -1)
	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "\\", "", -1)
	data = strings.Replace(data, "\"", "", -1)

	var minLen, maxLen int = 0, 0

	var err error = nil

	lenStr := strings.Split(tag.Get("len"), ",")
	if lenStr[0] != "" {
		minLen, err = strconv.Atoi(lenStr[0])
		if err != nil {
			return "", errors.New(tag.Get("name") + ":验证字符串的最小长度参数输入有误!")
		}
	}
	if len(lenStr) == 2 && lenStr[1] != "" {
		maxLen, err = strconv.Atoi(lenStr[1])
		if err != nil {
			return "", errors.New(tag.Get("name") + ":验证字符串的最大长度参数输入有误!")
		}
	}

	if minLen > maxLen {
		return "", errors.New(tag.Get("name") + ":最小长度和最大长度冲突!")
	}
	dataLen := len(data)
	library.FmtPrint("''''''''''data''''''", data)
	if dataLen < minLen {
		return data, errors.New(tag.Get("name") + ":字符串长度过短!")
	}

	if maxLen != 0 && dataLen > maxLen {
		return data, errors.New(tag.Get("name") + ":字符串过长!")
	}

	return data, nil
}

func (validator *Validator) CheckCardState(data interface{}, tag reflect.StructTag) error {
	if reflect.TypeOf(data).Name() != "int" {
		return errors.New(tag.Get("name") + ":类型只能是int")
	}
	cardState := data.(int)
	if cardState < 0 || cardState > 3 {
		return errors.New(tag.Get("name") + ":卡类类型出错")
	}
	return nil
}
