package monitor

import (
	"strings"
	"wikipay-admin/models/monitor"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//归集记录列表
func GetMonitorHistoryList(c *gin.Context) {
	var req monitor.MonitorHistoryReq

	if err := c.ShouldBindQuery(&req); err != nil {
		app.Custum(c, gin.H{
			"code": 111001,
			"msg":  "币种或地址不能为空",
		})
		return
	}
	utils.Pagination(&req.PageSize, &req.PageIndex)
	req.Coin = strings.ToUpper(req.Coin)

	result, count, err := new(monitor.MonitorHistory).GetPage(&req)
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.PageOK(c, result, count, int(req.PageIndex), int(req.PageSize), "")
}

//获取归集记录明细
func GetMonitorHistoryDetail(c *gin.Context) {
	var req monitor.MonitorHistoryReq

	if err := c.ShouldBindQuery(&req); err != nil {
		app.Custum(c, gin.H{
			"code": 111001,
			"msg":  "币种或地址不能为空",
		})
		return
	}
	utils.Pagination(&req.PageSize, &req.PageIndex)
	req.Coin = strings.ToUpper(req.Coin)

	result, count, err := new(monitor.MonitorHistory).GetDetails(&req)
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.PageOK(c, result, count, int(req.PageIndex), int(req.PageSize), "")
}
