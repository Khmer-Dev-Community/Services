package services

import (
	"fmt"
	"net/http"

	"telegram-service/lib/proxy/dtos"
	"telegram-service/lib/proxy/mappers"
	"telegram-service/lib/proxy/models"
	"telegram-service/lib/proxy/repository"

	"telegram-service/utils" // Assuming you have a utils package for logging/errors
)

// ProxyListServicer defines the interface for ProxyList business logic operations.
type ProxyListServicer interface {
	CreateProxyList(req dtos.CreateProxyListRequest) (*dtos.ProxyListResponse, error)
	GetProxyListByID(id uint) (*dtos.ProxyListResponse, error)
	GetAllProxyLists(page, limit int, query string) ([]*dtos.ProxyListResponse, int, error)
	UpdateProxyList(req dtos.UpdateProxyListRequest) (*dtos.ProxyListResponse, error)
	DeleteProxyList(id uint) (bool, error)
	HardDeleteProxyList(id uint) (bool, error)
	BulkCreateProxyLists(reqs []dtos.CreateProxyListRequest) ([]*dtos.ProxyListResponse, error)
}

// proxyListService is the concrete implementation of the ProxyListServicer interface.
type proxyListService struct {
	repo repository.ProxyListRepo
}

// NewProxyListService creates and returns a new ProxyListServicer instance.
func NewProxyListService(repo repository.ProxyListRepo) ProxyListServicer {
	return &proxyListService{repo: repo}
}

// CreateProxyList handles the creation of a new proxy list entry.
func (s *proxyListService) CreateProxyList(req dtos.CreateProxyListRequest) (*dtos.ProxyListResponse, error) {
	utils.InfoLog(req, "CreateProxyList Service: request received")
	proxyListModel := mappers.ToProxyListModelFromCreateRequest(&req)
	if err := s.repo.Create(proxyListModel); err != nil {
		utils.LoggerRepository(err, "CreateProxyList Service: failed to create proxy list in repository")
		return nil, utils.NewServiceError(err, "failed to create proxy list", http.StatusInternalServerError)
	}
	// Convert Model back to Response DTO
	resp := mappers.ToProxyListResponse(proxyListModel)
	utils.InfoLog(resp, "CreateProxyList Service: proxy list created successfully")
	return resp, nil
}

// BulkCreateProxyLists handles the bulk creation of new proxy list entries.
func (s *proxyListService) BulkCreateProxyLists(reqs []dtos.CreateProxyListRequest) ([]*dtos.ProxyListResponse, error) {
	utils.InfoLog(fmt.Sprintf("BulkCreateProxyLists Service: received %d requests", len(reqs)), "BulkCreateProxyLists Service")

	if len(reqs) == 0 {
		return []*dtos.ProxyListResponse{}, nil // Return empty slice if no requests
	}

	var proxyListModels []*models.ProxyList
	for _, req := range reqs {
		proxyListModels = append(proxyListModels, mappers.ToProxyListModelFromCreateRequest(&req))
	}

	// Persist to database in bulk
	if err := s.repo.BulkCreate(proxyListModels); err != nil {
		utils.LoggerRepository(err, "BulkCreateProxyLists Service: failed to bulk create proxy lists in repository")
		return nil, utils.NewServiceError(err, "failed to bulk create proxy lists", http.StatusInternalServerError)
	}

	// Convert created Models back to Response DTOs (they should now have their IDs)
	var responses []*dtos.ProxyListResponse
	for _, model := range proxyListModels {
		responses = append(responses, mappers.ToProxyListResponse(model))
	}

	utils.InfoLog(fmt.Sprintf("BulkCreateProxyLists Service: successfully created %d proxy lists", len(responses)), "BulkCreateProxyLists Service")
	return responses, nil
}

// GetProxyListByID retrieves a proxy list entry by its ID.
func (s *proxyListService) GetProxyListByID(id uint) (*dtos.ProxyListResponse, error) {
	utils.InfoLog(fmt.Sprintf("ProxyList ID: %d", id), "GetProxyListByID Service: request received")

	// Fetch from repository
	proxyListModel, err := s.repo.GetByID(id)
	if err != nil {
		utils.LoggerRepository(err, "GetProxyListByID Service: failed to fetch proxy list from repository")
		return nil, utils.NewServiceError(err, "failed to retrieve proxy list due to internal error", http.StatusInternalServerError)
	}
	if proxyListModel == nil {
		return nil, utils.NewServiceError(nil, fmt.Sprintf("proxy list with ID %d not found", id), http.StatusNotFound)
	}

	// Convert Model to Response DTO
	resp := mappers.ToProxyListResponse(proxyListModel)
	utils.InfoLog(resp, "GetProxyListByID Service: proxy list fetched successfully")
	return resp, nil
}

