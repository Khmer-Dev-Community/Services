package department

import (
	"log"
	"telegram-service/utils"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) GetDepartmentList(page int, limit int, query string) ([]*Department, int, error) {
	offset := (page - 1) * limit
	var Departments []*Department
	var total int64
	db := r.db                               // Assuming you have a db connection here
	log.Printf("Query Parameter: %s", query) // Log the query parameter to verify its value
	baseQuery := db.Model(&Department{})     //.Preload("Role")
	if query != "" {
		baseQuery = baseQuery.Where(" name LIKE ?", "%"+query+"%")
		log.Printf("Generated SQL Query: %v", baseQuery.Statement.SQL.String())
	}
	err := baseQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = baseQuery.Offset(offset).Limit(limit).Find(&Departments).Error
	if err != nil {
		return nil, 0, err
	}
	return Departments, int(total), nil
}

func (r *DepartmentRepository) CreateDepartment(deptDTO DepartmentDTO) (*Department, error) {
	dept := deptDTO.ToModel()
	err := r.db.Create(dept).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(dept, "Execute")
	return dept, nil
}

func (r *DepartmentRepository) UpdateDepartment(deptDTO DepartmentDTO) (*Department, error) {
	dept := deptDTO.ToModel()
	err := r.db.Save(dept).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(dept, "Execute")
	return dept, nil
}

func (r *DepartmentRepository) DeleteDepartmentByID(deptId int) (bool, error) {
	err := r.db.Delete(&Department{}, deptId).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}

func (r *DepartmentRepository) GetDepartmentByID(id uint) (*Department, error) {
	var dept Department
	err := r.db.First(&dept, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(dept, "Execute")
	return &dept, err
}
