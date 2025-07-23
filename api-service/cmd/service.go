package main

import (
	"github.com/Khmer-Dev-Community/Services/api-service/auth"
	"github.com/Khmer-Dev-Community/Services/api-service/auth02"
	"github.com/Khmer-Dev-Community/Services/api-service/config"
	clientAuthService "github.com/Khmer-Dev-Community/Services/api-service/lib/clientauth"
	menuService "github.com/Khmer-Dev-Community/Services/api-service/lib/menus"
	menusRepo "github.com/Khmer-Dev-Community/Services/api-service/lib/menus/repository"
	permissRepo "github.com/Khmer-Dev-Community/Services/api-service/lib/permission/repository"
	permissService "github.com/Khmer-Dev-Community/Services/api-service/lib/permission/service"
	"github.com/Khmer-Dev-Community/Services/api-service/lib/roles/repository"
	roleService "github.com/Khmer-Dev-Community/Services/api-service/lib/roles/services"
	userclientRepo "github.com/Khmer-Dev-Community/Services/api-service/lib/userclient"
	userRepo "github.com/Khmer-Dev-Community/Services/api-service/lib/users/repository"
	userService "github.com/Khmer-Dev-Community/Services/api-service/lib/users/services"

	"gorm.io/gorm"
)

type Services struct {
	Auth       *auth.AuthService
	Auth02     *auth02.ClientAuthController
	User       *userService.UserService
	Role       *roleService.RoleService
	Permission *permissService.PermissionService
	Menu       *menuService.MenuService
	Config     *config.Config
}

func InitServices(db *gorm.DB, cfg *config.Config) *Services {
	clientUserRepository := userclientRepo.NewClientUserRepository(db)
	clientAuthService := clientAuthService.NewClientAuthService(clientUserRepository) // clientauth_service
	userclientRepo.MigrateClientUsers(db)
	return &Services{
		Auth:       auth.NewAuthService(auth.NewUserRepository(db)),
		Auth02:     auth02.NewClientAuthController(clientAuthService),
		User:       userService.NewUserService(userRepo.NewUserRepository(db)),
		Role:       roleService.NewRoleService(repository.NewRoleRepository(db)),
		Permission: permissService.NewPermissionService(permissRepo.NewPermissionRepository(db)),
		Menu:       menuService.NewMenuService(menusRepo.NewMenuRepository(db)),
		Config:     cfg,
	}
}
