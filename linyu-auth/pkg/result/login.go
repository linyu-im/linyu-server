package result

type UserLoginInfoResult struct {
	UserID  string `json:"userId"`
	Account string `json:"account"`
	Token   string `json:"token"`
}
