package service

import (
	"context"
	"sort"

	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/id2"
	logger "git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	repo "git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/code"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// Menu menu
type Menu interface {
	// CreateMenu 创建菜单
	CreateMenu(ctx context.Context, req *CreateMenuReq) (*CreateMenuResp, error)

	// DeleteMenu 删除菜单
	DeleteMenu(ctx context.Context, req *DeleteMenuReq) (*DeleteMenuResp, error)

	// UpdateMenu 更新菜单
	UpdateMenu(ctx context.Context, req *UpdateMenuReq) (*UpdateMenuResp, error)

	// CreateGroup 创建分组
	CreateGroup(ctx context.Context, req *CreateGroupReq) (*CreateGroupResp, error)

	// DeleteGroup 删除分组
	DeleteGroup(ctx context.Context, req *DeleteGroupReq) (*DeleteGroupResp, error)

	// ListAllGroup 查询所有分组
	ListAllGroup(ctx context.Context, req *ListAllGroupReq) (*ListAllGroupResp, error)

	// ListAll 查询所有
	ListAll(ctx context.Context, req *ListAllReq) (*ListAllResp, error)

	// Transfer 移动菜单
	Transfer(ctx context.Context, req *TransferReq) (*TransferResp, error)

	// UserListAll 用户所有菜单
	UserListAll(ctx context.Context, req *UserListAllReq) (*UserListAllResp, error)

	// ModifyMenuType 更新类型为自定义页面
	ModifyMenuType(ctx context.Context, req *ModifyMenuTypeReq) (*ModifyMenuTypeResp, error)

	// FindByID 通过id获取菜单信息
	FindByID(ctx context.Context, req *FindByIDReq) (*FindByIDResp, error)

	// ListPage
	ListPage(ctx context.Context, req *ListPageReq) (*ListPageResp, error)
}

type menu struct {
	mongo *mongo.Database

	menuRepo models.MenuRepo

	dataBaseRepo models.DataBaseSchemaRepo

	schemaRepo models.TableRepo

	customRepo models.CustomPageRepo

	permission Permission
}

// New new a menu
func New(conf *config.Config, opts ...Options) (Menu, error) {
	menuRepo := repo.NewMenuRepo()
	dataBaseRepo := repo.NewDataBaseSchemaRepo()
	schemaRepo := repo.NewTableRepo()
	customRepo := repo.NewCustomPageRepo()
	m := &menu{
		menuRepo:     menuRepo,
		dataBaseRepo: dataBaseRepo,
		schemaRepo:   schemaRepo,
		customRepo:   customRepo,
	}

	for _, opt := range opts {
		opt(m)
	}
	return m, nil
}

func (m *menu) SetMongo(client *mongo.Client, dbName string) {
	m.mongo = client.Database(dbName)
}

// CreateMenuReq create menu request param
type CreateMenuReq struct {
	AppID    string `json:"appID"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Describe string `json:"describe"`
	GroupID  string `json:"groupID"`
}

// CreateMenuResp CreateMenuResp
type CreateMenuResp struct {
	ID string `json:"id"`
}

func (m *menu) CreateMenu(ctx context.Context, req *CreateMenuReq) (*CreateMenuResp, error) {
	menu := &models.Menu{
		ID:       id2.GenUpperID(),
		AppID:    req.AppID,
		Name:     req.Name,
		Icon:     req.Icon,
		Describe: req.Describe,
		GroupID:  req.GroupID,
	}
	// 1、查找当前应用 当前分组是否存在相同名称的菜单
	id, err := m.menuRepo.FindSameMenuName(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if id != "" {
		return nil, error2.NewError(code.ErrRepeatMenuName)
	}
	// 2、 如果不存在相同名称，找到当前组中的最大sort
	sort, err := m.menuRepo.FindMaxSortFromGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if sort == -1 {
		return nil, error2.NewError(code.ErrRepeatMenuName)
	}

	// 3、插入数据
	menu.Sort = sort + 1
	menu.BindingState = models.Unbound
	err = m.menuRepo.InsertMenu(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	// 4、返回结果
	resp := &CreateMenuResp{
		ID: menu.ID,
	}
	return resp, err
}

// DeleteMenuReq DeleteMenuReq
type DeleteMenuReq struct {
	ID      string `json:"id"`
	Sort    int    `json:"sort"`
	GroupID string `json:"groupID"`
	AppID   string `json:"appID"`
}

// DeleteMenuResp DeleteMenuResp
type DeleteMenuResp struct {
}

func (m *menu) DeleteMenu(ctx context.Context, req *DeleteMenuReq) (*DeleteMenuResp, error) {
	menu := &models.Menu{
		ID:      req.ID,
		Sort:    req.Sort,
		GroupID: req.GroupID,
		AppID:   req.AppID,
	}
	// 1、删除组中的当前菜单
	err := m.menuRepo.DeleteMenuFromGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	// 2、更新当前组 后面菜单的sort
	err = m.menuRepo.UpdateSortFromGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	err = m.schemaRepo.Delete(ctx, m.mongo, req.ID)
	if err != nil {
		return nil, err
	}
	err = m.dataBaseRepo.Delete(ctx, m.mongo, req.ID)
	if err != nil {
		return nil, err
	}

	// 3、返回结果
	resp := &DeleteMenuResp{}
	return resp, err
}

// UpdateMenuReq UpdateMenuReq
type UpdateMenuReq struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Describe string `json:"describe"`
	AppID    string `json:"appID"`
	GroupID  string `json:"groupID"`
}

// UpdateMenuResp UpdateMenuResp
type UpdateMenuResp struct {
}

func (m *menu) UpdateMenu(ctx context.Context, req *UpdateMenuReq) (*UpdateMenuResp, error) {

	menu := &models.Menu{
		ID:       req.ID,
		Name:     req.Name,
		Icon:     req.Icon,
		Describe: req.Describe,
		AppID:    req.AppID,
		GroupID:  req.GroupID,
	}
	// 1、查找当前应用 当前分组是否存在相同名称的菜单
	id, err := m.menuRepo.FindSameMenuName(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if id != menu.ID && id != "" {
		return nil, error2.NewError(code.ErrRepeatMenuName)
	}
	// 2、更新菜单
	err = m.menuRepo.UpdateMenuByID(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	resp := &UpdateMenuResp{}
	return resp, err
}

// CreateGroupReq CreateGroupReq
type CreateGroupReq struct {
	AppID   string `json:"appID"`
	Name    string `json:"name"`
	GroupID string `json:"groupID"`
}

// CreateGroupResp CreateGroupResp
type CreateGroupResp struct {
	ID string `json:"id"`
}

func (m *menu) CreateGroup(ctx context.Context, req *CreateGroupReq) (*CreateGroupResp, error) {
	menu := &models.Menu{
		ID:       id2.GenUpperID(),
		AppID:    req.AppID,
		Name:     req.Name,
		MenuType: models.GroupType,
		GroupID:  req.GroupID,
	}

	// 1、查找是否存在相同名称的分组[分组也是菜单]
	n, err := m.menuRepo.FindSameMenuName(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if n != "" {
		return nil, error2.NewError(code.ErrRepeatMenuName)
	}

	// 2、不存在，找到当前组最大的Sort
	sort, err := m.menuRepo.FindMaxSortFromGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if sort == -1 {
		return nil, error2.NewError(code.ErrRepeatMenuName)
	}

	// 3、插入数据
	menu.Sort = sort + 1
	menu.BindingState = models.Bound
	err = m.menuRepo.InsertMenu(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	// 4、返回结果
	resp := &CreateGroupResp{
		ID: menu.ID,
	}
	return resp, err
}

// DeleteGroupReq DeleteGroupReq
type DeleteGroupReq struct {
	AppID   string `json:"appID"`
	ID      string `json:"id"`
	Sort    int    `json:"sort"`
	GroupID string `json:"groupID"`
}

// DeleteGroupResp DeleteGroupResp
type DeleteGroupResp struct {
}

func (m *menu) DeleteGroup(ctx context.Context, req *DeleteGroupReq) (*DeleteGroupResp, error) {
	menu := &models.Menu{
		AppID:   req.AppID,
		ID:      req.ID,
		Sort:    req.Sort,
		GroupID: req.GroupID,
	}
	// 1、查找当前分组是否存在数据
	menus, err := m.menuRepo.FindAllFromGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if len(menus) != 0 {
		return nil, error2.NewError(code.ErrDeleteMenu)
	}

	// 2、不存在,删除当前分组,并更新后面的sort
	err = m.menuRepo.DeleteGroupByID(ctx, m.mongo, menu.ID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	err = m.menuRepo.UpdateSortFromGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	// 3、返回结果
	resp := &DeleteGroupResp{}
	return resp, err

}

// ListAllGroupReq ListAllGroupReq
type ListAllGroupReq struct {
	AppID string `json:"appID"`
}

// ListAllGroupResp ListAllGroupResp
type ListAllGroupResp struct {
	Group []*ListAllGroupVO `json:"group"`
	Count int               `json:"count"`
}

// ListAllGroupVO ListAllGroupVO
type ListAllGroupVO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (m *menu) ListAllGroup(ctx context.Context, req *ListAllGroupReq) (*ListAllGroupResp, error) {
	menu := &models.Menu{
		AppID:    req.AppID,
		MenuType: models.None,
	}
	groups, err := m.menuRepo.ListAllGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	result := make([]*ListAllGroupVO, 0)
	for _, element := range groups {
		if element.GroupID == "" && element.MenuType == models.GroupType {
			temp := &ListAllGroupVO{
				ID:   element.ID,
				Name: element.Name,
			}
			result = append(result, temp)
		}
	}
	resp := new(ListAllGroupResp)
	for len(result) > 0 {
		r := result[0]
		for _, element := range groups {
			if element.GroupID == r.ID && element.MenuType == models.GroupType {
				temp := &ListAllGroupVO{
					ID:   element.ID,
					Name: element.Name,
				}
				result = append(result, temp)
			}
		}
		temp := &ListAllGroupVO{
			ID:   r.ID,
			Name: r.Name,
		}
		resp.Group = append(resp.Group, temp)
		result = result[1:]
	}
	resp.Count = len(resp.Group)
	return resp, err
}

// ListAllReq ListAllReq
type ListAllReq struct {
	AppID string `json:"appID"`
}

// ListAllResp ListAllResp
type ListAllResp struct {
	Menu  []*ListAllVO `json:"menu"`
	Count int          `json:"count"`
}

// ListAllVO ListAllVO
type ListAllVO struct {
	ID           string          `json:"id" `
	AppID        string          `json:"appID"`
	Name         string          `json:"name"`
	Icon         string          `json:"icon"`
	Sort         int             `json:"sort"`
	Describe     string          `json:"describe"`
	MenuType     models.Type     `json:"menuType"`
	GroupID      string          `json:"groupID"`
	BindingState models.BindType `json:"bindingState"`
	Child        []*ListAllVO    `json:"child"`
	ChildCount   int             `json:"childCount"`
}

func (m *menu) ListAll(ctx context.Context, req *ListAllReq) (*ListAllResp, error) {
	menu := &models.Menu{
		AppID:    req.AppID,
		MenuType: models.None,
	}
	// 查找顶层菜单
	menus, err := m.menuRepo.ListAllGroup(ctx, m.mongo, menu)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	top := make([]*models.Menu, 0)
	for _, m := range menus {
		if m.GroupID == "" {
			top = append(top, m)
		}
	}

	result, count := findAllMenus(top, menus)
	resp := &ListAllResp{
		Count: count,
		Menu:  result,
	}
	return resp, err
}

// TransferReq TransferReq
type TransferReq struct {
	ID          string `json:"id"`
	AppID       string `json:"appID"`
	Name        string `json:"name"`
	FromSort    int    `json:"fromSort"`
	ToSort      int    `json:"toSort"`
	FromGroupID string `json:"fromGroupID"`
	ToGroupID   string `json:"toGroupID"`
}

// TransferResp TransferResp
type TransferResp struct {
}

func (m *menu) Transfer(ctx context.Context, req *TransferReq) (*TransferResp, error) {

	resp := new(TransferResp)
	// 相同分组中移动
	if req.FromGroupID == req.ToGroupID {
		return transferSameGroup(ctx, m, req)
	}

	// 不同分组中移动
	if req.FromGroupID != req.ToGroupID {
		return transferDifferentGroup(ctx, m, req)
	}
	return resp, nil
}

// UserListAllReq UserListAllReq
type UserListAllReq struct {
	AppID   string         `json:"appID"`
	FormID  []string       `json:"formID"`
	PerType models.PerType `json:"perType"`
}

// UserListAllResp UserListAllResp
type UserListAllResp struct {
	ErrCode int          `json:"-"`
	Menu    []*ListAllVO `json:"menu"`
	Count   int          `json:"count"`
}

func (m *menu) UserListAll(ctx context.Context, req *UserListAllReq) (*UserListAllResp, error) {

	// 1、查找顶层菜单
	allMenus, err := m.menuRepo.ListAllGroup(ctx, m.mongo, &models.Menu{AppID: req.AppID, MenuType: models.None})
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	var userMenus []*models.Menu
	if req.PerType == models.InitType {
		userMenus = allMenus
	} else {
		ids := req.FormID
		// 2、根据ids查询所有的菜单
		userMenus, err = m.menuRepo.BatchFindMenus(ctx, m.mongo, ids)
	}
	if err != nil {
		return nil, err
	}

	// 剔除没有绑定的菜单
	var menus []*models.Menu
	for _, usermenu := range userMenus {
		if usermenu.BindingState != models.Bound {
			continue
		}
		menus = append(menus, usermenu)
	}

	top := make([]*models.Menu, 0)
	for _, m := range allMenus {
		if m.GroupID == "" && m.MenuType == models.GroupType {
			top = append(top, m)
		}
	}
	for _, userMenu := range menus {
		if userMenu.GroupID == "" && userMenu.MenuType != models.GroupType {
			top = append(top, userMenu)
		}
	}

	// 对顶层排序
	sort.SliceStable(top, func(i, j int) bool {
		return top[i].Sort < top[j].Sort
	})

	temp, count := findAllMenus(top, menus)

	// 剔除空分组
	result := make([]*ListAllVO, 0)
	for _, element := range temp {
		if element.ChildCount == 0 && element.MenuType == models.GroupType {
			count--
			continue
		}
		result = append(result, element)
	}
	resp := &UserListAllResp{
		Count: count,
		Menu:  result,
	}
	return resp, err
}

// ModifyMenuTypeReq ModifyMenuTypeReq
type ModifyMenuTypeReq struct {
	ID string `json:"id"`
}

// ModifyMenuTypeResp ModifyMenuTypeResp
type ModifyMenuTypeResp struct {
}

func (m *menu) ModifyMenuType(ctx context.Context, req *ModifyMenuTypeReq) (*ModifyMenuTypeResp, error) {
	err := m.menuRepo.ModifyMenuTypeByID(ctx, m.mongo, req.ID, models.CustomType)
	if err != nil {
		return nil, err
	}
	return &ModifyMenuTypeResp{}, nil
}

func findAllMenus(top, menus []*models.Menu) ([]*ListAllVO, int) {
	result := make([]*ListAllVO, 0)
	var count int
	for _, element := range top {
		if element.MenuType == models.MenuType || element.MenuType == models.CustomType {
			temp, num := findChildrenMenus(count, element)
			count = num
			result = append(result, temp)
		}
		if element.MenuType == models.GroupType {
			temp, num := findChildrenGroups(count, element, menus)
			count = num
			result = append(result, temp)
		}
	}
	return result, count
}
func findChildrenMenus(count int, element *models.Menu) (*ListAllVO, int) {
	count++
	temp := &ListAllVO{
		ID:           element.ID,
		AppID:        element.AppID,
		Sort:         element.Sort,
		Name:         element.Name,
		Describe:     element.Describe,
		Icon:         element.Icon,
		MenuType:     element.MenuType,
		GroupID:      element.GroupID,
		BindingState: element.BindingState,
	}

	return temp, count
}
func findChildrenGroups(count int, element *models.Menu, menus []*models.Menu) (*ListAllVO, int) {
	childMenus := make([]*models.Menu, 0)
	for _, m2 := range menus {
		if m2.GroupID == element.ID {
			childMenus = append(childMenus, m2)
		}
	}
	children, n := findAllMenus(childMenus, menus)
	temp := &ListAllVO{
		ID:           element.ID,
		AppID:        element.AppID,
		Sort:         element.Sort,
		Name:         element.Name,
		GroupID:      element.GroupID,
		MenuType:     element.MenuType,
		BindingState: element.BindingState,
		Child:        children,
		ChildCount:   n,
	}
	count += n + 1
	return temp, count
}
func transferSameGroup(ctx context.Context, m *menu, req *TransferReq) (*TransferResp, error) {
	resp := new(TransferResp)
	if req.ToSort == req.FromSort {
		return nil, nil
	}
	// 1、查询范围内的数据
	menus, err := m.menuRepo.FindAllFromRange(ctx, m.mongo, req.FromSort, req.ToSort, req.AppID, req.ToGroupID)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	// 2、修改sort
	err = m.menuRepo.UpdateSortByID(ctx, m.mongo, &models.Menu{
		Sort:    req.ToSort,
		ID:      req.ID,
		GroupID: req.ToGroupID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	// 3、修改范围内的数据
	ids := make([]string, len(menus))
	// 3、更新 [tSort,fSort) 范围的sort
	for i, temp := range menus {
		id := temp.ID
		ids[i] = id
	}
	// 向上移动
	if req.FromSort > req.ToSort {
		err := m.menuRepo.BatchUpdateSortByID(ctx, m.mongo, 1, ids)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
			return nil, err
		}
	}
	// 向下移动
	if req.FromSort < req.ToSort {
		err := m.menuRepo.BatchUpdateSortByID(ctx, m.mongo, -1, ids)
		if err != nil {
			logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
			return nil, err
		}
	}
	return resp, err
}

func transferDifferentGroup(ctx context.Context, m *menu, req *TransferReq) (*TransferResp, error) {
	resp := new(TransferResp)
	// 1、判断移动的地方是否存在相同名称菜单
	sameMenu, err := m.menuRepo.FindSameMenuName(ctx, m.mongo, &models.Menu{
		Name:    req.Name,
		GroupID: req.ToGroupID,
		AppID:   req.AppID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	if sameMenu != "" {
		return nil, error2.NewError(code.ErrRepeatMenuName)
	}
	// 2、查找 tGroup中的 [tGroup, +∞)中的数据并修改
	var sort int
	if req.ToSort == 0 {
		sort = req.ToSort + 1
	} else {
		sort = req.ToSort
	}
	menus, err := m.menuRepo.FindAllBySortAndGroup(ctx, m.mongo, &models.Menu{
		AppID:   req.AppID,
		GroupID: req.ToGroupID,
		Sort:    sort,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	toIDs := make([]string, len(menus))
	for i, temp := range menus {
		id := temp.ID
		toIDs[i] = id
	}
	err = m.menuRepo.BatchUpdateSortByID(ctx, m.mongo, 1, toIDs)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	// 3、修改移动的数据
	err = m.menuRepo.UpdateSortByID(ctx, m.mongo, &models.Menu{
		Sort:    sort,
		ID:      req.ID,
		GroupID: req.ToGroupID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	// 4、查找 fGroup中的 (fSort, +∞)中的数据并修改
	menus, err = m.menuRepo.FindAllBySortAndGroup(ctx, m.mongo, &models.Menu{
		AppID:   req.AppID,
		Sort:    req.FromSort,
		GroupID: req.FromGroupID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}
	fromIDs := make([]string, len(menus))
	for i, temp := range menus {
		id := temp.ID
		fromIDs[i] = id
	}

	err = m.menuRepo.BatchUpdateSortByID(ctx, m.mongo, -1, fromIDs)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(ctx))
		return nil, err
	}

	return resp, err
}

// FindByIDReq FindByIDReq
type FindByIDReq struct {
	ID string `json:"menuId" binding:"required"`
}

// FindByIDResp FindByIDResp
type FindByIDResp struct {
	ID       string      `json:"_id"`
	AppID    string      `json:"app_id"`
	Name     string      `json:"name"`
	Icon     string      `json:"icon"`
	Sort     int         `json:"sort"`
	MenuType models.Type `json:"menu_type"`
	Describe string      `bson:"describe"`
	GroupID  string      `bson:"group_id"`
}

func (m *menu) FindByID(ctx context.Context, req *FindByIDReq) (*FindByIDResp, error) {
	menu, err := m.menuRepo.FindByID(ctx, m.mongo, req.ID)
	if err != nil {
		return nil, err
	}
	return &FindByIDResp{
		ID:       menu.ID,
		AppID:    menu.ID,
		Name:     menu.Name,
		Icon:     menu.Icon,
		Sort:     menu.Sort,
		MenuType: menu.MenuType,
		Describe: menu.Describe,
		GroupID:  menu.GroupID,
	}, nil
}

// ListPageReq ListPageReq
type ListPageReq struct {
	AppID string `json:"appID"`
}

// ListPageResp ListPageResp
type ListPageResp struct {
	Pages []*ListPageVO `json:"pages"`
}

// ListPageVO ListPageVO
type ListPageVO struct {
	ID   string `json:"id" `
	Name string `json:"name"`
}

func (m *menu) ListPage(ctx context.Context, req *ListPageReq) (*ListPageResp, error) {
	menus, err := m.menuRepo.ListAllGroup(ctx, m.mongo, &models.Menu{AppID: req.AppID, MenuType: models.MenuType})
	if err != nil {
		return nil, err
	}
	var pages = make([]*ListPageVO, 0)
	for _, menu := range menus {
		if menu.MenuType == models.MenuType && menu.BindingState == models.Bound {
			pages = append(pages, &ListPageVO{
				ID:   menu.ID,
				Name: menu.Name,
			})
		}
	}
	return &ListPageResp{
		Pages: pages,
	}, nil
}
