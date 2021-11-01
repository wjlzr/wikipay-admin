package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"wikipay-admin/models/user"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"
)

//列表
func GetFeedbackTypeList(c *gin.Context) {
	var (
		data      user.FeedbackType
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
	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

//获取
func GetFeedbackType(c *gin.Context) {
	var data user.FeedbackType

	data.Code, _ = tools.StringToInt(c.Param("code"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加
func InsertFeedbackType(c *gin.Context) {
	var data user.FeedbackType

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

//更新
func UpdateFeedbackType(c *gin.Context) {
	var data user.FeedbackType

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)

	result, err := data.Update(data.Code)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除
func DeleteFeedbackType(c *gin.Context) {
	var data user.FeedbackType

	IDS := tools.IdsStrToIdsIntGroup("code", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)

	app.OK(c, nil, msg.DeletedSuccess)
}
