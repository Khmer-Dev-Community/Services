package menus

import (
	"log"

	dto "github.com/Khmer-Dev-Community/Services/api-service/lib/menus/dto"
	models "github.com/Khmer-Dev-Community/Services/api-service/lib/menus/model"
	"github.com/Khmer-Dev-Community/Services/api-service/utils"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) GetMenuList(page int, limit int, query string) ([]*models.Menu, int, error) {
	offset := (page - 1) * limit
	var users []*models.Menu
	var total int64
	db := r.db
	baseQuery := db.Model(&models.Menu{})
	if query != "" {
		baseQuery = baseQuery.Where("name LIKE ?", "%"+query+"%")
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
func (r *MenuRepository) CreateMenu(menuDTO dto.MenuDTO) (*models.Menu, error) {
	menu := menuDTO.ToModel()
	err := r.db.Create(menu).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}

	utils.LoggerRepository(menu, "Execute")
	return menu, nil
}

func (r *MenuRepository) UpdateMenu(menuDTO *dto.MenuDTO) (*models.Menu, error) {
	menu := menuDTO.ToModel()
	err := r.db.Save(menu).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	utils.LoggerRepository(menu, "Execute")
	return menu, nil
}
func (r *MenuRepository) DeleteMenuByID(deptId int) (bool, error) {
	err := r.db.Delete(&models.Menu{}, deptId).Error
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return false, err
	}
	utils.LoggerRepository(true, "Execute")
	return true, nil
}

