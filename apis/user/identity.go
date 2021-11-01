package user

import (
	"wikipay-admin/models/user"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
)

//获取认证列表
func GetIdentityList(c *gin.Context) {
	var req user.IdentityStatusReq

	err := c.BindQuery(&req)
	tools.HasError(err, "数据解析失败", -1)
	utils.Pagination(&req.PageSize, &req.PageNum)

	info, count, err := user.FindIdentitys(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	//	app.OK(c, info, "")

	app.PageOK(c, info, count, int(req.PageNum), int(req.PageSize), "")
}

//从序列号获取身体信息
func GetIdentityFromId(c *gin.Context) {
	id, _ := tools.StringToInt(c.Param("id"))

	req := user.IdentityStatusReq{
		Id: id,
	}
	info, _, err := user.FindIdentitys(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, info, "")
}

//更新用户认证
func AuditIdentity(c *gin.Context) {
	var req user.IdentityReq

	err := c.ShouldBindJSON(&req)
	tools.HasError(err, "", 500)

	msg, err := user.UpdateIdentity(&req)
	tools.HasError(err, "", -1)
	if msg != "" {
		app.OK(c, gin.H{
			"success": false,
		}, msg)
		return
	}
	app.OK(c, gin.H{
		"success": true,
	}, "")
}
