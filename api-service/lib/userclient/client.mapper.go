package userclient

// ToClientUserResponseDTO maps a ClientUser model to a ClientUserResponseDTO.
// This is used when sending user data to the client, ensuring sensitive information is excluded.
func ToClientUserResponseDTO(user *ClientUser) *ClientUserResponseDTO {
	if user == nil {
		return nil
	}
	return &ClientUserResponseDTO{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
		Sex:         user.Sex,
		Status:      user.Status,
		Description: user.Description,
		GitHubID:    user.GitHubID,
		AvatarURL:   user.AvatarURL,
		Name:        user.Name,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// ToClientUserFromRegisterDTO maps a ClientRegisterRequestDTO to a ClientUser model.
// This is used during new user registration. Password hashing should happen after this mapping.
func ToClientUserFromRegisterDTO(dto *ClientRegisterRequestDTO) *ClientUser {
	if dto == nil {
		return nil
	}
	return &ClientUser{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Username:    dto.Username,
		Email:       dto.Email,
		Phone:       dto.Phone,
		Sex:         dto.Sex,
		Description: dto.Description,
		Status:      true, // Default status
	}
}

// UpdateClientUserFromDTO updates an existing ClientUser model with fields from ClientUpdateRequestDTO.
// Only non-nil (provided) fields in the DTO will update the model.
func UpdateClientUserFromDTO(user *ClientUser, dto *ClientUpdateRequestDTO) {
	if user == nil || dto == nil {
		return
	}

	if dto.FirstName != nil {
		user.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Phone != nil {
		user.Phone = *dto.Phone
	}
	if dto.Sex != nil {
		user.Sex = *dto.Sex
	}
	if dto.Status != nil {
		user.Status = *dto.Status
	}
	if dto.Description != nil {
		user.Description = *dto.Description
	}
}
