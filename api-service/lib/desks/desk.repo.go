package desks

import (
	"telegram-service/utils"
	"time"

	"gorm.io/gorm"
)

type DeskSettingRepository struct {
	db *gorm.DB
}

// NewDeskSettingRepository initializes a new DeskSettingRepository.
func NewDeskSettingRepository(db *gorm.DB) *DeskSettingRepository {
	return &DeskSettingRepository{db: db}
}

// DeskSetting initializes for redis
func (r *DeskSettingRepository) IntDeskSetting() ([]*DeskSetting, error) {
	var deskSettings []*DeskSetting
	err := r.db.Find(&deskSettings).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(deskSettings, "Execute")
	return deskSettings, nil
}

// CreateDeskSetting inserts a new desk setting record into the database.
func (r *DeskSettingRepository) CreateDeskSetting(deskSettingDTO DeskSettingCreateDTO) (*DeskSetting, error) {
	deskSetting := deskSettingDTO.ToModel()
	err := r.db.Create(deskSetting).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(deskSetting, "Execute")
	return deskSetting, nil
}

// UpdateDeskSetting updates an existing desk setting record.
func (r *DeskSettingRepository) UpdateDeskSetting(id uint, deskSettingDTO DeskSettingUpdateDTO) (*DeskSetting, error) {
	var existing DeskSetting
	err := r.db.First(&existing, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	updatedDeskSetting := deskSettingDTO.ToModelForUpdate(&existing)
	// Save the updated record
	err = r.db.Save(updatedDeskSetting).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(updatedDeskSetting, "Execute")
	return updatedDeskSetting, nil
}

// DeleteDeskSettingByID deletes a desk setting record by its ID.
func (r *DeskSettingRepository) DeleteDeskSettingByID(id uint) (bool, error) {
	err := r.db.Delete(&DeskSetting{}, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}

// GetDeskSettingByID retrieves a desk setting record by its ID.
func (r *DeskSettingRepository) GetDeskSettingByID(id uint) (*DeskSetting, error) {
	var deskSetting DeskSetting
	err := r.db.First(&deskSetting, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(deskSetting, "Execute")
	return &deskSetting, nil
}

// GetDeskSettingList retrieves a list of desk settings with pagination and optional search by setting name.
func (r *DeskSettingRepository) GetDeskSettingList(page int, limit int, query string) ([]*DeskSetting, int, error) {
	offset := (page - 1) * limit
	var deskSettings []*DeskSetting
	var total int64
	db := r.db

	utils.LoggerRepository(query, "Execute")
	baseQuery := db.Model(&DeskSetting{}).Preload("StreamServer")
	if query != "" {
		baseQuery = baseQuery.Where("setting_name LIKE ?", "%"+query+"%")
		utils.LoggerRepository(baseQuery.Statement.SQL.String(), "Execute")
	}
	err := baseQuery.Count(&total).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	err = baseQuery.Offset(offset).Limit(limit).Find(&deskSettings).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	utils.LoggerRepository(deskSettings, "Execute")
	return deskSettings, int(total), nil
}

// GetDeskSettingByDateRange retrieves desk settings within a date range.
func (r *DeskSettingRepository) GetDeskSettingByDateRange(startDate, endDate time.Time) ([]*DeskSetting, error) {
	var deskSettings []*DeskSetting
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&deskSettings).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(deskSettings, "Execute")
	return deskSettings, nil
}