func (r *MenuRepository) GetMenuByID(menuID uint) (*dto.MenuDTO, error) {
	var menu models.Menu
	baseQuery := r.db.Model(&models.Menu{})
	err := baseQuery.Where("id = ?", menuID).First(&menu).Error
	log.Printf("FindById SQL Query: %v", baseQuery.Statement.SQL.String())
	if err != nil {
		utils.LoggerRepository(err, "Execute")
		return nil, err
	}
	menuDTO := &dto.MenuDTO{}
	menuDTO.FromModel(&menu)
	utils.LoggerRepository(menuDTO, "Execute")
	return menuDTO, nil
}
func (r *MenuRepository) SelectMenusByTenantRoleID(roleId int) (*[]dto.MenuTree, error) {
	var menus []models.Menu
	query := r.db.
		Table("system_menu").
		Select("DISTINCT system_menu.*").
		Joins("JOIN system_menu_role r ON r.menu_id = system_menu.id").
		Where("system_menu.menu_type IN (?) ", []string{"M", "C", "F"}).
		Where("r.role_id = ?", roleId).
		//Where("r.deleted_at IS NULL").
		Order("system_menu.order_num").
		Find(&menus)
	if query.Error != nil {
		return nil, query.Error
	}
	//utils.LoggerRepository(query, "Execute")
	var buildMenuTree func(subOf int) []dto.MenuTree
	buildMenuTree = func(subOf int) []dto.MenuTree {
		var result []dto.MenuTree
		for _, menu := range menus {
			if menu.SubOf == subOf {
				children := buildMenuTree(menu.ID)
				// Only append if children is not empty
				if len(children) > 0 {
					result = append(result, dto.MenuTree{
						ID:         menu.ID,
						Name:       menu.Name,
						Path:       menu.Path,
						Redirect:   menu.Redirect,
						Component:  menu.Component,
						AlwaysShow: menu.AlwaysShow,
						Hidden:     menu.Hidden,
						Children:   children,
						MenuType:   menu.MenuType,
						Perms:      menu.Perms,
						APIURL:     menu.APIURL,
						Meta: dto.Meta{
							Title:       menu.Title,
							Icon:        menu.Icon,
							IconColor:   menu.IconColor,
							NoCache:     menu.NoCache,
							TitleKey:    menu.TitleKey,
							Link:        menu.Link,
							IsLink:      menu.Link,
							IsHide:      menu.Hidden,
							IsKeepAlive: true,
							IsAffix:     false,
							IsIframe:    false,
						},
					})
				} else {
					// Append only the menu itself if it has no children
					result = append(result, dto.MenuTree{
						ID:         menu.ID,
						Name:       menu.Name,
						Title:      menu.Title,
						Path:       menu.Path,
						Redirect:   menu.Redirect,
						Component:  menu.Component,
						AlwaysShow: menu.AlwaysShow,
						Hidden:     menu.Hidden,
						MenuType:   menu.MenuType,
						Perms:      menu.Perms,
						APIURL:     menu.APIURL,
						Meta: dto.Meta{
							Title:       menu.Title,
							Icon:        menu.Icon,
							IconColor:   menu.IconColor,
							NoCache:     menu.NoCache,
							TitleKey:    menu.TitleKey,
							Link:        menu.Link,
							IsLink:      menu.Link,
							IsHide:      menu.Hidden,
							IsKeepAlive: true,
							IsAffix:     menu.Hidden,
							IsIframe:    false,
						},
					})
				}
			}
		}
		return result
	}
	menuTree := buildMenuTree(0)
	//utils.LoggerRepository(menuTree, "Execute")
	return &menuTree, nil
}
func (r *MenuRepository) SelectMenus(roleId int) (*[]dto.MenuTree, error) {
	var menus []models.Menu
	query := r.db.
		Table("system_menu").
		Select("DISTINCT system_menu.*").
		Joins("JOIN system_menu_role r ON r.menu_id = system_menu.id").
		Where("system_menu.menu_type IN (?) ", []string{"M", "C"}).
		Where("r.role_id = ?", roleId).
		//Where("r.deleted_at IS NULL").
		Order("system_menu.order_num").
		Find(&menus)
	if query.Error != nil {
		return nil, query.Error
	}
	//utils.LoggerRepository(query, "Execute")
	var buildMenuTree func(subOf int) []dto.MenuTree
	buildMenuTree = func(subOf int) []dto.MenuTree {
		var result []dto.MenuTree
		for _, menu := range menus {
			if menu.SubOf == subOf {
				children := buildMenuTree(menu.ID)
				// Only append if children is not empty
				if len(children) > 0 {
					result = append(result, dto.MenuTree{
						ID:         menu.ID,
						Name:       menu.Name,
						Path:       menu.Path,
						Redirect:   menu.Redirect,
						Component:  menu.Component,
						AlwaysShow: menu.AlwaysShow,
						Hidden:     menu.Hidden,
						Children:   children,
						MenuType:   menu.MenuType,
						Perms:      menu.Perms,
						APIURL:     menu.APIURL,
						Meta: dto.Meta{
							Title:       menu.Title,
							Icon:        menu.Icon,
							IconColor:   menu.IconColor,
							NoCache:     menu.NoCache,
							TitleKey:    menu.TitleKey,
							Link:        menu.Link,
							IsLink:      menu.Link,
							IsHide:      menu.Hidden,
							IsKeepAlive: true,
							IsAffix:     false,
							IsIframe:    false,
						},
					})
				} else {
					// Append only the menu itself if it has no children
					result = append(result, dto.MenuTree{
						ID:         menu.ID,
						Name:       menu.Name,
						Title:      menu.Title,
						Path:       menu.Path,
						Redirect:   menu.Redirect,
						Component:  menu.Component,
						AlwaysShow: menu.AlwaysShow,
						Hidden:     menu.Hidden,
						MenuType:   menu.MenuType,
						Perms:      menu.Perms,
						APIURL:     menu.APIURL,
						Meta: dto.Meta{
							Title:       menu.Title,
							Icon:        menu.Icon,
							IconColor:   menu.IconColor,
							NoCache:     menu.NoCache,
							TitleKey:    menu.TitleKey,
							Link:        menu.Link,
							IsLink:      menu.Link,
							IsHide:      menu.Hidden,
							IsKeepAlive: true,
							IsAffix:     menu.Hidden,
							IsIframe:    false,
						},
					})
				}
			}
		}
		return result
	}
	menuTree := buildMenuTree(0)
	//utils.LoggerRepository(menuTree, "Execute")
	return &menuTree, nil
}

