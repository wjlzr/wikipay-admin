package handler

import (
	"github.com/gin-gonic/gin"
	"wikipay-admin/tools/app"
)

func HelloGoAdmin(c *gin.Context) {
	app.OK(c,"欢迎使用wikipay-admin 中后台脚手架","")
}