package tools

import (
	"fmt"
	"log"
	"time"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	var fileName, path string
	var data []string
	// 上传文件至指定目录
	fileName = tools.Int64ToString(time.Now().UnixNano()/1e6) + "-" + file.Filename
	path = "static/logo/"
	if err := c.SaveUploadedFile(file, path+fileName); err != nil {
		fmt.Println("UploadFile error", err)
		log.Println("UploadFile error", err)
		tools.HasError(err, "上传失败", -1)
	}
	data = append(data, utils.GetImage(utils.Logo, fileName))
	app.Success(c, data, "ok")
}
