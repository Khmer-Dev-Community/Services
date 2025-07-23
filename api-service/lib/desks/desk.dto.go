package desks

import (
	"telegram-service/lib/videosnap/setting"
	"time"
)

// DeskSettingCreateDTO is the DTO for creating a new DeskSetting.
type DeskSettingCreateDTO struct {
	DeskNo           string `json:"desk_no" binding:"required"`            // Desk number
	DeskStreamKey    string `json:"desk_streamkey" binding:"required"`     // Desk stream key
	DeskStreamServer uint   `json:"desk_stream_server" binding:"required"` // Desk stream server
}

// DeskSettingUpdateDTO is the DTO for updating an existing DeskSetting.
type DeskSettingUpdateDTO struct {
	ID               uint    `json:"id" binding:"required"`        // ID of the DeskSetting
	DeskNo           *string `json:"desk_no,omitempty"`            // Optional: Desk number
	DeskStreamKey    *string `json:"desk_streamkey,omitempty"`     // Optional: Desk stream key
	DeskStreamServer *uint   `json:"desk_stream_server,omitempty"` // Optional: Desk stream server
}

// DeskSettingResponseDTO is the DTO for the response of DeskSetting.
type DeskSettingResponseDTO struct {
	ID               uint                                `json:"id"`                 // ID of the DeskSetting
	DeskNo           string                              `json:"desk_no"`            // Desk number
	DeskStreamKey    string                              `json:"desk_streamkey"`     // Desk stream key
	DeskStreamServer uint                                `json:"desk_stream_server"` // Desk stream server
	StreamServer     setting.VideoSnapSettingResponseDTO `json:"stream_server"`
	CreatedAt        time.Time                           `json:"created_at"` // Created timestamp
	UpdatedAt        time.Time                           `json:"updated_at"` // Updated timestamp
}

// FromModel populates DeskSettingResponseDTO from DeskSetting model.
func (dto *DeskSettingResponseDTO) FromModel(setting *DeskSetting) {
	if setting == nil {
		return
	}
	dto.ID = setting.ID
	dto.DeskNo = setting.DeskNo
	dto.DeskStreamKey = setting.DeskStreamKey
	dto.DeskStreamServer = setting.DeskStreamServer
	dto.StreamServer = *setting.StreamServer.FromModelX()
	dto.CreatedAt = setting.CreatedAt
	dto.UpdatedAt = setting.UpdatedAt
}

// FromModel creates a new DeskSettingResponseDTO from DeskSetting model.
func FromModel(setting *DeskSetting) *DeskSettingResponseDTO {
	if setting == nil {
		return nil
	}
	return &DeskSettingResponseDTO{
		ID:               setting.ID,
		DeskNo:           setting.DeskNo,
		DeskStreamKey:    setting.DeskStreamKey,
		DeskStreamServer: setting.DeskStreamServer,
		StreamServer:     *setting.StreamServer.FromModelX(),
		CreatedAt:        setting.CreatedAt,
		UpdatedAt:        setting.UpdatedAt,
	}
}

// ToModel converts DeskSettingCreateDTO to DeskSetting model.
func (dto DeskSettingCreateDTO) ToModel() *DeskSetting {
	return &DeskSetting{
		DeskNo:           dto.DeskNo,
		DeskStreamKey:    dto.DeskStreamKey,
		DeskStreamServer: dto.DeskStreamServer,
	}
}

// ToModelForUpdate updates an existing DeskSetting model with values from DeskSettingUpdateDTO.
func (dto DeskSettingUpdateDTO) ToModelForUpdate(existing *DeskSetting) *DeskSetting {
	if dto.DeskNo != nil {
		existing.DeskNo = *dto.DeskNo
	}
	if dto.DeskStreamKey != nil {
		existing.DeskStreamKey = *dto.DeskStreamKey
	}
	if dto.DeskStreamServer != nil {
		existing.DeskStreamServer = *dto.DeskStreamServer
	}
	return existing
}
