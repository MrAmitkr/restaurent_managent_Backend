package routes

import (
	controller "golang-res/controllers"

	"github.com/gin-gonic/gin"
	//	"golang-res/controllers"
)

func TableRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.GET("/tables", controller.GetTables())
	incomingRoutes.GET("/tables/:table_id", controller.GetTable())
	incomingRoutes.POST("/tables", controller.CreateTable())
	incomingRoutes.PATCH("/tables/:table_id", controller.UpdateTable())

}
