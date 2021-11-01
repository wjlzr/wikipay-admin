package financial

import (
	"strings"
	"wikipay-admin/models/financial"
	"wikipay-admin/models/home"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//从时间区间获取财务信息
func GetFinancialInfoWithDateTime(c *gin.Context) {
	var req home.FinancialReq

	err := c.BindQuery(&req)
	tools.HasError(err, "数据解析失败", -1)

	resp, err := home.GetFinancialInfoWithDateTime(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, resp, "")
}

//交易账单
func GetTransactionFinancialInfo(c *gin.Context) {
	var req financial.TransactionReq

	err := c.BindQuery(&req)
	tools.HasError(err, "数据解析失败", -1)

	utils.Pagination(&req.PageSize, &req.PageNum)

	req.Coin = strings.ToUpper(req.Coin)
	resp, count, err := financial.GetTransactionInfo(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, gin.H{
		"financiallist": resp,
		"count":         count,
	}, "")
}


