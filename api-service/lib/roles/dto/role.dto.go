package roles

type RoleDTO struct {
	ID          int    `json:"id"`
	RoleName    string `json:"roleName"`
	RoleStatus  int    `json:"roleStatus"`
	RoleKey     string `json:"roleKey"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      bool   `json:"status"`
	CompanyID   int    `json:"companyId"`
}

type RoleCreateDTO struct {
	RoleName    string `json:"roleName" binding:"required"`
	RoleStatus  int    `json:"roleStatus" binding:"required"`
	RoleKey     string `json:"roleKey" binding:"required"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      bool   `json:"status"`
	CompanyID   int    `json:"companyId"`
}

type RoleUpdateDTO struct {
	ID          int    `json:"id"`
	RoleName    string `json:"roleName"`
	RoleStatus  int    `json:"roleStatus"`
	RoleKey     string `json:"roleKey"`
	Description string `json:"description"`
	Sort        int    `json:"sort"`
	Status      bool   `json:"status"`
	CompanyID   int    `json:"companyId"`
}
