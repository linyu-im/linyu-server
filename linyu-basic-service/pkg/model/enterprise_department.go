package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&EnterpriseDepartment{})
}

// EnterpriseDepartment 企业部门（企业组织结构）
type EnterpriseDepartment struct {
	ID           string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:部门id" json:"id"`
	EnterpriseID string              `gorm:"size:64;index;not null;comment:所属企业id" json:"enterpriseId"`
	ParentID     string              `gorm:"size:64;index;default:'';comment:父部门id(为空表示顶级部门)" json:"parentId"`
	Name         string              `gorm:"size:128;not null;comment:部门名称" json:"name"`
	Describe     string              `gorm:"type:text;comment:部门描述" json:"describe"`
	LeaderUserID string              `gorm:"size:64;comment:部门负责人用户id" json:"leaderUserId"`
	Sort         int                 `gorm:"default:0;comment:排序值(越小越靠前)" json:"sort"`
	Level        int                 `gorm:"default:1;comment:部门层级(顶级为1)" json:"level"`
	MemberNum    int                 `gorm:"default:0;comment:部门成员数" json:"memberNum"`
	CreatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt    localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deletedAt"`

	LeaderUsername string                  `gorm:"->;-:migration" json:"leaderUsername"`
	Children       []*EnterpriseDepartment `gorm:"-" json:"children,omitempty"`
	Members        []*EnterpriseMember     `gorm:"-" json:"members,omitempty"`
}

func (EnterpriseDepartment) TableName() string {
	return "t_enterprise_department"
}

func (EnterpriseDepartment) TableComment() string {
	return "企业部门表"
}
