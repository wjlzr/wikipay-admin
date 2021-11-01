package operate

import (
	"wikipay-admin/models/operate"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//列表
func GetOperateList(c *gin.Context) {

	var (
		err       error
		pageSize  = 10
		pageIndex = 1
	)

	if size := c.Request.FormValue("pageSize"); size != "" {
		pageSize = tools.StrToInt(err, size)
	}
	if index := c.Request.FormValue("pageIndex"); index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	result, count, err := operate.OperateManage{}.GetPage(pageSize, pageIndex)

	tools.HasError(err, "抱歉未找到相关信息", -1)

	app.PageOK(c, result, count, pageIndex, pageSize, "")
}

//新增
func InsertOperate(c *gin.Context) {
	var data operate.OperateManage
	err := c.ShouldBindJSON(&data)
	//时间判断
	if data.StartTime >= data.EndTime {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "结束时间必须大于开始时间",
		})
		return
	}
	//
	if xfloat64.FromFloatCmp(data.Cny, 0) < 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "输入金额不能小于0",
		})
		return
	}
	tools.HasError(err, "", 500)
	result, err := data.Create()
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//更新
func UpdateOperate(c *gin.Context) {
	var data operate.OperateManage
	err := c.MustBindWith(&data, binding.JSON)
	//时间判断
	if data.StartTime >= data.EndTime {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "结束时间必须大于开始时间",
		})
		return
	}
	//
	if xfloat64.FromFloatCmp(data.Cny, 0) < 0 {
		app.Custum(c, gin.H{
			"code": 10001,
			"msg":  "输入金额不能小于0",
		})
		return
	}

	tools.HasError(err, "数据解析失败", -1)
	result, err := data.Update(data.Id)
	tools.HasError(err, "", -1)
	app.OK(c, result, "")
}

//获取单个详情
func GetOperate(c *gin.Context) {
	var data operate.OperateManage

	data.Id, _ = tools.StringToInt64(c.Param("id"))
	result, err := data.Get()
	tools.HasError(err, "抱歉未找到相关信息", -1)
	app.OK(c, result, "")
}

//删除
func DeleteOperate(c *gin.Context) {
	var data operate.OperateManage

	IDS := tools.IdsStrToIdsIntGroup("id", c)
	_, err := data.BatchDelete(IDS)
	tools.HasError(err, msg.DeletedFail, 500)
	app.OK(c, nil, msg.DeletedSuccess)
}
