package monitor

import (
	"wikipay-admin/models/monitor"
	"wikipay-admin/tools/app"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

func GetGasPrice(c *gin.Context) {
	var req monitor.MonitorReq

	if err := c.ShouldBindQuery(&req); err != nil {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "输入币种不能为空",
		})
		return
	}

	gasPrice := utils.GetGasPrice(req.Coin)
	app.OK(c, gin.H{
		"coin":     req.Coin,
		"gasPrice": gasPrice,
	}, "")
}
