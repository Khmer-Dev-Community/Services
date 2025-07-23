package repository

import (
	"errors"
	"fmt"
	"telegram-service/lib/proxy/models"

	"gorm.io/gorm"
)

// ProxyListRepo defines the interface for ProxyList data access operations.
type ProxyListRepo interface {
	Create(proxyList *models.ProxyList) error
	GetByID(id uint) (*models.ProxyList, error)
	GetAll(page, limit int, query string) ([]*models.ProxyList, int64, error)
	Update(proxyList *models.ProxyList) error
	SoftDelete(id uint) (bool, error)
	HardDelete(id uint) (bool, error)
	BulkCreate(proxyLists []*models.ProxyList) error
}

// gormProxyListRepository implements ProxyListRepo using GORM.
type gormProxyListRepository struct {
	db *gorm.DB
}

// NewGormProxyListRepository creates a new instance of gormProxyListRepository.
func NewGormProxyListRepository(db *gorm.DB) ProxyListRepo {
	return &gormProxyListRepository{db: db}
}

// Create creates a new ProxyList record.
func (r *gormProxyListRepository) Create(proxyList *models.ProxyList) error {
	return r.db.Create(proxyList).Error
}

// GetByID retrieves a ProxyList record by its ID.
func (r *gormProxyListRepository) GetByID(id uint) (*models.ProxyList, error) {
	var proxyList models.ProxyList
	err := r.db.First(&proxyList, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if record not found
		}
		return nil, fmt.Errorf("failed to get proxy list by ID: %w", err)
	}
	return &proxyList, nil
}

// GetAll retrieves a paginated and searchable list of ProxyList records.
func (r *gormProxyListRepository) GetAll(page, limit int, query string) ([]*models.ProxyList, int64, error) {
	var proxyLists []*models.ProxyList
	var total int64

	offset := (page - 1) * limit
	db := r.db.Model(&models.ProxyList{})

	if query != "" {
		// Example: search by proxy_address or description
		db = db.Where("proxy_address LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count total proxy lists: %w", err)
	}

	if err := db.Offset(offset).Limit(limit).Find(&proxyLists).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get all proxy lists: %w", err)
	}

	return proxyLists, total, nil
}

// Update updates an existing ProxyList record.
func (r *gormProxyListRepository) Update(proxyList *models.ProxyList) error {
	return r.db.Save(proxyList).Error
}

// SoftDelete soft deletes a ProxyList record by its ID.
func (r *gormProxyListRepository) SoftDelete(id uint) (bool, error) {
	result := r.db.Delete(&models.ProxyList{}, id)
	if result.Error != nil {
		return false, fmt.Errorf("failed to soft delete proxy list: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}

// HardDelete permanently deletes a ProxyList record by its ID.
func (r *gormProxyListRepository) HardDelete(id uint) (bool, error) {
	result := r.db.Unscoped().Delete(&models.ProxyList{}, id)
	if result.Error != nil {
		return false, fmt.Errorf("failed to hard delete proxy list: %w", result.Error)
	}
	return result.RowsAffected > 0, nil
}
func (r *gormProxyListRepository) BulkCreate(proxyLists []*models.ProxyList) error {
	if len(proxyLists) == 0 {
		return nil // Nothing to create
	}
	// GORM will automatically batch if the slice is large, or insert individually.
	return r.db.Create(&proxyLists).Error
}
