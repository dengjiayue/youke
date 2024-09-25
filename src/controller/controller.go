package controller

import (
	"fmt"
	"youke/global"
	"youke/public_func"
	"youke/src/ocr_server"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

// 查看某天的订单
func (c *Controller) SelectOrderByYmd(ctx *gin.Context) {

}

// 顾客简略搜索
func (c *Controller) SelectCostomerSimple(ctx *gin.Context) {

}

// 查看顾客详情
func (c *Controller) SelectCostomerById(ctx *gin.Context) {

}

// 身份证识别
func (c *Controller) IdCardRecognitionAndCreateCostomer(ctx *gin.Context) {
	// idCardImg, err := ctx.FormFile("IdCardImg")
	// if err != nil {
	// 	global.Global.Logger.Error(err)
	// 	public_func.Fail(ctx, public_func.CommonERR, err)
	// }
	// idcardFile, err := idCardImg.Open()
	// defer idcardFile.Close()
	req := &struct {
		IdCardBase64 string //身份证图片base64编码
		PhoneNumber  string //电话号码(用户唯一标记)
	}{}

	err := ctx.Bind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
	}

	idCardInfo, err := ocr_server.IdCardOCR(req.IdCardBase64)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
	}
	fmt.Printf("=%#v\n", idCardInfo)
}

// 新增顾客登记
func (c *Controller) CreatOrderAndUpdateCostomer(ctx *gin.Context) {

}

// 一键登记顾客
func (c *Controller) CreatOrder(ctx *gin.Context) {

}

// 更新顾客与订单
func (c *Controller) UpdateOrderAndUpdateCostomer(ctx *gin.Context) {

}

// 更新顾客信息
func (c *Controller) UpdateCostomer(ctx *gin.Context) {

}

// 一键修改订单(换房/价格)
func (c *Controller) UpdateOrder(ctx *gin.Context) {

}
