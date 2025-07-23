package dtos

import "time"

// ProxyListResponse represents the data structure for sending ProxyList information to the client.
type ProxyListResponse struct {
	ID             uint      `json:"id"`
	Address        string    `json:"proxy_address"`
	Port           string    `json:"proxy_port"`
	GroupId        int       `json:"proxy_group"`
	Sort           int       `json:"sort"`
	CompanyID      int       `json:"companyId"`
	Description    string    `json:"description"`
	Status         bool      `json:"status"`
	Session        string    `json:"proxy_seesion"`
	Password       string    `json:"proxy_password"`
	TimeOut        int       `json:"proxy_timeout"`
	TragetLocation string    `json:"proxy_traget_location"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CreateProxyListRequest struct {
	Address        string `json:"proxy_address" validate:"required,max=250"`
	Port           string `json:"proxy_port" validate:"required,min=1,max=65535"` // Standard port range
	GroupId        int    `json:"proxy_group" validate:"required"`
	Sort           int    `json:"sort"`
	CompanyID      int    `json:"companyId" validate:"required"`  // Assuming CompanyID is provided by frontend or set by middleware
	Description    string `json:"description" validate:"max=255"` // Assuming max length for description
	Status         bool   `json:"status" validate:"required"`
	Session        string `json:"proxy_seesion"`
	Password       string `json:"proxy_password"`
	TimeOut        int    `json:"proxy_timeout"`
	TragetLocation string `json:"proxy_traget_location"`
}

type UpdateProxyListRequest struct {
	ID             uint    `json:"id" validate:"required"` // ID is required for update
	Address        *string `json:"proxy_address,omitempty" validate:"omitempty,max=250"`
	Port           *string `json:"proxy_port,omitempty" validate:"omitempty,min=1,max=65535"`
	GroupId        *int    `json:"proxy_group,omitempty"`
	Sort           *int    `json:"sort,omitempty"`
	CompanyID      *int    `json:"companyId,omitempty"`
	Description    *string `json:"description,omitempty" validate:"omitempty,max=255"`
	Status         *bool   `json:"status,omitempty"`
	Session        string  `json:"proxy_seesion"`
	Password       string  `json:"proxy_password"`
	TimeOut        int     `json:"proxy_timeout"`
	TragetLocation string  `json:"proxy_traget_location"`
}
