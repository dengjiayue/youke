package controller

import (
	"mime/multipart"
	"time"
	"youke/global"
	"youke/global/cos"
	model_customer "youke/model/customer"
	model_order "youke/model/order"
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

	req := &struct {
		IdCardBase64 string //身份证图片base64编码
		PhoneNumber  string //电话号码(用户唯一标记)
	}{}

	err := ctx.Bind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	idCardInfo, err := ocr_server.IdCardOCR(req.IdCardBase64)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	isAdult, err := public_func.IsAdultByID(idCardInfo.IdNum)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	resp := &struct {
		HeadImg []byte
		Name    string `json:"Name"`
		Age     int
		Address string `json:"Address"`
		IdNum   string `json:"IdNum"`
		IsAdult bool
	}{
		idCardInfo.HeadImg,
		idCardInfo.Name,
		idCardInfo.Age,
		idCardInfo.Address,
		idCardInfo.IdNum,
		isAdult,
	}

	public_func.Success(ctx, resp)
	// fmt.Printf("=%#v\n", idCardInfo)
}

// 新增顾客登记
func (c *Controller) CreatOrderAndUpdateCostomer(ctx *gin.Context) {
	req := struct {
		ChildId      int64                 `form:"child_id"`                       //被监护人id(只需要绑定监护人时提供)
		RoomNumber   string                `form:"room_number" binding:"required"` //房间号为 "0",表示只登记, 不产生订单(只登记接口/监护人登记接口使用)
		Name         string                `form:"name" binding:"required"`
		PhoneNumber  string                `form:"phone_number" binding:"required"`
		FaceImg      *multipart.FileHeader `form:"face_img" binding:"required"`
		IdcardImg    *multipart.FileHeader `form:"idcard_img" binding:"required"`
		IdcardNumber string                `form:"idcard_number" binding:"required"`
		Address      string                `form:"address" binding:"required"`
	}{}

	err := ctx.ShouldBind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	//检查输入
	if len(req.RoomNumber) == 0 || len(req.Name) == 0 || len(req.PhoneNumber) == 0 || req.FaceImg == nil || req.IdcardImg == nil || len(req.IdcardNumber) == 0 || len(req.Address) == 0 {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, "个人信息不完整,请填写完整提交")
		return
	}

	//对象存储
	var faceUrl, idcardUrl string
	{
		r, err := req.IdcardImg.Open()
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}

		faceUrl, err = cos.UploadFile(global.Global.Cos, r, req.IdcardImg.Size, req.IdcardNumber+"-c")
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
	}
	{
		r, err := req.FaceImg.Open()
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}

		idcardUrl, err = cos.UploadFile(global.Global.Cos, r, req.FaceImg.Size, req.IdcardNumber+"-f")
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
	}

	//记录顾客数据
	customer := &model_customer.Model{
		Name:         &req.Name,
		PhoneNumber:  &req.PhoneNumber,
		FaceImg:      &faceUrl,
		IdcardImg:    &idcardUrl,
		IdcardNumber: &req.IdcardNumber,
		Address:      &req.Address,
	}

	err = customer.CreateOrUpdateByIdcardNumber(global.Global.Db)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	//订单记录
	if req.RoomNumber != "0" {
		ymd := time.Now().Truncate(24 * time.Hour)
		order := &model_order.Model{RoomNumber: &req.RoomNumber, CustomerId: customer.Id, CustomerName: &req.Name, PhoneNumber: &req.PhoneNumber, Ymd: &ymd}
		err := order.Create(global.Global.Db)
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
	}

	//绑定监护人
	if req.ChildId != 0 {
		global.Global.Db.Model(&model_customer.Model{}).Where("id = ?", req.ChildId).Update("guardian_id", customer.Id)
	}

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
