package models

import (
	"errors"
	"math"
)

const (
	MaxLimit     = 100
	DefaultLimit = 20
)

var (
	ErrInvalidLimit           = errors.New("max limit is 100")
	ErrFailedToGetFromContext = errors.New("failed to get pagination from context")
)

type Pagination struct {
	Page      uint64 `query:"page"`
	PageLimit uint64 `query:"page_limit"`
}

type FullPagination struct {
	Page      uint64 `json:"page"`
	PageLimit uint64 `json:"page_limit"`
	PageCount uint64 `json:"page_count"`
	Total     uint64 `json:"total"`
}

func NewPagination(page, pageLimit uint64) (Pagination, error) {
	if pageLimit > MaxLimit {
		return Pagination{}, ErrInvalidLimit
	}

	return Pagination{
		Page:      page,
		PageLimit: pageLimit,
	}, nil
}

func (p *Pagination) Offset() uint64 {
	return (p.Page - 1) * p.PageLimit
}

func (p *Pagination) Limit() uint64 {
	return p.PageLimit
}

func (p *Pagination) GetFull(total uint64) FullPagination {
	pageCount := uint64(math.Ceil(float64(total / p.PageLimit)))
	if pageCount == 0 {
		pageCount = 1
	}

	return FullPagination{
		Page:      p.Page,
		PageLimit: p.PageLimit,
		PageCount: pageCount,
		Total:     total,
	}
}

type DBPagination struct {
	Offset uint64
	Limit  uint64
}
