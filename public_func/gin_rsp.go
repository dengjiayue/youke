package public_func

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code int         //0: 无错误(默认); 其他: 有错误
	Msg  string      //返回描述
	Data interface{} //返回数据
}

const (
	CommonERR    = 1 //一般错误
	NoPermission = 2 //没有权限
)

func ResponseData(ctx *gin.Context, code int, msg string, data interface{}) {
	ctx.JSON(200, &Resp{Code: code, Msg: msg, Data: data})
}

func Success(ctx *gin.Context, data interface{}) {
	ResponseData(ctx, 0, "ok", data)
}

func Fail(ctx *gin.Context, code int, errMsg interface{}) {
	ResponseData(ctx, code, fmt.Sprintf("%v", errMsg), nil)
}
