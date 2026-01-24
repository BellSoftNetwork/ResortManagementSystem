package dto

type PaginationQuery struct {
	Page int    `form:"page,default=0" binding:"min=0"`
	Size int    `form:"size,default=20" binding:"min=1,max=2000"`
	Sort string `form:"sort"` // 정렬 파라미터 추가
}

func (p *PaginationQuery) GetOffset() int {
	// Spring Boot와 동일하게 0-based 페이징
	return p.Page * p.Size
}

func (p *PaginationQuery) GetLimit() int {
	return p.Size
}
