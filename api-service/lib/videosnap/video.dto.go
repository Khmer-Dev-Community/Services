package videosnap

import "time"

type VideosnapStartDTO struct {
	GameNo    string `json:"period" binding:"required"`  // Game number
	Streamkey string `json:"rtmpurl" binding:"required"` // Stream key

}
type VideoSnapTaskResponseDTO struct {
	ImageURL string `json:"image_url" ` // Image URL
	VideoURL string `json:"video_url" ` // Video URL
	Status   int    `json:"status"`     // Status
}
type VideosnapCreateDTO struct {
	ProcessId uint   `json:"process_id" `               // Associated process ID
	Code      string `json:"tran_id"`                   // Transaction code
	GameNo    string `json:"period" binding:"required"` // Game number
	Streamkey string `json:"rtmpurl"`                   // Stream key
	Rtmp      string `json:"rtmp" `                     // RTMP URL
	ImageURL  string `json:"image_url" `                // Image URL
	VideoURL  string `json:"video_url" `                // Video URL
	StorePath string `json:"store_path"`                // Store path
	Status    bool   `json:"status"`                    // Status of the videosnap
}
type VideosnapUpdateDTO struct {
	ID        uint    `json:"id" binding:"required"` // ID of the videosnap
	ProcessId *uint   `json:"process_id,omitempty"`  // Optional: Associated process ID
	Code      *string `json:"tran_id,omitempty"`     // Optional: Transaction code
	GameNo    *string `json:"period,omitempty"`      // Optional: Game number
	Streamkey *string `json:"rtmpurl" `              // Stream key
	Rtmp      *string `json:"rtmp,omitempty"`        // Optional: RTMP URL
	ImageURL  *string `json:"image_url,omitempty"`   // Optional: Image URL
	VideoURL  *string `json:"video_url,omitempty"`   // Optional: Video URL
	StorePath *string `json:"store_path,omitempty"`  // Store path
	Status    *bool   `json:"status,omitempty"`      // Optional: Status of the videosnap
}
type VideosnapResponseDTO struct {
	ID        uint      `json:"id"`         // ID of the videosnap
	ProcessId uint      `json:"process_id"` // Associated process ID
	Code      string    `json:"tran_id"`    // Transaction code
	GameNo    string    `json:"period"`     // Game number
	Rtmp      string    `json:"rtmp"`       // RTMP URL
	Streamkey string    `json:"rtmpurl" `
	ImageURL  string    `json:"image_url"`  // Image URL
	VideoURL  string    `json:"video_url"`  // Video URL
	StorePath string    `json:"store_path"` // Store path
	Status    bool      `json:"status"`     // Status of the videosnap
	CreatedAt time.Time `json:"created_at"` // Created timestamp
	UpdatedAt time.Time `json:"updated_at"` // Updated timestamp
}

func (dto *VideosnapResponseDTO) FromModel(videosnap *Videosnap) {
	if videosnap == nil {
		return
	}
	dto.ID = videosnap.ID
	dto.ProcessId = videosnap.ProcessId
	dto.Code = videosnap.Code
	dto.GameNo = videosnap.GameNo
	dto.Rtmp = videosnap.Rtmp
	dto.Streamkey = videosnap.Streamkey
	dto.ImageURL = videosnap.ImageURL
	dto.VideoURL = videosnap.VideoURL
	dto.StorePath = videosnap.StorePath
	dto.Status = videosnap.Status
	dto.CreatedAt = videosnap.CreatedAt
	dto.UpdatedAt = videosnap.UpdatedAt
}
func FromModel(videosnap *Videosnap) *VideosnapResponseDTO {
	if videosnap == nil {
		return nil
	}
	return &VideosnapResponseDTO{
		ID:        videosnap.ID,
		ProcessId: videosnap.ProcessId,
		Code:      videosnap.Code,
		GameNo:    videosnap.GameNo,
		Rtmp:      videosnap.Rtmp,
		Streamkey: videosnap.Streamkey,
		ImageURL:  videosnap.ImageURL,
		VideoURL:  videosnap.VideoURL,
		StorePath: videosnap.StorePath,
		Status:    videosnap.Status,
		CreatedAt: videosnap.CreatedAt,
		UpdatedAt: videosnap.UpdatedAt,
	}
}
func (dto VideosnapCreateDTO) ToModel() *Videosnap {
	return &Videosnap{
		ProcessId: dto.ProcessId,
		Code:      dto.Code,
		GameNo:    dto.GameNo,
		Rtmp:      dto.Rtmp,
		Streamkey: dto.Streamkey,
		ImageURL:  dto.ImageURL,
		VideoURL:  dto.VideoURL,
		StorePath: dto.StorePath,
		Status:    dto.Status,
	}
}

func (dto VideosnapUpdateDTO) ToModelForUpdate(existing *Videosnap) *Videosnap {
	if dto.ProcessId != nil {
		existing.ProcessId = *dto.ProcessId
	}
	if dto.Code != nil {
		existing.Code = *dto.Code
	}
	if dto.GameNo != nil {
		existing.GameNo = *dto.GameNo
	}
	if dto.Rtmp != nil {
		existing.Rtmp = *dto.Rtmp
	}
	if dto.Streamkey != nil {
		existing.Streamkey = *dto.Streamkey
	}
	if dto.ImageURL != nil {
		existing.ImageURL = *dto.ImageURL
	}
	if dto.VideoURL != nil {
		existing.VideoURL = *dto.VideoURL
	}
	if dto.StorePath != nil {
		existing.StorePath = *dto.StorePath
	}
	if dto.Status != nil {
		existing.Status = *dto.Status
	}
	return existing
}
