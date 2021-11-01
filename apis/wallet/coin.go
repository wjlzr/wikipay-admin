package wallet

import (
	"wikipay-admin/models/wallet"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//
func GetCoinList(c *gin.Context) {
	var (
		data      wallet.Coin
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
func GetCoin(c *gin.Context) {
	var data wallet.Coin

	data.Code, _ = tools.StringToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加
func InsertCoin(c *gin.Context) {
	var data wallet.Coin

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)

	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新
func UpdateCoin(c *gin.Context) {
	var data wallet.Coin

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.Code)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

//删除
func DeleteCoin(c *gin.Context) {
	var data wallet.Coin

	IDS := tools.IdsStrToIdsIntGroup("code", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)

	app.OK(c, nil, msg.DeletedSuccess)
}
