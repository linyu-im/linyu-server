package dao

import (
	basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/constant"
	"gorm.io/gorm"
)

var EnterpriseMemberDao = newEnterpriseMemberDao()

func newEnterpriseMemberDao() *enterpriseMemberDao {
	return &enterpriseMemberDao{}
}

type enterpriseMemberDao struct{}

func (d *enterpriseMemberDao) GetByEnterpriseIdAndUserId(db *gorm.DB, enterpriseId string, userId string) *basicModel.EnterpriseMember {
	result := &basicModel.EnterpriseMember{}
	if err := db.First(result, "enterprise_id = ? AND user_id = ? AND status = ?", enterpriseId, userId, constant.UserStatus.Active).Error; err != nil {
		return nil
	}
	return result
}

func (d *enterpriseMemberDao) ListByEnterpriseId(db *gorm.DB, enterpriseId string) ([]*basicModel.EnterpriseMember, error) {
	var list []*basicModel.EnterpriseMember
	err := d.BuildEnterpriseMemberQuery(db, enterpriseId).
		Where("em.status = ?", constant.UserStatus.Active).
		Order("em.created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (d *enterpriseMemberDao) BuildEnterpriseMemberQuery(db *gorm.DB, enterpriseId string) *gorm.DB {
	return db.Table("t_enterprise_member AS em").
		Select("em.*, u.username AS username, u.user_level AS user_level, e.emotion_name AS emotion_name, e.url AS emotion_url, d.name AS department_name").
		Joins("LEFT JOIN t_user u ON u.id = em.user_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN t_emotion e ON u.emotion_id = e.id").
		Joins("LEFT JOIN t_enterprise_department d ON d.id = em.department_id AND d.deleted_at IS NULL").
		Where("em.enterprise_id = ? AND em.deleted_at IS NULL", enterpriseId)
}

func (d *enterpriseMemberDao) ListByEnterpriseIdAndUserId(
	db *gorm.DB,
	enterpriseId string,
	userId string,
) ([]*basicModel.EnterpriseMember, error) {
	var list []*basicModel.EnterpriseMember
	err := db.Where("enterprise_id = ? AND user_id = ? AND status = ?", enterpriseId, userId, constant.UserStatus.Active).
		Order("joined_at ASC, created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
