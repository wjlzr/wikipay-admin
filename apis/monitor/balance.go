package monitor

import (
	"strings"
	"wikipay-admin/common"
	"wikipay-admin/rpc/client"
	monitorservice "wikipay-admin/service/monitor"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//获取余额
func GetBalance(c *gin.Context) {
	var req client.BalanceReq

	err := c.ShouldBindQuery(&req)
	tools.HasError(err, "", 500)

	if strings.ToUpper(req.Coin) != common.BTC {
		if req.Address == "" {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "地址不能为空",
			})
		}
	}
	rpcClient := client.NewClient()
	if req.Address != "" {
		isValid := rpcClient.ValidateAddress(&client.BalanceReq{
			Coin:    req.Coin,
			Address: req.Address,
		})
		if !isValid {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "输入的地址不正确",
			})
			return
		}
	}

	var amount float64
	if req.Coin == "BTC" {
		if req.Address != "" {
			amount = utils.GetBtcBalance(req.Address)
			app.OK(c, gin.H{
				"amount": amount,
			}, "")
			return
		}
	}
	amount, err = rpcClient.GetBalance(&req)
	tools.HasError(err, "", -1)
	app.OK(c, gin.H{
		"amount": amount,
	}, "")
}

//获取网络状态
func GetNetworkInfo(c *gin.Context) {
	var req client.BaseReq

	err := c.ShouldBindQuery(&req)
	tools.HasError(err, "", 500)

	rpcClient := client.NewClient()
	resp, err := rpcClient.GetNetworkInfo(&req)

	tools.HasError(err, "", -1)
	app.OK(c, gin.H{
		"networkInfo": resp,
	}, "")
}

//统计
func GetSatistical(c *gin.Context) {
	req, err := monitorservice.CalcSatistical()

	tools.HasError(err, "", -1)
	app.OK(c, req, "")
}

//资产对比数据统计
func AssetComparison(c *gin.Context) {

	req, err := monitorservice.AssetComparison()
	tools.HasError(err, "", -1)
	app.OK(c, req, "")
}
