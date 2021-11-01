package exchange

import (
	"wikipay-admin/models/exchange"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//列表
func GetCoinBidAskList(c *gin.Context) {
	var (
		data         exchange.CoinBidAsk
		err          error
		pageSize     = 10
		pageIndex    = 1
		info         string
		searchParams exchange.SearchParams
	)

	//绑定不校验参数
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

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex, info, searchParams)
	tools.HasError(err, "", -1)
	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

//获取
func GetCoinBidAsk(c *gin.Context) {
	var data exchange.CoinBidAsk

	data.OrderNumber = c.Param("orderNumber")
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加
func InsertCoinBidAsk(c *gin.Context) {
	var data exchange.CoinBidAsk

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新
func UpdateCoinBidAsk(c *gin.Context) {
	var data exchange.CoinBidAsk

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.OrderNumber)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除
func DeleteCoinBidAsk(c *gin.Context) {
	var data exchange.CoinBidAsk

	IDS := tools.IdsStrToIdsIntGroup("orderNumber", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
