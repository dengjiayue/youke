package router

import (
	"youke/src/controller"
	"youke/src/middleware"

	"github.com/gin-gonic/gin"
)

func initRouter(engin *gin.Engine) {
	controller := controller.New()
	engin.Use(middleware.CORSMiddleware())
	engin.GET("ping", controller.Ping)
	engin.POST("CreatOrder", controller.CreatOrder)
	engin.POST("CreatOrderAndUpdateCostomer", controller.CreatOrderAndUpdateCostomer)
	engin.POST("SelectCostomerById", controller.SelectCostomerById)
	engin.POST("SelectOrderByYmd", controller.SelectOrderByYmd)
	engin.POST("SelectCostomerSimple", controller.SelectCostomerSimple)
	engin.POST("IdCardRecognition", controller.IdCardRecognition)
}

func Run() {
	engin := gin.Default()
	initRouter(engin)
	engin.Run(":11000")
}
