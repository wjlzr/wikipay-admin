package user

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

type Advert struct {
	Id          int    `json:"id" gorm:"type:int(10);primary_key"`    //
	Priority    int    `json:"priority" gorm:"type:int(10);"`         // 排列码
	Name        string `json:"name" gorm:"type:varchar(100);"`        // 名称
	Description string `json:"description" gorm:"type:varchar(300);"` // 描述
	ImageUrl    string `json:"imageUrl" gorm:"type:varchar(400);"`    // 图片链接
	LinkUrl     string `json:"linkUrl" gorm:"type:varchar(400);"`     // 广告链接
	Lang        string `json:"lang" gorm:"type:varchar(50);"`         // 语言
	Status      int    `json:"status" gorm:"type:tinyint(1);"`        // 状态
	CreateAt    int64  `json:"createAt" gorm:"type:bigint(13);"`      // 创建时间
	Type        int    `json:"type" gorm:"type:tinyint(1);"`          //
	DataScope   string `json:"dataScope" gorm:"-"`
	Params      string `json:"params"  gorm:"-"`
}

//
func (Advert) TableName() string {
	return "advert"
}

// 创建Advert
func (e *Advert) Create() (Advert, error) {
	var doc Advert

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Advert
func (e *Advert) Get() (Advert, error) {
	var doc Advert

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Advert带分页
func (e *Advert) GetPage(pageSize, pageIndex int, info string) ([]Advert, int, error) {
	var doc []Advert
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
		where = fmt.Sprintf("concat(priority) like '%s%s%s'", "%", info, "%")
	}
	if err := table.Where(where).
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Count(&count)
	return doc, count, nil
}

// 更新Advert
func (e *Advert) Update(id int) (update Advert, err error) {
	if err = orm.Eloquent.Table(e.TableName()).First(&update, id).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Advert
func (e *Advert) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).
		Delete(&Advert{}, id).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Advert) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Advert{}).Error; err != nil {
		return
	}
	Result = true
	return
}
