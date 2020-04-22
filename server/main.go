package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Name   string `json:"name" form:"name" xml:"name" binding:"required"`
	Number int    `json:"number" form:"number" xml:"number" binding:"required"`
}

func main() {
	r := gin.Default()

	// 模拟第三方提供添加产品API接口
	r.POST("/product", func(ctx *gin.Context) {
		var pro Product
		if err := ctx.ShouldBind(&pro); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, pro)
	})

	log.Fatal(r.Run(":7586"))
}
