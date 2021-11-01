package monitor

import (
	"strings"
	"wikipay-admin/common"
	"wikipay-admin/models/monitor"
	"wikipay-admin/rpc/client"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//获取网络状态
func SyncAccountAndBalance(c *gin.Context) {
	var req client.BaseReq

	err := c.ShouldBind(&req)
	tools.HasError(err, "", 500)

	req.Coin = utils.GetUpperCoin(req.Coin)
	if req.Coin == "" {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "输入币种不正确",
		})
		return
	}
	resp, err := monitor.SyncAccountAndBalance(&req)

	tools.HasError(err, "", -1)
	app.OK(c, gin.H{
		"balances": resp,
	}, "")
}

//获取列表
func GetMonitorUserAddressList(c *gin.Context) {
	var data monitor.MonitorUserAddress

	var (
		err       error
		pageSize  = 10
		pageIndex = 1
	)

	data.Coin = common.ETH
	if coin := c.Request.FormValue("coin"); coin != "" {
		data.Coin = strings.ToUpper(coin)
	}
	if info := c.Request.FormValue("info"); info != "" {
		data.Info = info
	}

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}
	if index := c.Request.FormValue("pageNum"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}
	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex)
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}
