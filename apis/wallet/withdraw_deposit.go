package wallet

import (
	"wikipay-admin/models/wallet"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

///
func GetWithdrawDepositList(c *gin.Context) {
	var (
		data         wallet.WithdrawDeposit
		err          error
		info         string
		status       string
		tradeType    string
		pageSize     = 10
		pageIndex    = 1
		searchParams wallet.SearchParams
	)

	//绑定无需校验的参数
	_ = c.ShouldBindQuery(&searchParams)

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}
	if s := c.Request.FormValue("info"); s != "" {
		info = s
	}
	if st := c.Request.FormValue("status"); st != "" {
		status = st
	}
	if t := c.Request.FormValue("type"); t != "" {
		tradeType = t
	}
	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex, info, status, tradeType, searchParams)
	tools.HasError(err, "", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

///
func GetWithdrawDeposit(c *gin.Context) {
	var data wallet.WithdrawDeposit
	data.OrderNumber = c.Param("orderNumber")
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

//
func InsertWithdrawDeposit(c *gin.Context) {
	var data wallet.WithdrawDeposit
	err := c.ShouldBindJSON(&data)

	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新
func UpdateWithdrawDeposit(c *gin.Context) {
	var data wallet.WithdrawDeposit
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)

	orderNumber, _ := tools.StringToInt64(data.OrderNumber)
	result, err := data.Update(orderNumber)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

//更新
func UpdateWithdrawDepositStatus(c *gin.Context) {
	var data wallet.WithdrawDeposit
	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)

	orderNumber, _ := tools.StringToInt64(data.OrderNumber)
	result, err := data.UpdateStatus(orderNumber)
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

//删除
func DeleteWithdrawDeposit(c *gin.Context) {
	var data wallet.WithdrawDeposit

	IDS := tools.IdsStrToIdsInt64Group("orderNumber", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