func (r *MenuRepository) SelectMenuWithChildren() (*[]dto.MenuTreeChild, error) {
	var menus []models.Menu
	query := r.db.
		Table("system_menu").
		Select("DISTINCT system_menu.*").
		Order("system_menu.order_num").
		Find(&menus)
	if query.Error != nil {
		utils.LoggerRepository(query.Error, "Execute")
		return nil, query.Error
	}

	var buildMenuTree func(subOf int) []dto.MenuTreeChild
	buildMenuTree = func(subOf int) []dto.MenuTreeChild {
		var result []dto.MenuTreeChild
		for _, menu := range menus {
			if menu.SubOf == subOf {
				children := buildMenuTree(menu.ID)
				menuTreeChild := dto.MenuTreeChild{
					ID:         menu.ID,
					Name:       menu.Name,
					Path:       menu.Path,
					Redirect:   menu.Redirect,
					Component:  menu.Component,
					AlwaysShow: menu.AlwaysShow,
					Hidden:     menu.Hidden,
					SubOf:      menu.SubOf,
					OrderNum:   menu.OrderNum,
					IsFrame:    menu.IsFrame,
					MenuType:   menu.MenuType,
					Perms:      menu.Perms,
					MenuStatus: menu.MenuStatus,
					APIURL:     menu.APIURL,
					Children:   children,
					Title:      menu.Title,
					Icon:       menu.Icon,
					IconColor:  menu.Icon,
					NoCache:    menu.NoCache,
					TitleKey:   menu.TitleKey,
					Link:       menu.Link,
				}
				result = append(result, menuTreeChild)
			}
		}
		return result
	}

	// Build the menu tree starting from the root (subOf = 0)
	menuTree := buildMenuTree(0)
	//utils.LoggerRepository(menuTree, "Execute")
	return &menuTree, nil
}

/*
	func (r *MenuRepository) SelectMenuWithChildrenByID1(menuID int) (*[]dto.MenuTreeChild, error) {
		var mainMenu []models.Menu
		if err := r.db.First(&mainMenu, menuID).Where("Menu.subof = ?", menuID).Error; err != nil {
			return nil, err
		}
		var chlids []models.Menu
		query := r.db.
			Table("Menus").
			Select("DISTINCT Menus.*").
			//Where("Menu.subof = ?", menuID).
			Order("Menus.orderNum").
			Find(&chlids)

		if query.Error != nil {
			return nil, query.Error
		}
		var buildMenuTree func(subOf int) []dto.MenuTreeChild
		buildMenuTree = func(subOf int) []dto.MenuTreeChild {
			var result []dto.MenuTreeChild
			for _, menu := range chlids {
				if menu.SubOf == subOf {
					children := buildMenuTree(menu.ID)
					menuTreeChild := dto.MenuTreeChild{
						ID:         menu.ID,
						Name:       menu.Name,
						Path:       menu.Path,
						Redirect:   menu.Redirect,
						Component:  menu.Component,
						AlwaysShow: menu.AlwaysShow,
						Hidden:     menu.Hidden,
						SubOf:      menu.SubOf,
						OrderNum:   menu.OrderNum,
						IsFrame:    menu.IsFrame,
						MenuType:   menu.MenuType,
						Perms:      menu.Perms,
						MenuStatus: menu.MenuStatus,
						APIURL:     menu.APIURL,
						Children:   children,
						Title:      menu.Title,
						Icon:       menu.Icon,
						NoCache:    menu.NoCache,
						TitleKey:   menu.TitleKey,
						Link:       menu.Link,
					}

					result = append(result, menuTreeChild)
				}
			}
			return result
		}

		// Build the menu tree starting from the root menu
		menuTree := dto.MenuTreeChild{
			ID:         mainMenu.ID,
			Name:       mainMenu.Name,
			Path:       mainMenu.Path,
			Redirect:   mainMenu.Redirect,
			Component:  mainMenu.Component,
			AlwaysShow: mainMenu.AlwaysShow,
			Hidden:     mainMenu.Hidden,
			SubOf:      mainMenu.SubOf,
			OrderNum:   mainMenu.OrderNum,
			IsFrame:    mainMenu.IsFrame,
			MenuType:   mainMenu.MenuType,
			Perms:      mainMenu.Perms,
			MenuStatus: mainMenu.MenuStatus,
			APIURL:     mainMenu.APIURL,
			Children:   buildMenuTree(mainMenu.ID),
			Title:      mainMenu.Title,
			Icon:       mainMenu.Icon,
			NoCache:    mainMenu.NoCache,
			TitleKey:   mainMenu.TitleKey,
			Link:       mainMenu.Link,
		}

		return &menuTree, nil
	}
*/
func (r *MenuRepository) SelectMenuWithChildrenByID(menuID int) (*[]dto.MenuTreeChild, error) {
	// Fetch main menu item with the given menuID
	var mainMenu models.Menu

	if err := r.db.First(&mainMenu, menuID).Where("system_menu.subof = ?", menuID).Error; err != nil {
		return nil, err
	}
	// Fetch all child menus
	var chlids []models.Menu
	query := r.db.
		Table("system_menu").
		Select("DISTINCT system_menu.*").
		//	Where("Menus.subof = ?", menuID).
		Order("system_menu.order_num").
		Find(&chlids)

	if query.Error != nil {
		utils.LoggerRepository(query.Error, "Execute")
		return nil, query.Error
	}

	// Recursive function to build the menu tree
	var buildMenuTree func(subOf int) []dto.MenuTreeChild
	buildMenuTree = func(subOf int) []dto.MenuTreeChild {
		var result []dto.MenuTreeChild
		for _, menu := range chlids {
			if menu.SubOf == subOf {
				children := buildMenuTree(menu.ID)
				menuTreeChild := dto.MenuTreeChild{
					ID:         menu.ID,
					Name:       menu.Name,
					Path:       menu.Path,
					Redirect:   menu.Redirect,
					Component:  menu.Component,
					AlwaysShow: menu.AlwaysShow,
					Hidden:     menu.Hidden,
					SubOf:      menu.SubOf,
					OrderNum:   menu.OrderNum,
					IsFrame:    menu.IsFrame,
					MenuType:   menu.MenuType,
					Perms:      menu.Perms,
					MenuStatus: menu.MenuStatus,
					APIURL:     menu.APIURL,
					Children:   children,
					Title:      menu.Title,
					Icon:       menu.Icon,
					IconColor:  menu.IconColor,
					NoCache:    menu.NoCache,
					TitleKey:   menu.TitleKey,
					Link:       menu.Link,
				}
				result = append(result, menuTreeChild)
			}
		}
		return result
	}

	// Build the menu tree starting from the root menu item
	menuTree := dto.MenuTreeChild{
		ID:         mainMenu.ID,
		Name:       mainMenu.Name,
		Path:       mainMenu.Path,
		Redirect:   mainMenu.Redirect,
		Component:  mainMenu.Component,
		AlwaysShow: mainMenu.AlwaysShow,
		Hidden:     mainMenu.Hidden,
		SubOf:      mainMenu.SubOf,
		OrderNum:   mainMenu.OrderNum,
		IsFrame:    mainMenu.IsFrame,
		MenuType:   mainMenu.MenuType,
		Perms:      mainMenu.Perms,
		MenuStatus: mainMenu.MenuStatus,
		APIURL:     mainMenu.APIURL,
		Children:   buildMenuTree(mainMenu.ID),
		Title:      mainMenu.Title,
		Icon:       mainMenu.Icon,
		IconColor:  mainMenu.IconColor,
		NoCache:    mainMenu.NoCache,
		TitleKey:   mainMenu.TitleKey,
		Link:       mainMenu.Link,
	}

	// Wrap menuTree in a slice to match the expected return type
	result := []dto.MenuTreeChild{menuTree}
	//utils.LoggerRepository(result, "Execute")
	return &result, nil
}

