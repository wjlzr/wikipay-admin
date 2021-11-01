package wallet

import (
	"wikipay-admin/models/wallet"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//获取收益列表
func GetProfitList(c *gin.Context) {
	var (
		data      wallet.Profit
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

//获取收益
func GetProfit(c *gin.Context) {
	var data wallet.Profit

	data.OrderNumber = c.Param("orderNumber")
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//新增收益
func InsertProfit(c *gin.Context) {
	var data wallet.Profit

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新收益
func UpdateProfit(c *gin.Context) {
	var data wallet.Profit

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(tools.StrToInt(err, data.OrderNumber))
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除收益
func DeleteProfit(c *gin.Context) {
	var data wallet.Profit

	IDS := tools.IdsStrToIdsIntGroup("orderNumber", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
