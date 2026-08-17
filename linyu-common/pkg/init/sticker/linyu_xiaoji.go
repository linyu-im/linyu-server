package sticker

import (
	"strings"

	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/storage"
)

func LinyuXiaoji() []StickerData {
	items := []struct {
		ID   string
		Name string
	}{
		{"1001", "OK"},
		{"1002", "安全守护"},
		{"1003", "傲娇"},
		{"1004", "拜托"},
		{"1005", "包我身上"},
		{"1006", "饱了"},
		{"1007", "抱歉"},
		{"1008", "暴怒"},
		{"1009", "崩溃"},
		{"1010", "比心"},
		{"1011", "不方便"},
		{"1012", "吃瓜"},
		{"1013", "打招呼"},
		{"1014", "大哭"},
		{"1015", "大笑"},
		{"1016", "到了"},
		{"1017", "等回复"},
		{"1018", "等我忙完"},
		{"1019", "点赞"},
		{"1020", "电话我"},
		{"1021", "调皮"},
		{"1022", "饿了"},
		{"1023", "发呆"},
		{"1024", "发送链接"},
		{"1025", "发送图片"},
		{"1026", "发送文件"},
		{"1027", "发送消息"},
		{"1028", "尴尬"},
		{"1029", "感动"},
		{"1030", "干饭"},
		{"1031", "鼓掌"},
		{"1032", "鬼脸"},
		{"1033", "害羞"},
		{"1034", "好吃"},
		{"1035", "喝水"},
		{"1036", "加油"},
		{"1037", "交给我"},
		{"1038", "惊讶"},
		{"1039", "拒绝"},
		{"1040", "咖啡"},
		{"1041", "开工"},
		{"1042", "开会中"},
		{"1043", "开心"},
		{"1044", "可以"},
		{"1045", "哭泣"},
		{"1046", "困困"},
		{"1047", "累瘫了"},
		{"1048", "冷"},
		{"1049", "离线"},
		{"1050", "裂开"},
		{"1051", "流汗"},
		{"1052", "麻了"},
		{"1053", "马上来"},
		{"1054", "满血复活"},
		{"1055", "忙碌"},
		{"1056", "没问题"},
		{"1057", "免打扰"},
		{"1058", "摸鱼"},
		{"1059", "奶茶"},
		{"1060", "陪伴"},
		{"1061", "佩服"},
		{"1062", "破防"},
		{"1063", "期待"},
		{"1064", "清空消息"},
		{"1065", "庆祝"},
		{"1066", "热"},
		{"1067", "撒花"},
		{"1068", "稍等"},
		{"1069", "生病"},
		{"1070", "生气"},
		{"1071", "失败了"},
		{"1072", "视频通话"},
		{"1073", "收到"},
		{"1074", "收到消息"},
		{"1075", "睡觉"},
		{"1076", "思考"},
		{"1077", "太强了"},
		{"1078", "偷看"},
		{"1079", "偷笑"},
		{"1080", "吐舌"},
		{"1081", "晚安"},
		{"1082", "委屈"},
		{"1083", "未读消息"},
		{"1084", "稳了"},
		{"1085", "无语"},
		{"1086", "午安"},
		{"1087", "下班啦"},
		{"1088", "消息提醒"},
		{"1089", "谢谢"},
		{"1090", "辛苦了"},
		{"1091", "休息一下"},
		{"1092", "疑问"},
		{"1093", "已读不回"},
		{"1094", "语音消息"},
		{"1095", "晕了"},
		{"1096", "在线"},
		{"1097", "早安"},
		{"1098", "震惊"},
		{"1099", "置顶"},
		{"1100", "走你"},
	}

	var iconPrefix string
	switch config.C.Storage.Type {
	case config.LocalStorageType:
		// http://{local.base-url}/api/basic/v1/local/storage/sticker/xiaoji/256x256/{name}.png
		iconPrefix = strings.TrimRight(config.C.Storage.LocalStorage.BaseURL, "/") +
			config.C.Server.RoutePrefix + storage.LocalStorageUrl + "/sticker/xiaoji/256x256/"
	case config.S3StorageType:
		// http://{s3.endpoint}/{bucket}/sticker/xiaoji/256x256/{name}.png
		s3 := config.C.Storage.S3Storage
		if s3.BaseURL != "" {
			iconPrefix = strings.TrimRight(s3.BaseURL, "/") + "/sticker/xiaoji/256x256/"
		} else {
			iconPrefix = strings.TrimRight(s3.Endpoint, "/") + "/" + s3.BucketName + "/sticker/xiaoji/256x256/"
		}
	default:
		iconPrefix = "/sticker/xiaoji/256x256/"
	}

	list := make([]StickerData, 0, len(items))
	for _, item := range items {
		list = append(list, StickerData{
			ID:            item.ID,
			Name:          item.Name,
			Type:          "image",
			IconUrl:       iconPrefix + item.Name + ".png",
			StickerPackID: "1",
		})
	}
	return list
}
