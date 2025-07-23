package employee

import (
	"encoding/json"
	"fmt"
	"net/http"
	"telegram-service/utils"

	"github.com/sirupsen/logrus"
)

// EmployeeServiceInterface defines the methods that the EmployeeService should implement.
type EmployeeServiceInterface interface {
	EmployeeServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	EmployeeServiceCreate(obj []byte) (*EmployeeUpdateDTO, error)
	EmployeeServiceUpdate(obj []byte) (*EmployeeUpdateDTO, error)
	EmployeeServiceGetById(obj []byte) (*EmployeeUpdateDTO, error)
	EmployeeServiceDelete(obj []byte) (bool, error)
}

type EmployeeService struct {
	repo *EmployeeRepository
}

// NewEmployeeService creates a new instance of EmployeeService.
func NewEmployeeService(repo *EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

// EmployeeServiceGetList retrieves a paginated list of employees.
func (us *EmployeeService) EmployeeServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	employee, total, err := us.repo.GetEmployeeList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
	}

	employeeDTOs := make([]*EmployeeUpdateDTO, len(employee))
	for i, employee_ := range employee {
		employeeDTO := &EmployeeUpdateDTO{}
		employeeDTO.FromModel(employee_)
		employeeDTOs[i] = employeeDTO
	}
	return utils.NewServicePaginationResponse(employeeDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "EmployeeService [EmployeeServiceGetList]")
}

// EmployeeServiceCreate creates a new employee.
func (us *EmployeeService) EmployeeServiceCreate(obj []byte) (*EmployeeUpdateDTO, error) {
	var createDto EmployeeUpdateDTO

	err := json.Unmarshal(obj, &createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	createdemp, err := us.repo.CreateEmployee(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	empDTO := &EmployeeUpdateDTO{}
	empDTO.FromModel(createdemp)

	utils.InfoLog(empDTO, string(utils.SuccessMessage))
	return empDTO, nil
}

// EmployeeServiceUpdate updates an existing employee.
func (us *EmployeeService) EmployeeServiceUpdate(obj []byte) (*EmployeeUpdateDTO, error) {
	var updateDTO EmployeeUpdateDTO
	if err := json.Unmarshal(obj, &updateDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto, err := us.repo.UpdateEmployee(updateDTO)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto_ := &EmployeeUpdateDTO{}
	updatedDto_.FromModel(updatedDto)
	utils.InfoLog(updatedDto, string(utils.SuccessMessage))
	return updatedDto_, nil
}

// EmployeeServiceGetById retrieves an employee by ID.
func (us *EmployeeService) EmployeeServiceGetById(obj []byte) (*EmployeeUpdateDTO, error) {
	var empDTO EmployeeUpdateDTO
	if err := json.Unmarshal(obj, &empDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	Employee, err := us.repo.GetEmployeeByID(uint(empDTO.ID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto_ := &EmployeeUpdateDTO{}
	updatedDto_.FromModel(Employee)
	utils.InfoLog(updatedDto_, string(utils.SuccessMessage))
	return updatedDto_, nil
}

// EmployeeServiceDelete deletes an employee by ID.
func (us *EmployeeService) EmployeeServiceDelete(obj []byte) (bool, error) {
	var deptDTO EmployeeUpdateDTO

	// Unmarshal the JSON input
	if err := json.Unmarshal(obj, &deptDTO); err != nil {
		return false, logAndReturnError("failed to unmarshal input", err)
	}

	// Attempt to delete the employee by ID
	deleted, err := us.repo.DeleteEmployeeByID(int(deptDTO.ID))
	if err != nil {
		return false, logAndReturnError("failed to delete employee", err)
	}

	// Log success message
	utils.InfoLog(&deptDTO, string(utils.SuccessMessage))
	return deleted, nil
}

func logAndReturnError(message string, err error) error {
	utils.ErrorLog(nil, err.Error())
	return fmt.Errorf("%s: %v", message, err)
}
