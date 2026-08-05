package dao

import (
	aiModel "github.com/linyu-im/linyu-server/linyu-ai/pkg/model"
	aiParam "github.com/linyu-im/linyu-server/linyu-ai/pkg/param"
	"gorm.io/gorm"
)

var SkillDao = newSkillDao()

func newSkillDao() *skillDao {
	return &skillDao{}
}

type skillDao struct{}

// List 查询技能列表，支持关键字/分类筛选
func (d *skillDao) List(db *gorm.DB, param *aiParam.SkillListParam) ([]*aiModel.Skill, error) {
	tx := db.Model(&aiModel.Skill{})
	if param.Keyword != "" {
		like := "%" + param.Keyword + "%"
		tx = tx.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if param.Category != "" {
		tx = tx.Where("category = ?", param.Category)
	}

	var list []*aiModel.Skill
	if err := tx.Order("featured DESC, created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
