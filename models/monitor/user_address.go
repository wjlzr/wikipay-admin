package monitor

import (
	"fmt"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/models/xfloat64"
	"wikipay-admin/tools"
)

//
type MonitorUserAddress struct {
	Id        int     `json:"id" gorm:"type:int(10) unsigned;primary_key"`
	Coin      string  `json:"coin" gorm:"type:varchar(32);"` // 币种名称
	Address   string  `json:"address" gorm:"type:varchar(64);"`
	Amount    float64 `json:"amount" gorm:"type:decimal(32,8);"`
	UserId    int64   `json:"userId" gorm:"type:bigint(10);"`
	Info      string  `json:"-" form:"info" gorm:"-"`
	DataScope string  `json:"-" gorm:"-"`
	models.BaseModel
}

//表名
func (e *MonitorUserAddress) TableName() string {
	return "monitor_user_address"
}

//创建
func CreateAndUpdate(e *MonitorUserAddress) (bool, error) {
	var mAddress MonitorAddress
	orm.Eloquent.Raw(fmt.Sprintf(`
			SELECT * FROM monitor_address
			WHERE coin = '%s' AND address = '%s'`,
		e.Coin,
		e.Address)).Scan(&mAddress)

	if mAddress.Coin != "" {
		return true, nil
	}

	var addr MonitorUserAddress
	orm.Eloquent.Raw(fmt.Sprintf(`
			SELECT * FROM monitor_user_address
			WHERE coin = '%s' AND address = '%s' AND TRUNCATE(amount,8) = %v`,
		e.Coin,
		e.Address,
		e.Amount)).Scan(&addr)

	if addr.Coin != "" {
		return true, nil
	}

	orm.Eloquent.Raw(fmt.Sprintf(`
			SELECT * FROM monitor_user_address 
			WHERE coin = '%s' AND address = '%s'`,
		e.Coin,
		e.Address)).
		Scan(&addr)

	if addr.Coin != "" {
		err := orm.Eloquent.Exec("UPDATE monitor_user_address SET amount = ?, user_id = ?, updated_at = NOW() WHERE coin = ? AND address = ?",
			e.Amount,
			e.UserId,
			e.Coin,
			e.Address).Error
		if err == nil {
			return true, nil
		}
	}

	result := orm.Eloquent.Exec("INSERT INTO monitor_user_address(`coin`,`address`,`amount`,`user_id`,`updated_at`) VALUES(?,?,?,?,now())",
		e.Coin,
		e.Address,
		e.Amount,
		e.UserId)
	if result.Error != nil {
		err := result.Error
		return false, err
	}
	return true, nil
}

//获取列表带分页
func (e *MonitorUserAddress) GetPage(pageSize int, pageIndex int) ([]MonitorUserAddress, int, error) {
	var doc []MonitorUserAddress

	table := orm.Eloquent.Select("*").Table(e.TableName())

	if e.Coin != "" {
		table = table.Where("coin = ?", e.Coin)
	}
	if e.Info != "" {
		table = table.Where(fmt.Sprintf("concat(address,user_id) like '%s%s%s'", "%", e.Info, "%"))
	}

	// 数据权限控制
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("monitor_user_address", table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Order("amount DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

//
func (e *MonitorUserAddress) Get() ([]MonitorUserAddress, error) {
	var addresses []MonitorUserAddress

	table := orm.Eloquent.Table(e.TableName())
	if e.Coin != "" {
		table = table.Where("coin = ?", e.Coin)
	}
	if e.Address != "" {
		table = table.Where("address = ?", e.Address)
	}

	if err := table.Order("created_at DESC").
		Find(&addresses).
		Error; err != nil {
		return addresses, err
	}
	return addresses, nil
}

//
func (e *MonitorUserAddress) GetDatas(min, max float64) ([]MonitorUserAddress, error) {
	var infos []MonitorUserAddress

	table := orm.Eloquent.Table(e.TableName())
	if e.Coin != "" {
		table = table.Where("coin = ?", e.Coin)
	}
	if xfloat64.FromFloatCmp(min, 0) > 0 {
		table = table.Where("amount >= ?", min)
	}
	if xfloat64.FromFloatCmp(max, 0) > 0 {
		table = table.Where("amount <= ?", max)
	}

	if err := table.Order("amount DESC").
		Find(&infos).
		Error; err != nil {
		return infos, err
	}
	return infos, nil
}
