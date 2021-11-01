package monitor

import (
	monitorservice "wikipay-admin/service/monitor"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"

	"github.com/gin-gonic/gin"
)

//
func GetWithdrawWithNow(c *gin.Context) {
	result, err := monitorservice.WithDrawWithNow()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

//获取资产
func GetAssetsGroup(c *gin.Context) {
	result, err := monitorservice.AssetsGroup()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

//归集资 产
func GetCollectAssetsGroup(c *gin.Context) {
	result, err := monitorservice.CollectAssetsGroup()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}
