package param

type FeedbackCreateParam struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}
