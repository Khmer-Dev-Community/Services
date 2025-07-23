package setting

import (
	"encoding/json"
	"fmt"
	"net/http"
	"telegram-service/utils"

	redis "telegram-service/config"

	"github.com/sirupsen/logrus"
)

// VideoSnapSettingServiceInterface defines the methods for the VideoSnapSettingService
type VideoSnapSettingServiceInterface interface {
	VideoSnapSettingServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	VideoSnapSettingServiceCreate(obj []byte) (*VideoSnapSettingResponseDTO, error)
	VideoSnapSettingServiceUpdate(obj []byte) (*VideoSnapSettingResponseDTO, error)
	VideoSnapSettingServiceGetById(obj []byte) (*VideoSnapSettingResponseDTO, error)
	VideoSnapSettingServiceDelete(obj []byte) (bool, error)
}

type VideoSnapSettingService struct {
	repo *VideoSnapSettingRepository
}

// NewVideoSnapSettingService initializes a new VideoSnapSettingService
func NewVideoSnapSettingService(repo *VideoSnapSettingRepository) *VideoSnapSettingService {
	return &VideoSnapSettingService{repo: repo}
}

// VideoSnapSettingServiceGetList retrieves a paginated list of video snap settings
func (vss *VideoSnapSettingService) VideoSnapSettingServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1 // Default to first page if invalid
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10 // Default limit if invalid
	}
	videoSnapSettings, total, err := vss.repo.GetVideoSnapSettingList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "VideoSnapSettingService [GetVideoSnapSettingList]")
	}
	videoSnapSettingDTOs := make([]*VideoSnapSettingResponseDTO, len(videoSnapSettings))
	for i, setting := range videoSnapSettings {
		settingDTO := &VideoSnapSettingResponseDTO{}
		settingDTO.FromModel(setting)
		videoSnapSettingDTOs[i] = settingDTO
	}
	return utils.NewServicePaginationResponse(videoSnapSettingDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "VideoSnapSettingService [GetVideoSnapSettingList]")
}

// VideoSnapSettingServiceCreate creates a new video snap setting
func (vss *VideoSnapSettingService) VideoSnapSettingServiceCreate(obj []byte) (*VideoSnapSettingResponseDTO, error) {
	var createDto VideoSnapSettingCreateDTO
	err := json.Unmarshal(obj, &createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	createdSetting, err := vss.repo.CreateVideoSnapSetting(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	redisKey := fmt.Sprintf("videoSetting:server:%d", createdSetting.ID)
	// Store in Redis without expiration
	if err := redis.SetWithoutExpired(redisKey, createdSetting); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to store setting in Redis: %v", err)
	}
	settingDTO := &VideoSnapSettingResponseDTO{}
	settingDTO.FromModel(createdSetting)
	utils.InfoLog(settingDTO, string(utils.SuccessMessage))
	return settingDTO, nil
}

// VideoSnapSettingServiceUpdate updates an existing video snap setting
func (vss *VideoSnapSettingService) VideoSnapSettingServiceUpdate(obj []byte) (*VideoSnapSettingResponseDTO, error) {
	var updateDto VideoSnapSettingUpdateDTO
	if err := json.Unmarshal(obj, &updateDto); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	utils.InfoLog(updateDto, "Data to be stored in Redis")
	redisKey := fmt.Sprintf("videoSetting:server:%d", updateDto.ID)
	if err := redis.SetWithoutExpired(redisKey, updateDto); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to store setting in Redis: %v", err)
	}

	updatedSetting, err := vss.repo.UpdateVideoSnapSetting(updateDto.ID, updateDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto := &VideoSnapSettingResponseDTO{}
	updatedDto.FromModel(updatedSetting)
	utils.InfoLog(updatedDto, string(utils.SuccessMessage))
	return updatedDto, nil
}

// VideoSnapSettingServiceGetById retrieves a video snap setting by ID
func (vss *VideoSnapSettingService) VideoSnapSettingServiceGetById(obj []byte) (*VideoSnapSettingResponseDTO, error) {
	var settingDTO VideoSnapSettingUpdateDTO
	if err := json.Unmarshal(obj, &settingDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	videoSnapSetting, err := vss.repo.GetVideoSnapSettingByID(settingDTO.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	resultDto := &VideoSnapSettingResponseDTO{}
	resultDto.FromModel(videoSnapSetting)
	utils.InfoLog(resultDto, string(utils.SuccessMessage))
	return resultDto, nil
}

// VideoSnapSettingServiceDelete deletes a video snap setting by ID
func (vss *VideoSnapSettingService) VideoSnapSettingServiceDelete(obj []byte) (bool, error) {
	var settingDTO VideoSnapSettingResponseDTO
	if err := json.Unmarshal(obj, &settingDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}
	success, err := vss.repo.DeleteVideoSnapSettingByID(settingDTO.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}

	utils.InfoLog(settingDTO, string(utils.SuccessMessage))
	return success, nil
}
