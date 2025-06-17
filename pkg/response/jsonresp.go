package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type JsonResponse struct {
	Code       int         `json:"code"`        // 错误码(0:成功, 1:失败, >1:错误码)
	Msg        string      `json:"msg"`         // 提示信息
	Data       interface{} `json:"data"`        // 返回数据(业务接口定义具体数据结构)
	ServerTime uint64      `json:"server_time"` // 服务器时间
}

type JsonResponseFunc func(obj *JsonResponse)

func WithCode(code int) JsonResponseFunc {
	return func(obj *JsonResponse) {
		obj.Code = code
	}
}

func WithMsg(msg string) JsonResponseFunc {
	return func(obj *JsonResponse) {
		obj.Msg = msg
	}
}

func WithData(data interface{}) JsonResponseFunc {
	return func(obj *JsonResponse) {
		obj.Data = data
	}
}

func WithErr(err error) JsonResponseFunc {
	return func(obj *JsonResponse) {
		if err != nil {
			obj.Code = -1
			obj.Msg = err.Error()
		} else {
			return
		}
	}
}

func Json(c *gin.Context, opts ...JsonResponseFunc) {
	jsonResp := &JsonResponse{
		Code:       0,
		Msg:        "success",
		ServerTime: uint64(time.Now().Unix()),
	}
	// 对jsonResp进行自定义赋值
	for _, opt := range opts {
		opt(jsonResp)
	}

	c.JSON(http.StatusOK, jsonResp)
}
