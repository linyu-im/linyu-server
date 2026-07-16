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
		{ID: "1", EmotionName: "爱你呦", Url: "/emotion/loveyou.png"},
		{ID: "2", EmotionName: "俺是100分", Url: "/emotion/score100.png"},
		{ID: "3", EmotionName: "熬夜冠军", Url: "/emotion/nightowl.png"},
		{ID: "4", EmotionName: "出差ing", Url: "/emotion/trip.png"},
		{ID: "5", EmotionName: "出去浪啦", Url: "/emotion/play.png"},
		{ID: "6", EmotionName: "度假ing", Url: "/emotion/holiday.png"},
		{ID: "7", EmotionName: "多云转晴", Url: "/emotion/cloudy.png"},
		{ID: "8", EmotionName: "乏了乏了", Url: "/emotion/tired.png"},
		{ID: "9", EmotionName: "呱~", Url: "/emotion/frog.png"},
		{ID: "10", EmotionName: "咖啡续命", Url: "/emotion/coffee.png"},
		{ID: "11", EmotionName: "开饭开饭", Url: "/emotion/meal.png"},
		{ID: "12", EmotionName: "快乐拉满", Url: "/emotion/happy.png"},
		{ID: "13", EmotionName: "恋爱ing", Url: "/emotion/inlove.png"},
		{ID: "14", EmotionName: "忙成陀螺", Url: "/emotion/busy.png"},
		{ID: "15", EmotionName: "摸鱼ing", Url: "/emotion/chill.png"},
		{ID: "16", EmotionName: "脑袋空空", Url: "/emotion/emptyhead.png"},
		{ID: "17", EmotionName: "破防ing", Url: "/emotion/broken.png"},
		{ID: "18", EmotionName: "秋意凉凉", Url: "/emotion/autumn.png"},
		{ID: "19", EmotionName: "求喝奶茶", Url: "/emotion/milktea.png"},
		{ID: "20", EmotionName: "燃起来了", Url: "/emotion/fire.png"},
		{ID: "21", EmotionName: "人生得意须尽欢", Url: "/emotion/cheers.png"},
		{ID: "22", EmotionName: "闪亮登场", Url: "/emotion/debut.png"},
		{ID: "23", EmotionName: "生活如此多彩", Url: "/emotion/colorful.png"},
		{ID: "24", EmotionName: "躺平躺平", Url: "/emotion/lieflat.png"},
		{ID: "25", EmotionName: "我超猛", Url: "/emotion/strong.png"},
		{ID: "26", EmotionName: "我是国宝你别动", Url: "/emotion/panda.png"},
		{ID: "27", EmotionName: "我有点G了", Url: "/emotion/gg.png"},
		{ID: "28", EmotionName: "我只是一种草", Url: "/emotion/grass.png"},
		{ID: "29", EmotionName: "一脸问号", Url: "/emotion/confused.png"},
		{ID: "30", EmotionName: "元气满满", Url: "/emotion/energy.png"},
		{ID: "31", EmotionName: "宅家ing", Url: "/emotion/home.png"},
		{ID: "32", EmotionName: "真的会谢", Url: "/emotion/thanks.png"},
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
