package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/config"
)

var (
	node *snowflake.Node
	once sync.Once
)

func InitSnowflake() {
	once.Do(func() {
		var err error
		var finalID int64
		if config.C.Server.NodeId > 0 {
			finalID = config.C.Server.NodeId
		} else {
			finalID = autoDeriveNodeID()
		}
		node, err = snowflake.NewNode(finalID)
		if err != nil {
			panic("init snowflake node error: " + err.Error())
		}
	})
}

func autoDeriveNodeID() int64 {
	hostname, err := os.Hostname()
	if err != nil {
		return 1
	}
	re := regexp.MustCompile(`(\d+)$`)
	match := re.FindStringSubmatch(hostname)
	if len(match) > 1 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err == nil {
			return id
		}
	}
	return int64(sumString(hostname) % 1024)
}

func sumString(s string) int {
	v := 0
	for _, r := range s {
		v += int(r)
	}
	return v
}

// GenerateSfID 生成唯一ID
func GenerateSfID() int64 {
	return node.Generate().Int64()
}

// GenerateSfIDString 生成唯一ID字符串
func GenerateSfIDString() string {
	return node.Generate().String()
}

func GenerateUuid() string {
	return uuid.NewString()
}

func GenerateOnlyNumber(accountPrefix string, checkFun func(account string) bool) string {
	const maxRetry = 5
	var account string
	for i := 0; i < maxRetry; i++ {
		account = GenerateAccount(accountPrefix)
		if checkFun(account) {
			return account
		}
		if i == maxRetry-1 {
			account = accountPrefix + GenerateSfIDString()
		}
	}
	return account
}

var groupNumberMu sync.Mutex

func GenerateGroupNumber(maxNumber string, checkFun func(number string) bool) string {
	groupNumberMu.Lock()
	defer groupNumberMu.Unlock()

	const (
		minVal   = 10000
		maxVal   = 99999999999
		fallback = 1000000000000
	)

	var current int64
	if maxNumber == "" {
		current = minVal
	} else {
		n, err := strconv.ParseInt(maxNumber, 10, 64)
		if err != nil || n < minVal {
			current = minVal
		} else if n >= maxVal {
			current = maxVal
		} else {
			current = n + 1
		}
	}

	for i := 0; i < 10; i++ {
		if current > maxVal {
			break
		}
		number := strconv.FormatInt(current, 10)
		if checkFun(number) {
			return number
		}
		current++
	}
	return fmt.Sprintf("%012d", GenerateSfID()%fallback)
}

func Generate1v1SessionID(id1, id2 string) string {
	if id1 < id2 {
		return fmt.Sprintf("%s_%s", id1, id2)
	}
	return fmt.Sprintf("%s_%s", id2, id1)
}

func Split1v1SessionID(sessionID string) []string {
	return strings.Split(sessionID, "_")
}
