package monitor

import (
	"errors"
	"fmt"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type MonitorReq struct {
	Coin string `json:"coin" form:"coin" binding:"required"` // 币种名称
}

//
type MonitorAddress struct {
	Id         int     `json:"id" gorm:"type:int(10) unsigned;primary_key"`
	Coin       string  `json:"coin" form:"coin" gorm:"type:varchar(32);"` // 币种名称
	Type       int     `json:"type" form:"type" gorm:"type:tinyint(1);"`  //1、热钱包 2、冷钱包
	Address    string  `json:"address" form:"address" gorm:"type:varchar(64);"`
	Amount     float64 `json:"amount" form:"amount" gorm:"type:decimal(32,8);"`
	Reason     string  `json:"reason" form:"reason" gorm:"type:varchar(1024);"`
	Ratio      float64 `json:"ratio" form:"ratio" gorm:"type:decimal(6,2);;"`
	Status     int     `json:"status" form:"status" gorm:"type:tinyint(1);"`     //1、正常 2、错误
	Comment    string  `json:"comment" form:"comment" gorm:"type:varchar(128);"` // 备注
	DataScope  string  `json:"dataScope" gorm:"-"`
	Info       string  `json:"info" form:"info" gorm:"-"`
	FeeAddress string  `json:"-" form:"-" gorm:"-"`
	models.BaseModel
}

//
type SatisticaInfo struct {
	Coin             string  `json:"coin"`
	HoldWalletAmount float64 `json:"holdWalletAmount"`
	ColdWalletAmount float64 `json:"coldWalletAmount"`
	Total            float64 `json:"total"`
	TotalCny         float64 `json:"totalCny"`
}

//资产对比
type AssetComparisonInfo struct {
	HotWallet     []HotWallet     `json:"hotWallet"`
	AddressAssets []AddressAssets `json:"addressAssets"`
}

//热钱包地址资产
type HotWallet struct {
	Coin             string  `json:"coin"`
	AccountAssets    float64 `json:"accountAssets"`    //账户资产
	HoldWalletAmount float64 `json:"holdWalletAmount"` //热钱包地址资产
	Diff             float64 `json:"diff"`             //差额
	DiffCny          float64 `json:"diffCny"`          //差额换算人民币
}

//地址资产合计
type AddressAssets struct {
	Coin                string  `json:"coin"`
	AccountAssets       float64 `json:"accountAssets"`       //账户资产
	AddressAssetsAmount float64 `json:"addressAssetsAmount"` //地址资产合计
	Diff                float64 `json:"diff"`                //差额
	DiffCny             float64 `json:"diffCny"`             //差额换算人民币
}

//比率
type RatioReq struct {
	Coin   string    `json:"coin"`
	Type   int       `json:"type"`
	Ids    []int     `json:"ids"`
	Ratios []float64 `json:"ratios"`
}

//
type MonitorAddressReq struct {
	Coin string `json:"coin" binding:"required"`
	Type int    `json:"type" binding:"required"`
}

//表名
func (e *MonitorAddress) TableName() string {
	return "monitor_address"
}

//创建
func (e *MonitorAddress) Create() (MonitorAddress, error) {
	var data MonitorAddress
	i := 0
	orm.Eloquent.Table(e.TableName()).
		Where("type=? AND coin = ? AND address = ? AND deleted_at IS NULL",
			e.Type,
			e.Coin,
			e.Address,
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
func (e *MonitorAddress) UpdateAddress() error {
	err := orm.Eloquent.Exec("UPDATE monitor_address SET amount = ?, updated_at = NOW() WHERE type=? AND coin = ? AND address = ? AND deleted_at IS NULL",
		e.Amount,
		e.Type,
		e.Coin,
		e.Address).Error
	return err
}

//获取列表不带分页
func (e *MonitorAddress) Get() ([]MonitorAddress, error) {
	var address []MonitorAddress

	table := orm.Eloquent.Table(e.TableName())
	if e.Type != 0 {
		table = table.Where("type = ?", e.Type)
	}
	if e.Coin != "" {
		table = table.Where("coin = ?", e.Coin)
	}
	if e.Ratio > 0 {
		table = table.Where("ratio > 0")
	}

	if e.Address != "" {
		table = table.Where("address = ?", e.Address)
	}
	if e.Comment != "" {
		table = table.Where("comment = ?", e.Comment)
	}

	if err := table.Order("amount DESC").
		Find(&address).
		Error; err != nil {
		return address, err
	}
	return address, nil
}

//获取列表带分页
func (e *MonitorAddress) GetPage(pageSize int, pageIndex int) ([]MonitorAddress, int, error) {
	var doc []MonitorAddress

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if e.Type != 0 {
		table = table.Where("type = ?", e.Type)
	}
	if e.Coin != "" {
		table = table.Where("coin = ?", e.Coin)
	}
	if e.Info != "" {
		table = table.Where(fmt.Sprintf("concat(address,comment) like '%s%s%s'", "%", e.Info, "%"))
	}

	//数据权限控制
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope("monitor_address", table)
	if err != nil {
		return nil, 0, err
	}

	var count int
	if err := table.Order("created_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)
	return doc, count, nil
}

//更新
func (e *MonitorAddress) Update(id int) (update MonitorAddress, err error) {
	if err = orm.Eloquent.Table(e.TableName()).
		Where("id = ?", id).
		First(&update).Error; err != nil {
		return
	}

	//
	if err = orm.Eloquent.
		Table(e.TableName()).
		Model(&update).
		Updates(&e).Error; err != nil {
		return
	}
	return
}

//更新占比
func (e *MonitorAddress) UpdateRatio(ratio RatioReq) (bool, error) {
	//	err := orm.Eloquent.Transaction(func(tx *gorm.DB) error {
	if err := orm.Eloquent.Exec("UPDATE monitor_address SET ratio = 0  WHERE coin = ? AND type = ?",
		ratio.Coin,
		ratio.Type).Error; err != nil {
		return false, err
	}

	for k, v := range ratio.Ids {
		fmt.Println(ratio.Ratios[k], v)
		if err := orm.Eloquent.Exec("UPDATE monitor_address SET ratio = ?  WHERE id = ? AND coin = ? AND type = ?",
			ratio.Ratios[k],
			v,
			ratio.Coin,
			ratio.Type).Error; err != nil {
			return false, err
		}
	}
	return true, nil
}

//删除
func (e *MonitorAddress) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).
		Where("id = ?", id).Delete(&MonitorAddress{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//统计
func (e *MonitorAddress) Satistical() ([]MonitorAddress, error) {
	var infos []MonitorAddress
	err := orm.Eloquent.Table(e.TableName()).
		Select("any_value(coin) AS coin,any_value(type) AS type,IFNULL(any_value(SUM(amount)),0) AS amount").
		Group("coin,type").
		Find(&infos).Error
	return infos, err
}