func (r *MenuRepository) SelectMenuWithChildrenToPermissionAuthorize() (*[]dto.MenuLableChild, error) {
	var menus []models.Menu
	query := r.db.
		Table("system_menu").
		Select("DISTINCT system_menu.*").
		Order("system_menu.order_num").
		Find(&menus)
	if query.Error != nil {
		utils.LoggerRepository(query.Error, "Execute")
		return nil, query.Error
	}

	var buildMenuTree func(subOf int) []dto.MenuLableChild
	buildMenuTree = func(subOf int) []dto.MenuLableChild {
		var result []dto.MenuLableChild
		for _, menu := range menus {
			if menu.SubOf == subOf {
				children := buildMenuTree(menu.ID)
				menuTreeChild := dto.MenuLableChild{
					ID:       menu.ID,
					Label:    menu.Name,
					Children: children,
				}
				result = append(result, menuTreeChild)
			}
		}
		return result
	}

	// Build the menu tree starting from the root (subOf = 0)
	menuTree := buildMenuTree(0)
	//utils.LoggerRepository(menuTree, "Execute")
	return &menuTree, nil
}

func (r *MenuRepository) GetPermissionMenuIDs(roleID int) ([]int, error) {
	var menuIDs []int
	utils.LoggerRepository(roleID, "Execute")
	// Corrected the WHERE clause
	//AND deleted_at IS NULL
	if err := r.db.Table("system_menu_role").Where("role_id = ? ", roleID).Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}

	//utils.LoggerRepository(menuIDs, "Execute")
	return menuIDs, nil
}
