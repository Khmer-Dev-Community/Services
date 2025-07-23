package department

type DepartmentCreateDTO struct {
	Code      string `json:"dept_code"`
	Name      string `json:"name"`
	CompanyID int    `json:"companyId"`
	Status    string `json:"status"`
}

type DepartmentUpdateDTO struct {
	ID        int    `json:"id"`
	Code      string `json:"dept_code"`
	Name      string `json:"name"`
	CompanyID int    `json:"companyId"`
	Status    string `json:"status"`
}

type DepartmentDTO struct {
	ID        int    `json:"id"`
	Code      string `json:"dept_code"`
	Name      string `json:"name"`
	CompanyID int    `json:"companyId"`
	Status    string `json:"status"`
}

// FromModel maps
func (dto *DepartmentDTO) FromModel(dept *Department) {
	dto.ID = int(dept.ID)
	dto.Code = dept.Code
	dto.Name = dept.Name
	dto.CompanyID = dept.CompanyID
	dto.Status = dept.Status

}

// ToModel maps
func (dto *DepartmentDTO) ToModel() *Department {
	return &Department{
		ID:        uint(dto.ID),
		Code:      dto.Code,
		Name:      dto.Name,
		CompanyID: dto.CompanyID,
		Status:    dto.Status,
	}
}
