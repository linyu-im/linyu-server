package result

type UserLoginInfoResult struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}
