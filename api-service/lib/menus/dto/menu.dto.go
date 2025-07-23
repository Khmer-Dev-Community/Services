package menus

import (
	models "telegram-service/lib/menus/model"
	"time"
)

type MenuDTO struct {
	ID         int       `json:"id"`
	Name       string    `json:"name" binding:"required"`
	Path       string    `json:"path"`
	Redirect   string    `json:"redirect"`
	Hidden     bool      `json:"hidden"`
	Title      string    `json:"title"`
	Icon       string    `json:"icon"`
	IconColor  string    `json:"icon_color"`
	NoCache    bool      `json:"noCache"`
	TitleKey   string    `json:"titleKey"`
	Link       string    `json:"link"`
	SubOf      int       `json:"subof"`
	Component  string    `json:"component"`
	OrderNum   int       `json:"orderNum"`
	IsFrame    int       `json:"isFrame"`
	MenuType   string    `json:"menuType"`
	Perms      string    `json:"perms"`
	AlwaysShow bool      `json:"alwaysShow"`
	MenuStatus bool      `json:"menuStatus"`
	APIURL     string    `json:"apiURL"`
	CreateTime time.Time `gorm:"column:created_time" json:"createTime"`
}

// MenuCreateDTO represents the DTO for creating a menu
type MenuCreateDTO struct {
	Name       string    `json:"name" binding:"required"`
	Path       string    `json:"path"`
	Redirect   string    `json:"redirect"`
	Hidden     bool      `json:"hidden"`
	Title      string    `json:"title"`
	Icon       string    `json:"icon"`
	IconColor  string    `json:"icon_color"`
	NoCache    bool      `json:"noCache"`
	TitleKey   string    `json:"titleKey"`
	Link       string    `json:"link"`
	SubOf      int       `json:"subof"`
	Component  string    `json:"component"`
	OrderNum   int       `json:"orderNum"`
	IsFrame    int       `json:"isFrame"`
	MenuType   string    `json:"menuType"`
	Perms      string    `json:"perms"`
	AlwaysShow bool      `json:"alwaysShow"`
	MenuStatus bool      `json:"menuStatus"`
	APIURL     string    `json:"apiURL"`
	CreateTime time.Time `gorm:"column:created_time" json:"createTime"`
}

type MenuUpdateDTO struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`
	Path       string `json:"path"`
	Redirect   string `json:"redirect"`
	Hidden     bool   `json:"hidden"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	IconColor  string `json:"icon_color"`
	NoCache    bool   `json:"noCache"`
	TitleKey   string `json:"titleKey"`
	Link       string `json:"link"`
	SubOf      int    `json:"subof"`
	Component  string `json:"component"`
	OrderNum   int    `json:"orderNum"`
	IsFrame    int    `json:"isFrame"`
	MenuType   string `json:"menuType"`
	Perms      string `json:"perms"`
	AlwaysShow bool   `json:"alwaysShow"`
	MenuStatus bool   `json:"menuStatus"`
	APIURL     string `json:"apiURL"`
	//CreateTime time.Time `gorm:"column:created_time" json:"createTime"`
}

type MenuTreeChild struct {
	ID         int             `json:"id"`
	Name       string          `json:"name" binding:"required"`
	Path       string          `json:"path"`
	Redirect   string          `json:"redirect"`
	Hidden     bool            `json:"hidden"`
	Title      string          `json:"title"`
	Icon       string          `json:"icon"`
	IconColor  string          `json:"icon_color"`
	NoCache    bool            `json:"noCache"`
	TitleKey   string          `json:"titleKey"`
	Link       string          `json:"link"`
	SubOf      int             `json:"subof"`
	Component  string          `json:"component"`
	OrderNum   int             `json:"orderNum"`
	IsFrame    int             `json:"isFrame"`
	MenuType   string          `json:"menuType"`
	Perms      string          `json:"perms"`
	AlwaysShow bool            `json:"alwaysShow"`
	MenuStatus bool            `json:"menuStatus"`
	APIURL     string          `json:"apiURL"`
	Children   []MenuTreeChild `json:"children,omitempty"`
}
type MenuChecked struct {
	Menu []MenuLableChild `json:"menu,omitempty"`
}
type MenuLableChild struct {
	ID       int              `json:"id"`
	Label    string           `json:"label"`
	Children []MenuLableChild `json:"children,omitempty"`
}
type MenuTreeselectDTO struct {
	Menu        []MenuLableChild `json:"menu,omitempty"`
	CheckedKeys []int            `json:"checkedKeys,omitempty"`
}
type MenuTree struct {
	Title      string     `json:"title"`
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Redirect   string     `json:"redirect"`
	Component  string     `json:"component,omitempty"`
	AlwaysShow bool       `json:"alwaysShow"`
	Hidden     bool       `json:"hidden"`
	Children   []MenuTree `json:"children,omitempty"`
	Meta       Meta       `json:"meta"`
	MenuType   string     `json:"menuType"`
	Perms      string     `json:"perms"`
	APIURL     string     `json:"apiURL"`
}

type Meta struct {
	Title       string `json:"title"`
	MenuType    string `json:"menuType"`
	Icon        string `json:"icon"`
	IconColor   string `json:"icon_color"`
	NoCache     bool   `json:"noCache"`
	TitleKey    string `json:"titleKey"`
	Link        string `json:"link"`
	IsLink      string `json:"isLink"`
	IsHide      bool   `json:"isHide"`
	IsKeepAlive bool   `json:"isKeepAlive"`
	IsAffix     bool   `json:"isAffix"`
	IsIframe    bool   `json:"isIframe"`
}

// FromModel populates the MenuDTO from a Menu model
func (dto *MenuDTO) FromModel(menu *models.Menu) {
	dto.ID = menu.ID
	dto.Name = menu.Name
	dto.Path = menu.Path
	dto.Redirect = menu.Redirect
	dto.Hidden = menu.Hidden
	dto.Title = menu.Title
	dto.Icon = menu.Icon
	dto.IconColor = menu.IconColor
	dto.NoCache = menu.NoCache
	dto.TitleKey = menu.TitleKey
	dto.Link = menu.Link
	dto.SubOf = menu.SubOf
	dto.Component = menu.Component
	dto.OrderNum = menu.OrderNum
	dto.IsFrame = menu.IsFrame
	dto.MenuType = menu.MenuType
	dto.Perms = menu.Perms
	dto.AlwaysShow = menu.AlwaysShow
	dto.MenuStatus = menu.MenuStatus
	dto.APIURL = menu.APIURL
	// If you have CreateTime, uncomment and set it
	dto.CreateTime = menu.CreateTime
}

// ToModel converts the MenuDTO to a Menu model
func (dto *MenuDTO) ToModel() *models.Menu {
	return &models.Menu{
		ID:         dto.ID,
		Name:       dto.Name,
		Path:       dto.Path,
		Redirect:   dto.Redirect,
		Hidden:     dto.Hidden,
		Title:      dto.Title,
		Icon:       dto.Icon,
		IconColor:  dto.IconColor,
		NoCache:    dto.NoCache,
		TitleKey:   dto.TitleKey,
		Link:       dto.Link,
		SubOf:      dto.SubOf,
		Component:  dto.Component,
		OrderNum:   dto.OrderNum,
		IsFrame:    dto.IsFrame,
		MenuType:   dto.MenuType,
		Perms:      dto.Perms,
		AlwaysShow: dto.AlwaysShow,
		MenuStatus: dto.MenuStatus,
		APIURL:     dto.APIURL,
		// If you have CreateTime, uncomment and set it
		CreateTime: dto.CreateTime,
	}
}
