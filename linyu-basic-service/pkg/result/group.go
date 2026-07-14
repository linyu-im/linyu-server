package result

import basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"

type GroupInfoResult struct {
	Info          *basicModel.Group         `json:"info"`
	CurrentMember *basicModel.GroupMember   `json:"currentMember"`
	Tops          []*basicModel.GroupMember `json:"tops"`
}
