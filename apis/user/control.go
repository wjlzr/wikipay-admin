package user

import (
	"wikipay-admin/models/user"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"
	"wikipay-admin/tools/app/msg"
	"wikipay-admin/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//获取用户风控
func GetUserControlInfos(c *gin.Context) {
	var req user.ControlReq
	err := c.BindWith(&req, binding.Query)
	tools.HasError(err, "数据解析失败", -1)

	utils.Pagination(&req.PageSize, &req.PageNum)
	info, count, err := user.GetUserControlInfos(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, gin.H{
		"count":            count,
		"maxEmailSendNum":  user.MaxEmailSendNum,
		"maxPaypwdErrNum":  user.MaxPaypwdErrNum,
		"maxSmsSendNum":    user.MaxSmsSendNum,
		"userControlInfos": info,
	}, "")
}

//冻结、解锁
func UpdateUserStatus(c *gin.Context) {
	var req user.UserStatusReq
	err := c.ShouldBind(&req)
	tools.HasError(err, "数据解析失败", -1)

	err = user.UpdateStatus(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, gin.H{
		"success": true,
	}, "")
}

//重置
func ResetUser(c *gin.Context) {
	var req user.UserStatusReq
	err := c.ShouldBind(&req)
	tools.HasError(err, "数据解析失败", -1)

	err = user.Reset(&req)
	tools.HasError(err, msg.ServerInternalError, 500)
	app.OK(c, gin.H{
		"success": true,
	}, "")
}
