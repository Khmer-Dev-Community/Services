package repository

import (
	"log"
	"telegram-service/telegram/dtos"
	"telegram-service/telegram/models"
	"telegram-service/utils"

	"gorm.io/gorm"
)

// TelegramGroupRepository handles database operations for TelegramGroup models.
type TelegramGroupRepository struct {
	db *gorm.DB
}

// NewTelegramGroupRepository creates a new instance of TelegramGroupRepository.
func NewTelegramGroupRepository(db *gorm.DB) *TelegramGroupRepository {
	return &TelegramGroupRepository{db: db}
}

// GetTelegramGroupList retrieves a paginated list of Telegram groups based on a filter.
func (r *TelegramGroupRepository) GetTelegramGroupList(page int, limit int, filter *dtos.TelegramGroupFilter) ([]*models.TelegramGroup, int, error) {
	offset := (page - 1) * limit
	var telegramGroups []*models.TelegramGroup
	var total int64

	db := r.db.Model(&models.TelegramGroup{})

	// Apply filters if provided
	if filter != nil {
		if filter.GroupID != nil {
			db = db.Where("group_id = ?", *filter.GroupID)
		}
		if filter.GroupName != nil {
			db = db.Where("group_name LIKE ?", "%"+*filter.GroupName+"%")
		}
		if filter.GameType != nil {
			db = db.Where("game_type LIKE ?", "%"+*filter.GameType+"%")
		}
	}

	// Count total records
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Retrieve paginated records
	err = db.Offset(offset).Limit(limit).Find(&telegramGroups).Error
	if err != nil {
		return nil, 0, err
	}

	log.Printf("Retrieved %d Telegram Groups (Total: %d)", len(telegramGroups), total)
	return telegramGroups, int(total), nil
}

// CreateTelegramGroup creates a new Telegram group in the database.
func (r *TelegramGroupRepository) CreateTelegramGroup(group *models.TelegramGroup) (*models.TelegramGroup, error) {
	err := r.db.Create(group).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(group, "Execute")
	return group, nil
}

// UpdateTelegramGroup updates an existing Telegram group in the database.
func (r *TelegramGroupRepository) UpdateTelegramGroup(group *models.TelegramGroup) (*models.TelegramGroup, error) {
	err := r.db.Save(group).Error // Save will update if primary key exists, otherwise insert
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(group, "Execute")
	return group, nil
}

// DeleteTelegramGroupByID deletes a Telegram group by its ID.
func (r *TelegramGroupRepository) DeleteTelegramGroupByID(id uint) (bool, error) {
	err := r.db.Delete(&models.TelegramGroup{}, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}

// GetTelegramGroupByID retrieves a single Telegram group by its primary key ID.
func (r *TelegramGroupRepository) GetTelegramGroupByID(id uint) (*models.TelegramGroup, error) {
	var telegramGroup models.TelegramGroup
	err := r.db.First(&telegramGroup, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(telegramGroup, "Execute")
	return &telegramGroup, err
}

// GetTelegramGroupByGroupID retrieves a single Telegram group by its GroupID (Telegram's chat ID).
func (r *TelegramGroupRepository) GetTelegramGroupByGroupID(groupID int64) (*models.TelegramGroup, error) {
	var telegramGroup models.TelegramGroup
	err := r.db.Where("group_id = ?", groupID).First(&telegramGroup).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(telegramGroup, "Execute")
	return &telegramGroup, nil
}

func (r *TelegramGroupRepository) GetTelegramGroupName(filter *dtos.TelegramGroupFilter) (*models.TelegramGroup, error) {
	var telegramGroup models.TelegramGroup
	err := r.db.Where("group_name = ?", filter.GroupName).First(&telegramGroup).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	return &telegramGroup, err
}
