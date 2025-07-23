package main

import (
	"telegram-service/auth"
	"telegram-service/config"
	"telegram-service/lib/department"
	"telegram-service/lib/employee"
	menuService "telegram-service/lib/menus"
	menusRepo "telegram-service/lib/menus/repository"
	permissRepo "telegram-service/lib/permission/repository"
	permissService "telegram-service/lib/permission/service"
	"telegram-service/lib/roles/repository"
	roleService "telegram-service/lib/roles/services"
	userRepo "telegram-service/lib/users/repository"
	userService "telegram-service/lib/users/services"

	// --- PROXY IMPORTS (Make sure these are present if you use proxy) ---

	proxyService "telegram-service/lib/proxy/services"
	// --- END PROXY IMPORTS ---

	telegrambot "telegram-service/telegram/bot"
	telegramRepo "telegram-service/telegram/repository"
	telegramService "telegram-service/telegram/services" // Import your existing telegram service

	"gorm.io/gorm"
)

type Services struct {
	Auth               *auth.AuthService
	User               *userService.UserService
	Role               *roleService.RoleService
	Permission         *permissService.PermissionService
	Menu               *menuService.MenuService
	Employee           *employee.EmployeeService
	Department         *department.DepartmentService
	Proxy              proxyService.ProxyListServicer
	ActiveUserBots     map[string]*telegrambot.BotAccount
	Config             *config.Config
	TransactionService telegramService.TransactionService // <--- CORRECTED THIS LINE (removed '*')
	TgGroup            telegramService.TelegramGroupService
}

func InitServices(db *gorm.DB, cfg *config.Config, activeUserBots map[string]*telegrambot.BotAccount) *Services {
	tgRepo := telegramRepo.NewTransactionRepository(db)
	tgService := telegramService.NewTransactionService(tgRepo) // This returns the interface type

	tgGroupRepo := telegramRepo.NewTelegramGroupRepository(db)
	tgGroupService := telegramService.NewTelegramGroupService(tgGroupRepo)

	return &Services{
		Auth:               auth.NewAuthService(auth.NewUserRepository(db)),
		User:               userService.NewUserService(userRepo.NewUserRepository(db)),
		Role:               roleService.NewRoleService(repository.NewRoleRepository(db)),
		Permission:         permissService.NewPermissionService(permissRepo.NewPermissionRepository(db)),
		Menu:               menuService.NewMenuService(menusRepo.NewMenuRepository(db)),
		Employee:           employee.NewEmployeeService(employee.NewEmployeeRepository(db)),
		Department:         department.NewDepartmentService(department.NewDepartmentRepository(db)),
		ActiveUserBots:     activeUserBots,
		Config:             cfg,
		TransactionService: tgService, // This now correctly matches the field type
		TgGroup:            *tgGroupService,
	}
}
