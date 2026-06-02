package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&EnterpriseMember{})
}

// EnterpriseMember 企业成员（企业通讯录）
// Roles 可存多个身份，逗号分隔，取值见 constant.EnterpriseMemberRole
type EnterpriseMember struct {
	ID           string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	EnterpriseID string              `gorm:"size:64;index;not null;comment:所属企业id" json:"enterpriseId"`
	DepartmentID string              `gorm:"size:64;index;default:'';comment:所属部门id" json:"departmentId"`
	UserID       string              `gorm:"size:64;index;not null;comment:用户id" json:"userId"`
	MemberName   string              `gorm:"size:128;comment:真实姓名" json:"memberName"`
	JobNumber    string              `gorm:"size:64;comment:工号" json:"jobNumber"`
	JobTitle     string              `gorm:"size:128;comment:职位" json:"jobTitle"`
	Email        string              `gorm:"size:128;comment:企业邮箱" json:"email"`
	Phone        string              `gorm:"size:32;comment:联系电话" json:"phone"`
	Roles        string              `gorm:"size:512;default:'member';comment:身份列表" json:"roles"`
	Status       string              `gorm:"size:32;default:'active';comment:成员状态" json:"status"`
	JoinedAt     localtime.LocalTime `gorm:"type:timestamp(3);comment:加入企业时间" json:"joinedAt"`
	CreatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	Username       string `gorm:"->;-:migration" json:"username"`
	UserLevel      int    `gorm:"->;-:migration" json:"userLevel"`
	EmotionName    string `gorm:"->;-:migration" json:"emotionName"`
	EmotionUrl     string `gorm:"->;-:migration" json:"emotionUrl"`
	DepartmentName string `gorm:"->;-:migration" json:"departmentName"`
}

func (EnterpriseMember) TableName() string {
	return "t_enterprise_member"
}

func (EnterpriseMember) TableComment() string {
	return "企业成员表"
}
