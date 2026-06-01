package result

import basicModel "github.com/linyu-im/linyu-server/linyu-basic-service/pkg/model"

type GroupInfoResult struct {
	Info *basicModel.Group         `json:"info"`
	Tops []*basicModel.GroupMember `json:"tops"`
}
