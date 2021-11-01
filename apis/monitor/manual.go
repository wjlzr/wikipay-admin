package monitor

import (
	"strings"
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/rpc/client"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//添加
func InsertMonitorManualSetting(c *gin.Context) {
	var data monitor.MonitorManualSetting

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)

	ma := monitor.MonitorUserAddress{
		Coin:    data.Coin,
		Address: data.FromAddress,
	}
	infos, err := ma.Get()
	if len(infos) < 1 || err != nil {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "地址不存在",
		})
		return
	}

	grpcClient := client.NewClient()
	isValid := grpcClient.ValidateAddress(&client.BalanceReq{
		Coin:    data.Coin,
		Address: data.ToAddress,
	})
	if !isValid {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "输入的归集地址不正确",
		})
		return
	}
	//
	if data.GasStatus == 2 && data.GasPrice <= 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "gas费用不能小于等于0",
		})
		return
	}

	balance, _ := grpcClient.GetBalance(&client.BalanceReq{
		Coin:    data.Coin,
		Address: data.FromAddress,
	})
	if balance <= 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "被归集地址数量小于等于0",
		})
		return
	}

	data.Coin = strings.ToUpper(data.Coin)

	var fee float64
	if data.GasPrice > 0 {
		fee = data.GasPrice
	} else {
		if data.Coin == "ETH" {
			fee = utils.GetGasPrice(data.Coin)
		} else {
			fee = 0
		}
	}

	//fmt.Println(balance, data.Amount)
	if xfloat64.FromFloatCmp(balance, xfloat64.Add(data.Amount, fee)) < 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "输入的归集数量大于本身数量",
		})
		return
	}

	if data.Coin == "USDT-OMNI" {
		if data.FeeAddress == "" {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "手续费地址不能为空",
			})
			return
		}
	}

	result, err := data.Create()
	tools.HasError(err, "", -1)
	// go func() {
	// 	monitorservice.MonitorSend(&monitorservice.MonitorReq{
	// 		Coin:        data.Coin,
	// 		FromAddress: data.FromAddress,
	// 		ToAddress:   data.ToAddress,
	// 		Amount:      data.Amount,
	// 		FeeAddress:  data.FeeAddress,
	// 	})
	// }()
	app.OK(c, result, "")
}
