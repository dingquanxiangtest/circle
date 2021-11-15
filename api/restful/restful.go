package restful

import (
	"net/http"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mongo"
	"git.internal.yunify.com/qxp/molecule/internal/comet"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

const (
	managerPath = "manager"
	homePath    = "home"
	flowPath    = "flow"
	processPath = "process"
	basePath    = "base"
)

// Router 路由
type Router struct {
	c *config.Config

	engine *gin.Engine

	// 流程引擎与表单的交互，单独开端口，校验逻辑不同
	processEngine *gin.Engine
}
type router func(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error

var routers = []router{
	cometRouter,
	tableRouter,
	//权限
	permissionRouter,
	menuRouter,
	flowRouter,
	datasetRouter,
	organizationsRouter,
	formulaRouter,
	customPageRouter,
}

// NewRouter 开启路由
func NewRouter(c *config.Config) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	processEngine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	//appCenterClient := client.NewAppCenterClient(c)
	var router = map[string]*gin.RouterGroup{
		//appCenterClient.CheckIsAppAdmin
		managerPath: engine.Group("/api/v1/structor/:appID/m"),
		homePath:    engine.Group("/api/v1/structor/:appID/home"),
		flowPath:    engine.Group("/api/v1/structor/:appID/flow"),
		processPath: processEngine.Group("/api/v1/structor/"),
		basePath:    engine.Group("/api/v1/structor/:appID/base"),
	}
	client, err := mongo.New(&c.Mongo)
	if err != nil {
		return nil, err
	}
	opt := service.WithMongo(client, c.Service.DB)

	for _, f := range routers {
		err = f(c, router, opt)
		if err != nil {
			return nil, err
		}
	}
	// 流程相关（前端调用的走外网）
	return &Router{
		c:             c,
		engine:        engine,
		processEngine: processEngine,
	}, nil
}

func newRouter(c *config.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(logger.GinLogger(),
		logger.GinRecovery())

	return engine, nil
}

func flowRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	process, err := NewProcess(c, opts...)
	if err != nil {
		return err
	}
	processGroup := r[managerPath].Group("/process")
	{
		processGroup.POST("/getByID", process.GetSchema)

	}

	// 流程相关接口只内网暴露服务，单独处理
	processGroupIn := r[processPath].Group("/process")
	{
		processGroupIn.POST("/getByID", process.GetRawSchema)
		processGroupIn.POST("/getSubTable", process.GetSubTable)

		processGroupIn.POST("/getData", process.GetData)
		processGroupIn.POST("/batchGetData", process.BatchGetData)
		processGroupIn.POST("/saveProcessData", process.UpdateProcessData)
		processGroupIn.POST("/createProcessData", process.CreateProcessData)

	}
	// 不需要鉴权的接口
	noAuth, err := comet.NewNoAuth(c, opts...)
	if err != nil {
		return err
	}
	//noAuthGroup := r[processPath].Group("/:appID/form/:tableName/Search",comet.HandleSchema(auth))
	noAuthGroup := r[processPath].Group("/noAuth/:appID/form/:tableName")
	{
		noAuthGroup.POST("/search", comet.SearchData(noAuth))
		noAuthGroup.POST("/get", comet.SearchData(noAuth))
		noAuthGroup.POST("/create", comet.CreateData(noAuth))
		noAuthGroup.POST("/update", comet.UpdateData(noAuth))
		noAuthGroup.POST("/delete", comet.SearchData(noAuth))
	}
	return nil
}

func menuRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	menu, err := NewMenu(c, opts...)
	if err != nil {
		return err
	}
	menuManger := r[managerPath].Group("/menu")
	{
		menuManger.POST("/create", menu.CreateMenu)
		menuManger.POST("/delete", menu.DeleteMenu)
		menuManger.POST("/update", menu.UpdateMenu)
		menuManger.POST("/list", menu.ListAll)
		menuManger.POST("/transfer", menu.Transfer)
		menuManger.POST("/listPage", menu.ListPage)
	}
	menuHome := r[homePath].Group("/menu")
	{
		menuHome.POST("/list", menu.UserListAll)
	}

	groupManger := r[managerPath].Group("/group")
	{

		groupManger.POST("/create", menu.CreateGroup)
		groupManger.POST("/delete", menu.DeleteGroup)
		groupManger.POST("/list", menu.ListAllGroup)

	}

	return nil
}

func datasetRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	dataset, err := NewDataSet(c, opts...)
	if err != nil {
		return err
	}
	datasetHome := r[homePath].Group("", func(c *gin.Context) {
		if c.Param("appID") != "dataset" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})
	datasetManager := r[managerPath].Group("", func(c *gin.Context) {
		if c.Param("appID") != "dataset" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})
	{
		datasetHome.POST("/get", dataset.GetDataSet)
		datasetManager.POST("/create", dataset.CreateDataSet)
		datasetManager.POST("/get", dataset.GetDataSet)
		datasetManager.POST("/update", dataset.UpdateDataSet)
		datasetManager.POST("/getByCondition", dataset.GetByConditionSet)
		datasetManager.POST("/delete", dataset.DeleteDataSet)

	}
	return nil
}
func cometRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	comet1, err := comet.New(c, opts...)
	if err != nil {
		return err
	}
	cometHome := r[homePath].Group("/form/:tableName")
	{
		cometHome.POST("", comet1.Handle)
		//cometHome.POST("/search", comet1.Search)
	}

	subCometHome := r[homePath].Group("/subForm/:tableName")
	{
		subCometHome.POST("", comet1.HandleWithoutAuth)
	}

	// 不需要鉴权的接口
	auth, err := comet.NewAuth(c, opts...)
	if err != nil {
		return err
	}
	{
		cometHome.POST("/search", comet.SearchData(auth))
		cometHome.POST("/get", comet.SearchData(auth))
		cometHome.POST("/create", comet.CreateData(auth))
		cometHome.POST("/update", comet.UpdateData(auth))
		cometHome.POST("/delete", comet.DeleteData(auth))
	}

	return nil
}

func tableRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	table, err := comet.NewTable(c, opts...)
	if err != nil {
		return err
	}

	schemaGroup := r[homePath].Group("/schema/:tableName")
	{
		schemaGroup.POST("", table.Handle)
	}

	kernel, err := NewKernel(c, opts...)
	if err != nil {
		return err
	}
	schemaManager := r[managerPath].Group("/table")
	{
		schemaManager.POST("/create", kernel.CreateSchema)
		schemaManager.POST("/getByID", kernel.GetSchema)
		schemaManager.POST("/delete", kernel.DeleteSchema)
		schemaManager.POST("/createBlank", kernel.CreateBlankSchema)
		schemaManager.POST("/search", kernel.SearchSchema)
		schemaManager.POST("/getXName", kernel.GetXName)
		schemaManager.POST("/checkRepeat", kernel.CheckRepeat)
		schemaManager.POST("/getInfo", kernel.GetModelDataByMenu)

	}

	managerConfig := r[managerPath].Group("/config")
	{
		managerConfig.POST("/create", kernel.CreateOrUpdateConfig)
		managerConfig.POST("/update", kernel.CreateOrUpdateConfig)
	}

	// sub table
	subTable, err := NewSubTable(c, opts...)
	if err != nil {
		return err
	}
	subTableGroup := r[managerPath].Group("/subTable")
	{
		subTableGroup.POST("/create", subTable.CreateSubTable)
		subTableGroup.POST("/getByID", subTable.GetSubTable)
		subTableGroup.POST("/getByCondition", subTable.GetSubTables)
		subTableGroup.POST("/update", subTable.UpdateSubTable)
		subTableGroup.POST("/delete", subTable.DeleteSubTable)
	}
	subTableHomeGroup := r[homePath].Group("/subTable")
	{
		subTableHomeGroup.POST("/getByCondition", subTable.GetSubTables)
		subTableHomeGroup.POST("/getByType", subTable.GetSubTablesByType)
		subTableHomeGroup.POST("/getByID", subTable.GetSubTable)
	}

	return nil
}

func permissionRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	permission, err := NewPermission(c, opts...)
	if err != nil {
		return err
	}

	baseGroup := r[basePath].Group("/permission")
	{
		baseGroup.POST("/perGroup/create", permission.CreatePermissionGroup) //  创建权限组
		baseGroup.POST("/perGroup/update", permission.UpdatePermissionGroup, permission.AddAppDepUser)

	}
	manager := r[managerPath].Group("/permission")
	{
		manager.POST("/perGroup/create", permission.CreatePermissionGroup) //  创建权限组
		manager.POST("/perGroup/updateName", permission.UpdatePerName)     //  更新权限组
		manager.POST("/perGroup/update", permission.UpdatePermissionGroup, permission.AddAppDepUser)
		manager.POST("/perGroup/delete", permission.DeletePermissionGroup)          // 删除权限组
		manager.POST("/perGroup/getByID", permission.GetByIDPermissionGroup)        // 根据id 获取权限组
		manager.POST("/perGroup/getList", permission.GetListPermissionGroup)        // 根据条件获取 权限组列表
		manager.POST("/perGroup/getByCondition", permission.GetByConditionPerGroup) // 根据用户id ，或者部门ID ，得到权限组id

		manager.POST("/perGroup/getPerData", permission.GetPerData)
		manager.POST("/perGroup/saveForm", permission.SaveForm)     // 保存表单权限
		manager.POST("/perGroup/deleteForm", permission.DeleteForm) // 删除表单权限
		manager.POST("/perGroup/getForm", permission.GetForm)
		manager.POST("/perGroup/updatePage", permission.SavePagePermission)
		manager.POST("/perGroup/pageList", permission.GetPagePermission)
		manager.POST("/perGroup/updatePagePer", permission.ModifyPagePermission)
		manager.POST("/perGroup/getPerGroupByMenu", permission.GetPerGroupsByMenu)
	}

	home := r[homePath].Group("/permission")
	{
		home.POST("/operatePer/getOperate", permission.GetOperate) // 跟据用户id 和 部门ID，得到操作权限
		home.POST("/perGroup/getPerOption", permission.GetPerOption)
		home.POST("/perGroup/saveUserPerMatch", permission.SaveUserPerMatch) // 保存用户匹配的权限组
	}

	return nil
}

func organizationsRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	organizations := NewOrganizations(c)
	home := r[homePath].Group("/org")
	{
		home.POST("/DEPTree", organizations.DEPTree)                // 部门树
		home.POST("/depByIDs", organizations.SelectDepByIDs)        // 根据IDs查询部门信息
		home.POST("/onlineUserDep", organizations.OnlineUserDep)    // 查询在线用户部门信息
		home.POST("/userList", organizations.SelectUserByCondition) // 根据条件获取用户信息
		home.POST("/usersInfoByIDs", organizations.UserUsersInfo)   // 根据IDs查询用户信息
		home.POST("/onlineUserInfo", organizations.OnlineUserInfo)  //  查询在线用户信息
	}
	return nil
}

func formulaRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	formula, err := NewFormula(c, opts...)
	if err != nil {
		return err
	}
	formulaHome := r[homePath].Group("", func(c *gin.Context) {
		if c.Param("appID") != "formula" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})
	formulaManager := r[managerPath].Group("", func(c *gin.Context) {
		if c.Param("appID") != "formula" {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	})
	{
		formulaHome.POST("/calculation", formula.Calculation)
		formulaManager.POST("/calculation", formula.Calculation)
	}
	return nil
}

func customPageRouter(c *config.Config, r map[string]*gin.RouterGroup, opts ...service.Options) error {
	customPage, err := NewCustomPage(c, opts...)
	if err != nil {
		return err
	}
	manager := r[managerPath].Group("/page")
	{
		manager.POST("/create", customPage.CreateCustomPage)
		manager.POST("/update", customPage.UpdateCustomPage)
		manager.POST("/getByMenu", customPage.GetByMenuID)
	}
	return nil
}

// Run 启动服务
func (r *Router) Run() {
	go r.processEngine.Run(r.c.ProcessPort)
	r.engine.Run(r.c.Port)
}

// Close 关闭服务
func (r *Router) Close() {
}
