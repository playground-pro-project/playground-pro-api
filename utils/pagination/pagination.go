package pagination

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort,omitempty" validate:"sort,omitempty"`
	TotalRows  int64  `json:"total_rows,omitempty"`
	TotalPages int    `json:"total_pages,omitempty"`
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 6
	}
	return p.Limit
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetPage() int {
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

func CalculateTotalPages(totalRows int64, limit int) int {
	return int(math.Ceil(float64(totalRows) / float64(limit)))
}

// Clousure function implementation
func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
