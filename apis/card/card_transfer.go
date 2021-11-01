package card

import (
	"wikipay-admin/models/card"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//列表
func GetCardTransferList(c *gin.Context) {
	var (
		data         card.CardTransfer
		info         string
		status       int
		cardType     int
		err          error
		pageSize     = 10
		pageIndex    = 1
		searchParams card.SearchParams
	)

	_ = c.ShouldBindQuery(&searchParams)

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}
	if s := c.Request.FormValue("info"); s != "" {
		info = s
	}
	if t := c.Request.FormValue("type"); t != "" {
		cardType = tools.StrToInt(err, t)
	}
	if st := c.Request.FormValue("status"); st != "" {
		status = tools.StrToInt(err, st)
	}

	data.DataScope = tools.GetUserIdStr(c)
	result, count, err := data.GetPage(pageSize, pageIndex, status, cardType, info, searchParams)
	tools.HasError(err, "", -1)
	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

//获取
func GetCardTransfer(c *gin.Context) {
	var data card.CardTransfer

	data.Id, _ = tools.StringToInt(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//添加
func InsertCardTransfer(c *gin.Context) {
	var data card.CardTransfer

	err := c.ShouldBindJSON(&data)
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新
func UpdateCardTransfer(c *gin.Context) {
	var data card.CardTransfer

	err := c.BindWith(&data, binding.JSON)
	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//删除
func DeleteCardTransfer(c *gin.Context) {
	var data card.CardTransfer

	IDS := tools.IdsStrToIdsIntGroup("id", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}

//审核
func AuditCardTransfer(c *gin.Context) {
	var (
		data card.CardTransfer
		req  card.CardTransferReq
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		tools.HasError(err, "数据解析失败", -1)
	}
	err := data.CardAudit(&req)
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, gin.H{
		"success": true,
	}, "")
}
