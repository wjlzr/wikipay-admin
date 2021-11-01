package monitor

import (
	"strings"
	"wikipay-admin/common"
	"wikipay-admin/models/monitor"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/rpc/client"
	monitorservice "wikipay-admin/service/monitor"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//获取
func GetMonitorAddress(c *gin.Context) {
	var data monitor.MonitorAddress

	data.Id, _ = tools.StringToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.OK(c, result, "")
}

//获取列表
func GetMonitorAddressList(c *gin.Context) {
	var data monitor.MonitorAddress

	var (
		err       error
		pageSize  = 10
		pageIndex = 1
	)

	data.Type = 1
	if strType := c.Request.FormValue("type"); strType != "" {
		data.Type = tools.StrToInt(err, strType)
	}
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

//添加
func InsertMonitorAddress(c *gin.Context) {
	var data monitor.MonitorAddress

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)

	//检查地址
	if data.Coin != "" {
		rpcClient := client.NewClient()
		isValid := rpcClient.ValidateAddress(&client.BalanceReq{
			Coin:    data.Coin,
			Address: data.Address,
		})
		if !isValid {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "输入的地址不正确",
			})
			return
		}
	}
	data.Coin = strings.ToUpper(data.Coin)
	result, err := data.Create()
	tools.HasError(err, "", -1)

	app.OK(c, result, "")
}

//更新比率
func UpdateMonitorRatio(c *gin.Context) {
	var data monitor.RatioReq

	err := c.ShouldBind(&data)
	tools.HasError(err, "数据解析失败", -1)

	var total float64
	for _, ratio := range data.Ratios {
		total = total + ratio
	}

	if xfloat64.FromFloatCmp(total, 1) != 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "比率设置错误，不等于100%",
		})
		return
	}

	data.Coin = strings.ToUpper(data.Coin)
	isUpdate, err := new(monitor.MonitorAddress).UpdateRatio(data)
	tools.HasError(err, "", -1)

	app.OK(c, gin.H{
		"success": isUpdate,
	}, "")
}

//更新
func UpdateMonitorAddress(c *gin.Context) {
	var data monitor.MonitorAddress

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)

	//检查地址
	if data.Coin != "" {
		rpcClient := client.NewClient()
		isValid := rpcClient.ValidateAddress(&client.BalanceReq{
			Coin:    data.Coin,
			Address: data.Address,
		})
		if !isValid {
			app.Custum(c, gin.H{
				"code": 10001,
				"msg":  "输入的地址不正确",
			})
			return
		}
	}

	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除
func DeleteMonitorAddress(c *gin.Context) {
	var data monitor.MonitorAddress

	id, _ := tools.StringToInt(c.Param("id"))
	_, err := data.Delete(id)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}

//同步
func SyncMonitorAddress(c *gin.Context) {
	var req monitor.MonitorAddressReq
	if err := c.ShouldBind(&req); err != nil {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "请选择钱包类型或币种",
		})
		return
	}

	err := monitorservice.SyncAddress(&req)
	if err != nil {
		app.Custum(c, gin.H{
			"code": 10002,
			"msg":  "获取或更新数据失败",
		})
		return
	}

	app.OK(c, nil, "同步成功")
}
