package param

type UserEmotionSetParam struct {
	EmotionId string `json:"emotionId"`
}

type GetUserAvatarParam struct {
	UserId string `json:"userId"`
}
