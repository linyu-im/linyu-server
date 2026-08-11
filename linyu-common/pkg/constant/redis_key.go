package constant

type redisKey struct {
	UserToken         string //用户登录token
	UserOnline        string //用户在线设备
	UserActiveSession string //用户当前激活的会话
	UserCode          string //用户验证码
	UserCodeLock      string //用户验证码锁（限制频繁发送）
	UploadChunk       string //上传的分片内容
	UploadMultipart   string //分片上传会话（uploadId/fileKey）
	UploadPartETag    string //分片 ETag
	GroupMaxNumber    string //当前最大群号
}

var RedisKey = redisKey{
	UserToken:         "user:%s:%s",                  //（user:用户id:登录设备）
	UserOnline:        "user:online:%s",              //（user:online:用户id，值为在线设备集合）
	UserActiveSession: "user:active-session:%s:%s",   //（user:active-session:用户id:设备）
	UserCode:          "user:code:%s",                //(user:code:手机号/邮箱)
	UserCodeLock:      "user:code:lock:%s",           //(user:code:lock:手机号/邮箱)
	UploadChunk:       "upload:chunk:%s",             //(upload:chunk:文件hash)
	UploadMultipart:   "upload:multipart:%s",         //(upload:multipart:文件hash)
	UploadPartETag:    "upload:multipart:%s:part:%s", //(upload:multipart:文件hash:part:分片序号)
	GroupMaxNumber:    "group:max-number",            //(group:max-number)
}
