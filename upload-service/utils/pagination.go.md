package utils

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

type MetaData struct {
	Total    int64  `json:"total"`
	Page     string `json:"page"`
	LastPage int    `json:"last_page"`
}

// PaginateResult represents the paginated result
type PaginateResult struct {
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Meta    MetaData    `json:"meta"`
}
type PaginationFilter struct {
	page         int
	limit        int
	filtersearch string
}

// Paginate executes pagination logic
func Paginate(c *gin.Context, query *gorm.DB, model interface{}, preloadOptions map[string]bool) (*PaginateResult, error) {
	// Parse pagination parameters from query string
	pagination, err := ParsePaginationParams(c)
	if err != nil {
		return nil, err
	}

	// Retrieve total count
	var totalCount int64
	if err := query.Model(model).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Ensure pagination limit is valid
	if pagination.Limit <= 0 {
		return nil, errors.New("invalid limit parameter")
	}

	// Convert pagination limit to int64
	limit := int64(pagination.Limit)

	// Calculate last page
	lastPage := totalCount / limit
	if totalCount%limit != 0 {
		lastPage++
	}

	// Calculate offset based on page and limit
	offset := (pagination.Page - 1) * pagination.Limit

	// Create a slice of the appropriate type dynamically
	sliceType := reflect.SliceOf(reflect.TypeOf(model).Elem())
	data := reflect.New(sliceType).Interface()

	// Preload relationships based on the preload options
	for preloadField, shouldPreload := range preloadOptions {
		if shouldPreload {
			query = query.Preload(preloadField)
		}
	}

	// Retrieve data with pagination and preloaded relationships
	if err := query.Offset(offset).Limit(pagination.Limit).Find(data).Error; err != nil {
		fmt.Println("Pagination Query Error:", err.Error())
		return nil, err
	}

	// Create paginated result
	result := &PaginateResult{
		Data:    data,
		Status:  http.StatusText(200),
		Message: "Data retrieved successfully",
		Meta: MetaData{
			Total:    totalCount,
			Page:     strconv.Itoa(pagination.Page),
			LastPage: int(lastPage),
		},
	}
	return result, nil
}

// Paginate executes pagination logic with preload
func PaginateWithPayload(Ispage int, Islimit int, filter string, query *gorm.DB, model interface{}, preload []string, payload map[string][]string) (*PaginateResult, error) {
	// Parse pagination parameters from query string
	pagination, err := ParsePaginationParamsQuery(Ispage, Islimit)
	if err != nil {
		return nil, err
	}

	// Retrieve total count
	var totalCount int64
	if err := query.Model(model).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Ensure pagination limit is valid
	if pagination.Limit <= 0 {
		return nil, errors.New("invalid limit parameter")
	}

	// Convert pagination limit to int64
	limit := int64(pagination.Limit)

	// Calculate last page
	lastPage := totalCount / limit
	if totalCount%limit != 0 {
		lastPage++
	}

	// Calculate offset based on page and limit
	offset := (pagination.Page - 1) * pagination.Limit

	// Apply preload with specific columns if provided
	for _, p := range preload {
		if cols, ok := payload[p]; ok {
			query = query.Preload(p, func(db *gorm.DB) *gorm.DB {
				return db.Select(cols)
			})
		} else {
			query = query.Preload(p)
		}
	}

	// Create a slice of the appropriate type dynamically
	sliceType := reflect.SliceOf(reflect.TypeOf(model).Elem())
	data := reflect.New(sliceType).Interface()

	// Retrieve data with pagination
	if err := query.Offset(offset).Limit(pagination.Limit).Find(data).Error; err != nil {
		fmt.Println("Pagination Query Error:", err.Error())
		return nil, err
	}

	// Filter data to include only selected columns from payload
	filteredData := filterData(data, payload)

	// Create paginated result
	result := &PaginateResult{
		Data:    filteredData,
		Status:  "success",
		Message: "Data retrieved successfully",
		Meta: MetaData{
			Total:    totalCount,
			Page:     strconv.Itoa(pagination.Page),
			LastPage: int(lastPage),
		},
	}
	return result, nil
}

func filterData(data interface{}, payload map[string][]string) interface{} {
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice {
		return data
	}

	// Create a new slice to hold the filtered data
	filteredData := reflect.MakeSlice(dataValue.Type(), dataValue.Len(), dataValue.Cap())

	// Iterate through the original data slice
	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)

		// Create a new map to hold the filtered fields for the current item
		filteredItem := reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))

		// Iterate through the fields of the item
		for k := range payload {
			fieldValue := item.FieldByName(k)
			if fieldValue.IsValid() {
				filteredItem.SetMapIndex(reflect.ValueOf(k), fieldValue)
			}
		}

		// Set the filtered item in the new slice
		filteredData.Index(i).Set(filteredItem)
	}

	return filteredData.Interface()
}

// ParsePaginationParams parses pagination parameters from query string
func ParsePaginationParams(c *gin.Context) (PaginationParams, error) {
	const defaultPage = 1
	const defaultLimit = 10

	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	if pageStr == "" || limitStr == "" {
		return PaginationParams{}, errors.New("page and limit parameters are required")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return PaginationParams{}, errors.New("invalid page parameter")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return PaginationParams{}, errors.New("invalid limit parameter")
	}

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}, nil
}

// RespondWithPagination responds with paginated data
func RespondWithPagination(c *gin.Context, data interface{}, status, message string, meta MetaData) {
	result := &PaginateResult{
		Data:    data,
		Status:  status,
		Message: message,
		Meta:    meta,
	}

	c.JSON(http.StatusOK, result)
}

func ParsePaginationParamsQuery(page int, limit int) (PaginationParams, error) {

	return PaginationParams{
		Page:  page,
		Limit: limit,
	}, nil
}
