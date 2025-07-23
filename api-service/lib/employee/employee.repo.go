package employee

import (
	"telegram-service/utils"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) GetEmployeeList(page int, limit int, query string) ([]*Employee, int, error) {
	offset := (page - 1) * limit
	var Employees []*Employee
	var total int64
	db := r.db // Assuming you have a db connection here

	utils.LoggerRepository(query, "Execute")
	baseQuery := db.Model(&Employee{}) //.Preload("Role")
	if query != "" {
		baseQuery = baseQuery.Where(" employee_id LIKE ?", "%"+query+"%")
		utils.LoggerRepository(baseQuery.Statement.SQL.String(), "Execute")

	}
	err := baseQuery.Count(&total).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	err = baseQuery.Offset(offset).Limit(limit).Find(&Employees).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, 0, err
	}
	utils.LoggerRepository(Employees, "Execute")
	return Employees, int(total), nil
}

func (r *EmployeeRepository) CreateEmployee(empDTO EmployeeUpdateDTO) (*Employee, error) {
	emp := empDTO.ToModel()
	err := r.db.Create(emp).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(emp, "Execute")
	return &emp, nil
}
func (r *EmployeeRepository) UpdateEmployee(empDTO EmployeeUpdateDTO) (*Employee, error) {
	emp := empDTO.ToModel()
	err := r.db.Save(emp).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(emp, "Execute")
	return &emp, nil
}

func (r *EmployeeRepository) DeleteEmployeeByID(Id int) (bool, error) {
	err := r.db.Delete(&Employee{}, Id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	return true, nil
}
func (r *EmployeeRepository) GetEmployeeByID(id uint) (*Employee, error) {
	var emp Employee
	err := r.db.First(&emp, id).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(emp, "Execute")
	return &emp, err
}
