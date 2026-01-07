package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/logger"
	"go.uber.org/zap"
	"sync"

	"github.com/sony/sonyflake/v2"
)

var (
	idGen *sonyflake.Sonyflake //雪花
	once  sync.Once
)

func init() {
	once.Do(func() {
		settings := sonyflake.Settings{
			MachineID: func() (int, error) {
				return 1, nil // 机器ID
			},
		}
		sf, err := sonyflake.New(settings)
		if err != nil {
			panic("init Sonyflake error:" + err.Error())
		}
		idGen = sf
	})
}

// GenerateSfID 生成唯一ID
func GenerateSfID() (int64, error) {
	id, err := idGen.NextID()
	if err != nil {
		logger.Log.Error("[GenerateSfID] 生成SonyflakeID失败:", zap.Error(err))
		return 0, err
	}
	return id, nil
}

// GenerateSfIDString 生成唯一ID字符串
func GenerateSfIDString() string {
	id, err := GenerateSfID()
	if err != nil {
		return uuid.New().String()
	}
	return fmt.Sprintf("%d", id)
}

func GenerateUuid() string {
	return uuid.New().String()
}
