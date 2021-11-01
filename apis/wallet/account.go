package wallet

import (
	"wikipay-admin/models/wallet"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"wikipay-admin/tools/app/msg"
)

//账户列表
func GetAccountList(c *gin.Context) {
	var (
		data      wallet.Account
		err       error
		info      string
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

//获取账户
func GetAccount(c *gin.Context) {
	var data wallet.Account

	data.Id, _ = tools.StringToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加账户
func InsertAccount(c *gin.Context) {
	var data wallet.Account

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//修改账户
func UpdateAccount(c *gin.Context) {
	var data wallet.Account

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除账户
func DeleteAccount(c *gin.Context) {
	var data wallet.Account

	IDS := tools.IdsStrToIdsIntGroup("id", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
