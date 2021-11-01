package wallet

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//
type Transfer struct {
	OrderNumber     string `json:"orderNumber" gorm:"type:bigint(18) unsigned;primary_key"` // 订单号
	TransFlowNumber string `json:"transFlowNumber" gorm:"type:varchar(32);"`                // 交易流水号
	FromId          int    `json:"fromId" gorm:"type:bigint(10);"`                          //
	ToId            int    `json:"toId" gorm:"type:bigint(10);"`                            //
	Money           string `json:"money" gorm:"type:decimal(32,8);"`                        //
	Status          int    `json:"status" gorm:"type:tinyint(1);"`                          // 状态(0-失败，1-成功)
	Comment         string `json:"comment" gorm:"type:varchar(100);"`                       // 备注
	CreateAt        int    `json:"createAt" gorm:"type:bigint(13);"`                        // 交易时间
	Trmo            string `json:"trmo" gorm:"type:varchar(1000);"`                         //
	Coin            string `json:"coin" gorm:"type:varchar(20);"`                           //
	CoinPrice       string `json:"coinPrice" gorm:"type:decimal(16,4);"`                    //
	Source          string `json:"source" gorm:"type:varchar(50);"`                         //
	DataScope       string `json:"dataScope" gorm:"-"`
	Params          string `json:"params"  gorm:"-"`
}

//
func (Transfer) TableName() string {
	return "transfer"
}

// 创建Transfer
func (e *Transfer) Create() (Transfer, error) {
	var doc Transfer
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Transfer
func (e *Transfer) Get() (Transfer, error) {
	var doc Transfer
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Transfer带分页
func (e *Transfer) GetPage(pageSize, pageIndex int, info string, req SearchParams) ([]Transfer, int, error) {
	var doc []Transfer
	table := orm.Eloquent.Select("*").Table(e.TableName())

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string
	if info != "" {
		where = fmt.Sprintf("concat(order_number) like '%s%s%s'", "%", info, "%")
	}

	if req.StartTime != "" {
		table = table.Where("create_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		table = table.Where("create_at <= ?", req.EndTime)
	}

	if err := table.Select("order_number,trans_flow_number,coin,coin_price,source,from_id,to_id,0+CAST(money AS char(32)) AS money,status,comment,create_at,trmo").
		Where(where).
		Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Where(where).Count(&count)
	return doc, count, nil
}

// 更新Transfer
func (e *Transfer) Update(orderNumber int64) (update Transfer, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", orderNumber).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Transfer
func (e *Transfer) Delete(orderNumber int64) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number = ?", orderNumber).Delete(&Transfer{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Transfer) BatchDelete(orderNumbers []int64) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("order_number in (?)", orderNumbers).Delete(&Transfer{}).Error; err != nil {
		return
	}
	Result = true
	return
}
