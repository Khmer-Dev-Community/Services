package clientfollow

// ToNotificationDTO converts a Notification to a Notification
func ToNotificationDTO(n Notification) NotificationDTO {
	return NotificationDTO{
		ID:         n.ID,
		FromUserID: n.FromUserID,
		PostID:     n.PostID,
		Type:       string(n.Type), // Convert NotificationType to a string
		Message:    n.Message,
		IsRead:     n.IsRead,
		CreatedAt:  n.CreatedAt,
	}
}

// ToNotificationDTOs converts a slice of Notification to a slice of NotificationDTOs.
func ToNotificationDTOs(notifications []Notification) []NotificationDTO {
	dtos := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = ToNotificationDTO(n)
	}
	return dtos
}

// ToAccountFollowDTO converts a AccountFollow to a AccountFollow
func ToAccountFollowDTO(a AccountFollow) AccountFollowDTO {
	return AccountFollowDTO{
		FollowerID: a.FollowerID,
		FollowedID: a.FollowedID,
		CreatedAt:  a.CreatedAt,
	}
}

// ToAccountFollowDTOs converts a slice of AccountFollow to a slice of AccountFollowDTOs.
func ToAccountFollowDTOs(follows []AccountFollow) []AccountFollowDTO {
	dtos := make([]AccountFollowDTO, len(follows))
	for i, f := range follows {
		dtos[i] = ToAccountFollowDTO(f)
	}
	return dtos
}

// ToPostFollowDTO converts a PostFollow to a PostFollow
func ToPostFollowDTO(p PostFollow) PostFollowDTO {
	return PostFollowDTO{
		UserID:    p.UserID,
		PostID:    p.PostID,
		CreatedAt: p.CreatedAt,
	}
}

// ToPostFollowDTOs converts a slice of PostFollow to a slice of PostFollowDTOs.
func ToPostFollowDTOs(follows []PostFollow) []PostFollowDTO {
	dtos := make([]PostFollowDTO, len(follows))
	for i, f := range follows {
		dtos[i] = ToPostFollowDTO(f)
	}
	return dtos
}
