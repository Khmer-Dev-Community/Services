package utils

import "math"

type Meta struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	LastPage int `json:"lastPage"`
}

type PaginationResponse struct {
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
	Status  uint        `json:"status"`
	Message string      `json:"message"`
}

func NewPaginationResponse(data interface{}, total, page, limit int, status int, message string) PaginationResponse {
	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	return PaginationResponse{
		Data: data,
		Meta: Meta{
			Total:    total,
			Page:     page,
			LastPage: lastPage,
		},
		Status:  uint(status),
		Message: message,
	}
}
