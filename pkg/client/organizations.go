package client

import (
	"context"
	"net/http"

	"git.internal.yunify.com/qxp/misc/client"
)

const (
	organizationsHost = "http://org/api/v1/org"
	depTree           = "/DEPTree"
	depByIDs          = "/depByIDs"
	usersInfo         = "/usersInfo"
	adminUserList     = "/adminUserList"
)

// Organizations 对外接口
type Organizations interface {
	DEPTree(ctx context.Context, req *DEPTreeReq) (*DEPTreeResp, error)
	SelectDepByIDs(ctx context.Context, req *SelectDepByIDsReq) (*SelectDepByIDsResp, error)
	UserSelectByIDs(ctx context.Context, req *UserSelectByIDsReq) (*UserSelectByIDsResp, error)
	AdminSelectByCondition(ctx context.Context, req *AdminSelectByConditionReq) (*AdminSelectByConditionResp, error)
}
type organizations struct {
	client http.Client
}

// NewOrganizations create organization
func NewOrganizations(conf client.Config) Organizations {
	return &organizations{
		client: client.New(conf),
	}
}

// DEPTreeReq DEPTreeReq
type DEPTreeReq struct {
}

// DEPTreeResp DEPTreeResp
type DEPTreeResp struct {
	ID                 string        `json:"id,omitempty"`
	DepartmentName     string        `json:"departmentName,omitempty"`
	DepartmentLeaderID string        `json:"departmentLeaderID,omitempty"`
	UseStatus          int           `json:"useStatus,omitempty"`
	PID                string        `json:"pid,omitempty"`       //上层ID
	SuperPID           string        `json:"superID,omitempty"`   //最顶层父级ID
	CompanyID          string        `json:"companyID,omitempty"` //所属公司id
	Grade              int           `json:"grade,omitempty"`     //部门等级
	Child              []DEPTreeResp `json:"child"`
}

// DEPTree 部门树
func (o *organizations) DEPTree(ctx context.Context, req *DEPTreeReq) (*DEPTreeResp, error) {

	departmentTree := &DEPTreeResp{}
	err := client.POST(ctx, &o.client, organizationsHost+depTree, req, departmentTree)
	return departmentTree, err
}

// SelectDepByIDsReq SelectDepByIDsReq
type SelectDepByIDsReq struct {
	IDs []string `json:"ids"`
}

// SelectDepByIDsResp SelectDepByIDsResp
type SelectDepByIDsResp struct {
	UserDepartments []*UserDepartment
}

// UserDepartment 返回前端用户查看的应用结构体
type UserDepartment struct {
	ID                 string `json:"id,omitempty"`
	DepartmentName     string `json:"departmentName,omitempty"`
	DepartmentLeaderID string `json:"departmentLeaderID,omitempty"`
	UseStatus          int    `json:"useStatus,omitempty"`
	PID                string `json:"pid,omitempty"`       //上层ID
	SuperPID           string `json:"superID,omitempty"`   //最顶层父级ID
	CompanyID          string `json:"companyID,omitempty"` //所属公司id
	Grade              int    `json:"grade,omitempty"`     //部门等级
}

func (o *organizations) SelectDepByIDs(ctx context.Context, req *SelectDepByIDsReq) (*SelectDepByIDsResp, error) {
	params := req
	userDepartments := make([]*UserDepartment, 0)
	err := client.POST(ctx, &o.client, organizationsHost+depByIDs, params, &userDepartments)

	return &SelectDepByIDsResp{
		UserDepartments: userDepartments,
	}, err
}

// UserSelectByIDsReq UserSelectByIDsReq
type UserSelectByIDsReq struct {
	IDs []string `json:"ids"`
}

// UserSelectByIDsResp UserSelectByIDsResp
type UserSelectByIDsResp struct {
	UserUsersInfo []*UserUser
}

// UserUser 用户可见字段
type UserUser struct {
	ID          string    `json:"id,omitempty"`
	UserName    string    `json:"userName,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Email       string    `json:"email,omitempty"`
	Address     string    `json:"address,omitempty"`
	LeaderID    string    `json:"leaderID,omitempty"`
	CompanyID   string    `json:"companyID,omitempty"`   //所属公司id
	Avatar      string    `json:"avatar"`                //用户头像
	IsDEPLeader int       `json:"isDEPLeader,omitempty"` //是否部门领导 1是，-1不是
	Status      int       `json:"status"`                //第一位：密码是否需要重置
	DEP         DEPOnline `json:"dep,omitempty"`         //用户所在部门
}

// DEPOnline 用于用户部门层级线索
type DEPOnline struct {
	ID                 string     `json:"id,omitempty"`
	DepartmentName     string     `json:"departmentName"`
	DepartmentLeaderID string     `json:"departmentLeaderID"`
	UseStatus          int        `json:"useStatus,omitempty"`
	PID                string     `json:"pid"`                 //上层ID
	SuperPID           string     `json:"superID,omitempty"`   //最顶层父级ID
	CompanyID          string     `json:"companyID,omitempty"` //所属公司id
	Grade              int        `json:"grade,omitempty"`     //部门等级
	Child              *DEPOnline `json:"child"`
}

func (o *organizations) UserSelectByIDs(ctx context.Context, req *UserSelectByIDsReq) (*UserSelectByIDsResp, error) {
	params := req
	userUsersInfo := make([]*UserUser, 0)
	err := client.POST(ctx, &o.client, organizationsHost+usersInfo, params, &userUsersInfo)

	return &UserSelectByIDsResp{
		UserUsersInfo: userUsersInfo,
	}, err
}

// AdminSelectByConditionReq AdminSelectByConditionReq
type AdminSelectByConditionReq struct {
	DepID                string `json:"depID" binding:"max=64"`
	IncludeChildDEPChild int    `json:"includeChildDEPChild" ` //是否包含当前部门的子部门 1是包含，其它都不包含
	UserName             string `json:"userName" binding:"max=80,excludesall=0x2C!@#$?.%:*&^+><=；;"`
	Phone                string `json:"phone" binding:"max=64"`
	Email                string `json:"email" binding:"max=64"`
	IDCard               string `json:"idCard" binding:"max=64"`
	LeaderID             string `json:"leaderID" binding:"max=64"`
	UseStatus            int    `json:"useStatus"`                   //状态：1正常，-2禁用，-1删除 （与账号库相同）
	CompanyID            string `json:"companyID"  binding:"max=64"` //所属公司id
	Page                 int    `json:"page"`
	Limit                int    `json:"limit"`
}

// AdminSelectByConditionResp AdminSelectByConditionReq
type AdminSelectByConditionResp struct {
	PageSize    int         `json:"-"`
	TotalCount  int64       `json:"total_count"`
	TotalPage   int         `json:"-"`
	CurrentPage int         `json:"-"`
	StartIndex  int         `json:"-"`
	Data        interface{} `json:"data"`
}

func (o *organizations) AdminSelectByCondition(ctx context.Context, req *AdminSelectByConditionReq) (*AdminSelectByConditionResp, error) {
	params := req
	resp := &AdminSelectByConditionResp{}
	err := client.POST(ctx, &o.client, organizationsHost+adminUserList, params, resp)
	return resp, err
}
