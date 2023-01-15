package paginator

import (
	"github.com/sirupsen/logrus"
	"math"
	"strings"
)

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 40
	DefaultOrder    = DESC
)

type Order string

type Input struct {
	Page       uint
	PageSize   uint
	SystemName string
	EventID    uint
	OrderField string
	Order      Order
}

type Output struct {
	TotalItems      uint `json:"total_items"`
	TotalPages      uint `json:"total_pages"`
	Page            uint `json:"page"`
	ItemsPerPage    uint `json:"items_per_page"`
	HasNextPage     bool `json:"has_next_page"`
	HasPreviousPage bool `json:"has_previous_page"`
}

func NewInput(page, pageSize, eventID uint, systemName, orderField string, order string) *Input {
	if page == 0 {
		page = DefaultPage
	}

	if pageSize == 0 {
		pageSize = DefaultPageSize
	}

	return &Input{
		Page:       page,
		PageSize:   pageSize,
		SystemName: systemName,
		EventID:    eventID,
		OrderField: orderField,
		Order:      getOrder(order),
	}
}

func getOrder(order string) Order {
	if strings.ToLower(order) == "asc" {
		return ASC
	}

	return DefaultOrder
}

func (i *Input) LogFields() logrus.Fields {
	return logrus.Fields{
		"page":        i.Page,
		"page_size":   i.PageSize,
		"system_name": i.SystemName,
		"event_id":    i.EventID,
		"order_field": i.OrderField,
		"order":       i.Order,
	}
}

func NewOutput(input *Input, total uint) *Output {
	totalPages := uint(math.Ceil(float64(total) / float64(input.PageSize)))

	return &Output{
		TotalItems:      total,
		TotalPages:      totalPages,
		Page:            input.Page,
		ItemsPerPage:    input.PageSize,
		HasNextPage:     input.Page < totalPages,
		HasPreviousPage: input.Page > 1,
	}
}
