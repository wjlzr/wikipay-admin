package user

import (
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"

	_ "time"
)

type FeedbackType struct {
	Code      int    `json:"code" gorm:"type:tinyint(1) unsigned zerofill;primary_key"` // 编码
	ZhCN      string `json:"zhCN" gorm:"column:zh-CN;type:varchar(100);"`               // 中文
	ZhHK      string `json:"zhHK" gorm:"column:zh-HK;type:varchar(100);"`               // 香港
	ZhTW      string `json:"zhTW" gorm:"column:zh-TW;type:varchar(100);"`               // 台湾
	En        string `json:"en" gorm:"type:varchar(200);"`                              // 英文
	Vi        string `json:"vi" gorm:"type:varchar(200);"`                              // 越南语
	Th        string `json:"th" gorm:"type:varchar(200);"`                              // 泰语
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

func (FeedbackType) TableName() string {
	return "feedback_type"
}

// 创建FeedbackType
func (e *FeedbackType) Create() (FeedbackType, error) {
	var doc FeedbackType

	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

//获取
func (e *FeedbackType) Get() (FeedbackType, error) {
	var doc FeedbackType

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "code = ?", e.Code).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取带分页
func (e *FeedbackType) GetPage(pageSize int, pageIndex int) ([]FeedbackType, int, error) {
	var doc []FeedbackType

	table := orm.Eloquent.Select("*").Table(e.TableName())
	// 数据权限控制(如果不需要数据权限请将此处去掉)
	dataPermission := new(models.DataPermission)
	dataPermission.UserId, _ = tools.StringToInt(e.DataScope)
	table, err := dataPermission.GetDataScope(e.TableName(), table)
	if err != nil {
		return nil, 0, err
	}
	var count int
	if err := table.Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&doc).
		Error; err != nil {
		return nil, 0, err
	}
	table.Count(&count)

	return doc, count, nil
}

// 更新
func (e *FeedbackType) Update(id int) (update FeedbackType, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("code = ?", id).First(&update).Error; err != nil {
		return
	}

	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除
func (e *FeedbackType) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("code = ?", id).Delete(&FeedbackType{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *FeedbackType) BatchDelete(id []int) (Result bool, err error) {

	if err = orm.Eloquent.Table(e.TableName()).Where("code in (?)", id).Delete(&FeedbackType{}).Error; err != nil {
		return
	}
	Result = true
	return
}
