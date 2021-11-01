package utils

import (
	"wikipay-admin/tools"

	"github.com/spf13/viper"
)

//获取身份证路径
const (
	Feedback = "feedback"
	Logo     = "logo"
	Identity = "identity"

	testPath = "http://18.162.243.214:81"
	prodPath = "http://www.wikipay.net"
)

//获取图片地址
func GetImage(name, img string) string {
	var basePath string
	if viper.GetString("settings.application.mode") == string(tools.ModeTest) {
		basePath = testPath
	} else if viper.GetString("settings.application.mode") == string(tools.ModeProd) {
		basePath = prodPath
	}

	return basePath + "/" + name + "/" + img
}
