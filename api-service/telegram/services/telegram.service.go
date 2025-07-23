package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"telegram-service/telegram/dtos"
	"telegram-service/telegram/mappers"
	"telegram-service/telegram/repository"
	"telegram-service/utils" // Assuming utils package is available and contains PaginationRequestDTO, ServicePaginationResponse, ErrorLog, InfoLog etc.

	"github.com/sirupsen/logrus" // Assuming logrus is used for logging
)

// TelegramGroupServiceInterface defines the methods for the TelegramGroupService
type TelegramGroupServiceInterface interface {
	TelegramGroupServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	TelegramGroupServiceCreate(obj []byte) (*dtos.TelegramGroupResponseDTO, error)
	TelegramGroupServiceUpdate(obj []byte) (*dtos.TelegramGroupResponseDTO, error)
	TelegramGroupServiceGetByID(id uint) (*dtos.TelegramGroupResponseDTO, error)
	TelegramGroupServiceDelete(id uint) (bool, error)
	TelegramGroupServiceGetByName(filter *dtos.TelegramGroupFilter) (*dtos.TelegramGroupResponseDTO, error)
}

// TelegramGroupService implements the TelegramGroupServiceInterface
type TelegramGroupService struct {
	repo *repository.TelegramGroupRepository
}

// NewTelegramGroupService creates a new instance of TelegramGroupService
func NewTelegramGroupService(repo *repository.TelegramGroupRepository) *TelegramGroupService {
	return &TelegramGroupService{repo: repo}
}

// TelegramGroupServiceGetList retrieves a paginated list of Telegram groups
func (s *TelegramGroupService) TelegramGroupServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}

	var filter dtos.TelegramGroupFilter
	if requestDto.Query != "" {
		if err := json.Unmarshal([]byte(requestDto.Query), &filter); err != nil {
			utils.ErrorLog(nil, fmt.Sprintf("Failed to unmarshal filter: %v", err))
			return utils.NewServicePaginationResponse(nil, 0, int(requestDto.Page), int(requestDto.Limit), http.StatusBadRequest, "Invalid filter format", logrus.ErrorLevel, "TelegramGroupService [GetList]")
		}
	}

	telegramGroups, total, err := s.repo.GetTelegramGroupList(int(requestDto.Page), int(requestDto.Limit), &filter)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error retrieving Telegram group list: %v", err))
		return utils.NewServicePaginationResponse(nil, 0, int(requestDto.Page), int(requestDto.Limit), http.StatusInternalServerError, string(utils.ErrorMessage), logrus.ErrorLevel, "TelegramGroupService [GetList]")
	}

	groupDTOs := make([]*dtos.TelegramGroupResponseDTO, len(telegramGroups))
	for i, group := range telegramGroups {
		groupDTOs[i] = mappers.ToTelegramGroupResponseDTO(group)
	}

	return utils.NewServicePaginationResponse(groupDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "TelegramGroupService [GetList]")
}

// TelegramGroupServiceCreate creates a new Telegram group
func (s *TelegramGroupService) TelegramGroupServiceCreate(obj []byte) (*dtos.TelegramGroupResponseDTO, error) {
	var createDto dtos.TelegramGroupCreateDTO
	err := json.Unmarshal(obj, &createDto)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Failed to unmarshal create DTO: %v", err))
		return nil, fmt.Errorf("invalid request payload: %w", err)
	}

	telegramGroup := mappers.ToTelegramGroupModelCreate(&createDto)
	if telegramGroup == nil {
		utils.ErrorLog(nil, "Failed to map create DTO to model: received nil model")
		return nil, fmt.Errorf("internal server error: failed to prepare data")
	}

	createdGroup, err := s.repo.CreateTelegramGroup(telegramGroup)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error creating Telegram group: %v", err))
		return nil, fmt.Errorf("failed to create Telegram group: %w", err)
	}

	responseDTO := mappers.ToTelegramGroupResponseDTO(createdGroup)
	utils.InfoLog(responseDTO, string(utils.SuccessMessage))
	return responseDTO, nil
}

