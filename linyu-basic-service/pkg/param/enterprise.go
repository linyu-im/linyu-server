package param

type EnterpriseInfoParam struct {
	EnterpriseID string `json:"enterpriseId" binding:"required"`
}

type GetEnterpriseAvatarParam struct {
	EnterpriseID string `json:"enterpriseId" binding:"required"`
}
