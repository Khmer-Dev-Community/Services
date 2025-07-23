package permission

import (
	models "telegram-service/lib/permission/model"
)

type PermissionDTO struct {
	ID                int    `json:"id"`
	RoleID            int    `json:"roleId"`
	MenuID            int    `json:"menuId"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	RoleKey           string `json:"roleKey"`
	RoleName          string `json:"roleName"`
}

type RolePermissionsDTO struct {
	Permissions []string `json:"permissions"`
	RoleName    []string `json:"roles"`
}
type PermissionCreateDTO struct {
	RoleID            int    `json:"roleId"`
	MenuID            int    `json:"menuId"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	RoleKey           string `json:"roleKey"`
	RoleName          string `json:"roleName"`
}

type PermissionUpdateDTO struct {
	ID                int    `json:"id"`
	RoleID            int    `json:"roleId"`
	MenuID            int    `json:"menuId"`
	MenuCheckStrictly bool   `json:"menuCheckStrictly"`
	RoleKey           string `json:"roleKey"`
	RoleName          string `json:"roleName"`
}
type UpdateRolePermissionsMessage struct {
	RoleID  int   `json:"id"`
	MenuIDs []int `json:"menuIds"`
}

func (permission *PermissionDTO) FromModel(permis *models.Permission) {
	permission.ID = int(permis.ID)
	permission.RoleID = permis.RoleID
	permission.MenuID = permis.MenuID
	permission.MenuCheckStrictly = permis.MenuCheckStrictly
	permission.RoleKey = permis.RoleKey
	permission.RoleName = permis.RoleName

}

func (permission *PermissionDTO) ToModel() *models.Permission {
	return &models.Permission{
		ID:                permission.ID,
		RoleID:            permission.RoleID,
		MenuID:            permission.MenuID,
		MenuCheckStrictly: permission.MenuCheckStrictly,
		RoleKey:           permission.RoleKey,
		RoleName:          permission.RoleName,
	}
}
