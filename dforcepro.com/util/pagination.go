package util

import (
	"math"

	mgo "gopkg.in/mgo.v2"
)

type PaginationRes struct {
	Raws     interface{} `json:"raws,omitempty"`
	Total    int         `json:"total,omitempty"`
	AllPages int         `json:"allPages,omitempty"`
	Page     int         `json:"page,omitempty"`
	Limit    int         `json:"limit,omitempty"`
}

const (
	MaxLimit = 300
)

func MongoPagination(query *mgo.Query, limit int, page int) (*PaginationRes, error) {
	total, err := query.Count()

	if err != nil {
		return nil, err
	}
	if limit < 1 || limit > 100 {
		limit = 100
	}
	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	if page > totalPage {
		page = totalPage
	} else if page < 1 {
		page = 1
	}
	return &PaginationRes{Total: total, AllPages: totalPage, Page: page, Limit: limit}, nil
}
