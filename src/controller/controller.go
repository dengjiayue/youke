package controller

import (
	"fmt"
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

func (c *Controller) Ping(ctx *gin.Context) {
	fmt.Printf("ping seccess")
	ctx.String(200, "pong")
}

// 查看某天的订单
func (c *Controller) SelectOrderByYmd(ctx *gin.Context) {
	req := &struct {
		Ymd string
	}{}
	err := ctx.Bind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	resp, err := model_order.SelectByYmd(global.Global.Db, req.Ymd)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	public_func.Success(ctx, resp)

}

// 顾客简略搜索
func (c *Controller) SelectCostomerSimple(ctx *gin.Context) {
	req := &model_customer.CostomerSimple{}
	err := ctx.Bind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	data, count, err := model_customer.SelectCostomerSimple(global.Global.Db, req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	resp := gin.H{"Data": data, "Count": count}
	public_func.Success(ctx, resp)

}

// 查看顾客详情
func (c *Controller) SelectCostomerById(ctx *gin.Context) {
	req := &struct{ Id int64 }{}
	err := ctx.Bind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	resp, err := model_customer.SelectById(global.Global.Db, req.Id)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	public_func.Success(ctx, resp)

}

// 上传头像图片
func (c *Controller) UploadFaceImg(ctx *gin.Context) {
	req := &struct {
		IdcardNumber string                `form:"idcard_number" binding:"required"`
		FaceImg      *multipart.FileHeader `form:"face_img" binding:"required"`
	}{}

	err := ctx.ShouldBind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	if req.FaceImg == nil || req.IdcardNumber == "" {
		global.Global.Logger.Error("参数有误", req)
		public_func.Fail(ctx, public_func.CommonERR, "参数有误")
		return
	}

	//对象存储
	var faceUrl string
	{
		r, err := req.FaceImg.Open()
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}

		faceUrl, err = cos.UploadFile(global.Global.Cos, r, req.FaceImg.Size, req.IdcardNumber+"-f.jpg")
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
	}
	// fmt.Printf("faceurl=%#v\n", faceUrl)
	public_func.Success(ctx, gin.H{"FaceImg": faceUrl})

}

// 身份证识别
func (c *Controller) IdCardRecognition(ctx *gin.Context) {

	req := &struct {
		IdCardBase64 string //身份证图片base64编码
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

	//对象存储
	var faceUrl, idcardUrl string
	{
		r, err := ocr_server.ConvertBase64ToReader(req.IdCardBase64)
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}

		idcardUrl, err = cos.UploadFile(global.Global.Cos, r, 0, idCardInfo.IdNum+"-c.jpg")
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
	}
	{
		faceUrl, err = cos.UploadFile(global.Global.Cos, idCardInfo.HeadImg, 0, idCardInfo.IdNum+"-f.jpg")
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
	}
	fmt.Printf("faceurl=%#v\n", faceUrl)

	resp := &struct {
		IdcardImg string
		HeadImg   string
		Name      string `json:"Name"`
		// Age     int
		Address string `json:"Address"`
		IdNum   string `json:"IdNum"`
		IsAdult bool
	}{
		idcardUrl,
		faceUrl,
		idCardInfo.Name,
		// idCardInfo.Age,
		idCardInfo.Address,
		idCardInfo.IdNum,
		isAdult,
	}

	public_func.Success(ctx, resp)
	// fmt.Printf("=%#v\n", idCardInfo)
}

// 新增顾客登记
func (c *Controller) CreatOrderAndUpdateCostomer(ctx *gin.Context) {
	fmt.Printf("登记\n")
	req := &struct {
		ChildIdNumber string `form:"child_id_number"`                //被监护人身份证号(只需要绑定监护人时提供)
		RoomNumber    string `form:"room_number" binding:"required"` //房间号为 "0",表示只登记, 不产生订单(只登记接口/监护人登记接口使用)
		Name          string `form:"name" binding:"required"`
		PhoneNumber   string `form:"phone_number" binding:"required"`
		FaceImg       string `form:"face_img" binding:"required"`
		IdcardImg     string `form:"idcard_img" binding:"required"`
		IdcardNumber  string `form:"idcard_number" binding:"required"`
		Address       string `form:"address" binding:"required"`
		GuardianId    int64  `form:"guardian_id"`
	}{}

	err := ctx.ShouldBind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}
	fmt.Printf("req=%#v\n", req)

	//检查输入
	if len(req.RoomNumber) == 0 || len(req.Name) == 0 || len(req.PhoneNumber) == 0 || len(req.IdcardNumber) == 0 || len(req.Address) == 0 {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, "个人信息不完整,请填写完整提交")
		return
	}

	// //对象存储
	// var faceUrl, idcardUrl string
	// {
	// 	r, err := req.IdcardImg.Open()
	// 	if err != nil {
	// 		global.Global.Logger.Error(err)
	// 		public_func.Fail(ctx, public_func.CommonERR, err)
	// 		return
	// 	}

	// 	idcardUrl, err = cos.UploadFile(global.Global.Cos, r, req.IdcardImg.Size, req.IdcardNumber+"-c.jpg")
	// 	if err != nil {
	// 		global.Global.Logger.Error(err)
	// 		public_func.Fail(ctx, public_func.CommonERR, err)
	// 		return
	// 	}
	// }
	// {
	// 	r, err := req.FaceImg.Open()
	// 	if err != nil {
	// 		global.Global.Logger.Error(err)
	// 		public_func.Fail(ctx, public_func.CommonERR, err)
	// 		return
	// 	}

	// 	faceUrl, err = cos.UploadFile(global.Global.Cos, r, req.FaceImg.Size, req.IdcardNumber+"-f.jpg")
	// 	if err != nil {
	// 		global.Global.Logger.Error(err)
	// 		public_func.Fail(ctx, public_func.CommonERR, err)
	// 		return
	// 	}
	// }
	// fmt.Printf("faceurl=%#v\n", faceUrl)

	//记录顾客数据
	customer := &model_customer.Model{
		Name:         &req.Name,
		PhoneNumber:  &req.PhoneNumber,
		FaceImg:      &req.FaceImg,
		IdcardImg:    &req.IdcardImg,
		IdcardNumber: &req.IdcardNumber,
		Address:      &req.Address,
		GuardianId:   &req.GuardianId,
	}

	err = customer.CreateOrUpdateByIdcardNumber(global.Global.Db)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}
	fmt.Printf("登记成功\n")

	//订单记录
	if req.RoomNumber != "0" && req.RoomNumber != "" {
		ymd := time.Now().Truncate(24 * time.Hour)
		order := &model_order.Model{RoomNumber: &req.RoomNumber, CustomerId: customer.Id, CustomerName: &req.Name, PhoneNumber: &req.PhoneNumber, Ymd: &ymd}
		err := order.Create(global.Global.Db)
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
		fmt.Printf("订单已生成\n")
	}

	//绑定监护人
	if req.ChildIdNumber != "" {
		err = global.Global.Db.Model(&model_customer.Model{}).Where("idcard_number = ?", req.ChildIdNumber).Update("guardian_id", customer.Id).Error
		if err != nil {
			global.Global.Logger.Error(err)
			public_func.Fail(ctx, public_func.CommonERR, err)
			return
		}
		fmt.Printf("绑定监护人成功\n")
	}

	public_func.Success(ctx, "ok")

}

