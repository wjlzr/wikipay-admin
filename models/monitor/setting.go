package monitor

import (
	"errors"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
)

//
type MonitorSettingReq struct {
	Coin string `json:"coin" form:"coin" binding:"required"` // 币种名称
	Type int    `json:"type" form:"type" binding:"required" `
}

//设置
type MonitorSetting struct {
	Id         int     `json:"id"  gorm:"type:int(10) unsigned;primary_key"`
	Coin       string  `json:"coin" binding:"required" gorm:"type:varchar(32);"`     // 币种名称
	Type       int     `json:"type" binding:"required" gorm:"type:tinyint(1);"`      //类型 1、热钱包 2、冷钱包
	Min        float64 `json:"min" binding:"required" gorm:"type:decimal(32,8);"`    //最小数量
	Max        float64 `json:"max" gorm:"type:decimal(32,8);"`                       //最大数量
	GasPrice   float64 `json:"gasPrice" gorm:"type:decimal(32,8);"`                  //gas费用
	GasStatus  int     `json:"gasStatus" binding:"required" gorm:"type:tinyint(1);"` //gas状态 1、自动  2、手动
	Status     int     `json:"status" gorm:"type:tinyint(1);"`                       //归集状态 1、可用 2、信用
	Address1   string  `json:"address1" gorm:"type:varchar(128);"`                   //打手续费地址
	Address2   string  `json:"address2" gorm:"type:varchar(128);"`                   //打手续费地址
	Address3   string  `json:"address3" gorm:"type:varchar(128);"`
	FeeAddress string  `json:"feeAddress" gorm:"type:varchar(128);"` //手续费地址
	Day        int     `json:"day" gorm:"type:tinyint(1);"`          //天
	Week       int     `json:"week" gorm:"type:tinyint(1);"`         //周
	Hour       int     `json:"hour" gorm:"type:tinyint(1);"`         //时
	Comment    string  `json:"comment" gorm:"type:varchar(256);"`
	models.BaseModel
}

//表名
func (e *MonitorSetting) TableName() string {
	return "monitor_setting"
}

//创建
func (e *MonitorSetting) Create() (MonitorSetting, error) {
	var data MonitorSetting
	i := 0
	orm.Eloquent.Table(e.TableName()).
		Where("type=? AND coin = ? AND deleted_at IS NULL",
			e.Type,
			e.Coin,
		).Count(&i)

	if i > 0 {
		return data, errors.New("记录已经存在！")
	}

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return data, err
	}
	data = *e
	return data, nil
}

//更新
func (e *MonitorSetting) Update(id int) (update MonitorSetting, err error) {
	if err = orm.Eloquent.Table(e.TableName()).
		Where("id = ?", id).
		First(&update).Error; err != nil {
		return
	}

	err = orm.Eloquent.Exec("UPDATE monitor_setting SET coin = ?,type = ?, max = ?, min = ?, day = ?, week=?, hour = ?, gas_status=?, gas_price=?, address1=?, address2=?, address3=?, fee_address = ?, comment = ? WHERE id = ?",
		e.Coin,
		e.Type,
		e.Max,
		e.Min,
		e.Day,
		e.Week,
		e.Hour,
		e.GasStatus,
		e.GasPrice,
		e.Address1,
		e.Address2,
		e.Address3,
		e.Comment,
		e.FeeAddress,
		e.Id,
	).Error

	return
}

//获取规则
func (e *MonitorSetting) Get() (MonitorSetting, error) {
	var setting MonitorSetting

	table := orm.Eloquent.Table(e.TableName())
	if e.Type != 0 {
		table = table.Where("type = ?", e.Type)
	}
	if e.Coin != "" {
		table = table.Where("coin = ?", e.Coin)
	}
	if err := table.Scan(&setting).
		Error; err != nil {
		return setting, err
	}
	return setting, nil
}
