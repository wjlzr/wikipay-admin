package monitor

import (
	orm "wikipay-admin/database"
	"wikipay-admin/models"
)

//
type MonitorManualSetting struct {
	Id          int     `json:"id"  gorm:"type:int(10) unsigned;primary_key"`
	Coin        string  `json:"coin" form:"coin" gorm:"type:varchar(32);"`
	FromAddress string  `json:"fromAddress" binding:"required" gorm:"type:varchar(128);"`
	ToAddress   string  `json:"toAddress" binding:"required" gorm:"type:varchar(128);"`
	FeeAddress  string  `json:"feeAddress"  gorm:"type:varchar(128);"`
	Amount      float64 `json:"amount" binding:"required" gorm:"type:decimal(32,8);"` //数量
	GasPrice    float64 `json:"gasPrice" gorm:"type:decimal(32,8);"`                  //gas费用
	GasStatus   int     `json:"gasStatus" binding:"required" gorm:"type:tinyint(1);"` //gas状态 1、自动  2、手动
	Comment     string  `json:"comment" gorm:"type:varchar(256);"`
	models.BaseModel
}

//表名
func (e *MonitorManualSetting) TableName() string {
	return "monitor_manual_setting"
}

//创建
func (e *MonitorManualSetting) Create() (MonitorManualSetting, error) {
	var data MonitorManualSetting

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return data, err
	}
	data = *e
	return data, nil
}
