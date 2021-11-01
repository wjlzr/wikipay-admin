package wallet

import (
	"wikipay-admin/models/wallet"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//列表
func GetTransferTempList(c *gin.Context) {
	var (
		data      wallet.TransferTemp
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

//获取
func GetTransferTemp(c *gin.Context) {
	var data wallet.TransferTemp

	data.TransactionNumber = c.Param("transactionNumber")
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加
func InsertTransferTemp(c *gin.Context) {
	var data wallet.TransferTemp

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//修改
func UpdateTransferTemp(c *gin.Context) {
	var data wallet.TransferTemp

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	transactionNumber, _ := tools.StringToInt64(data.TransactionNumber)

	result, err := data.Update(transactionNumber)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除
func DeleteTransferTemp(c *gin.Context) {
	var data wallet.TransferTemp

	IDS := tools.IdsStrToIdsInt64Group("transactionNumber", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
