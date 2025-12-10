package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Query struct {
	Page     int
	PageSize int
}

func Get(c *gin.Context) Query {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	s, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if p < 1 {
		p = 1
	}
	if s < 1 {
		s = 10
	}
	if s > 100 {
		s = 100
	}
	return Query{Page: p, PageSize: s}
}

func Scope(page, size int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

type Info struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func Do[T any](db *gorm.DB, q Query, out *[]T) (Info, error) {
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return Info{}, err
	}
	if err := db.Scopes(Scope(q.Page, q.PageSize)).Find(out).Error; err != nil {
		return Info{}, err
	}
	return Info{Page: q.Page, PageSize: q.PageSize, Total: total}, nil
}
