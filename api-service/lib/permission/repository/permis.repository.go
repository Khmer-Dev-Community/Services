package permission

import (
	"fmt"
	"log"
	dto "telegram-service/lib/permission/dto"
	models "telegram-service/lib/permission/model"
	role "telegram-service/lib/roles/models"
	"telegram-service/utils"
	"time"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GetPermissionList(page int, limit int, query string) ([]*models.Permission, int, error) {
	offset := (page - 1) * limit
	var permision []*models.Permission
	var total int64
	db := r.db
	baseQuery := db.Model(&models.Permission{})
	if query != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+query+"%")
		log.Printf("Generated SQL Query: %v", baseQuery.Statement.SQL.String())
	}
	// Count the total number of records matching the query
	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Then, retrieve the paginated list
	err = baseQuery.Offset(offset).Limit(limit).Find(&permision).Error
	if err != nil {
		return nil, 0, err
	}

	return permision, int(total), nil
}

func (r *PermissionRepository) CreatePermission(dto *dto.PermissionCreateDTO) (*models.Permission, error) {
	Permission := &models.Permission{

		RoleID:            dto.RoleID,
		MenuID:            dto.MenuID,
		MenuCheckStrictly: dto.MenuCheckStrictly,
		RoleKey:           dto.RoleKey,
		RoleName:          dto.RoleName,
		CreatedAt:         time.Now(),
	}
	if err := r.db.Create(Permission).Error; err != nil {
		return nil, err
	}
	return Permission, nil
}

func (r *PermissionRepository) UpdatePermission(Id uint, dto *dto.PermissionUpdateDTO) (*models.Permission, error) {
	var permission models.Permission

	if err := r.db.First(&permission, Id).Error; err != nil {
		return nil, err
	}
	permission.RoleID = dto.RoleID
	permission.MenuID = dto.MenuID
	permission.MenuCheckStrictly = dto.MenuCheckStrictly
	permission.RoleKey = dto.RoleKey
	permission.RoleName = dto.RoleName
	permission.CreatedAt = time.Now()

	// Save the updated permission
	if err := r.db.Save(&permission).Error; err != nil {
		return nil, err
	}

	// Return the updated permission
	return &permission, nil
}

func (r *PermissionRepository) DeletePermissionByID(id uint) error {
	result := r.db.Delete(&models.Permission{}, id)
	return result.Error
}
func (r *PermissionRepository) GetPermissionByID(PermissionID uint) (*models.Permission, error) {
	var Permission models.Permission
	baseQuery := r.db.Model(&models.Permission{})
	err := baseQuery.Where("id = ?", PermissionID).First(&Permission).Error
	log.Printf("FindById SQL Query: %v", baseQuery.Statement.SQL.String())
	if err != nil {
		return nil, err
	}
	return &Permission, nil
}

