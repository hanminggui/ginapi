package ginapi

// Uint 接受gin.Context.Get(key) 的返回值作为参数。返回uint类型结果
func Uint(v interface{}, exists bool) (u uint) {
	if exists && v != nil {
		u, _ = v.(uint)
	}
	return
}
