package user

import (
	"wikipay-admin/models/user"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//
func GetAdvertList(c *gin.Context) {
	var (
		data      user.Advert
		info      string
		err       error
		pageSize  = 10
		pageIndex = 1
	)

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}
	if s := c.Request.FormValue("info"); s != "" {
		info = s
	}

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex, info)
	tools.HasError(err, "", -1)
	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

//
func GetAdvert(c *gin.Context) {
	var data user.Advert
	data.Id, _ = tools.StringToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//
func InsertAdvert(c *gin.Context) {
	var data user.Advert
	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//
func UpdateAdvert(c *gin.Context) {
	var data user.Advert
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//
func DeleteAdvert(c *gin.Context) {
	var data user.Advert
	IDS := tools.IdsStrToIdsIntGroup("id", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
