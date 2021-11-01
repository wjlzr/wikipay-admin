package topupcard

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type TopupCard struct {
	Id        int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` //
	CardNo    string `json:"cardNo" gorm:"type:varchar(18);"`             // 卡号
	Money     string `json:"money" gorm:"type:decimal(9,2);"`             // 面值
	Type      string `json:"type" gorm:"type:tinyint(1);"`                // 类型
	Status    string `json:"status" gorm:"type:tinyint(1);"`              // 状态
	ChannelNo string `json:"channelNo" gorm:"type:varchar(10);"`          // 渠道编号
	CreateAt  string `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	ExpireAt  string `json:"expireAt" gorm:"type:bigint(13);"`            // 过期时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (TopupCard) TableName() string {
	return "topup_card"
}

// 创建
func (e *TopupCard) Create() (TopupCard, error) {
	var doc TopupCard
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e

	return doc, nil
}

//获取
func (e *TopupCard) Get() (TopupCard, error) {
	var doc TopupCard

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取
func (e *TopupCard) GetPage(pageSize, pageIndex int, info string) ([]TopupCard, int, error) {

	var doc []TopupCard
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
		where = fmt.Sprintf("concat(card_no) like '%s%s%s'", "%", info, "%")
	}
	if err := table.Where(where).
		Order("id DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Count(&count)
	return doc, count, nil
}

//更新
func (e *TopupCard) Update(id int) (update TopupCard, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

//删除
func (e *TopupCard) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&TopupCard{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *TopupCard) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&TopupCard{}).Error; err != nil {
		return
	}
	Result = true
	return
}
