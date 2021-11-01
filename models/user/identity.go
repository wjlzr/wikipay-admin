package user

import (
	"bytes"
	"fmt"
	orm "wikipay-admin/database"
	"wikipay-admin/models/financial"
	"wikipay-admin/utils"

	"github.com/jinzhu/gorm"
)

// const identityPath = "http://18.162.243.214:81/identity/"

type Identity struct {
	Id             int    `xorm:"int(10)" json:"id"`
	UserId         int64  `xorm:"int(13)" json:"-"`
	CountryCode    string `xorm:"varchar(32)" json:"countryCode"`
	IdNumber       string `xorm:"varchar(128)" json:"idNumber"`
	FirstName      string `xorm:"varchar(60)" json:"firstName"`
	LastName       string `xorm:"varchar(60)" json:"lastName"`
	Phone          string `xorm:"varchar(100)" json:"phone"`
	ImgUrl         string `xorm:"varchar(256)" json:"imgUrl"`
	Content        string `xorm:"varchar(256)" json:"content"`
	Status         int    `xorm:"tinyint(1)" json:"status"`
	CreateAt       int64  `xorm:"bigint(13)" json:"createAt"`
	Email          string `xorm:"-" json:"email"`
	NickName       string `xorm:"-" json:"nickName"`
	Name           string `xorm:"-" json:"name"`
	OldCountryCode string `xorm:"-" json:"oldCountryCode"`
}

//
type IdentityReq struct {
	UserId  string `json:"userId" binding:"required"`
	Status  int    `json:"status" binding:"required"`
	Content string `json:"content"`
}

//
type IdentityStatusReq struct {
	financial.Pagination
	Id   int    `form:"id"`
	Info string `form:"info"`
}

//更新用户状态
func UpdateIdentity(req *IdentityReq) (string, error) {
	if req.Status < 2 || req.Status > 3 {
		return "状态不正确", nil
	}
	var info User
	//if err := orm.Eloquent.Table(User{}.TableName()).Where("id = ? AND auth_status = 1", req.UserId).First(&info); err == nil {
	orm.Eloquent.Raw(fmt.Sprintf("SELECT id FROM user WHERE id = %s", req.UserId)).Find(&info)
	if info.Id > 0 {
		return "", orm.Eloquent.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("UPDATE user SET auth_status = ? WHERE id = ?", req.Status, req.UserId).Error; err != nil {
				return err
			}

			if err := tx.Exec("UPDATE identity SET status = ?, content = ? WHERE user_id = ?",
				req.Status,
				req.Content,
				req.UserId).Error; err != nil {
				return err
			}
			return nil
		})
	}
	return "此用户未认证或已认证", nil
}

//查找认证信息
func FindIdentitys(req *IdentityStatusReq) ([]Identity, int, error) {
	//var where string
	var where bytes.Buffer

	//	where.WriteString(" WHERE a.country_code <> '156'" )
	if req.Id > 0 {
		where.WriteString(fmt.Sprintf(" AND a.id = %d", req.Id))
	} else {
		if req.Info != "" {
			where.WriteString(fmt.Sprintf(" AND concat(a.email,a.phone) LIKE '%s' ", "%"+req.Info+"%"))
		}
		where.WriteString(" ORDER BY b.create_at DESC")
		where.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d",
			req.PageSize,
			req.PageNum))

	}

	var (
		identitys []Identity
		count     int
	)

	//sql :=
	sql := `
		SELECT b.create_at,a.id,a.email,a.phone,a.country_code AS old_country_code, CONCAT(a.last_name,' ',a.first_name) AS name ,a.email,a.nick_name,a.phone,b.content,b.id_number,a.auth_status AS status,b.create_at,b.country_code,b.first_name,b.last_name,b.img_url
		FROM user a
		LEFT JOIN identity b ON a.id = b.user_id 
		WHERE a.country_code<>''
		`
	orm.Eloquent.Table(new(User).TableName()).Where("country_code<>''").Count(&count)
	//orm.Eloquent.Raw(fmt.Sprintf("%s", sql)).Find(&identitys)
	//count := len(identitys)

	sql = fmt.Sprintf("%s %s", sql, where.String())
	err := orm.Eloquent.Raw(sql).Find(&identitys).Error

	if err == nil {
		for i, v := range identitys {
			oldIdentity := utils.FindFlag(v.OldCountryCode)
			if oldIdentity != nil {
				identitys[i].OldCountryCode = oldIdentity.Name + " " + oldIdentity.Code
			}
			//	log.Log(identitys[i].OldCountryCode)
			newIdentity := utils.FindFlag(v.CountryCode)
			if newIdentity != nil {
				identitys[i].CountryCode = newIdentity.Name + " " + newIdentity.Code
			}
			if v.ImgUrl != "" {
				identitys[i].ImgUrl = utils.GetImage(utils.Identity, "") + v.ImgUrl
			}
			//	log.Log(identitys[i].CountryCode)
		}
	}
	return identitys, count, err
}