func (r *PermissionRepository) GetPermissionsForRole11(roleID uint) ([]string, error) {
	var permissions []string
	// Fetch permissions associated with the given role
	err := r.db.Model(&models.Permission{}).
		Select("perms").
		Joins("JOIN role_accesses ON permissions.menu_id = role_accesses.menu_id").
		Where("role_accesses.role_id = ?", roleID).
		Pluck("perms", &permissions).
		Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
func (r *PermissionRepository) GetPermissionsForRole(roleId int) ([]string, error) {
	var permissions []string
	//AND system_menu_role.deleted_at IS NULL
	err := r.db.Table("system_menu").
		Select("system_menu.perms").
		Joins("JOIN system_menu_role ON system_menu_role.menu_id = system_menu.id").
		Where("system_menu_role.role_id = ? AND system_menu.perms IS NOT NULL ", roleId).
		Pluck("system_menu.perms", &permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r *PermissionRepository) GetRoleName(roleId int) ([]string, error) {
	var roles []string
	err := r.db.Table("system_role").
		Select("system_role.role_name").
		Where("system_role.id = ?", roleId).
		Pluck("system_role.role_name", &roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
func (repo *PermissionRepository) GetRoleNames(roleID int) ([]string, error) {
	var roleNames []string
	err := repo.db.Model(&role.Role{}). // Assuming `Role` is your role model
						Where("id = ?", roleID).
						Pluck("role_name", &roleNames).Error // Assuming `roleName` holds the role name
	if err != nil {
		return nil, err
	}
	return roleNames, nil
}

func (r *PermissionRepository) GetPermissionMenuIDs(roleID int) ([]int, error) {
	var menuIDs []int
	if err := r.db.Table("system_menu_role").Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}
	return menuIDs, nil
}

func (r *PermissionRepository) UpdateUserPermissions11(roleID int, newMenuIDs []int) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	var existingMenuIDs []int
	// Fetch existing permissions
	if err := tx.Table("RoleAccesse").Where("role_id = ?", roleID).Pluck("MenuId", &existingMenuIDs).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Determine which menuIDs to delete
	toDelete := difference(existingMenuIDs, newMenuIDs)
	// Determine which menuIDs to add
	toAdd := difference(newMenuIDs, existingMenuIDs)
	// Delete permissions not in new menu IDs
	if len(toDelete) > 0 {
		if err := tx.Where("role_id = ? AND MenuId IN ?", roleID, toDelete).Delete(&models.Permission{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// Add new permissions
	for _, menuID := range toAdd {
		permission := models.Permission{
			RoleID: roleID,
			MenuID: menuID,
		}
		if err := tx.Create(&permission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (r *PermissionRepository) UpdateUserPermissions2(roleID int, menuIds []int) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var existingMenuIDs []int
	if err := tx.Table("system_menu_role").Where("role_id = ?", roleID).Pluck("menu_id", &existingMenuIDs).Error; err != nil {
		tx.Rollback()
		return err
	}
	// Delete permissions not in new menuIds
	if len(menuIds) > 0 {
		if err := tx.Exec("DELETE FROM system_menu_role WHERE role_id = ? AND menu_id NOT IN (?)", roleID, menuIds).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// If menuIds slice is empty, delete all permissions for the role
		if err := tx.Where("role_id = ?", roleID).Delete(&models.Permission{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	toAdd := difference(menuIds, existingMenuIDs)
	for _, menuID := range toAdd {
		permission := models.Permission{
			RoleID: roleID,
			MenuID: menuID,
		}
		if err := tx.Create(&permission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (r *PermissionRepository) UpdateUserPermissions(roleID int, menuIds []int) error {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		utils.LoggerRepository("Failed to begin transaction: "+tx.Error.Error(), "Error")
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			utils.LoggerRepository("Transaction rolled back due to panic: "+fmt.Sprint(r), "Error")
		}
	}()

	utils.LoggerRepository("Starting UpdateUserPermissions for roleID: "+fmt.Sprintf("%d", roleID), "Debug")

	// Retrieve existing menu IDs for the role
	var existingMenuIDs []int
	if err := tx.Table("system_menu_role").Where("role_id = ?", roleID).Pluck("menu_id", &existingMenuIDs).Error; err != nil {
		tx.Rollback()
		utils.LoggerRepository("Error retrieving existing menu IDs: "+err.Error(), "Error")
		return err
	}
	utils.LoggerRepository("Existing menu IDs for roleID "+fmt.Sprintf("%d", roleID)+": "+fmt.Sprintf("%v", existingMenuIDs), "Debug")

	// Delete permissions not in new menuIds
	if len(menuIds) > 0 {
		utils.LoggerRepository("Deleting permissions not in new menuIds for roleID: "+fmt.Sprintf("%d", roleID), "Debug")
		if err := tx.Where("role_id = ? AND menu_id NOT IN ?", roleID, menuIds).Delete(&models.Permission{}).Error; err != nil {
			tx.Rollback()
			utils.LoggerRepository("Error deleting permissions: "+err.Error(), "Error")
			return err
		}
	} else {
		// If menuIds slice is empty, delete all permissions for the role
		utils.LoggerRepository("menuIds is empty, deleting all permissions for roleID: "+fmt.Sprintf("%d", roleID), "Debug")
		if err := tx.Where("role_id = ?", roleID).Delete(&models.Permission{}).Error; err != nil {
			tx.Rollback()
			utils.LoggerRepository("Error deleting all permissions: "+err.Error(), "Error")
			return err
		}
	}

	// Identify new permissions to add
	toAdd := difference(menuIds, existingMenuIDs)
	utils.LoggerRepository("Identified new permissions to add: "+fmt.Sprintf("%v", toAdd), "Debug")
	for _, menuID := range toAdd {
		permission := models.Permission{
			RoleID:    roleID,
			MenuID:    menuID,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&permission).Error; err != nil {
			tx.Rollback()
			utils.LoggerRepository("Error adding new permission for menuID "+fmt.Sprintf("%d", menuID)+": "+err.Error(), "Error")
			return err
		}
		utils.LoggerRepository("Added new permission for menuID: "+fmt.Sprintf("%d", menuID), "Debug")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utils.LoggerRepository("Error committing transaction: "+err.Error(), "Error")
		return err
	}
	utils.LoggerRepository("Successfully updated user permissions for roleID: "+fmt.Sprintf("%d", roleID), "Debug")

	return nil
}

func (r *PermissionRepository) UpdateUserPermissions89(roleID int, menuIds []int) error {
	// Start a transaction

	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Retrieve existing menu IDs for the role
	var existingMenuIDs []int
	if err := tx.Table("system_menu_role").Where("role_id = ?", roleID).Pluck("menu_id", &existingMenuIDs).Error; err != nil {
		tx.Rollback()
		utils.LoggerRepository(err.Error(), "Execute")
		return err
	}

	// Delete permissions not in new menuIds
	if len(menuIds) > 0 {
		// Use GORM Delete method
		if err := tx.Where("role_id = ? AND menu_id NOT IN ?", roleID, menuIds).Delete(&models.Permission{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// If menuIds slice is empty, delete all permissions for the role
		if err := tx.Where("role_id = ?", roleID).Delete(&models.Permission{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Identify new permissions to add
	toAdd := difference(menuIds, existingMenuIDs)
	for _, menuID := range toAdd {
		permission := models.Permission{
			RoleID:    roleID,
			MenuID:    menuID,
			CreatedAt: time.Now(),
			// Ensure that CreatedAt is not being set here
		}
		if err := tx.Create(&permission).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func difference(slice1, slice2 []int) []int {
	m := make(map[int]bool)
	for _, item := range slice2 {
		m[item] = true
	}
	var diff []int
	for _, item := range slice1 {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
