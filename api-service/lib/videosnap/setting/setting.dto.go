package setting

import "time"

type VideoSnapSettingCreateDTO struct {
	ID               uint   `json:"id"`
	Rtmp             string `json:"rtmp" binding:"required"` // RTMP URL
	ServerName       string `json:"server_name"`
	Duration         int    `json:"duration" binding:"required"`           // Duration in seconds
	VideoSize        string `json:"video_size" binding:"required"`         // Video size
	VideoType        string `json:"video_type" binding:"required"`         // Video type
	BeforeStart      int    `json:"before_start" binding:"required"`       // Time before start in seconds
	AfterEnd         int    `json:"after_end" binding:"required"`          // Time after end in seconds
	OutputPath       string `json:"output_path" binding:"required"`        // Output path for video
	ServiceAcccount  string `json:"service_account" binding:"required"`    // Service account
	AcceessDomain    string `json:"access_domain" binding:"required"`      // Access domain
	BucketName       string `json:"bucket_name" binding:"required"`        // Bucket name
	DeleteLocalStore bool   `json:"delete_local_store" binding:"required"` // Flag to delete local store
}

type VideoSnapSettingUpdateDTO struct {
	ID               uint    `json:"id" binding:"required"` // ID of the VideoSnapSetting
	ServerName       *string `json:"server_name"`
	Rtmp             *string `json:"rtmp,omitempty"`               // Optional: RTMP URL
	Duration         *int    `json:"duration,omitempty"`           // Optional: Duration in seconds
	VideoSize        *string `json:"video_size,omitempty"`         // Optional: Video size
	VideoType        *string `json:"video_type,omitempty"`         // Optional: Video type
	BeforeStart      *int    `json:"before_start,omitempty"`       // Optional: Time before start in seconds
	AfterEnd         *int    `json:"after_end,omitempty"`          // Optional: Time after end in seconds
	OutputPath       *string `json:"output_path,omitempty"`        // Optional: Output path for video
	ServiceAcccount  *string `json:"service_account,omitempty"`    // Optional: Service account
	BucketName       *string `json:"bucket_name,omitempty"`        // Optional: Bucket name
	AcceessDomain    *string `json:"access_domain,omitempty"`      // Optional: Access domain
	DeleteLocalStore *bool   `json:"delete_local_store,omitempty"` // Optional: Flag to delete local store
}

type VideoSnapSettingResponseDTO struct {
	ID               uint      `json:"id"` // ID of the VideoSnapSetting
	ServerName       string    `json:"server_name"`
	Rtmp             string    `json:"rtmp"`               // RTMP URL
	Duration         int       `json:"duration"`           // Duration in seconds
	VideoSize        string    `json:"video_size"`         // Video size
	VideoType        string    `json:"video_type"`         // Video type
	BeforeStart      int       `json:"before_start"`       // Time before start in seconds
	AfterEnd         int       `json:"after_end"`          // Time after end in seconds
	OutputPath       string    `json:"output_path"`        // Output path for video
	ServiceAcccount  string    `json:"service_account"`    // Service account
	BucketName       string    `json:"bucket_name"`        // Bucket name
	AcceessDomain    string    `json:"access_domain"`      // Access domain
	DeleteLocalStore bool      `json:"delete_local_store"` // Flag to delete local store
	CreatedAt        time.Time `json:"created_at"`         // Created timestamp
	UpdatedAt        time.Time `json:"updated_at"`         // Updated timestamp
}

