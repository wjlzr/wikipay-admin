package financial

import (
	"wikipay-admin/models/financial"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//获取首页用户统计信息
func GetWithdrawAuditInfo(c *gin.Context) {
	var req financial.WithdrawAuditReq
	err := c.BindQuery(&req)
	tools.HasError(err, "数据解析失败", -1)

	utils.Pagination(&req.PageSize, &req.PageNum)
	info, count, err := financial.FindWithdrawAuditInfo(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, gin.H{
		"count":             count,
		"withdrawAuditList": info,
	}, "")
}

//提币审核
func WithdrawAudit(c *gin.Context) {
	var req financial.WithdrawReq

	if err := c.ShouldBindJSON(&req); err != nil {
		tools.HasError(err, "数据解析失败", -1)
	}

	msgCode, _ := financial.WithdrawAuditOn(&req)
	//fmt.Println(msgCode, financial.TipMessage[msgCode])
	if msgCode == financial.WithdrawSuccess {
		app.OK(c, gin.H{
			"success": true}, "")
	} else {
		app.Custum(c, gin.H{
			"code": 10002,
			"msg":  financial.TipMessage[msgCode]})
	}
}

//点差列表
func ListBidAsk(c *gin.Context) {
	var req financial.TimeReq
	err := c.BindQuery(&req)
	tools.HasError(err, "数据解析失败", -1)

	info, err := financial.FindBidAsk(&req)

	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, info, "")
}

//获取手续费
func ListFee(c *gin.Context) {
	var req financial.TimeReq
	err := c.BindQuery(&req)
	tools.HasError(err, "数据解析失败", -1)

	info, err := financial.FindFee(&req)

	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, info, "")
}
