package videosnap

import (
	"telegram-service/utils"
	"time"

	"gorm.io/gorm"
)

type VideosnapRepository struct {
	db *gorm.DB
}

func NewVideosnapRepository(db *gorm.DB) *VideosnapRepository {
	return &VideosnapRepository{db: db}
}

// CreateVideosnap inserts a new videosnap record into the database.
func (r *VideosnapRepository) CreateVideosnap(videosnapDTO VideosnapCreateDTO) (*Videosnap, error) {
	videosnap := videosnapDTO.ToModel()
	err := r.db.Create(videosnap).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videosnap, "Execute")
	return videosnap, nil
}

// UpdateVideosnap updates an existing videosnap record.
func (r *VideosnapRepository) UpdateVideosnap(id uint, videosnapDTO VideosnapUpdateDTO) (*Videosnap, error) {
	var existing Videosnap
	err := r.db.First(&existing, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	updatedVideosnap := videosnapDTO.ToModelForUpdate(&existing)
	// Save the updated record
	err = r.db.Save(updatedVideosnap).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(updatedVideosnap, "Execute")
	return updatedVideosnap, nil
}

// DeleteVideosnapByID deletes a videosnap record by its ID.
func (r *VideosnapRepository) DeleteVideosnapByID(id uint) (bool, error) {
	err := r.db.Delete(&Videosnap{}, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}

// GetVideosnapByID retrieves a videosnap record by its ID.
func (r *VideosnapRepository) GetVideosnapByID(id uint) (*Videosnap, error) {
	var videosnap Videosnap
	err := r.db.First(&videosnap, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videosnap, "Execute")
	return &videosnap, nil
}

// GetVideosnapGameNo retrieves a videosnap record by its ID.
func (r *VideosnapRepository) GetVideosnapByGameNo(gameNo string) (*Videosnap, error) {
	var videosnap Videosnap
	err := r.db.Where("game_no = ?", gameNo).First(&videosnap).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videosnap, "Execute")
	return &videosnap, nil
}

// GetVideosnapList retrieves a list of videosnaps with pagination and optional search by game number.
func (r *VideosnapRepository) GetVideosnapList(page int, limit int, query string) ([]*Videosnap, int, error) {
	offset := (page - 1) * limit
	var videosnaps []*Videosnap
	var total int64
	db := r.db

	utils.LoggerRepository(query, "Execute")
	baseQuery := db.Model(&Videosnap{})
	if query != "" {
		baseQuery = baseQuery.Where("game_no LIKE ?", "%"+query+"%")
		utils.LoggerRepository(baseQuery.Statement.SQL.String(), "Execute")
	}
	// Add order by id
	baseQuery = baseQuery.Order("id DESC") // Change ASC to DESC for descending order if needed
	err := baseQuery.Count(&total).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	err = baseQuery.Offset(offset).Limit(limit).Find(&videosnaps).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	utils.LoggerRepository(videosnaps, "Execute")
	return videosnaps, int(total), nil
}

// GetVideosnapByProcessID retrieves videosnap records by process ID.
func (r *VideosnapRepository) GetVideosnapByProcessID(processID uint) ([]*Videosnap, error) {
	var videosnaps []*Videosnap
	err := r.db.Where("process_id = ?", processID).Find(&videosnaps).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videosnaps, "Execute")
	return videosnaps, nil
}

// GetVideosnapByDateRange retrieves videosnap records within a date range.
func (r *VideosnapRepository) GetVideosnapByDateRange(startDate, endDate time.Time) ([]*Videosnap, error) {
	var videosnaps []*Videosnap
	err := r.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&videosnaps).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(videosnaps, "Execute")
	return videosnaps, nil
}
