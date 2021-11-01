package card

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Card struct {
	Id          int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` // 序列号
	UserId      int    `json:"userId" gorm:"type:bigint(10);"`              // 用户序列号
	Currencies  string `json:"currencies" gorm:"type:varchar(50);"`         // 支持的法定货币列表
	Brand       string `json:"brand" gorm:"type:varchar(30);"`              // 品牌类型(visa....)
	CardId      string `json:"cardId" gorm:"type:varchar(100);"`            // 卡识别码
	CardNo      string `json:"cardNo" gorm:"type:varchar(50);"`             // 卡号
	CardLimit   string `json:"cardLimit" gorm:"type:decimal(16,2);"`        // 额度
	Cvv         string `json:"cvv" gorm:"type:varchar(10);"`                // cvv码
	CardStatus  string `json:"cardStatus" gorm:"type:varchar(20);"`         // /状态
	ExpiryYear  int    `json:"expiryYear" gorm:"type:int(4);"`              // 过期年份
	Name        string `json:"name" gorm:"type:varchar(50);"`               // 名字
	ExpiryMonth int    `json:"expiryMonth" gorm:"type:int(2);"`             // 过期月份
	CreateAt    int    `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
}

func (Card) TableName() string {
	return "card"
}

// 创建Card
func (e *Card) Create() (Card, error) {
	var doc Card
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Card
func (e *Card) Get() (Card, error) {
	var doc Card
	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Card带分页
func (e *Card) GetPage(pageSize, pageIndex int, info string) ([]Card, int, error) {
	var doc []Card
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
		where = fmt.Sprintf("concat(user_id,card_no) like '%s%s%s'", "%", info, "%")
	}
	if err := table.Where(where).
		Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}
	var count int
	table.Count(&count)
	return doc, count, nil
}

// 更新Card
func (e *Card) Update(id int) (update Card, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Card
func (e *Card) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Card{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Card) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Card{}).Error; err != nil {
		return
	}
	Result = true
	return
}
