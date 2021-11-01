package system

import (
	"wikipay-admin/models"
	"wikipay-admin/tools"
	"wikipay-admin/tools/app"

	"github.com/gin-gonic/gin"
)

func GetInfo(c *gin.Context) {

	var roles = make([]string, 1)
	roles[0] = tools.GetRoleName(c)

	var permissions = make([]string, 1)
	permissions[0] = "*:*:*"
	RoleMenu := models.RoleMenu{}
	RoleMenu.RoleId = tools.GetRoleId(c)

	var mp = make(map[string]interface{})
	mp["roles"] = roles
	//fmt.Println(tools.GetRoleName(c))
	if tools.GetRoleName(c) == "admin" || tools.GetRoleName(c) == "系统管理员" {
		mp["permissions"] = permissions
	} else {
		list, _ := RoleMenu.GetPermis()
		mp["permissions"] = list
	}
	//fmt.Println(tools.GetUserId(c))
	sysuser := models.SysUser{}
	sysuser.UserId = tools.GetUserId(c)
	user, err := sysuser.Get()
	tools.HasError(err, "", 600)

	mp["introduction"] = " am a super administrator"

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}
	mp["name"] = user.NickName

	app.OK(c, mp, "")
}
