package employee

type EmployeeCreateDTO struct {
	EmployeeID   string `json:"empCard"`
	EmployeeCard string `json:"ci_card"`
	FirstName    string `json:"fname"`
	LastName     string `json:"lname"`
	Username     string `json:"username"`
	Sex          string `json:"gender"`
	AccountName  string `json:"account_name"`
	Password     string `json:"password"`
	Position     string `json:"user_position"`
	RoleID       uint   `json:"roleId"`
	CompanyID    int    `json:"companyId"`
	DepartmentID int    `json:"departmentId"`
	JoinDate     string `json:"joinDate"`
	Level        string `json:"cardType"`
	Language     string `json:"default_lag"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Balanace     string `json:"balance"`
	Note         string `json:"note"`
	Status       string `json:"status"`
	Token        string `json:"token"`
}

type EmployeeUpdateDTO struct {
	ID           uint   `json:"id"`
	EmployeeID   string `json:"empCard"`
	EmployeeCard string `json:"ci_card"`
	FirstName    string `json:"fname"`
	LastName     string `json:"lname"`
	Username     string `json:"username"`
	Sex          string `json:"gender"`
	AccountName  string `json:"account_name"`
	Password     string `json:"password"`
	Position     string `json:"user_position"`
	RoleID       uint   `json:"roleId"`
	CompanyID    int    `json:"companyId"`
	DepartmentID int    `json:"departmentId"`
	JoinDate     string `json:"joinDate"`
	Level        string `json:"cardType"`
	Language     string `json:"default_lag"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Balanace     string `json:"balance"`
	Note         string `json:"note"`
	Status       string `json:"status"`
	Token        string `json:"token"`
}

// FromModel populates the EmployeeUpdateDTO from an Employee model
func (dto *EmployeeUpdateDTO) FromModel(emp *Employee) {
	dto.ID = emp.ID
	dto.EmployeeID = emp.EmployeeID
	dto.EmployeeCard = emp.EmployeeCard
	dto.FirstName = emp.FirstName
	dto.LastName = emp.LastName
	dto.Username = emp.Username
	dto.Sex = emp.Sex
	dto.AccountName = emp.AccountName
	dto.Password = emp.Password
	dto.Position = emp.Position
	dto.RoleID = emp.RoleID
	dto.CompanyID = emp.CompanyID
	dto.DepartmentID = emp.DepartmentID
	dto.JoinDate = emp.JoinDate
	dto.Level = emp.Level
	dto.Language = emp.Language
	dto.Phone = emp.Phone
	dto.Email = emp.Email
	dto.Balanace = emp.Balanace // Note: Adjust this to 'Balance' if there is a typo
	dto.Note = emp.Note
	dto.Status = emp.Status
	dto.Token = emp.Token
}

// ToModel converts the EmployeeUpdateDTO to an Employee model
func (dto *EmployeeUpdateDTO) ToModel() Employee {
	return Employee{
		ID:           dto.ID,
		EmployeeID:   dto.EmployeeID,
		EmployeeCard: dto.EmployeeCard,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Username:     dto.Username,
		Sex:          dto.Sex,
		AccountName:  dto.AccountName,
		Password:     dto.Password,
		Position:     dto.Position,
		RoleID:       dto.RoleID,
		CompanyID:    dto.CompanyID,
		DepartmentID: dto.DepartmentID,
		JoinDate:     dto.JoinDate,
		Level:        dto.Level,
		Language:     dto.Language,
		Phone:        dto.Phone,
		Email:        dto.Email,
		Balanace:     dto.Balanace, // Adjust if necessary
		Note:         dto.Note,
		Status:       dto.Status,
		Token:        dto.Token,
	}
}
