package mappers

import (
	"telegram-service/telegram/dtos"
	"telegram-service/telegram/models"
)

// ToTelegramGroupResponseDTO converts a models.TelegramGroup to a dtos.TelegramGroupResponseDTO.
func ToTelegramGroupResponseDTO(tg *models.TelegramGroup) *dtos.TelegramGroupResponseDTO {
	if tg == nil {
		return nil
	}
	return &dtos.TelegramGroupResponseDTO{
		ID:        tg.ID,
		GroupID:   tg.GroupID,
		GroupName: tg.GroupName,
		GameType:  tg.GameType,
		CreatedAt: tg.CreatedAt,
		UpdatedAt: tg.UpdatedAt,
	}
}

// ToTelegramGroupModelCreate converts a dtos.TelegramGroupCreateDTO to a models.TelegramGroup.
func ToTelegramGroupModelCreate(dto *dtos.TelegramGroupCreateDTO) *models.TelegramGroup {
	if dto == nil {
		return nil
	}
	return &models.TelegramGroup{
		GroupID:   dto.GroupID,
		GroupName: dto.GroupName,
		GameType:  dto.GameType,
	}
}

// UpdateTelegramGroupModel updates an existing models.TelegramGroup with values from a dtos.TelegramGroupUpdateDTO.
func UpdateTelegramGroupModel(dto *dtos.TelegramGroupUpdateDTO, tg *models.TelegramGroup) {
	if dto == nil || tg == nil {
		return // Do nothing if DTO or model is nil
	}

	if dto.GroupID != nil {
		tg.GroupID = *dto.GroupID
	}
	if dto.GroupName != nil {
		tg.GroupName = *dto.GroupName
	}
	if dto.GameType != nil {
		tg.GameType = *dto.GameType
	}
}
