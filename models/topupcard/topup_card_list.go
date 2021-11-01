package topupcard

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type TopupCardList struct {
	Id          int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` //
	UserId      string `json:"userId" gorm:"type:bigint(10);"`              // 用户序列号
	CardNo      string `json:"cardNo" gorm:"type:varchar(18);"`             // 卡号
	CreateAt    string `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	OrderNumber string `json:"orderNumber" gorm:"type:bigint(18);"`         //
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
}

//筛选条件扩展
type SearchParams struct {
	StartTime string `form:"startTime"`
	EndTime   string `form:"endTime"`
}

//
func (TopupCardList) TableName() string {
	return "topup_card_list"
}

// 创建
func (e *TopupCardList) Create() (TopupCardList, error) {
	var doc TopupCardList
	result := orm.Eloquent.Table(e.TableName()).Create(&e)

	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *TopupCardList) Get() (TopupCardList, error) {
	var doc TopupCardList

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取带分页
func (e *TopupCardList) GetPage(pageSize, pageIndex int, info string, req SearchParams) ([]TopupCardList, int, error) {
	var doc []TopupCardList
	table := orm.Eloquent.Select("*").Table(e.TableName())

	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}

	var where string
	if info != "" {
		where = fmt.Sprintf("concat(card_no,order_number,user_id) like '%s%s%s'", "%", info, "%")
	}

	if req.StartTime != "" {
		table = table.Where("create_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		table = table.Where("create_at <= ?", req.EndTime)
	}

	if err := table.Where(where).
		Order("create_at desc").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Where(where).Count(&count)
	return doc, count, nil
}

//更新
func (e *TopupCardList) Update(id int) (update TopupCardList, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除
func (e *TopupCardList) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&TopupCardList{}).Error; err != nil {
		success = false
		return
	}

	success = true
	return
}

//批量删除
func (e *TopupCardList) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&TopupCardList{}).Error; err != nil {
		return
	}

	Result = true
	return
}
