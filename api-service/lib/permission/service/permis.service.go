// Permission_service.go
package permissions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	redis "github.com/Khmer-Dev-Community/Services/api-service/config"
	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/permission/dto"
	repository "github.com/Khmer-Dev-Community/Services/api-service/lib/permission/repository"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"github.com/sirupsen/logrus"
)

type PermissionService struct {
	repo *repository.PermissionRepository
	// PermissionService *PermissionService
}

func NewPermissionService(repo *repository.PermissionRepository) *PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) HandlePermissionListRequest1(page int, limit int, query string) ([]*dto.PermissionDTO, int, error) {
	permissions, total, err := s.repo.GetPermissionList(page, limit, query)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving permission list: %v", err)
	}

	permisDTOs := make([]*dto.PermissionDTO, len(permissions))
	for i, permission := range permissions {
		permisDTO := &dto.PermissionDTO{}
		permisDTO.FromModel(permission)
		permisDTOs[i] = permisDTO
	}
	return permisDTOs, total, nil
}
func (s *PermissionService) HandlePermissionListRequest(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}

	permiss, total, err := s.repo.GetPermissionList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())

	}

	permission_ := make([]*dto.PermissionDTO, len(permiss))
	for i, obj := range permiss {
		permission := &dto.PermissionDTO{}
		permission.FromModel(obj) // Use method to map model to DTO
		permission_[i] = permission
	}
	return utils.NewServicePaginationResponse(permission_, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "")
}
func (s *PermissionService) HandlePermissionCreateRequest(obj []byte) (*dto.PermissionDTO, error) {
	var createRoleRequest dto.PermissionCreateDTO
	if err := json.Unmarshal(obj, &createRoleRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	createdPermission, err := s.repo.CreatePermission(&createRoleRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	permisDTO := &dto.PermissionDTO{}
	permisDTO.FromModel(createdPermission)

	return permisDTO, nil
}
func (s *PermissionService) HandlePermissionUpdateRequest(permissID uint, obj []byte) (*dto.PermissionDTO, error) {
	var updateRoleDto dto.PermissionUpdateDTO
	if err := json.Unmarshal(obj, &updateRoleDto); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedPermission, err := s.repo.UpdatePermission(permissID, &updateRoleDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	permisDTO := &dto.PermissionDTO{}
	permisDTO.FromModel(updatedPermission)
	return permisDTO, nil
}
func (s *PermissionService) HandlePermissionByIdRequest(PermissionID uint) (*dto.PermissionDTO, error) {
	permis, err := s.repo.GetPermissionByID(PermissionID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	permisDTO := &dto.PermissionDTO{}
	permisDTO.FromModel(permis)

	return permisDTO, nil
}
func (us *PermissionService) HandleRolePermissionById(roleID int) (*dto.RolePermissionsDTO, error) {
	utils.InfoLog(roleID, fmt.Sprintf("Fetching permissions for role ID: %d", roleID))
	permissionStrings, err := us.repo.GetPermissionsForRole(roleID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to get permissions: %s", err)
	}
	// Fetch role associated with the role by role id
	roleNames, err := us.repo.GetRoleNames(roleID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to get role names: %s", err)
	}
	// Ensure we have at least one role name
	if len(roleNames) == 0 {
		utils.ErrorLog(roleNames, "no role names found for role id : %s")
		return nil, fmt.Errorf("no role names found for role ID: %d", roleID)
	}
	// Construct the RolePermissionsDTO
	rolePermissionsDTO := &dto.RolePermissionsDTO{
		Permissions: permissionStrings,
		RoleName:    roleNames,
	}
	// 20 nimutes expiration will be return reqtest status 401
	// client auto logout
	expiration := time.Minute * 0
	key := fmt.Sprintf("role:%d", roleID)
	if err := redis.SetWithExpiration(key, rolePermissionsDTO, expiration); err != nil {
		utils.ErrorLog("Failed Redis Store Role Permission", err.Error())
		return nil, fmt.Errorf("failed redis store Role permission: %s", err)
	}
	return rolePermissionsDTO, nil
}
func (us *PermissionService) UpdateRolePermissionByRoleID(msgBody []byte) (*dto.PermissionDTO, error) {
	utils.InfoLog(msgBody, "UpdateRolePermissionByRoleID")
	// Define the request data structure
	var updateRolePermission dto.UpdateRolePermissionsMessage

	// Unmarshal the JSON message into requestData struct
	if err := json.Unmarshal(msgBody, &updateRolePermission); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	err := us.repo.UpdateUserPermissions(updateRolePermission.RoleID, updateRolePermission.MenuIDs)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	return nil, nil
}

func (us *PermissionService) HandlePermissionDeleteRequest(obj []byte) (*dto.PermissionDTO, error) {
	var deleteRequest dto.PermissionDTO
	if err := json.Unmarshal(obj, &deleteRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	err := us.repo.DeletePermissionByID(uint(deleteRequest.ID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	return nil, nil

}
