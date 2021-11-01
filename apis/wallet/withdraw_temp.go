package wallet

import (
	"github.com/gin-gonic/gin"

	"wikipay-admin/models/wallet"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"
)

//列表
func GetWithdrawTempList(c *gin.Context) {
	var (
		data      wallet.WithdrawTemp
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

//获取
func GetWithdrawTemp(c *gin.Context) {
	var data wallet.WithdrawTemp

	data.TransactionNumber = c.Param("transactionNumber")
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加
func InsertWithdrawTemp(c *gin.Context) {
	var data wallet.WithdrawTemp

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新
func UpdateWithdrawTemp(c *gin.Context) {
	// var data wallet.WithdrawTemp

	// err := c.BindWith(&data, binding.JSON)
	// tools.HasError(err, "数据解析失败", -1)
	// result, err := data.Update(data.TransactionNumber)
	// tools.HasError(err, "", -1)
	// app.OK(c, result, "")
}

//删除
func DeleteWithdrawTemp(c *gin.Context) {
	var data wallet.WithdrawTemp

	IDS := tools.IdsStrToIdsIntGroup("transactionNumber", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
