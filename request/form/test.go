package form

type TestValidator struct {
	IsCount    int64  `valid:"PosNO" requried:"true" name:"isCount" len:""` // 反射为什么只能是int64
	ItemString string `valid:"String" len:"1,3" name:"felix"`               // 方法必须大写
}
