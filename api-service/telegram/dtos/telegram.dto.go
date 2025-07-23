package dtos

import (
	"time"
)

// TelegramGroupResponseDTO represents the response structure for a TelegramGroup.
type TelegramGroupResponseDTO struct {
	ID        uint       `json:"id"`
	GroupID   int64      `json:"group_id"`
	GroupName string     `json:"group_name"`
	GameType  string     `json:"game_type"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// TelegramGroupCreateDTO represents the request structure for creating a TelegramGroup.
type TelegramGroupCreateDTO struct {
	GroupID   int64  `json:"group_id" binding:"required"`
	GroupName string `json:"group_name" binding:"required"`
	GameType  string `json:"game_type" binding:"required"`
}

// TelegramGroupUpdateDTO represents the request structure for updating an existing TelegramGroup.
type TelegramGroupUpdateDTO struct {
	GroupID   *int64  `json:"group_id"`
	GroupName *string `json:"group_name"`
	GameType  *string `json:"game_type"`
}

// TelegramGroupFilter represents the filter structure for querying TelegramGroups.
type TelegramGroupFilter struct {
	GroupID   *int64  `json:"group_id"`
	GroupName *string `json:"group_name"`
	GameType  *string `json:"game_type"`
}
