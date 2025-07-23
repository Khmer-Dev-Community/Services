package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"telegram-service/lib/proxy/dtos"
	"telegram-service/lib/proxy/services"
	"telegram-service/utils"

	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	proxyUrlTemplate = "http://%s:%s@%s:%s"
)

const (
	ADDRESS  = "proxy.geonode.io"
	USER     = "geonode_99xssXgUIa-type-residential"
	PASSWORD = "0f5ba465-1d87-4066-ba69-f30912d2a64e"
	PORT     = "9000"
)

// proxy.geonode.io:10000:geonode_99xssXgUIa-type-residential-country-pk-lifetime-360-session-zhJJXr:5c4be9c4-7fb5-433c-b77d-9a88c5159cd4

// http://geonode_99xssXgUIa-type-residential-country-pk-lifetime-session-X9YpQz1N:0f5ba465-1d87-4066-ba69-f30912d2a64e@proxy.geonode.io:9000
type ProxyList struct {
	IP       string `json:"ip,omitempty"`
	Port     int    `json:"port,omitempty"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	TimeOut  string `json:"timeout,omitempty"`
	Session  int    `json:"session,omitempty"`
	Address  string `json:"address"`
}
type RequestFilter struct {
	Limit       int    `json:"limit" validate:"required,max=5000"`
	Port        string `json:"port" validate:"required,min=1,max=65535"` // Standard port range
	GroupId     string `json:"group" validate:"required"`
	Location    string `json:"proxy_traget_location"`
	CompanyID   int    `json:"companyId" validate:"required"`  // Assuming CompanyID is provided by frontend or set by middleware
	Description string `json:"description" validate:"max=255"` // Assuming max length for description
	Status      bool   `json:"status" validate:"required"`
}

func PgetProxyByTargetLocation(w http.ResponseWriter, r *http.Request) {
	var proxies []ProxyList
	Qlimit := r.URL.Query().Get("limit")
	QLocation := r.URL.Query().Get("location")

	limitInt, err := strconv.Atoi(Qlimit)
	if err != nil || limitInt <= 0 {
		limitInt = 100
	}

	for i := 0; i < limitInt; i++ {
		sessionID := randString(8)
		addressUser := fmt.Sprintf("%s-country-%s-lifetime-session-%s", USER, QLocation, sessionID)
		formattedProxy := fmt.Sprintf(proxyUrlTemplate, addressUser, PASSWORD, ADDRESS, PORT)

		proxies = append(proxies, ProxyList{
			UserName: addressUser,
			Password: PASSWORD,
			Address:  formattedProxy,
			Session:  i + 1,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proxies)
}

// randString generates a random alphanumeric string of given length
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// ProxyListController handles HTTP requests for proxy list management.
type ProxyListController struct {
	proxyListService services.ProxyListServicer
}

// NewProxyListController initializes a new ProxyListController.
func NewProxyListController(svr services.ProxyListServicer) *ProxyListController {
	return &ProxyListController{proxyListService: svr}
}
func writeErrorResponse(w http.ResponseWriter, err error, context string) {
	var se *utils.ServiceError
	if errors.As(err, &se) {
		utils.ErrorLog(se.Err, fmt.Sprintf("%s: %s", context, se.Message))
		utils.HttpSuccessResponse(w, nil, se.StatusCode, se.Message) // HttpSuccessResponse can send errors too based on message and status
	} else {
		utils.ErrorLog(err, fmt.Sprintf("%s: %v", context, err))
		utils.HttpSuccessResponse(w, nil, http.StatusInternalServerError, "an unexpected error occurred")
	}
}

// ProxyListControllerGetList handles GET requests to retrieve a paginated list of ProxyList records.
// @Summary Get all proxy list records
// @Description Get a list of proxy list records with pagination and fuzzy search
// @Tags ProxyLists
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)" default(1)
// @Param limit query int false "Number of records per page (default: 10)" default(10)
// @Param query query string false "Search query for proxy address or description"
// @Success 200 {object} utils.PaginationResponse{data=[]dtos.ProxyListResponse} "Successfully retrieved proxy list"
// @Failure 500 {object} utils.ErrorResponseDTO "Internal server error"
// @Router /proxylists [get]
func (ctrl *ProxyListController) ProxyListControllerGetList(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	query := r.URL.Query().Get("query")

	pageInt, err := utils.ParseIntOrDefault(pageStr, 1)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("invalid 'page' parameter: %w", err), "ProxyListControllerGetList")
		return
	}

	limitInt, err := utils.ParseIntOrDefault(limitStr, 10)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("invalid 'limit' parameter: %w", err), "ProxyListControllerGetList")
		return
	}

	proxyListResponseData, total, err := ctrl.proxyListService.GetAllProxyLists(pageInt, limitInt, query)
	if err != nil {
		writeErrorResponse(w, err, "ProxyListControllerGetList: failed to fetch proxy list from service")
		return
	}

	lastPage := 1
	if limitInt > 0 && total > 0 {
		lastPage = (int(total) + limitInt - 1) / limitInt
	}

	// Assuming utils.NewHttpPaginationResponse exists and works with the provided data
	utils.NewHttpPaginationResponse(
		w,
		proxyListResponseData,
		total,
		pageInt,
		lastPage,
		http.StatusOK,
		"Proxy lists retrieved successfully",
	)
}

// ProxyListControllerGetByID handles GET requests to retrieve a single ProxyList record by its ID.
// @Summary Get a proxy list record by ID
// @Description Get a proxy list record by its unique ID
// @Tags ProxyLists
// @Accept json
// @Produce json
// @Param id path int true "ProxyList ID"
// @Success 200 {object} dtos.ProxyListResponse "Successfully retrieved proxy list record"
// @Failure 400 {object} utils.ErrorResponseDTO "Invalid ID parameter"
// @Failure 404 {object} utils.ErrorResponseDTO "ProxyList record not found"
// @Failure 500 {object} utils.ErrorResponseDTO "Internal server error"
// @Router /proxylists/{id} [get]
func (ctrl *ProxyListController) ProxyListControllerGetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := utils.ParseUintID(idStr)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("invalid ID format: %w", err), "ProxyListControllerGetByID")
		return
	}

	proxyListResponseDTO, err := ctrl.proxyListService.GetProxyListByID(id)
	if err != nil {
		writeErrorResponse(w, err, fmt.Sprintf("ProxyListControllerGetByID: error fetching proxy list ID %d", id))
		return
	}

	utils.HttpSuccessResponse(w, proxyListResponseDTO, http.StatusOK, "Proxy list retrieved successfully")
}

// ProxyListControllerCreate handles POST requests to create a new ProxyList record.
// @Summary Create a new proxy list record
// @Description Create a new proxy list record.
// @Tags ProxyLists
// @Accept json
// @Produce json
// @Param body body dtos.CreateProxyListRequest true "ProxyList creation details"
// @Success 201 {object} dtos.ProxyListResponse "Successfully created proxy list record"
// @Failure 400 {object} utils.ErrorResponseDTO "Invalid request payload or validation error"
// @Failure 500 {object} utils.ErrorResponseDTO "Internal server error"
// @Router /proxylists [post]
func (ctrl *ProxyListController) ProxyListControllerCreate(w http.ResponseWriter, r *http.Request) {
	var createDTO dtos.CreateProxyListRequest
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		utils.ErrorLog(err, "ProxyListControllerCreate: failed to decode request body")
		writeErrorResponse(w, fmt.Errorf("invalid request payload: %w", err), "ProxyListControllerCreate")
		return
	}
	createdResponseDTO, err := ctrl.proxyListService.CreateProxyList(createDTO)
	if err != nil {
		writeErrorResponse(w, err, "ProxyListControllerCreate: failed to create proxy list via service")
		return
	}
	utils.HttpSuccessResponse(w, createdResponseDTO, http.StatusCreated, "Proxy list created successfully")
}

func (ctrl *ProxyListController) ProxyListControllerGenerateBulk(w http.ResponseWriter, r *http.Request) {
	var createDTO RequestFilter
	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		utils.ErrorLog(err, "ProxyListControllerCreate: failed to decode request body")
		writeErrorResponse(w, fmt.Errorf("invalid request payload: %w", err), "ProxyListControllerCreate")
		return
	}
	limitInt := createDTO.Limit
	QLocation := createDTO.Location
	if limitInt <= 0 {
		limitInt = 100
	}
	if limitInt > 1000 { // Prevent excessively large requests
		limitInt = 1000
		utils.InfoLog(fmt.Sprintf("Bulk create limit capped at %d", limitInt), "ProxyListControllerGenerateBulk")
	}
	if QLocation == "" {
		QLocation = "unknown" // Default location
	}

	var proxiesToCreate []dtos.CreateProxyListRequest
	for i := 0; i < limitInt; i++ {
		sessionID := randString(8)
		// Ensure USER, PASSWORD, ADDRESS, PORT, proxyUrlTemplate are defined globally or passed
		addressUser := fmt.Sprintf("%s-country-%s-lifetime-session-%s", USER, QLocation, sessionID)
		formattedProxy := fmt.Sprintf(proxyUrlTemplate, addressUser, PASSWORD, ADDRESS, PORT)

		proxiesToCreate = append(proxiesToCreate, dtos.CreateProxyListRequest{
			Address:        formattedProxy,
			GroupId:        1,
			CompanyID:      1,
			Description:    fmt.Sprintf("Generated proxy for %s - Session %s", QLocation, sessionID),
			Status:         true,
			Session:        sessionID,
			Password:       PASSWORD,
			Port:           PORT,
			TragetLocation: QLocation,
		})
	}
	createdResponses, err := ctrl.proxyListService.BulkCreateProxyLists(proxiesToCreate)
	if err != nil {
		writeErrorResponse(w, err, "ProxyListControllerGenerateBulk: failed to bulk create proxy lists via service")
		return
	}

	utils.HttpSuccessResponse(w, createdResponses, http.StatusCreated, fmt.Sprintf("Successfully generated and created %d proxy lists", len(createdResponses)))
}

// ProxyListControllerUpdate handles PUT requests to update an existing ProxyList record.
// @Summary Update an existing proxy list record
// @Description Update an existing proxy list record by ID with provided details.
// @Tags ProxyLists
// @Accept json
// @Produce json
// @Param body body dtos.UpdateProxyListRequest true "ProxyList update details including ID"
// @Success 200 {object} dtos.ProxyListResponse "Successfully updated proxy list record"
// @Failure 400 {object} utils.ErrorResponseDTO "Invalid request payload or validation error"
// @Failure 404 {object} utils.ErrorResponseDTO "ProxyList record not found"
// @Failure 500 {object} utils.ErrorResponseDTO "Internal server error"
// @Router /proxylists [put]
func (ctrl *ProxyListController) ProxyListControllerUpdate(w http.ResponseWriter, r *http.Request) {
	var updateDTO dtos.UpdateProxyListRequest
	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		utils.ErrorLog(err, "ProxyListControllerUpdate: failed to decode request body")
		writeErrorResponse(w, fmt.Errorf("invalid request payload: %w", err), "ProxyListControllerUpdate")
		return
	}

	if updateDTO.ID == 0 {
		writeErrorResponse(w, fmt.Errorf("proxy list ID is required for update"), "ProxyListControllerUpdate")
		return
	}

	// You would typically validate the DTO here
	// if err := validator.New().Struct(updateDTO); err != nil {
	// 	writeErrorResponse(w, utils.NewServiceError(err, "validation failed", http.StatusBadRequest), "ProxyListControllerUpdate")
	// 	return
	// }

	updatedResponseDTO, err := ctrl.proxyListService.UpdateProxyList(updateDTO)
	if err != nil {
		writeErrorResponse(w, err, "ProxyListControllerUpdate: failed to update proxy list via service")
		return
	}

	utils.HttpSuccessResponse(w, updatedResponseDTO, http.StatusOK, "Proxy list updated successfully")
}

// ProxyListControllerDelete handles DELETE requests to soft delete a ProxyList record by ID.
// @Summary Delete a proxy list record
// @Description Soft delete a proxy list record by its ID
// @Tags ProxyLists
// @Accept json
// @Produce json
// @Param id path int true "ProxyList ID"
// @Success 200 {object} object{deleted=boolean} "Successfully soft deleted proxy list record"
// @Failure 400 {object} utils.ErrorResponseDTO "Invalid ID parameter"
// @Failure 404 {object} utils.ErrorResponseDTO "ProxyList record not found"
// @Failure 500 {object} utils.ErrorResponseDTO "Internal server error"
// @Router /proxylists/{id} [delete]
func (ctrl *ProxyListController) ProxyListControllerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := utils.ParseUintID(idStr)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("invalid proxy list ID format: %w", err), "ProxyListControllerDelete")
		return
	}

	success, err := ctrl.proxyListService.DeleteProxyList(id)
	if err != nil {
		writeErrorResponse(w, err, fmt.Sprintf("ProxyListControllerDelete: failed to soft delete proxy list ID %d", id))
		return
	}

	utils.HttpSuccessResponse(w, map[string]bool{"deleted": success}, http.StatusOK, "Proxy list soft deleted successfully")
}

// ProxyListControllerHardDelete handles DELETE requests to permanently remove a ProxyList record by ID.
// @Summary Permanently delete a proxy list record
// @Description Permanently delete a proxy list record by its ID. Use with extreme caution.
// @Tags ProxyLists
// @Accept json
// @Produce json
// @Param id path int true "ProxyList ID"
// @Success 200 {object} object{deleted=boolean} "Successfully permanently deleted proxy list record"
// @Failure 400 {object} utils.ErrorResponseDTO "Invalid ID parameter"
// @Failure 404 {object} utils.ErrorResponseDTO "ProxyList record not found"
// @Failure 500 {object} utils.ErrorResponseDTO "Internal server error"
// @Router /proxylists/hard/{id} [delete]
func (ctrl *ProxyListController) ProxyListControllerHardDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := utils.ParseUintID(idStr)
	if err != nil {
		writeErrorResponse(w, fmt.Errorf("invalid proxy list ID format: %w", err), "ProxyListControllerHardDelete")
		return
	}

	success, err := ctrl.proxyListService.HardDeleteProxyList(id)
	if err != nil {
		writeErrorResponse(w, err, fmt.Sprintf("ProxyListControllerHardDelete: failed to hard delete proxy list ID %d", id))
		return
	}

	utils.HttpSuccessResponse(w, map[string]bool{"deleted": success}, http.StatusOK, "Proxy list permanently deleted successfully")
}
