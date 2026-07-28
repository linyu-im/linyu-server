package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&Space{})
}

// Space 空间容量表
type Space struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	SpaceType   string              `gorm:"size:32;not null;uniqueIndex:uk_space_target,priority:1;comment:空间类型user/group/org" json:"spaceType"`
	TargetID    string              `gorm:"size:64;not null;uniqueIndex:uk_space_target,priority:2;comment:关联目标id，用户id/群id/组织id" json:"targetId"`
	OwnerID     string              `gorm:"size:64;index;not null;comment:空间所有者或负责人id" json:"ownerId"`
	SpaceName   string              `gorm:"size:256;not null;comment:空间名称" json:"spaceName"`
	QuotaBytes  int64               `gorm:"type:bigint;not null;default:0;comment:空间总容量，0表示不限制" json:"quotaBytes"`
	UsedBytes   int64               `gorm:"type:bigint;not null;default:0;comment:已使用容量" json:"usedBytes"`
	FileCount   int64               `gorm:"type:bigint;not null;default:0;comment:文件数量" json:"fileCount"`
	FolderCount int64               `gorm:"type:bigint;not null;default:0;comment:文件夹数量" json:"folderCount"`
	Status      string              `gorm:"size:32;index;not null;default:active;comment:状态 active/disabled/readonly" json:"status"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   localtime.LocalTime `gorm:"type:timestamp(3);autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Space) TableName() string {
	return "t_space"
}

func (Space) TableComment() string {
	return "空间容量表"
}
