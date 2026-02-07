package event

import "github.com/linyu-im/linyu-server/linyu-common/pkg/constant"

type WsData struct {
	SeqId   string      `json:"seqId"`
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type WsDataEvent struct {
	FromUserId string   `json:"fromUserId"`
	ToUserIds  []string `json:"ToUserIds"`
	Data       *WsData  `json:"data"`
}

func (m WsDataEvent) EventName() string {
	return constant.EventType.WsDataEvent
}

func (m WsDataEvent) GetSequenceId() string {
	return m.Data.SeqId
}
