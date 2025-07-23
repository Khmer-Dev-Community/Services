package repository

import (
	"log"
	dto "telegram-service/lib/roles/dto"
	models "telegram-service/lib/roles/models"
	"telegram-service/utils"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetRoleList(page int, limit int, query string) ([]*models.Role, int, error) {
	offset := (page - 1) * limit
	var users []*models.Role
	var total int64
	db := r.db
	baseQuery := db.Model(&models.Role{})
	if query != "" {
		baseQuery = baseQuery.Where("role_name LIKE ?", "%"+query+"%")
		log.Printf("Generated SQL Query: %v", baseQuery.Statement.SQL.String())
	}

	// Count the total number of records matching the query
	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Then, retrieve the paginated list
	err = baseQuery.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (r *RoleRepository) CreateRole(dto *dto.RoleCreateDTO) (*models.Role, error) {
	role := &models.Role{
		RoleName:   dto.RoleName,
		RoleStatus: int(dto.RoleStatus),
		RoleKey:    dto.RoleKey,
		Decription: dto.Description,
		Sort:       dto.Sort,
		Status:     dto.Status,
		CompanyID:  dto.CompanyID,
	}
	if err := r.db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) UpdateRole(userID uint, dto *dto.RoleUpdateDTO) (*models.Role, error) {

	var existingRole models.Role
	if err := r.db.First(&existingRole, userID).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	existingRole.RoleName = dto.RoleName
	existingRole.RoleKey = dto.RoleKey
	existingRole.Decription = dto.Description
	existingRole.RoleStatus = int(dto.RoleStatus)
	existingRole.Sort = dto.Sort
	existingRole.Status = dto.Status
	existingRole.CompanyID = dto.CompanyID
	if err := r.db.Save(&existingRole).Error; err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	return &existingRole, nil
}

func (r *RoleRepository) DeleteRoleByID(id uint) error {
	result := r.db.Delete(&models.Role{}, id)
	return result.Error
}

func (r *RoleRepository) GetRoleByID(roleID uint) (*models.Role, error) {
	var role models.Role
	baseQuery := r.db.Model(&models.Role{})
	err := baseQuery.Where("id = ?", roleID).First(&role).Error

	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(baseQuery.Statement.SQL.String(), "Execute")
	return &role, nil
}
