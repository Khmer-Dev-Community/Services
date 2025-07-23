// Menu_service.go
package menu

import (
	"encoding/json"
	"fmt"
	"net/http"
	dto "telegram-service/lib/menus/dto"
	repository "telegram-service/lib/menus/repository"
	"telegram-service/utils"
	"time"

	"github.com/sirupsen/logrus"
)

type MenuServiceInterface interface {
	GetMenuTree(roleId int) (*[]dto.MenuTree, error)
	HandleMenuListSub(roleId int) (*[]dto.MenuTree, error)
	HandleSelectMenusByTenantRoleID(roleId int) (*[]dto.MenuTree, error)
	HandleMenuListRequest(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	HandleMenuCreateRequest(msgBody []byte) (*dto.MenuDTO, error)
	HandleMenuUpdateRequest(msgBody []byte) (*dto.MenuDTO, error)
	HandleMenuDeleteRequest(msgBody []byte) (bool, error)
	HandleMenuListWithChild() (*[]dto.MenuTreeChild, error)
	HandleMenuListWithChildByID(menuId int) (*[]dto.MenuTreeChild, error)
}
type MenuService struct {
	repo *repository.MenuRepository
}

var _ MenuServiceInterface = (*MenuService)(nil)

func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (ms *MenuService) GetMenuTree(RoleId int) (*[]dto.MenuTree, error) {
	return ms.repo.SelectMenus(RoleId)
}

func (ms *MenuService) HandleSelectMenusByTenantRoleID(roleId int) (*[]dto.MenuTree, error) {
	utils.InfoLog(roleId, "HandleSelectMenusByTenantRoleID")
	menuPermissions, err := ms.repo.SelectMenusByTenantRoleID(int(roleId))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to fetch menu permissions: %v", err)
	}
	utils.InfoLog(menuPermissions, string(utils.ResponseMessage))

	return menuPermissions, nil
}

func (ms *MenuService) HandleMenuListSub(roleId int) (*[]dto.MenuTree, error) {
	utils.InfoLog(roleId, "Fetching menu permissions")
	menuPermissions, err := ms.repo.SelectMenus(int(roleId))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to fetch menu permissions: %v", err)
	}

	return menuPermissions, nil
}

func (us *MenuService) HandleMenuListRequest(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	utils.InfoLog(requestDto, "HandleMenuListRequest")
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	permiss, total, err := us.repo.GetMenuList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())

	}
	permission_ := make([]*dto.MenuDTO, len(permiss))
	for i, obj := range permiss {
		permission := &dto.MenuDTO{}
		permission.FromModel(obj) // Use method to map model to DTO
		permission_[i] = permission
	}
	return utils.NewServicePaginationResponse(permission_, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "")
}
func (us *MenuService) HandleMenuListWithChildToPermissionAuthorize(roleId int) (*dto.MenuTreeselectDTO, error) {
	utils.InfoLog(roleId, "HandleMenuListWithChildToPermissionAuthorize Role ID")
	MenusList, err := us.repo.SelectMenuWithChildrenToPermissionAuthorize()
	if err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
	}

	checkedKeys, err := us.repo.GetPermissionMenuIDs(int(roleId))
	if err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err
	}
	responseData := &dto.MenuTreeselectDTO{
		Menu:        *MenusList,
		CheckedKeys: checkedKeys,
	}
	utils.InfoLog(responseData, string(utils.SuccessMessage))
	return responseData, nil
	//retrun utils.ServiceResponse(responseData, http.StatusOK, "Success", logrus.InfoLevel)

}
func (us *MenuService) HandleMenuListWithChild() (*[]dto.MenuTreeChild, error) {
	utils.InfoLog(nil, "HandleMenuListWithChild")
	MenusList, err := us.repo.SelectMenuWithChildren()
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err
	}
	utils.InfoLog(MenusList, string(utils.SuccessMessage))
	return MenusList, nil
}

func (us *MenuService) HandleMenuListWithChildByID(menuId int) (*[]dto.MenuTreeChild, error) {
	utils.InfoLog(menuId, "HandleMenuListWithChildByID")
	MenusList, err := us.repo.SelectMenuWithChildrenByID(int(menuId))
	if err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err
	}
	utils.InfoLog(MenusList, string(utils.SuccessMessage))
	return MenusList, nil
	//utils.ServiceResponse(MenusList, http.StatusOK, "Success", logrus.InfoLevel)

}
func (us *MenuService) HandleMenuCreateRequest(msgBody []byte) (*dto.MenuDTO, error) {

	var createMenuRequest dto.MenuDTO
	if err := json.Unmarshal(msgBody, &createMenuRequest); err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err
	}
	utils.InfoLog(createMenuRequest, "HandleMenuCreateRequest")
	MenusList, err := us.repo.CreateMenu(createMenuRequest)
	if err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err
	}
	MenusList_ := &dto.MenuDTO{}
	MenusList_.FromModel(MenusList)
	utils.InfoLog(MenusList, string(utils.SuccessMessage))
	return MenusList_, nil
}

func (us *MenuService) HandleMenuUpdateRequest(msgBody []byte) (*dto.MenuDTO, error) {
	var updateMenuRequest dto.MenuDTO
	if err := json.Unmarshal(msgBody, &updateMenuRequest); err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err

	}
	utils.InfoLog(updateMenuRequest, "HandleMenuUpdateRequest")
	updateMenuRequest.CreateTime = time.Now()
	updatedDto, err := us.repo.UpdateMenu(&updateMenuRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto_ := &dto.MenuDTO{}
	updatedDto_.FromModel(updatedDto)
	utils.InfoLog(updatedDto, string(utils.SuccessMessage))
	return updatedDto_, nil
}

func (us *MenuService) HandleMenuByIdRequest(MenuID int) (*dto.MenuDTO, error) {
	utils.InfoLog(MenuID, "HandleMenuByIdRequest")
	MenusList, err := us.repo.GetMenuByID(uint(MenuID))
	if err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return nil, err
	}

	utils.InfoLog(MenusList, string(utils.SuccessMessage))
	return MenusList, nil
}

func (us *MenuService) HandleMenuDeleteRequest(msgBody []byte) (bool, error) {
	var req dto.MenuDTO
	if err := json.Unmarshal(msgBody, &req); err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return false, err

	}
	utils.InfoLog(req, "HandleMenuDeleteRequest")
	deleted, err := us.repo.DeleteMenuByID(req.ID)
	if err != nil {
		utils.ServiceResponse(nil, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel)
		return false, err
	}
	utils.InfoLog(deleted, string(utils.SuccessMessage))
	return deleted, nil
}
