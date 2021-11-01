package user

import (
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

///
type Version struct {
	Id        int    `json:"id" gorm:"type:int(10) unsigned;primary_key"` // 序列号
	Platform  string `json:"platform" gorm:"type:varchar(20);"`           // 平台类型
	Version   string `json:"version" gorm:"type:varchar(20);"`            // 版本号
	Forced    int    `json:"forced" gorm:"type:tinyint(1);"`              // 1、是 2、否
	Url       string `json:"url" gorm:"type:varchar(256);"`               // 链接
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13);"`            // 创建时间
	En        string `json:"en" gorm:"type:varchar(2000);"`               // 英文
	ZhCN      string `json:"zhCN" gorm:"type:varchar(2000);column:zh-CN"` // 中文
	ZhHK      string `json:"zhHK" gorm:"type:varchar(2000);column:zh-HK"` // 香港
	ZhTW      string `json:"zhTW" gorm:"type:varchar(2000);column:zh-TW"` // 台湾
	Vi        string `json:"vi" gorm:"type:varchar(2000);"`
	Th        string `json:"th" gorm:"type:varchar(300);"` // 越南语
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (Version) TableName() string {
	return "version"
}

// 创建Version
func (e *Version) Create() (Version, error) {
	var doc Version

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Version
func (e *Version) Get() (Version, error) {
	var doc Version
	table := orm.Eloquent.Table(e.TableName())

	if err := table.First(&doc, e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Version带分页
func (e *Version) GetPage(pageSize int, pageIndex int) ([]Version, int, error) {
	var doc []Version

	table := orm.Eloquent.Select("*").Table(e.TableName())

	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int

	if err := table.Order("id DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).
		Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)
	return doc, count, nil
}

// 更新Version
func (e *Version) Update(id int) (update Version, err error) {
	if err = orm.Eloquent.Table(e.TableName()).
		First(&update, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Version
func (e *Version) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).
		Delete(&Version{}, id).
		Error; err != nil {

		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Version) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Version{}).Error; err != nil {
		return
	}
	Result = true
	return
}
