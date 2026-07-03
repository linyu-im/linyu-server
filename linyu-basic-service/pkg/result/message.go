package result

type ForwardError struct {
	SessionId string `json:"sessionId"`
	SceneType string `json:"sceneType"`
	ToId      string `json:"toId"`
	ErrMsg    string `json:"errMsg"`
}