func (setting *VideoSnapSetting) FromModelX() *VideoSnapSettingResponseDTO {
	if setting == nil {
		return nil
	}
	return &VideoSnapSettingResponseDTO{
		ID:               setting.ID,
		ServerName:       setting.ServerName,
		Rtmp:             setting.Rtmp,
		Duration:         setting.Duration,
		VideoSize:        setting.VideoSize,
		VideoType:        setting.VideoType,
		BeforeStart:      setting.BeforeStart,
		AfterEnd:         setting.AfterEnd,
		OutputPath:       setting.OutputPath,
		BucketName:       setting.BucketName,
		AcceessDomain:    setting.AcceessDomain,
		DeleteLocalStore: setting.DeleteLocalStore,
	}
}
func (dto *VideoSnapSettingResponseDTO) FromModel(setting *VideoSnapSetting) {
	if setting == nil {
		return
	}
	dto.ID = setting.ID
	dto.ServerName = setting.ServerName
	dto.Rtmp = setting.Rtmp
	dto.Duration = setting.Duration
	dto.VideoSize = setting.VideoSize
	dto.VideoType = setting.VideoType
	dto.BeforeStart = setting.BeforeStart
	dto.AfterEnd = setting.AfterEnd
	dto.OutputPath = setting.OutputPath
	dto.ServiceAcccount = setting.ServiceAcccount
	dto.AcceessDomain = setting.AcceessDomain
	dto.BucketName = setting.BucketName
	dto.DeleteLocalStore = setting.DeleteLocalStore
	dto.CreatedAt = setting.CreatedAt
	dto.UpdatedAt = setting.UpdatedAt
}

func FromModel(setting *VideoSnapSetting) *VideoSnapSettingResponseDTO {
	if setting == nil {
		return nil
	}
	return &VideoSnapSettingResponseDTO{
		ID:               setting.ID,
		ServerName:       setting.ServerName,
		Rtmp:             setting.Rtmp,
		Duration:         setting.Duration,
		VideoSize:        setting.VideoSize,
		VideoType:        setting.VideoType,
		BeforeStart:      setting.BeforeStart,
		AfterEnd:         setting.AfterEnd,
		OutputPath:       setting.OutputPath,
		ServiceAcccount:  setting.ServiceAcccount,
		BucketName:       setting.BucketName,
		AcceessDomain:    setting.AcceessDomain,
		DeleteLocalStore: setting.DeleteLocalStore,
		CreatedAt:        setting.CreatedAt,
		UpdatedAt:        setting.UpdatedAt,
	}
}

func (dto VideoSnapSettingCreateDTO) ToModel() *VideoSnapSetting {
	return &VideoSnapSetting{
		ID:               dto.ID,
		ServerName:       dto.ServerName,
		Rtmp:             dto.Rtmp,
		Duration:         dto.Duration,
		VideoSize:        dto.VideoSize,
		VideoType:        dto.VideoType,
		BeforeStart:      dto.BeforeStart,
		AfterEnd:         dto.AfterEnd,
		OutputPath:       dto.OutputPath,
		ServiceAcccount:  dto.ServiceAcccount,
		BucketName:       dto.BucketName,
		AcceessDomain:    dto.AcceessDomain,
		DeleteLocalStore: dto.DeleteLocalStore,
	}
}

func (dto VideoSnapSettingUpdateDTO) ToModelForUpdate(existing *VideoSnapSetting) *VideoSnapSetting {
	if dto.ServerName != nil {
		existing.ServerName = *dto.ServerName
	}
	if dto.Rtmp != nil {
		existing.Rtmp = *dto.Rtmp
	}
	if dto.Duration != nil {
		existing.Duration = *dto.Duration
	}
	if dto.VideoSize != nil {
		existing.VideoSize = *dto.VideoSize
	}
	if dto.VideoType != nil {
		existing.VideoType = *dto.VideoType
	}
	if dto.BeforeStart != nil {
		existing.BeforeStart = *dto.BeforeStart
	}
	if dto.AfterEnd != nil {
		existing.AfterEnd = *dto.AfterEnd
	}
	if dto.OutputPath != nil {
		existing.OutputPath = *dto.OutputPath
	}
	if dto.ServiceAcccount != nil {
		existing.ServiceAcccount = *dto.ServiceAcccount
	}
	if dto.BucketName != nil {
		existing.BucketName = *dto.BucketName
	}
	if dto.AcceessDomain != nil {
		existing.AcceessDomain = *dto.AcceessDomain
	}
	if dto.DeleteLocalStore != nil {
		existing.DeleteLocalStore = *dto.DeleteLocalStore
	}
	return existing
}
