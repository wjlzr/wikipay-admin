package monitor

import (
	orm "wikipay-admin/database"
	"wikipay-admin/models"
)

//
type MonitorHistory struct {
	Id               int     `json:"id" gorm:"type:int(10) unsigned;primary_key"`
	Coin             string  `json:"coin,omitempty" binding:"required" form:"coin" gorm:"type:varchar(32);"` // 币种名称
	FromAddress      string  `json:"fromAddress,omitempty" form:"fromAddress"  gorm:"type:varchar(64);"`     //被归集地址
	ToAddress        string  `json:"toAddress,omitempty" form:"toAddress"  gorm:"type:varchar(64);"`         //归集地址
	TxId             string  `json:"txId,omitempty" form:"txId"  gorm:"type:varchar(128);"`                  //交易哈希
	TxFee            float64 `json:"txFee,omitempty" form:"txFee"  gorm:"type:decimal(32,8);"`               //
	CollectionAmount float64 `json:"collectionAmount,omitempty" form:"collectionAmount"  gorm:"type:decimal(32,8);"`
	Balance          float64 `json:"balance,omitempty" form:"balance"  gorm:"type:decimal(32,8);"`
	SendAmount       float64 `json:"sendAmount,omitempty" form:"sendAmount"  gorm:"type:decimal(32,8);"` //打币数量
	SendFee          float64 `json:"sendFee,omitempty" form:"sendFee"  gorm:"type:decimal(32,8);"`       //打币费用
	//	Type             int     `gorm:"type:tinyint(1);"`                                                   //笔数                                             //                                                // 1、自动归集 2、主动归集                                             //                                               //
	Status    int    `json:"status,omitempty" form:"status"  gorm:"type:tinyint(1);"`
	Count     int    `json:"count,omitempty" gorm:"-"` //笔数
	DataScope string `json:"-" gorm:"-"`
	Address   string `json:"address,omitempty" form:"address" gorm:"-"` //状态
	models.BaseModel
}

//
type MonitorHistoryReq struct {
	Coin      string `form:"coin"`     ///币种
	Address   string `form:"address"`  //地址
	PageSize  int32  `form:"pageSize"` //页数
	PageIndex int32  `form:"pageIndex"`
	StartTime string `form:"startTime"` //开始时间
	EndTime   string `form:"endTime"`   //结束时间
	Status    int    `form:"status"`    //状态
}

//表名
func (e *MonitorHistory) TableName() string {
	return "monitor_history"
}

//创建
func (e *MonitorHistory) Create() error {
	err := orm.Eloquent.Table(e.TableName()).Create(&e).Error
	return err
}

//获取列表带分页
func (e *MonitorHistory) GetPage(req *MonitorHistoryReq) ([]MonitorHistory, int, error) {
	var doc []MonitorHistory

	table := orm.Eloquent.Select("ANY_VALUE(coin) AS coin, ANY_VALUE(from_address) AS address, ANY_VALUE(created_at) AS created_at, ANY_VALUE(SUM(collection_amount)) AS collection_amount, ANY_VALUE(COUNT(*)) AS count, ANY_VALUE(SUM(tx_fee)) AS tx_fee, ANY_VALUE(status) AS status").Table(e.TableName())
	if req.Coin != "" {
		table = table.Where("coin = ?", req.Coin)
	}
	if req.Address != "" {
		table = table.Where("from_address = ?", req.Address)
		//table = table.Where(fmt.Sprintf("concat(from_address,to_address) like '%s%s%s'", "%", e.Address, "%"))
	}
	if req.StartTime != "" {
		table = table.Where("DATE(created_at)>= ?", req.StartTime)
	}
	if req.EndTime != "" {
		table = table.Where("DATE(created_at)<= ?", req.EndTime)
	}

	//分页查询
	var count int
	if err := table.Offset(req.PageIndex).
		Limit(req.PageSize).
		Group("from_address").
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)

	return doc, count, nil
}

//获取明细带分页
func (e *MonitorHistory) GetDetails(req *MonitorHistoryReq) ([]MonitorHistory, int, error) {
	var details []MonitorHistory

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if req.Coin != "" {
		table = table.Where("coin = ?", req.Coin)
	}
	if req.Address != "" {
		table = table.Where("from_address = ?", req.Address)
	}
	if req.Status > 0 {
		table = table.Where("status = ?", req.Status)
	}

	//分页查询
	var count int
	if err := table.Order("created_at DESC").
		Offset(req.PageIndex).
		Limit(req.PageSize).
		Find(&details).Error; err != nil {
		return nil, 0, err
	}
	table.Where("`deleted_at` IS NULL").Count(&count)

	// if err := orm.Eloquent.Select("ANY_VALUE(coin) AS coin, ANY_VALUE(to_address) AS address,ANY_VALUE(SUM(collection_amount)) AS collection_amount,ANY_VALUE(COUNT(*)) AS count,ANY_VALUE(SUM(tx_fee)) AS tx_fee").
	// 	Table(e.TableName()).
	// 	Where("coin = ? AND from_address = ?", req.Coin, req.Address).
	// 	Group("to_address").
	// 	Find(&total).Error; err != nil {
	// 	return nil, nil, 0, err
	// }

	return details, count, nil
}
