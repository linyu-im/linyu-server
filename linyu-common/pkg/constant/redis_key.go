package constant

type redisKey struct {
	UserToken    string //用户登录token
	UserCode     string //用户验证码
	UserCodeLock string //用户验证码锁（限制频繁发送）
	UploadChunk  string //上传的分片内容
}

var RedisKey = redisKey{
	UserToken:    "user:%s:%s",        //（user:用户id:登录设备）
	UserCode:     "user:code:%s",      //(user:code:手机号/邮箱)
	UserCodeLock: "user:code:lock:%s", //(user:code:lock:手机号/邮箱)
	UploadChunk:  "upload:chunk:%s",   //(upload:chunk:文件hash)
}
