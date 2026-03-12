package request

type PageQuery struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}

func (p *PageQuery) SetDefault() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		p.SortOrder = "asc"
	}
}
