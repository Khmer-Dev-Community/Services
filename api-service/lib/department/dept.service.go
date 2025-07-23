package department

import (
	"encoding/json"
	"fmt"
	"net/http"
	"telegram-service/utils"

	"github.com/sirupsen/logrus"
)

// DepartmentServiceInterface defines the methods for the DepartmentService
type DepartmentServiceInterface interface {
	DepartmentServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	DepartmentServiceCreate(obj []byte) (*DepartmentDTO, error)
	DepartmentServiceUpdate(obj []byte) (*DepartmentDTO, error)
	DepartmentServiceGetById(obj []byte) (*DepartmentDTO, error)
	DepartmentServiceDelete(obj []byte) (bool, error)
}
type DepartmentService struct {
	repo *DepartmentRepository
}

func NewDepartmentService(repo *DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

// DepartmentServiceGetList retrieves a paginated list of departments
func (us *DepartmentService) DepartmentServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10
	}
	depts, total, err := us.repo.GetDepartmentList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
	}
	deptDTOs := make([]*DepartmentDTO, len(depts))
	for i, dept := range depts {
		deptDTO := &DepartmentDTO{}
		deptDTO.FromModel(dept)
		deptDTOs[i] = deptDTO
	}
	return utils.NewServicePaginationResponse(deptDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "DepartmentService [GetDepartmentList]")
}

// DepartmentServiceCreate creates a new department
func (us *DepartmentService) DepartmentServiceCreate(obj []byte) (*DepartmentDTO, error) {
	var createDto DepartmentDTO
	err := json.Unmarshal(obj, &createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	createdDept, err := us.repo.CreateDepartment(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	deptDTO := &DepartmentDTO{}
	deptDTO.FromModel(createdDept)
	utils.InfoLog(deptDTO, string(utils.SuccessMessage))
	return deptDTO, nil
}

// DepartmentServiceUpdate updates an existing department
func (us *DepartmentService) DepartmentServiceUpdate(obj []byte) (*DepartmentDTO, error) {
	var updateDepartmentRequest DepartmentDTO
	if err := json.Unmarshal(obj, &updateDepartmentRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto, err := us.repo.UpdateDepartment(updateDepartmentRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto_ := &DepartmentDTO{}
	updatedDto_.FromModel(updatedDto)
	utils.InfoLog(updatedDto, string(utils.SuccessMessage))
	return updatedDto_, nil
}

// DepartmentServiceGetById retrieves a department by ID
func (us *DepartmentService) DepartmentServiceGetById(obj []byte) (*DepartmentDTO, error) {
	var deptDTO DepartmentUpdateDTO
	if err := json.Unmarshal(obj, &deptDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	department, err := us.repo.GetDepartmentByID(uint(deptDTO.ID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	updatedDto_ := &DepartmentDTO{}
	updatedDto_.FromModel(department)
	utils.InfoLog(updatedDto_, string(utils.SuccessMessage))
	return updatedDto_, nil
}

func (us *DepartmentService) DepartmentServiceDelete(obj []byte) (bool, error) {
	var deptDTO DepartmentDTO
	if err := json.Unmarshal(obj, &deptDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}
	success, err := us.repo.DeleteDepartmentByID(int(deptDTO.ID))
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}

	utils.InfoLog(deptDTO, string(utils.SuccessMessage))
	return success, nil
}
