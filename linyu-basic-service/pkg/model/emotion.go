package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	db.AddMigrateTable(&Emotion{})

}

// Emotion 情绪表
type Emotion struct {
	ID          string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	EmotionName string              `gorm:"size:256;not null;comment:emo名称" json:"emotionName"`
	Url         string              `gorm:"size:512;not null;comment:emo地址" json:"url"`
	Type        string              `gorm:"size:64;comment:类型" json:"type"`
	CreatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt   localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Emotion) TableName() string {
	return "t_emotion"
}

func (Emotion) TableComment() string {
	return "情绪表"
}

func (Emotion) TableInit(db *gorm.DB) error {
	datas := []Emotion{
		{ID: "1", EmotionName: "自由万岁~", Url: "/emotion/freedom.png"},
		{ID: "2", EmotionName: "受伤的总是我", Url: "/emotion/injured.png"},
		{ID: "3", EmotionName: "闭麦谢谢", Url: "/emotion/shutup.png"},
		{ID: "4", EmotionName: "不知所措", Url: "/emotion/ataloss.png"},
		{ID: "5", EmotionName: "多云转晴", Url: "/emotion/cloudytofine.png"},
		{ID: "6", EmotionName: "晴转多云", Url: "/emotion/finetocloudy.png"},
		{ID: "7", EmotionName: "爱你呦", Url: "/emotion/loveyou.png"},
		{ID: "8", EmotionName: "暖暖的我", Url: "/emotion/warm.png"},
		{ID: "9", EmotionName: "闪亮登场", Url: "/emotion/debut.png"},
		{ID: "10", EmotionName: "人生如此多彩", Url: "/emotion/rainbow.png"},
		{ID: "11", EmotionName: "快乐拉满", Url: "/emotion/joyful.png"},
		{ID: "12", EmotionName: "人已被掏空", Url: "/emotion/empty.png"},
		{ID: "13", EmotionName: "沉默是金", Url: "/emotion/silence.png"},
		{ID: "14", EmotionName: "困了困了", Url: "/emotion/sleepy.png"},
		{ID: "15", EmotionName: "有点寄了", Url: "/emotion/beill.png"},
		{ID: "16", EmotionName: "纯小丑", Url: "/emotion/clown.png"},
		{ID: "17", EmotionName: "恋爱ing", Url: "/emotion/inlove.png"},
		{ID: "18", EmotionName: "失恋ing", Url: "/emotion/lovelorn.png"},
		{ID: "19", EmotionName: "俺是100昏", Url: "/emotion/100.png"},
		{ID: "20", EmotionName: "我超猛的", Url: "/emotion/strong.png"},
		{ID: "21", EmotionName: "我要悄悄变瘦", Url: "/emotion/thin.png"},
		{ID: "22", EmotionName: "出差ing", Url: "/emotion/onbusiness.png"},
		{ID: "23", EmotionName: "专注ing", Url: "/emotion/focus.png"},
		{ID: "24", EmotionName: "人生得意须尽欢", Url: "/emotion/complacent.png"},
		{ID: "25", EmotionName: "度假ing", Url: "/emotion/holiday.png"},
		{ID: "26", EmotionName: "晚安啦世界", Url: "/emotion/goodnight.png"},
		{ID: "27", EmotionName: "呱~", Url: "/emotion/croak.png"},
		{ID: "28", EmotionName: "秋意凉凉", Url: "/emotion/autumn.png"},
		{ID: "29", EmotionName: "我是国宝别动", Url: "/emotion/panda.png"},
		{ID: "30", EmotionName: "我好燃", Url: "/emotion/fire.png"},
		{ID: "31", EmotionName: "我只是一种草", Url: "/emotion/grass.png"},
	}
	for _, item := range datas {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		}).Create(&item).Error
		if err != nil {
			return err
		}
	}
	return nil
}
