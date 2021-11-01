package user

import (
	"fmt"
	_ "time"
	orm "wikipay-admin/database"
	"wikipay-admin/models"
	"wikipay-admin/tools"
)

//意见反馈
type Feedback struct {
	Id        int    `json:"id" gorm:"type:int(10) unsigned;primary_key"`        // 序列号
	UserId    int64  `json:"userId" gorm:"type:int(10);"`                        // 用户序列号
	Type      int    `json:"type" gorm:"type:tinyint(1) unsigned zerofill;"`     // 类型
	Suggest   string `json:"suggest" gorm:"type:varchar(500);"`                  // 建议
	Img1      string `json:"img1" gorm:"type:varchar(300);"`                     // l图片1
	Img2      string `json:"img2" gorm:"type:varchar(300);"`                     // l图片2
	Img3      string `json:"img3" gorm:"type:varchar(300);"`                     // l图片3
	Img4      string `json:"img4" gorm:"type:varchar(300);"`                     // l图片4
	Phone     string `json:"phone" gorm:"type:varchar(20);"`                     // 手机号
	CreateAt  int64  `json:"createAt" gorm:"type:bigint(13) unsigned zerofill;"` // 创建时间
	DataScope string `json:"dataScope" gorm:"-"`
	Params    string `json:"params"  gorm:"-"`
}

//
func (Feedback) TableName() string {
	return "feedback"
}

// 创建Feedback
func (e *Feedback) Create() (Feedback, error) {
	var doc Feedback
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取Feedback
func (e *Feedback) Get() (Feedback, error) {
	var doc Feedback

	table := orm.Eloquent.Table(e.TableName())
	if err := table.First(&doc, "id = ?", e.Id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取Feedback带分页
func (e *Feedback) GetPage(pageSize, pageIndex int, info string) ([]Feedback, int, error) {
	var doc []Feedback
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
		where = fmt.Sprintf("concat(user_id) like '%s%s%s'", "%", info, "%")
	}
	if err := table.Where(where).
		Order("create_at DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).Find(&doc).Error; err != nil {
		return nil, 0, err
	}

	var count int
	table.Count(&count)
	return doc, count, nil
}

// 更新Feedback
func (e *Feedback) Update(id int) (update Feedback, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id = ?", id).First(&update).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除Feedback
func (e *Feedback) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id = ?", id).Delete(&Feedback{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除
func (e *Feedback) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("Id in (?)", id).Delete(&Feedback{}).Error; err != nil {
		return
	}
	Result = true
	return
}