// TelegramGroupServiceUpdate updates an existing Telegram group
func (s *TelegramGroupService) TelegramGroupServiceUpdate(obj []byte) (*dtos.TelegramGroupResponseDTO, error) {
	var updateDto dtos.TelegramGroupUpdateDTO
	if err := json.Unmarshal(obj, &updateDto); err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Failed to unmarshal update DTO: %v", err))
		return nil, fmt.Errorf("invalid request payload: %w", err)
	}
	type TelegramGroupUpdateRequestBody struct {
		ID        uint    `json:"id" binding:"required"` // Assuming ID is passed in the body for update
		GroupID   *int64  `json:"group_id"`
		GroupName *string `json:"group_name"`
		GameType  *string `json:"game_type"`
	}

	var updateReqBody TelegramGroupUpdateRequestBody
	if err := json.Unmarshal(obj, &updateReqBody); err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Failed to unmarshal update request body: %v", err))
		return nil, fmt.Errorf("invalid request payload: %w", err)
	}

	existingGroup, err := s.repo.GetTelegramGroupByID(updateReqBody.ID)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Telegram group with ID %d not found for update: %v", updateReqBody.ID, err))
		return nil, fmt.Errorf("telegram group not found: %w", err)
	}

	// Apply updates using the mapper
	updateDTO := dtos.TelegramGroupUpdateDTO{
		GroupID:   updateReqBody.GroupID,
		GroupName: updateReqBody.GroupName,
		GameType:  updateReqBody.GameType,
	}
	mappers.UpdateTelegramGroupModel(&updateDTO, existingGroup)

	updatedGroup, err := s.repo.UpdateTelegramGroup(existingGroup)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error updating Telegram group with ID %d: %v", updateReqBody.ID, err))
		return nil, fmt.Errorf("failed to update Telegram group: %w", err)
	}

	responseDTO := mappers.ToTelegramGroupResponseDTO(updatedGroup)
	utils.InfoLog(responseDTO, string(utils.SuccessMessage))
	return responseDTO, nil
}
func (s *TelegramGroupService) TelegramGroupServiceGetByID(id uint) (*dtos.TelegramGroupResponseDTO, error) {
	telegramGroup, err := s.repo.GetTelegramGroupByID(id)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error retrieving Telegram group by ID %d: %v", id, err))
		return nil, fmt.Errorf("telegram group not found: %w", err)
	}

	responseDTO := mappers.ToTelegramGroupResponseDTO(telegramGroup)
	utils.InfoLog(responseDTO, string(utils.SuccessMessage))
	return responseDTO, nil
}
func (s *TelegramGroupService) TelegramGroupServiceDelete(id uint) (bool, error) {
	success, err := s.repo.DeleteTelegramGroupByID(id)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error deleting Telegram group by ID %d: %v", id, err))
		return false, fmt.Errorf("failed to delete Telegram group: %w", err)
	}

	utils.InfoLog(fmt.Sprintf("Deleted Telegram Group with ID: %d", id), string(utils.SuccessMessage))
	return success, nil
}
func (s *TelegramGroupService) TelegramGroupServiceGetByGroupID(groupID int64) (*dtos.TelegramGroupResponseDTO, error) {
	telegramGroup, err := s.repo.GetTelegramGroupByGroupID(groupID)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error retrieving Telegram group by Telegram Group ID %d: %v", groupID, err))
		return nil, fmt.Errorf("telegram group with group_id %d not found: %w", groupID, err)
	}

	responseDTO := mappers.ToTelegramGroupResponseDTO(telegramGroup)
	utils.InfoLog(responseDTO, string(utils.SuccessMessage))
	return responseDTO, nil
}

func (s *TelegramGroupService) TelegramGroupServiceGetByName(filter *dtos.TelegramGroupFilter) (*dtos.TelegramGroupResponseDTO, error) {
	telegramGroup, err := s.repo.GetTelegramGroupName(filter)
	if err != nil {
		utils.ErrorLog(nil, fmt.Sprintf("Error retrieving Telegram group by Telegram Group ID %d: %v", filter, err))
		return nil, fmt.Errorf("telegram group  %d not found: %w", filter, err)
	}

	responseDTO := mappers.ToTelegramGroupResponseDTO(telegramGroup)
	utils.InfoLog(responseDTO, string(utils.SuccessMessage))
	return responseDTO, nil
}
