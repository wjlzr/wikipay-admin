package monitor

import (
	"errors"
	"strings"
	"wikipay-admin/models/monitor"
	"wikipay-admin/rpc/client"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//获取
func GetMonitorSetting(c *gin.Context) {
	var (
		data monitor.MonitorSetting
		req  monitor.MonitorSettingReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		app.Error(c, 10001, errors.New("null"), "币种为空或类型未选择")
		return
	}

	data.Coin = strings.ToUpper(req.Coin)
	data.Type = req.Type

	result, err := data.Get()
	if err != nil {
		app.Custum(c, gin.H{
			"code": 55555,
			"msg":  "未找到相关记录",
		})
		return
	}
	app.OK(c, result, "")
}

//添加
func InsertMonitorSetting(c *gin.Context) {
	var data monitor.MonitorSetting

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)

	if data.Min <= 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "最小值不能小于等于0",
		})
		return
	}
	//检查地址
	rpcClient := client.NewClient()
	if data.Coin != "" {
		if data.Address1 != "" {
			isValid := rpcClient.ValidateAddress(&client.BalanceReq{
				Coin:    data.Coin,
				Address: data.Address1,
			})
			if !isValid {
				app.Custum(c, gin.H{
					"code": 10001,
					"msg":  "输入的地址不正确",
				})
				return
			}
		}

		if data.Address2 != "" {
			isValid := rpcClient.ValidateAddress(&client.BalanceReq{
				Coin:    data.Coin,
				Address: data.Address2,
			})
			if !isValid {
				app.Custum(c, gin.H{
					"code": 10001,
					"msg":  "输入的地址不正确",
				})
				return
			}
		}

		if data.Address3 != "" {
			isValid := rpcClient.ValidateAddress(&client.BalanceReq{
				Coin:    data.Coin,
				Address: data.Address3,
			})
			if !isValid {
				app.Custum(c, gin.H{
					"code": 10001,
					"msg":  "输入的地址不正确",
				})
				return
			}
		}
	}

	data.Coin = strings.ToUpper(data.Coin)
	if data.Coin == "USDT-OMNI" {
		if data.FeeAddress == "" {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "请输入手续费地址",
			})
			return
		}

		isValid := rpcClient.ValidateAddress(&client.BalanceReq{
			Coin:    data.Coin,
			Address: data.FeeAddress,
		})
		if !isValid {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "输入的手续费地址不正确",
			})
			return
		}
	}

	result, err := data.Create()
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

//更新
func UpdateMonitorSetting(c *gin.Context) {
	var data monitor.MonitorSetting

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)

	//检查最小值
	if data.Min <= 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "最小值不能小于等于0",
		})
		return
	}

	//检查地址
	rpcClient := client.NewClient()
	if data.Coin != "" {
		if data.Address1 != "" {
			isValid := rpcClient.ValidateAddress(&client.BalanceReq{
				Coin:    data.Coin,
				Address: data.Address1,
			})
			if !isValid {
				app.Custum(c, gin.H{
					"code": 10001,
					"msg":  "输入的地址不正确",
				})
				return
			}
		}

		if data.Address2 != "" {
			isValid := rpcClient.ValidateAddress(&client.BalanceReq{
				Coin:    data.Coin,
				Address: data.Address2,
			})
			if !isValid {
				app.Custum(c, gin.H{
					"code": 10001,
					"msg":  "输入的地址不正确",
				})
				return
			}
		}

		if data.Address3 != "" {
			isValid := rpcClient.ValidateAddress(&client.BalanceReq{
				Coin:    data.Coin,
				Address: data.Address3,
			})
			if !isValid {
				app.Custum(c, gin.H{
					"code": 10001,
					"msg":  "输入的地址不正确",
				})
				return
			}
		}
	}

	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}
