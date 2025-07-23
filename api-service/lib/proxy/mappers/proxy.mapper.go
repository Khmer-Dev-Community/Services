package mappers

import (
	"telegram-service/lib/proxy/dtos"
	"telegram-service/lib/proxy/models"
)

// ToProxyListResponse converts a models.ProxyList to a dtos.ProxyListResponse.
func ToProxyListResponse(model *models.ProxyList) *dtos.ProxyListResponse {
	if model == nil {
		return nil
	}
	return &dtos.ProxyListResponse{
		ID:          model.ID,
		Address:     model.Address,
		Port:        model.Port,
		GroupId:     model.GroupId,
		Sort:        model.Sort,
		CompanyID:   model.CompanyID,
		Description: model.Decription, // Mapped from original model's 'Decription'
		Status:      model.Status,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

// ToProxyListModelFromCreateRequest converts a dtos.CreateProxyListRequest to a models.ProxyList.
func ToProxyListModelFromCreateRequest(req *dtos.CreateProxyListRequest) *models.ProxyList {
	if req == nil {
		return nil
	}
	return &models.ProxyList{
		Address:    req.Address,
		Port:       req.Port,
		GroupId:    req.GroupId,
		Sort:       req.Sort,
		CompanyID:  req.CompanyID,
		Decription: req.Description, // Mapped to original model's 'Decription'
		Status:     req.Status,
	}
}

func UpdateProxyListModelFromUpdateRequest(model *models.ProxyList, req *dtos.UpdateProxyListRequest) {
	if model == nil || req == nil {
		return // Nothing to update if model or request is nil
	}

	if req.Address != nil {
		model.Address = *req.Address
	}
	if req.Port != nil {
		model.Port = *req.Port
	}
	if req.GroupId != nil {
		model.GroupId = *req.GroupId
	}
	if req.Sort != nil {
		model.Sort = *req.Sort
	}
	if req.CompanyID != nil {
		model.CompanyID = *req.CompanyID
	}
	if req.Description != nil {
		model.Decription = *req.Description // Mapped to original model's 'Decription'
	}
	if req.Status != nil {
		model.Status = *req.Status
	}
}
