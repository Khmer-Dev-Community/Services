// Role_service.go
package roles

import (
	"encoding/json"
	"fmt"
	"net/http"

	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/roles/dto"
	"github.com/Khmer-Dev-Community/Services/api-service/lib/roles/repository"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"github.com/sirupsen/logrus"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) HandleRoleListRequest(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {

	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}

	roles, total, err := s.repo.GetRoleList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())

	}

	var roleDTOs []*dto.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, &dto.RoleDTO{
			ID:          int(role.ID),
			RoleName:    role.RoleName,
			RoleStatus:  role.RoleStatus,
			RoleKey:     role.RoleKey,
			Description: role.Decription,
			Status:      role.Status,
			CompanyID:   role.CompanyID,
			Sort:        role.Sort,
		})
	}
	return utils.NewServicePaginationResponse(roleDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "")

}
func (s *RoleService) HandleRoleListRequest1(obj []byte) {
	var requestDto utils.PaginationRequestDTO
	err := json.Unmarshal(obj, &requestDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		//return nil, 0, fmt.Errorf("error retrieving user list: %v", err)
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	roles, total, err := s.repo.GetRoleList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {

		utils.ErrorLog(nil, err.Error())
	}
	// Map roles (models) to DTOs
	var roleDTOs []*dto.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, &dto.RoleDTO{
			ID:          int(role.ID),
			RoleName:    role.RoleName,
			RoleStatus:  role.RoleStatus,
			RoleKey:     role.RoleKey,
			Description: role.Decription,
			Status:      role.Status,
			Sort:        role.Sort,
			CompanyID:   role.CompanyID,
		})
	}
	utils.NewServicePaginationResponse(roleDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, "Success", logrus.InfoLevel, "RoleService [HandleRoleListRequest]")

}
func (us *RoleService) HandleRoleCreateRequest(obj []byte) (*dto.RoleDTO, error) {
	var createRoleRequest dto.RoleCreateDTO
	if err := json.Unmarshal(obj, &createRoleRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	createdRole, err := us.repo.CreateRole(&createRoleRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	roleDTO := &dto.RoleDTO{
		RoleName:    createdRole.RoleName,
		RoleStatus:  createdRole.RoleStatus,
		RoleKey:     createdRole.RoleKey,
		Description: createdRole.Decription,
		Status:      createdRole.Status,
		Sort:        createdRole.Sort,
		CompanyID:   createdRole.CompanyID,
	}
	utils.InfoLog(roleDTO, string(utils.SuccessMessage))
	return roleDTO, nil

}

func (us *RoleService) HandleRoleUpdateRequest(obj []byte) (*dto.RoleDTO, error) {
	var updateRoleRequest dto.RoleUpdateDTO
	if err := json.Unmarshal(obj, &updateRoleRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	RoleID := uint(updateRoleRequest.ID)
	updatedRole, err := us.repo.UpdateRole(RoleID, &updateRoleRequest)
	if err != nil {
		utils.ErrorLog(nil, "Failed to update role: "+err.Error())
		return nil, err
	}

	roleDTO := &dto.RoleDTO{
		ID:          int(updatedRole.ID),
		RoleName:    updatedRole.RoleName,
		RoleStatus:  updatedRole.RoleStatus,
		RoleKey:     updatedRole.RoleKey,
		Description: updatedRole.Decription,
		Status:      updatedRole.Status,
		Sort:        updatedRole.Sort,
		CompanyID:   updatedRole.CompanyID,
	}

	utils.InfoLog(roleDTO, string(utils.SuccessMessage))
	return roleDTO, nil
}

func (us *RoleService) HandleRoleByIdRequest(RoleID int) (*dto.RoleDTO, error) {

	RoleFiller, err := us.repo.GetRoleByID(uint(RoleID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	roleDTO := &dto.RoleDTO{
		ID:          int(RoleFiller.ID),
		RoleName:    RoleFiller.RoleName,
		RoleStatus:  RoleFiller.RoleStatus,
		RoleKey:     RoleFiller.RoleKey,
		Description: RoleFiller.Decription,
		Status:      RoleFiller.Status,
		Sort:        RoleFiller.Sort,
		CompanyID:   RoleFiller.CompanyID,
	}

	utils.InfoLog(roleDTO, string(utils.SuccessMessage))
	return roleDTO, nil

}

func (us *RoleService) HandleRoleDeleteRequest(obj []byte) (*dto.RoleDTO, error) {
	var deleteRoleRequest dto.RoleDTO
	if err := json.Unmarshal(obj, &deleteRoleRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	err := us.repo.DeleteRoleByID(uint(deleteRoleRequest.ID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	utils.InfoLog(nil, string(utils.SuccessMessage))
	return nil, nil
}
