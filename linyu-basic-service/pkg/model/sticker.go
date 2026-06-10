package model

import (
	"github.com/linyu-im/linyu-server/linyu-common/pkg/db"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/localtime"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	db.AddMigrateTable(&Sticker{})
}

// Sticker 表情表
type Sticker struct {
	ID            string              `gorm:"size:64;primaryKey;autoIncrement:false;comment:id" json:"id"`
	Name          string              `gorm:"size:128;not null;comment:表情名称" json:"name"`
	IconUrl       string              `gorm:"size:512;comment:表情图标地址" json:"iconUrl"`
	Type          string              `gorm:"size:64;comment:类型" json:"type"`
	IconValue     string              `gorm:"size:512;comment:表情值" json:"iconValue"`
	StickerPackID string              `gorm:"size:64;index;not null;comment:表情分组id" json:"stickerPackId"`
	CreatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt     localtime.LocalTime `gorm:"type:timestamp(3);not null;autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}

func (Sticker) TableName() string {
	return "t_sticker"
}

func (Sticker) TableComment() string {
	return "表情表"
}

func (Sticker) TableInit(db *gorm.DB) error {
	datas := []Sticker{
		// 笑脸
		{ID: "1", Name: "露齿笑", Type: "unicode", IconValue: "😀", StickerPackID: "default"},
		{ID: "2", Name: "大笑", Type: "unicode", IconValue: "😄", StickerPackID: "default"},
		{ID: "3", Name: "眯眼笑", Type: "unicode", IconValue: "😆", StickerPackID: "default"},
		{ID: "4", Name: "笑哭", Type: "unicode", IconValue: "😂", StickerPackID: "default"},
		{ID: "5", Name: "笑到打滚", Type: "unicode", IconValue: "🤣", StickerPackID: "default"},
		{ID: "6", Name: "害羞笑", Type: "unicode", IconValue: "😊", StickerPackID: "default"},
		{ID: "7", Name: "微笑", Type: "unicode", IconValue: "🙂", StickerPackID: "default"},
		{ID: "8", Name: "眨眼", Type: "unicode", IconValue: "😉", StickerPackID: "default"},
		{ID: "9", Name: "吐舌", Type: "unicode", IconValue: "😛", StickerPackID: "default"},
		{ID: "10", Name: "调皮", Type: "unicode", IconValue: "😜", StickerPackID: "default"},
		// 喜爱
		{ID: "11", Name: "花痴", Type: "unicode", IconValue: "😍", StickerPackID: "default"},
		{ID: "12", Name: "爱心脸", Type: "unicode", IconValue: "🥰", StickerPackID: "default"},
		{ID: "13", Name: "飞吻", Type: "unicode", IconValue: "😘", StickerPackID: "default"},
		{ID: "14", Name: "亲亲", Type: "unicode", IconValue: "😗", StickerPackID: "default"},
		{ID: "15", Name: "抱抱", Type: "unicode", IconValue: "🤗", StickerPackID: "default"},
		// 表情
		{ID: "16", Name: "思考", Type: "unicode", IconValue: "🤔", StickerPackID: "default"},
		{ID: "17", Name: "闭嘴", Type: "unicode", IconValue: "🤐", StickerPackID: "default"},
		{ID: "18", Name: "无语", Type: "unicode", IconValue: "😐", StickerPackID: "default"},
		{ID: "19", Name: "翻白眼", Type: "unicode", IconValue: "🙄", StickerPackID: "default"},
		{ID: "20", Name: "得意", Type: "unicode", IconValue: "😏", StickerPackID: "default"},
		{ID: "21", Name: "困", Type: "unicode", IconValue: "😴", StickerPackID: "default"},
		{ID: "22", Name: "流口水", Type: "unicode", IconValue: "🤤", StickerPackID: "default"},
		{ID: "23", Name: "墨镜酷", Type: "unicode", IconValue: "😎", StickerPackID: "default"},
		{ID: "24", Name: "派对", Type: "unicode", IconValue: "🥳", StickerPackID: "default"},
		// 难过
		{ID: "25", Name: "委屈", Type: "unicode", IconValue: "🥺", StickerPackID: "default"},
		{ID: "26", Name: "流泪", Type: "unicode", IconValue: "😢", StickerPackID: "default"},
		{ID: "27", Name: "大哭", Type: "unicode", IconValue: "😭", StickerPackID: "default"},
		{ID: "28", Name: "失望", Type: "unicode", IconValue: "😞", StickerPackID: "default"},
		{ID: "29", Name: "疲惫", Type: "unicode", IconValue: "😩", StickerPackID: "default"},
		// 生气
		{ID: "30", Name: "生气", Type: "unicode", IconValue: "😠", StickerPackID: "default"},
		{ID: "31", Name: "怒火", Type: "unicode", IconValue: "😡", StickerPackID: "default"},
		{ID: "32", Name: "骂人", Type: "unicode", IconValue: "🤬", StickerPackID: "default"},
		// 惊讶
		{ID: "33", Name: "惊讶", Type: "unicode", IconValue: "😮", StickerPackID: "default"},
		{ID: "34", Name: "震惊", Type: "unicode", IconValue: "😲", StickerPackID: "default"},
		{ID: "35", Name: "害怕", Type: "unicode", IconValue: "😱", StickerPackID: "default"},
		{ID: "36", Name: "脸红", Type: "unicode", IconValue: "😳", StickerPackID: "default"},
		// 手势
		{ID: "37", Name: "挥手", Type: "unicode", IconValue: "👋", StickerPackID: "default"},
		{ID: "38", Name: "点赞", Type: "unicode", IconValue: "👍", StickerPackID: "default"},
		{ID: "39", Name: "踩", Type: "unicode", IconValue: "👎", StickerPackID: "default"},
		{ID: "40", Name: "OK", Type: "unicode", IconValue: "👌", StickerPackID: "default"},
		{ID: "41", Name: "胜利", Type: "unicode", IconValue: "✌️", StickerPackID: "default"},
		{ID: "42", Name: "鼓掌", Type: "unicode", IconValue: "👏", StickerPackID: "default"},
		{ID: "43", Name: "握手", Type: "unicode", IconValue: "🤝", StickerPackID: "default"},
		{ID: "44", Name: "祈祷", Type: "unicode", IconValue: "🙏", StickerPackID: "default"},
		{ID: "45", Name: "加油", Type: "unicode", IconValue: "💪", StickerPackID: "default"},
		// 爱心
		{ID: "46", Name: "红心", Type: "unicode", IconValue: "❤️", StickerPackID: "default"},
		{ID: "47", Name: "橙心", Type: "unicode", IconValue: "🧡", StickerPackID: "default"},
		{ID: "48", Name: "黄心", Type: "unicode", IconValue: "💛", StickerPackID: "default"},
		{ID: "49", Name: "绿心", Type: "unicode", IconValue: "💚", StickerPackID: "default"},
		{ID: "50", Name: "蓝心", Type: "unicode", IconValue: "💙", StickerPackID: "default"},
		{ID: "51", Name: "紫心", Type: "unicode", IconValue: "💜", StickerPackID: "default"},
		{ID: "52", Name: "心碎", Type: "unicode", IconValue: "💔", StickerPackID: "default"},
		{ID: "53", Name: "闪亮", Type: "unicode", IconValue: "✨", StickerPackID: "default"},
		{ID: "54", Name: "庆祝", Type: "unicode", IconValue: "🎉", StickerPackID: "default"},
		{ID: "55", Name: "礼花", Type: "unicode", IconValue: "🎊", StickerPackID: "default"},
		// 动物
		{ID: "56", Name: "狗", Type: "unicode", IconValue: "🐶", StickerPackID: "default"},
		{ID: "57", Name: "猫", Type: "unicode", IconValue: "🐱", StickerPackID: "default"},
		{ID: "58", Name: "熊猫", Type: "unicode", IconValue: "🐼", StickerPackID: "default"},
		{ID: "59", Name: "兔子", Type: "unicode", IconValue: "🐰", StickerPackID: "default"},
		{ID: "60", Name: "猪头", Type: "unicode", IconValue: "🐷", StickerPackID: "default"},
		// 食物
		{ID: "61", Name: "啤酒", Type: "unicode", IconValue: "🍺", StickerPackID: "default"},
		{ID: "62", Name: "咖啡", Type: "unicode", IconValue: "☕", StickerPackID: "default"},
		{ID: "63", Name: "蛋糕", Type: "unicode", IconValue: "🎂", StickerPackID: "default"},
		{ID: "64", Name: "披萨", Type: "unicode", IconValue: "🍕", StickerPackID: "default"},
		{ID: "65", Name: "汉堡", Type: "unicode", IconValue: "🍔", StickerPackID: "default"},
		// 天气
		{ID: "66", Name: "太阳", Type: "unicode", IconValue: "☀️", StickerPackID: "default"},
		{ID: "67", Name: "月亮", Type: "unicode", IconValue: "🌙", StickerPackID: "default"},
		{ID: "68", Name: "下雨", Type: "unicode", IconValue: "🌧️", StickerPackID: "default"},
		{ID: "69", Name: "雪花", Type: "unicode", IconValue: "❄️", StickerPackID: "default"},
		{ID: "70", Name: "彩虹", Type: "unicode", IconValue: "🌈", StickerPackID: "default"},
		// 其他
		{ID: "71", Name: "玫瑰", Type: "unicode", IconValue: "🌹", StickerPackID: "default"},
		{ID: "72", Name: "枯萎", Type: "unicode", IconValue: "🥀", StickerPackID: "default"},
		{ID: "73", Name: "礼物", Type: "unicode", IconValue: "🎁", StickerPackID: "default"},
		{ID: "74", Name: "炸弹", Type: "unicode", IconValue: "💣", StickerPackID: "default"},
		{ID: "75", Name: "便便", Type: "unicode", IconValue: "💩", StickerPackID: "default"},
		{ID: "76", Name: "骷髅", Type: "unicode", IconValue: "💀", StickerPackID: "default"},
		{ID: "77", Name: "幽灵", Type: "unicode", IconValue: "👻", StickerPackID: "default"},
		{ID: "78", Name: "外星人", Type: "unicode", IconValue: "👽", StickerPackID: "default"},
		{ID: "79", Name: "机器人", Type: "unicode", IconValue: "🤖", StickerPackID: "default"},
		{ID: "80", Name: "小丑", Type: "unicode", IconValue: "🤡", StickerPackID: "default"},
	}
	miyouRabbitDatas := []Sticker{
		{ID: "1001", Name: "阿姬-倒地了", Type: "image", IconUrl: "https://bbs-static.miyoushe.com/static/2023/05/22/a2fa805a0a782cdb8ab44ddaf8d5343d_2757305461654236072.png", StickerPackID: "1"},
		{ID: "1002", Name: "阿姬-倒地", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2022/12/08/a2fa805a0a782cdb8ab44ddaf8d5343d_7408497524870727211.png", StickerPackID: "1"},
		{ID: "1003", Name: "阿姬-得意", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/b171cc5b6c8d8b6257a3bf8163fb2680_567796524796359289.png", StickerPackID: "1"},
		{ID: "1004", Name: "阿姬-低落", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/9e8bcae8c207bef5c02c176cdb6916ee_7651949051600440018.png", StickerPackID: "1"},
		{ID: "1005", Name: "阿姬-调查", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/43338db72e2bb898d13a84c2fc302cb8_7770278814004091890.png", StickerPackID: "1"},
		{ID: "1006", Name: "阿姬-惊讶", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/829f0a1d533765b80769dd1a38dc338c_9182108590963360703.png", StickerPackID: "1"},
		{ID: "1007", Name: "阿姬-开心", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/4b4fb49d4bd402c530600e758ebb06bd_5477506439621383475.png", StickerPackID: "1"},
		{ID: "1008", Name: "阿姬-灵感", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/62158b26331a045bbaabf48d1fc5c8eb_980401343708652819.png", StickerPackID: "1"},
		{ID: "1009", Name: "阿姬-期待", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/955210989c73f207db2268b43b95f7d5_2639242042836509216.png", StickerPackID: "1"},
		{ID: "1010", Name: "阿姬-疑问", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/d3f9322fcc499df3fca778e8c366f34d_4195564687951765835.png", StickerPackID: "1"},
		{ID: "1011", Name: "吃雪糕", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/e380fdc2668f3dd0e8c8b9420a3cf124_3493794542940321040.png", StickerPackID: "1"},
		{ID: "1012", Name: "米游兔-加油", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/838dba51bd47c248d1f02b2e24a28b18_6640243918864203211.png", StickerPackID: "1"},
		{ID: "1013", Name: "米游兔—飚汗", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/74442ea449af598301cd16ed2841d3b0_5330135795233750081.png", StickerPackID: "1"},
		{ID: "1014", Name: "阿君-得意", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/2b2c2b9c68222997c21867d0d2bf8e09_3346073008037045533.png", StickerPackID: "1"},
		{ID: "1015", Name: "米游君-喝可乐", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/06f3f4c54fc51d82ba3e575194504d80_6167770388617778929.png", StickerPackID: "1"},
		{ID: "1016", Name: "米游君-认真工作", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/dec64f96a41d05d0f0c73432e7636e7b_3857353910783327632.png", StickerPackID: "1"},
		{ID: "1017", Name: "米游君-杂耍", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/c932711f2b7c36254e7afde91596f358_7200038206758881023.png", StickerPackID: "1"},
		{ID: "1018", Name: "绝区姬-嗨", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/2645aa282eee6e34ff5fbd546476df2b.png", StickerPackID: "1"},
		{ID: "1019", Name: "绝区姬-哇吼", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/6e1e0ce67ffad5de3b21cb50d37a1eac.png", StickerPackID: "1"},
		{ID: "1020", Name: "米游君-认真工作", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/dec64f96a41d05d0f0c73432e7636e7b_5424104392001050566.png", StickerPackID: "1"},
		{ID: "1021", Name: "绝区姬-自拍", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/d6279aaf178f56d095a1abae4a5ba3c6.png", StickerPackID: "1"},
		{ID: "1022", Name: "绝区姬-自信", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/45f524b90875b85a461bafa9fe9b483d.png", StickerPackID: "1"},
		{ID: "1023", Name: "米游姬-抱抱", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/7a514b631b93abb424cc01ade18020d5.png", StickerPackID: "1"},
		{ID: "1024", Name: "米游姬-好耶", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/6abecce847497d0cca150543bbf14709.png", StickerPackID: "1"},
		{ID: "1025", Name: "米游姬-抱米游兔", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/402fcfc2bd296f8189e6c85df3227f73.png", StickerPackID: "1"},
		{ID: "1026", Name: "米游姬-期待哦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/5f401916a1ef1c133b827949ff765692.png", StickerPackID: "1"},
		{ID: "1027", Name: "米游姬-抛心心", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/a2016a74dc5c31f7ee259e646415407c.png", StickerPackID: "1"},
		{ID: "1028", Name: "米游姬-献花", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/ca9a10a02c1007bf2626e65df2df7381.png", StickerPackID: "1"},
		{ID: "1029", Name: "米游姬-休息", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/c35e095d32d704ed85302d6fac6e8ca8.png", StickerPackID: "1"},
		{ID: "1030", Name: "米游兔-OK", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/db9c29ac19e22a5d9a60195803c3276b.png", StickerPackID: "1"},
		{ID: "1031", Name: "米游姬-卖萌", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/3831379d9e92f4050e640080e73be1bb.gif", StickerPackID: "1"},
		{ID: "1032", Name: "米游姬-感谢", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/02/03/a16450f92757f9d45588a8541fbbab09_4901912883842187271.gif", StickerPackID: "1"},
		{ID: "1033", Name: "米游姬-干杯", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/8a1c5e858456b720f2c3c45982729628_5411981834900996336.gif", StickerPackID: "1"},
		{ID: "1034", Name: "米游姬-开心", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/67ac33a760e1b8a4d1ed346bd69a9467_6056156571882767008.gif", StickerPackID: "1"},
		{ID: "1035", Name: "偶像姬-WINK", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/13042dc38b444c3dab24bce0b83d35d7.png", StickerPackID: "1"},
		{ID: "1036", Name: "偶像姬-应援", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/96934755c76f0abe4ebbdfb34950beaa.png", StickerPackID: "1"},
		{ID: "1037", Name: "偶像姬-达咩", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/f289ff3ea6828d1ab47b8e66fef20815.png", StickerPackID: "1"},
		{ID: "1038", Name: "偶像姬-糖葫芦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/23cbc98f6908928ed4b4495421310664.png", StickerPackID: "1"},
		{ID: "1039", Name: "偶像姬-闪亮登场", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/53d2ad52039025a469d2c5c39d85b2a0.png", StickerPackID: "1"},
		{ID: "1040", Name: "偶像姬-期待", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/5cb972102e1492c61aa2a404953956bd.png", StickerPackID: "1"},
		{ID: "1041", Name: "偶像姬-大声bb", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/59fe234835e808c41b0208c50be44914.png", StickerPackID: "1"},
		{ID: "1042", Name: "阿姬-开心", Type: "image", IconUrl: "https://upload-bbs.mihoyo.com/upload/2022/11/14/4b4fb49d4bd402c530600e758ebb06bd_6955490723450463851.png", StickerPackID: "1"},
		{ID: "1043", Name: "偶像姬-超凶", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/99dad3dfda99a78517532cd36a5c01c5.png", StickerPackID: "1"},
		{ID: "1044", Name: "偶像姬-沮丧", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/6bdba3c03727d9e42ff5432d3e422751.png", StickerPackID: "1"},
		{ID: "1045", Name: "偶像姬-沮丧", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/045a3b4140d74bd4c091f7e557d26213.png", StickerPackID: "1"},
		{ID: "1046", Name: "偶像姬-鸽子", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/850ed1c7790db7beffd44156ddb3173e.png", StickerPackID: "1"},
		{ID: "1047", Name: "偶像姬-咕咕", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/059e5d62858227e9e9660ca669187cec.png", StickerPackID: "1"},
		{ID: "1048", Name: "偶像姬-撒花", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/7085182d437642634c95e6fecee9ad01.png", StickerPackID: "1"},
		{ID: "1049", Name: "米游姬-哼", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/047350694eb6b3b68fbb02d31c1a91d1.png", StickerPackID: "1"},
		{ID: "1050", Name: "米游姬-乖巧", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/2f2e3096743e864e75fa2ff50d36fa50.png", StickerPackID: "1"},
		{ID: "1051", Name: "米游姬-吃瓜", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/613a2b262af0319edde21587b88a9c6e.png", StickerPackID: "1"},
		{ID: "1052", Name: "米游姬-呆滞", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/03c78415a3318d527d307c52a188a5ad.png", StickerPackID: "1"},
		{ID: "1053", Name: "米游姬-疑问", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/cfc6b254655ee9b2a08202020d898f87.png", StickerPackID: "1"},
		{ID: "1054", Name: "米游姬-睡觉", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/dac9d3391a4a0b5efc0c0acd589d3b30.png", StickerPackID: "1"},
		{ID: "1055", Name: "米游姬-暗中观察", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/d4d609d5af6bc85a0c3bba6b4b59fffa.png", StickerPackID: "1"},
		{ID: "1056", Name: "米游姬-惊", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/9dec3e01e834018c63b1f6710d4f1b8e.png", StickerPackID: "1"},
		{ID: "1057", Name: "米游姬-点赞", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/d6ec212722aa079647bb233a3853f775.png", StickerPackID: "1"},
		{ID: "1058", Name: "米游姬-求求你啦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/1642da910ce6a64b72eb436a008d3af8.png", StickerPackID: "1"},
		{ID: "1059", Name: "米游姬-期待", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/6adaac5ed9b16311259d3bbb6c108125.png", StickerPackID: "1"},
		{ID: "1060", Name: "米游姬-打call", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/f9c9a65998f81afff034bfbe20087f4a.png", StickerPackID: "1"},
		{ID: "1061", Name: "米游姬-撒花", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/19cf3986f2c2e034dd8704641a430cf9.png", StickerPackID: "1"},
		{ID: "1062", Name: "米游姬-来了来了", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/e31ca4bc3b36d83ba0095505d4e46972.png", StickerPackID: "1"},
		{ID: "1063", Name: "米游姬-周五啦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/fce537cc087bab80640209c5a2b5c59f.png", StickerPackID: "1"},
		{ID: "1064", Name: "米游姬-累趴", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/968635e0f4b0b0fbb6cc3e19aa64ffb3.png", StickerPackID: "1"},
		{ID: "1065", Name: "星穹米游姬-得意", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/6520fc446c0d2daa7aa0bb0d4947392f_3450328440953691483.png", StickerPackID: "1"},
		{ID: "1066", Name: "星穹米游姬-哭唧唧", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/9b6cab9f6fe05d7d99888525338b8930_4389753021835428841.png", StickerPackID: "1"},
		{ID: "1067", Name: "星穹米游姬-你好", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/ca2646288114295698f89b9431db206f_45031739628247345.png", StickerPackID: "1"},
		{ID: "1068", Name: "星穹米游姬-耶！", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/a0ca769e6f3bfc4291e36287fc9d38ef_8994291476410828569.png", StickerPackID: "1"},
		{ID: "1069", Name: "吃糖葫芦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/ffb9724864fd36eff5df36a3e7145075.png", StickerPackID: "1"},
		{ID: "1070", Name: "春节咕咕", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/1f7ac4a1432f8cd2802ad2adb9162199.png", StickerPackID: "1"},
		{ID: "1071", Name: "放鞭炮", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/6e68fe14b046b4e0dcadb358f59b57c5.png", StickerPackID: "1"},
		{ID: "1072", Name: "福到啦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/97f16d18549d2649232931a12e4553fb.png", StickerPackID: "1"},
		{ID: "1073", Name: "恭贺新春", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/ffa54fd8e3cbf92405b8053462fdfe86.png", StickerPackID: "1"},
		{ID: "1074", Name: "米游姬-发福利", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/4a4990c1e1294e40e5078b4ff5518cc4.png", StickerPackID: "1"},
		{ID: "1075", Name: "米游姬-福袋", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/906db216f9a2c3bfba0f3c0705cfcc7e.png", StickerPackID: "1"},
		{ID: "1076", Name: "米游姬-哈欠", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/9a0d2f4bffa1913f52e19bd738beb41a.png", StickerPackID: "1"},
		{ID: "1077", Name: "米游姬-抢红包", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/1bee21e3f57024d72abd4aabeb411315.png", StickerPackID: "1"},
		{ID: "1078", Name: "米游姬-入欧", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/366803aa8452bd1a2b2b57b556c7a504.png", StickerPackID: "1"},
		{ID: "1079", Name: "米游姬-唢呐", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/3e54210658c07342924cc536caa22c99.png", StickerPackID: "1"},
		{ID: "1080", Name: "米游兔-好耶", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/422b384f970bf369dc0083d2fa17ede9.png", StickerPackID: "1"},
		{ID: "1081", Name: "米游兔-脱非", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/a7fa138da244faa48e57e81c656639be.png", StickerPackID: "1"},
		{ID: "1082", Name: "牛气冲天", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/02390a2b74d3a9c70b4ff78cdb9d608a.png", StickerPackID: "1"},
		{ID: "1083", Name: "受到惊吓", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/0b568ef852feeb12f72c1b7b6123461f.png", StickerPackID: "1"},
		{ID: "1084", Name: "我的心好痛", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/8e53cb4b90a9bbb67db92a9ca95cfd19.png", StickerPackID: "1"},
		{ID: "1085", Name: "魔游姬-哼哼", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/0fb1f674bffda4fa9ab8b078171c53d9.png", StickerPackID: "1"},
		{ID: "1086", Name: "魔游姬-得意", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/5a97d54c56fcb87e041c29eb0e6a5f27.png", StickerPackID: "1"},
		{ID: "1087", Name: "魔游姬-害怕", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/e907f9fe3609227517f70697835e0cae.png", StickerPackID: "1"},
		{ID: "1088", Name: "魔游姬-鸽子", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/03564ea172a338ac1e4e977792d31f3f.png", StickerPackID: "1"},
		{ID: "1089", Name: "魔游姬-生气", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/45ca6788363e15288193244bb5493a9b.png", StickerPackID: "1"},
		{ID: "1090", Name: "魔游姬-糖葫芦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/ebd96478a64bb7ac193806be15913198.png", StickerPackID: "1"},
		{ID: "1091", Name: "米游姬-比心", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/4d41193fa02ac31e2daf7436ba5e12cf.png", StickerPackID: "1"},
		{ID: "1092", Name: "米游姬糖-葫芦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/6744fdc365a234f152bc3fcc6443c7a9.png", StickerPackID: "1"},
		{ID: "1093", Name: "米游姬-咕咕", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/8568e2f26acbd1b86ae03dbf4b62f1b7.png", StickerPackID: "1"},
		{ID: "1094", Name: "米游姬-疑惑", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/611e5b672d31a4fd203b46253c1f2348.png", StickerPackID: "1"},
		{ID: "1095", Name: "米游姬-得意", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/068a8b611fae7cde48470ebfe5e21847.png", StickerPackID: "1"},
		{ID: "1096", Name: "米游姬-观察", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/251a9bb6d4a4f1a39a430159919ed5fc.png", StickerPackID: "1"},
		{ID: "1097", Name: "米游姬-抱抱", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/b37d75bce678aaed49a93a911299aa39.png", StickerPackID: "1"},
		{ID: "1098", Name: "米游姬-大哭", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/521b9f657d4265d2299641aa6552e7f3.png", StickerPackID: "1"},
		{ID: "1099", Name: "米游姬-喝茶", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/01ecd5b3b0f65c59ec46e2dc2538ba79.png", StickerPackID: "1"},
		{ID: "1100", Name: "米游姬流-鼻血", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/8bcf809e4bdc770476544b1276acece1.png", StickerPackID: "1"},
		{ID: "1101", Name: "米游姬-委屈", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/dc9830d58892a88244364ce51394011c.png", StickerPackID: "1"},
		{ID: "1102", Name: "米游兔-加油", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/5857b8a3d4023bd05954225b0d578845_8473504187038159665.png", StickerPackID: "1"},
		{ID: "1103", Name: "米游兔-无奈", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/3597ff62b268362f91576dda93e6d58b_2427201235588485119.png", StickerPackID: "1"},
		{ID: "1104", Name: "米游兔-烟雾弹", Type: "image", IconUrl: "https://upload-bbs.miyoushe.com/upload/2023/01/18/3d0feb3d67c59658d1a084a8eabd00af_2634705373803711026.png", StickerPackID: "1"},
		{ID: "1105", Name: "米游兔-笔心", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/2416bcd2bd669db70ea8e8db923a81eb.png", StickerPackID: "1"},
		{ID: "1106", Name: "米游兔-糖葫芦", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/befed8e6dd7a9d00918dd257cdd1a71d.png", StickerPackID: "1"},
		{ID: "1107", Name: "米游兔-害怕", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/a03798f27603dd3d94eb83786f75b71d.png", StickerPackID: "1"},
		{ID: "1108", Name: "米游兔-暗中观察", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/c81b4153b1b7e687aaea76fc1e4c5dc2.png", StickerPackID: "1"},
		{ID: "1109", Name: "米游兔-可怜兮兮", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/d57ac9b36df80b48b1acb5769ac349ac.png", StickerPackID: "1"},
		{ID: "1110", Name: "米游兔-星星", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/43e6636cada23e5411a4e82546ee9f22.png", StickerPackID: "1"},
		{ID: "1111", Name: "米游兔-递茶", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/ba72f75be849cf28e55c9dce548967f0.png", StickerPackID: "1"},
		{ID: "1112", Name: "米游兔-吃惊", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/378cdb8ccf15a49f0e7b3aba18acf2ff.png", StickerPackID: "1"},
		{ID: "1113", Name: "米游兔-乖巧", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/d93728df6c1ea1fb2a0ae381824a970d.png", StickerPackID: "1"},
		{ID: "1114", Name: "米游兔-自闭", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/77645095c66d46eb001231934c28be1a.png", StickerPackID: "1"},
		{ID: "1115", Name: "米游兔-安详", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/fd58c953ab1b2b886cc524ff1f636aa8.png", StickerPackID: "1"},
		{ID: "1116", Name: "米游兔-吨吨", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/2a935e275eaa9b7d7f50523f56adadb6.png", StickerPackID: "1"},
		{ID: "1117", Name: "米游兔-举牌子", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/577a81a5ede754faa4d92882c5525fff.png", StickerPackID: "1"},
		{ID: "1118", Name: "米游兔-吃瓜", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/5ea22ea42dc918b8263006ef3c2fd7e0.png", StickerPackID: "1"},
		{ID: "1119", Name: "米游兔-飞翔", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/cb81bf30de24cf79829961cca5b11428.png", StickerPackID: "1"},
		{ID: "1120", Name: "米游兔-期待", Type: "image", IconUrl: "https://img-static.mihoyo.com/communityweb/upload/8ae8647b41a08ffec4e5eb2239fc30c1.png", StickerPackID: "1"},
	}
	for _, stickers := range [][]Sticker{datas, miyouRabbitDatas} {
		for _, item := range stickers {
			err := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoNothing: true,
			}).Create(&item).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
