package tools

import (
	"strings"

	"github.com/gin-gonic/gin"
)

//获取URL中批量id并解析
func IdsStrToIdsIntGroup(key string, c *gin.Context) []int {
	return idsStrToIdsIntGroup(c.Param(key))
}

func idsStrToIdsIntGroup(keys string) []int {
	IDS := make([]int, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := StringToInt(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}

///
func IdsStrToIdsInt64Group(key string, c *gin.Context) []int64 {
	return idsStrToIdsInt64Group(c.Param(key))
}

func idsStrToIdsInt64Group(keys string) []int64 {
	IDS := make([]int64, 0)
	ids := strings.Split(keys, ",")
	for i := 0; i < len(ids); i++ {
		ID, _ := StringToInt64(ids[i])
		IDS = append(IDS, ID)
	}
	return IDS
}
