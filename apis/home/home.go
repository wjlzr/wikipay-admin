package home

import (
	"strings"
	"wikipay-admin/models/home"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//获取首页用户统计信息
func GetTotalInfo(c *gin.Context) {
	info, err := home.GetTotalInfo()
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, info, "")
}

//获取收益
func GetProfitInfo(c *gin.Context) {
	var data home.MonthReq
	err := c.ShouldBindWith(&data, binding.Query)
	tools.HasError(err, "数据解析失败", -1)

	infos := home.GetProfit(data.Month)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, infos, "")
}

//获取交易
func GetTradeInfo(c *gin.Context) {
	var data home.TradeReq
	err := c.ShouldBindWith(&data, binding.Query)
	tools.HasError(err, "数据解析失败", -1)

	depositInfo, withdrawInfo, err := home.GroupTradeInfo(strings.ToUpper(data.Coin), data.Month)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, gin.H{
		"deposits":  depositInfo,
		"withdraws": withdrawInfo,
	}, "")
}

//获取财务
func GetFinancialInfo(c *gin.Context) {
	info, err := home.GetFinancialInfo()
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, info, "")
}
