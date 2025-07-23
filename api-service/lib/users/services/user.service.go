package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"telegram-service/lib/users/dto"
	"telegram-service/lib/users/repository"
	"telegram-service/utils"

	"github.com/sirupsen/logrus"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserList1(page int, limit int, query string) ([]*dto.UserDTO, int, error) {
	log.Print("UserService Page, v%", page)
	log.Print("UserService Page, v%", limit)
	users, total, err := s.repo.GetUserList(page, limit, query)

	if err != nil {
		log.Printf("Error retrieving user list: %v", err)
		return nil, 0, fmt.Errorf("error retrieving user list: %v", err)
	}

	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		if user == nil {
			log.Printf("User is nil at index %d", i)
			continue // or handle appropriately (e.g., return an error)
		}
		userDTO := &dto.UserDTO{}
		userDTO.FromModel(user) // Use method to map model to DTO
		userDTOs[i] = userDTO
	}

	return userDTOs, total, nil
}

func (s *UserService) GetUserList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {

	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}

	users, total, err := s.repo.GetUserList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())

	}

	userDTOs := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		userDTO := &dto.UserDTO{}
		userDTO.FromModel(user) // Use method to map model to DTO
		userDTOs[i] = userDTO
	}
	//utils.Logger(userDTOs, http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "UserService")
	return utils.NewServicePaginationResponse(userDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "UserService [GetUserList]")

}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(userID uint) (*dto.UserDTO, error) {
	user, err := s.repo.FindOneById(userID)
	if err != nil {
		utils.WarnLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)

	}
	userDTO := &dto.UserDTO{}
	userDTO.FromModel(user)
	utils.InfoLog(userDTO, string(utils.SuccessMessage))
	return userDTO, nil
}
func (s *UserService) CreateUser(obj []byte) (*dto.UserDTO, error) {
	var createUserRequest dto.UserCreateDTO
	if err := json.Unmarshal(obj, &createUserRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	createdUser, err := s.repo.CreateUser(&createUserRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	userDTO := &dto.UserDTO{}
	userDTO.FromModel(createdUser)

	utils.InfoLog(userDTO, string(utils.SuccessMessage))
	return userDTO, nil
}

// UserService.go
func (s *UserService) UpdateUser(userID uint, obj []byte) (*dto.UserDTO, error) {
	var updateUserRequest dto.UserUpdateDTO
	if err := json.Unmarshal(obj, &updateUserRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedUser, err := s.repo.UpdateUser(userID, &updateUserRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	userDTO := &dto.UserDTO{}
	userDTO.FromModel(updatedUser)

	utils.InfoLog(userDTO, string(utils.SuccessMessage))
	return userDTO, nil
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(userID uint) (*dto.UserDTO, error) {
	err := s.repo.DeleteUserByID(userID)
	if err != nil {
		log.Printf("Error deleting user with ID %d: %v", userID, err)
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	utils.InfoLog(nil, string(utils.SuccessMessage))
	return nil, nil
}
