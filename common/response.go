package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse API 返回结果结构
type APIResponse struct {
	Code      int         `json:"-"`
	ErrorCode int         `json:"error_code"`
	Message   string      `json:"message"`
	Hint      string      `json:"hint"`
	Data      interface{} `json:"data"`
}

const (
	success        = 0    // 成功
	serverError    = 1000 // 系统错误
	notFound       = 1001 // 401错误
	unknownError   = 1002 // 未知错误​
	parameterError = 1003 // 参数错误
	authError      = 1004 // 错误​
)

func (e *APIResponse) Error() string {
	return e.Message
}

func newAPIResponse(code int, errorCode int, msg string) *APIResponse {
	return &APIResponse{
		Code:      code,
		ErrorCode: errorCode,
		Message:   msg,
	}
}

// ServerError 500 错误处理
func ServerError(c *gin.Context) {
	returnResponse(c, newAPIResponse(http.StatusInternalServerError, serverError, http.StatusText(http.StatusInternalServerError)))
}

// NotFound 404 错误
func NotFound(c *gin.Context) {
	returnResponse(c, newAPIResponse(http.StatusNotFound, notFound, http.StatusText(http.StatusNotFound)))
}

// UnknownError 未知错误
func UnknownError(c *gin.Context, message string) {
	returnResponse(c, newAPIResponse(http.StatusForbidden, unknownError, message))
}

// ParameterError 参数错误
func ParameterError(c *gin.Context, message string) {
	returnResponse(c, newAPIResponse(http.StatusBadRequest, parameterError, message))
}

// Success 正常请求
func Success(c *gin.Context, data interface{}) {
	response := newAPIResponse(http.StatusOK, success, "ok")
	response.Data = data
	returnResponse(c, response)
}

func returnResponse(c *gin.Context, response *APIResponse) {
	response.Hint = c.GetString("requestId")
	c.Set("response", response)
	c.JSON(response.Code, response)
}
