package controller

import (
	"intelliq/app/model"
	"intelliq/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AddMetaData adds meta data
func AddMetaData(ctx *gin.Context) {
	var metaData model.Meta
	ctx.BindJSON(&metaData)
	res := service.AddNewData(&metaData)
	ctx.JSON(http.StatusOK, res)
}

//UpdateMetaData updates meta data
func UpdateMetaData(ctx *gin.Context) {
	var metaData model.Meta
	ctx.BindJSON(&metaData)
	res := service.UpdateMetaData(&metaData)
	ctx.JSON(http.StatusOK, res)
}

//ReadMetaData fetches meta data
func ReadMetaData(ctx *gin.Context) {
	res := service.ReadMetaData()
	ctx.JSON(http.StatusOK, res)
}