// 一键登记顾客
func (c *Controller) CreatOrder(ctx *gin.Context) {
	req := &struct {
		CustomerId int64
		RoomNumber string
	}{}

	err := ctx.Bind(req)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	// 查询顾客信息
	var data model_customer.Model
	err = global.Global.Db.Select("name,phone_number").Where("id=?", req.CustomerId).First(&data).Error
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	// 记录订单
	ymd := time.Now().Truncate(24 * time.Hour)
	order := &model_order.Model{
		RoomNumber:   &req.RoomNumber,
		CustomerId:   &req.CustomerId,
		CustomerName: data.Name,
		PhoneNumber:  data.PhoneNumber,
		Ymd:          &ymd,
	}
	err = order.Create(global.Global.Db)
	if err != nil {
		global.Global.Logger.Error(err)
		public_func.Fail(ctx, public_func.CommonERR, err)
		return
	}

	public_func.Success(ctx, "ok")

}

// // 更新顾客与订单
// func (c *Controller) UpdateOrderAndUpdateCostomer(ctx *gin.Context) {

// }

// // 更新顾客信息
// func (c *Controller) UpdateCostomer(ctx *gin.Context) {

// }

// // 一键修改订单(换房/价格)
// func (c *Controller) UpdateOrder(ctx *gin.Context) {

// }