// GetAllProxyLists retrieves a paginated and searchable list of proxy lists.
func (s *proxyListService) GetAllProxyLists(page, limit int, query string) ([]*dtos.ProxyListResponse, int, error) {
	utils.InfoLog(fmt.Sprintf("Page: %d, Limit: %d, Query: %s", page, limit, query), "GetAllProxyLists Service: request received")

	// Fetch from repository
	proxyListModels, total, err := s.repo.GetAll(page, limit, query)
	if err != nil {
		utils.LoggerRepository(err, "GetAllProxyLists Service: failed to fetch proxy lists from repository")
		return nil, 0, utils.NewServiceError(err, "failed to retrieve proxy lists due to internal error", http.StatusInternalServerError)
	}

	// Convert Models to Response DTOs
	var responses []*dtos.ProxyListResponse
	for _, model := range proxyListModels {
		responses = append(responses, mappers.ToProxyListResponse(model))
	}
	utils.InfoLog(fmt.Sprintf("Fetched %d proxy lists, total %d", len(responses), total), "GetAllProxyLists Service: proxy lists fetched successfully")
	return responses, int(total), nil
}

// UpdateProxyList handles updating an existing proxy list entry.
func (s *proxyListService) UpdateProxyList(req dtos.UpdateProxyListRequest) (*dtos.ProxyListResponse, error) {
	utils.InfoLog(req, "UpdateProxyList Service: request received")

	// Retrieve existing model
	existingProxyListModel, err := s.repo.GetByID(req.ID)
	if err != nil {
		utils.LoggerRepository(err, "UpdateProxyList Service: failed to fetch proxy list for update")
		return nil, utils.NewServiceError(err, "failed to retrieve proxy list for update due to internal error", http.StatusInternalServerError)
	}
	if existingProxyListModel == nil {
		return nil, utils.NewServiceError(nil, fmt.Sprintf("proxy list with ID %d not found for update", req.ID), http.StatusNotFound)
	}

	// Update model from DTO
	mappers.UpdateProxyListModelFromUpdateRequest(existingProxyListModel, &req)

	// Persist changes to database
	if err := s.repo.Update(existingProxyListModel); err != nil {
		utils.LoggerRepository(err, "UpdateProxyList Service: failed to update proxy list in repository")
		return nil, utils.NewServiceError(err, "failed to update proxy list", http.StatusInternalServerError)
	}

	// Convert updated Model to Response DTO
	resp := mappers.ToProxyListResponse(existingProxyListModel)
	utils.InfoLog(resp, "UpdateProxyList Service: proxy list updated successfully")
	return resp, nil
}

// DeleteProxyList handles soft deleting a proxy list entry.
func (s *proxyListService) DeleteProxyList(id uint) (bool, error) {
	utils.InfoLog(fmt.Sprintf("ProxyList ID: %d", id), "DeleteProxyList Service: soft delete request received")

	deleted, err := s.repo.SoftDelete(id)
	if err != nil {
		utils.LoggerRepository(err, "DeleteProxyList Service: failed to soft delete proxy list in repository")
		return false, utils.NewServiceError(err, "failed to soft delete proxy list", http.StatusInternalServerError)
	}
	if !deleted {
		return false, utils.NewServiceError(nil, fmt.Sprintf("proxy list with ID %d not found for soft delete", id), http.StatusNotFound)
	}
	utils.InfoLog(fmt.Sprintf("ProxyList ID %d", id), "DeleteProxyList Service: proxy list soft deleted successfully")
	return true, nil
}

// HardDeleteProxyList handles permanently deleting a proxy list entry.
func (s *proxyListService) HardDeleteProxyList(id uint) (bool, error) {
	utils.InfoLog(fmt.Sprintf("ProxyList ID: %d", id), "HardDeleteProxyList Service: permanent delete request received")

	deleted, err := s.repo.HardDelete(id)
	if err != nil {
		utils.LoggerRepository(err, "HardDeleteProxyList Service: failed to hard delete proxy list in repository")
		return false, utils.NewServiceError(err, "failed to hard delete proxy list", http.StatusInternalServerError)
	}
	if !deleted {
		return false, utils.NewServiceError(nil, fmt.Sprintf("proxy list with ID %d not found for hard delete", id), http.StatusNotFound)
	}
	utils.InfoLog(fmt.Sprintf("ProxyList ID %d", id), "HardDeleteProxyList Service: proxy list permanently deleted successfully")
	return true, nil
}
