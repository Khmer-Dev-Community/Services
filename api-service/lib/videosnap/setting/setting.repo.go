package setting

import (
	"telegram-service/utils"
	"time"

	"gorm.io/gorm"
)

type VideoSnapSettingRepository struct {
	db *gorm.DB
}

// NewVideoSnapSettingRepository initializes a new VideoSnapSettingRepository.
func NewVideoSnapSettingRepository(db *gorm.DB) *VideoSnapSettingRepository {
	return &VideoSnapSettingRepository{db: db}
}

// CreateVideoSnapSetting inserts a new video snap setting record into the database.
func (r *VideoSnapSettingRepository) CreateVideoSnapSetting(videoSnapSettingDTO VideoSnapSettingCreateDTO) (*VideoSnapSetting, error) {
	videoSnapSetting := videoSnapSettingDTO.ToModel()
	err := r.db.Create(videoSnapSetting).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videoSnapSetting, "Execute")
	return videoSnapSetting, nil
}

// UpdateVideoSnapSetting updates an existing video snap setting record.
func (r *VideoSnapSettingRepository) UpdateVideoSnapSetting(id uint, videoSnapSettingDTO VideoSnapSettingUpdateDTO) (*VideoSnapSetting, error) {
	var existing VideoSnapSetting
	err := r.db.First(&existing, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	updatedVideoSnapSetting := videoSnapSettingDTO.ToModelForUpdate(&existing)
	// Save the updated record
	err = r.db.Save(updatedVideoSnapSetting).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(updatedVideoSnapSetting, "Execute")
	return updatedVideoSnapSetting, nil
}

// DeleteVideoSnapSettingByID deletes a video snap setting record by its ID.
func (r *VideoSnapSettingRepository) DeleteVideoSnapSettingByID(id uint) (bool, error) {
	err := r.db.Delete(&VideoSnapSetting{}, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}

// GetVideoSnapSettingByID retrieves a video snap setting record by its ID.
func (r *VideoSnapSettingRepository) GetVideoSnapSettingByID(id uint) (*VideoSnapSetting, error) {
	var videoSnapSetting VideoSnapSetting
	err := r.db.First(&videoSnapSetting, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videoSnapSetting, "Execute")
	return &videoSnapSetting, nil
}

// GetVideoSnapSettingList retrieves a list of video snap settings with pagination and optional search by setting name.
func (r *VideoSnapSettingRepository) GetVideoSnapSettingList(page int, limit int, query string) ([]*VideoSnapSetting, int, error) {
	offset := (page - 1) * limit
	var videoSnapSettings []*VideoSnapSetting
	var total int64
	db := r.db

	utils.LoggerRepository(query, "Execute")
	baseQuery := db.Model(&VideoSnapSetting{})
	if query != "" {
		baseQuery = baseQuery.Where("setting_name LIKE ?", "%"+query+"%")
		utils.LoggerRepository(baseQuery.Statement.SQL.String(), "Execute")
	}
	err := baseQuery.Count(&total).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	err = baseQuery.Offset(offset).Limit(limit).Find(&videoSnapSettings).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	utils.LoggerRepository(videoSnapSettings, "Execute")
	return videoSnapSettings, int(total), nil
}

// GetVideoSnapSettingByProcessID retrieves video snap settings by process ID.
func (r *VideoSnapSettingRepository) GetVideoSnapSettingByProcessID(processID uint) ([]*VideoSnapSetting, error) {
	var videoSnapSettings []*VideoSnapSetting
	err := r.db.Where("process_id = ?", processID).Find(&videoSnapSettings).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videoSnapSettings, "Execute")
	return videoSnapSettings, nil
}

// GetVideoSnapSettingByDateRange retrieves video snap settings within a date range.
func (r *VideoSnapSettingRepository) GetVideoSnapSettingByDateRange(startDate, endDate time.Time) ([]*VideoSnapSetting, error) {
	var videoSnapSettings []*VideoSnapSetting
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&videoSnapSettings).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videoSnapSettings, "Execute")
	return videoSnapSettings, nil
}
