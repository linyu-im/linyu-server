package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
)

func init() {
	db.AddMigrateTable(&SpaceFile{})
}

type SpaceFile struct {
	ID                  string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	SpaceID             string              `gorm:"size:64;index;not null;comment:空间id" json:"spaceId"`
	UserID              string              `gorm:"size:64;index;comment:创建人id" json:"userId"`
	PhysicalID          string              `gorm:"size:64;index;not null;comment:物理文件id" json:"physicalId"`
	PhysicalStoragePath string              `gorm:"type:text;not null;comment:物理的存储路径" json:"physicalStoragePath"`
	ParentID            string              `gorm:"size:64;index;comment:父目录id" json:"parentId"`
	Path                string              `gorm:"type:text;comment:完整路径 /id1/id2/id3" json:"path"`
	Level               int                 `gorm:"type:int;comment:目录层级" json:"level"`
	FileName            string              `gorm:"size:255;not null;comment:文件名称" json:"fileName"`
	IsDir               bool                `gorm:"default:0;comment:是否文件夹" json:"isDir"`
	FileType            string              `gorm:"size:64;comment:文件类型(后缀)" json:"fileType"`
	FileCategory        string              `gorm:"size:32;index;not null;default:other;comment:文件分类" json:"fileCategory"`
	FileSize            int64               `gorm:"type:bigint;comment:文件大小" json:"fileSize"`
	Status              string              `gorm:"size:32;comment:状态" json:"status"`
	CreatedAt           localtime.LocalTime `gorm:"type:timestamp(3);autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt           localtime.LocalTime `gorm:"type:timestamp(3);autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (SpaceFile) TableName() string {
	return "t_space_file"
}

func (SpaceFile) TableComment() string {
	return "空间（用户/组）文件表"
}
