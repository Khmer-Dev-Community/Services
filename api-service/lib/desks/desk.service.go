package desks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"telegram-service/utils"

	redis "telegram-service/config"

	"github.com/sirupsen/logrus"
)

// DeskSettingServiceInterface defines the methods for the DeskSettingService
type DeskSettingServiceInterface interface {
	InitializeDeskRedis() error
	DeskSettingServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	DeskSettingServiceCreate(obj []byte) (*DeskSettingResponseDTO, error)
	DeskSettingServiceUpdate(obj []byte) (*DeskSettingResponseDTO, error)
	DeskSettingServiceGetById(obj []byte) (*DeskSettingResponseDTO, error)
	DeskSettingServiceDelete(obj []byte) (bool, error)
}

type DeskSettingService struct {
	repo *DeskSettingRepository
}

// InitializeDesk To Redis
// InitializeDeskRedis initializes desk settings in Redis
func (dss *DeskSettingService) InitializeDeskRedis() error {
	deskSettings, err := dss.repo.IntDeskSetting()

	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return err // Return the error for further handling
	}
	utils.InfoLog(fmt.Sprintf("Initialize desk Data :%v", deskSettings), string(utils.SuccessMessage))
	for _, deskSetting := range deskSettings {
		redisKey := fmt.Sprintf("deskSetting:table:%s", deskSetting.DeskStreamKey)
		if err := redis.SetWithoutExpired(redisKey, deskSetting); err != nil {
			utils.ErrorLog(nil, err.Error())
			return err // Return the error if Redis operation fails
		}
	}
	utils.InfoLog("Initialize desk settings in Redis Success", string(utils.SuccessMessage))
	return nil // Return nil if all operations are successful
}

// NewDeskSettingService initializes a new DeskSettingService
func NewDeskSettingService(repo *DeskSettingRepository) *DeskSettingService {
	service := &DeskSettingService{repo: repo}
	// Initialize Redis data
	if err := service.InitializeDeskRedis(); err != nil {
		utils.ErrorLog(nil, "Failed to initialize desk settings in Redis: "+err.Error())
	}
	return service
	//return &DeskSettingService{repo: repo}
}

// DeskSettingServiceGetList retrieves a paginated list of desk settings
func (dss *DeskSettingService) DeskSettingServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1 // Default to first page if invalid
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10 // Default limit if invalid
	}
	deskSettings, total, err := dss.repo.GetDeskSettingList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "DeskSettingService [GetDeskSettingList]")
	}
	deskSettingDTOs := make([]*DeskSettingResponseDTO, len(deskSettings))
	for i, setting := range deskSettings {
		settingDTO := &DeskSettingResponseDTO{}
		settingDTO.FromModel(setting)
		deskSettingDTOs[i] = settingDTO
	}
	if err := dss.InitializeDeskRedis(); err != nil {
		utils.ErrorLog(nil, "Failed to initialize desk settings in Redis: "+err.Error())
	}
	return utils.NewServicePaginationResponse(deskSettingDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "DeskSettingService [GetDeskSettingList]")
}

// DeskSettingServiceCreate creates a new desk setting
func (dss *DeskSettingService) DeskSettingServiceCreate(obj []byte) (*DeskSettingResponseDTO, error) {
	var createDto DeskSettingCreateDTO
	err := json.Unmarshal(obj, &createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	createdSetting, err := dss.repo.CreateDeskSetting(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	redisKey := fmt.Sprintf("deskSetting:table:%s", createdSetting.DeskStreamKey)
	// Store in Redis without expiration
	if err := redis.SetWithoutExpired(redisKey, createdSetting); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to store setting in Redis: %v", err)
	}
	settingDTO := &DeskSettingResponseDTO{}
	settingDTO.FromModel(createdSetting)
	utils.InfoLog(settingDTO, string(utils.SuccessMessage))
	return settingDTO, nil
}

// DeskSettingServiceUpdate updates an existing desk setting
func (dss *DeskSettingService) DeskSettingServiceUpdate(obj []byte) (*DeskSettingResponseDTO, error) {
	var updateDto DeskSettingUpdateDTO
	if err := json.Unmarshal(obj, &updateDto); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	utils.InfoLog(updateDto, "Data to be stored in Redis")
	redisKey := fmt.Sprintf("deskSetting:table:%d", updateDto.ID)
	if err := redis.SetWithoutExpired(redisKey, updateDto); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to store setting in Redis: %v", err)
	}

	updatedSetting, err := dss.repo.UpdateDeskSetting(updateDto.ID, updateDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto := &DeskSettingResponseDTO{}
	updatedDto.FromModel(updatedSetting)
	utils.InfoLog(updatedDto, string(utils.SuccessMessage))
	return updatedDto, nil
}

// DeskSettingServiceGetById retrieves a desk setting by ID
func (dss *DeskSettingService) DeskSettingServiceGetById(obj []byte) (*DeskSettingResponseDTO, error) {
	var settingDTO DeskSettingUpdateDTO
	if err := json.Unmarshal(obj, &settingDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	deskSetting, err := dss.repo.GetDeskSettingByID(settingDTO.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	resultDto := &DeskSettingResponseDTO{}
	resultDto.FromModel(deskSetting)
	utils.InfoLog(resultDto, string(utils.SuccessMessage))
	return resultDto, nil
}

// DeskSettingServiceDelete deletes a desk setting by ID
func (dss *DeskSettingService) DeskSettingServiceDelete(obj []byte) (bool, error) {
	var settingDTO DeskSettingResponseDTO
	if err := json.Unmarshal(obj, &settingDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}
	success, err := dss.repo.DeleteDeskSettingByID(settingDTO.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}

	utils.InfoLog(settingDTO, string(utils.SuccessMessage))
	return success, nil
}
