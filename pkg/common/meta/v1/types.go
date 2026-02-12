package v1

// ListMeta describes metadata that synthetic resources must have, including lists and
type ListMeta struct {
	Page int `json:"totalCount,omitempty"`

	PageSize int `json:"offset,omitempty" form:"offset"`
}

func (p ListMeta) GetLimit() int {
	if p.PageSize <= 0 || p.PageSize >= 50 {
		p.PageSize = 10
	}
	return p.PageSize
}
func (p ListMeta) GetPage() int {
	if p.Page <= 0 || p.Page >= 20 {
		return 1
	}
	return p.Page
}
func (p ListMeta) GetOffset() int {
	offset := (p.GetPage() - 1) * p.GetLimit()
	return offset
}
