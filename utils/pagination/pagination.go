package pagination

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit uint `json:"limit,omitempty"`
	Offset uint `json:"offset,omitempty"`
	Page uint `json:"page,omitempty"`
	Sort string `json:"sort,omitempty" validate:"sort,omitempty"`
	TotalRows uint `json:"total_rows,omitempty"`
	TotalPages uint `json:"total_pages,omitempty"`
}

func (p *Pagination) GetLimit() uint {
	if p.Limit == 0 {
		p.Limit = 5
	}
	return p.Limit
}

func (p *Pagination) GetOffset() uint {
	if p.Offset == 0 {
		p.Offset = 0
	}
	return p.Offset
}

func (p *Pagination) GetPage() uint {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "updated_at DESC"
	}
	return p.Sort
}

// Clousure function implementation
func Paginate(val interface{}, p *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(val).Count(&totalRows)

	p.TotalRows = uint(totalRows)
	totalPages := math.Ceil(float64(totalRows) / float64(p.GetLimit()))
	p.TotalPages = uint(totalPages)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(p.GetOffset())).Limit(int(p.GetLimit())).Order(p.GetSort())
	}
}
