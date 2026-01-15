package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.MysqlAddMigrateTable(&User{})
}

// User 用户账号表
type User struct {
	ID        string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	Username  string              `gorm:"size:255;not null;comment:用户名" json:"username"`
	Account   string              `gorm:"size:64;uniqueIndex;not null;comment:账号" json:"account"`
	Password  string              `gorm:"size:255;not null;comment:密码" json:"-"`
	Phone     string              `gorm:"size:11;uniqueIndex:uniq_phone_deleted_at;comment:手机号" json:"phone"`
	Email     string              `gorm:"size:255;uniqueIndex:uniq_email_deleted_at;comment:邮箱" json:"email"`
	Gender    string              `gorm:"size:10;comment:性别" json:"gender"`
	Birthday  localtime.LocalTime `gorm:"comment:生日" json:"birthday"`
	Avatar    string              `gorm:"size:512;comment:头像URL" json:"avatar"`
	Status    string              `gorm:"size:64;default:'active';comment:用户状态" json:"status"`
	CreatedAt localtime.LocalTime `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt localtime.LocalTime `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt      `gorm:"uniqueIndex:uniq_phone_deleted_at;uniqueIndex:uniq_email_deleted_at;index" json:"deletedAt"`
}

func (User) TableName() string {
	return "t_user"
}

func (User) TableComment() string {
	return "用户账号表"
}
