package ginapi

var apiStucts []interface{} = make([]interface{}, 0)

// LoadAPI 预加载PAI载体
func LoadApi(list ...interface{}) {
	apiStucts = append(apiStucts, list...)
}
