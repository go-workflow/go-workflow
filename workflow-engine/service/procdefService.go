package service

import (
	"errors"
	"sync"
	"time"

	"github.com/go-workflow/go-workflow/workflow-engine/flow"

	"github.com/mumushuiding/util"

	"github.com/go-workflow/go-workflow/workflow-engine/model"
)

var saveLock sync.Mutex

// Procdef 流程定义表
type Procdef struct {
	Name string `json:"name"`
	// 流程定义json字符串
	Resource *flow.Node `json:"resource"`
	// 用户id
	Userid string `jsong:"userid"`
	// 用户所在公司
	Company   string `json:"company"`
	PageSize  int    `json:"pageSize"`
	PageIndex int    `json:"pageIndex"`
}

// GetProcdefByID 根据流程定义id获取流程定义
func GetProcdefByID(id int) (*model.Procdef, error) {
	return model.GetProcdefByID(id)
}

// GetProcdefLatestByNameAndCompany GetProcdefLatestByNameAndCompany
// 根据流程定义名字和公司查询流程定义
func GetProcdefLatestByNameAndCompany(name, company string) (*model.Procdef, error) {
	return model.GetProcdefLatestByNameAndCompany(name, company)
}

// GetResourceByNameAndCompany GetResourceByNameAndCompany
// 获取流程定义配置信息
func GetResourceByNameAndCompany(name, company string) (*flow.Node, int, string, error) {
	prodef, err := GetProcdefLatestByNameAndCompany(name, company)
	if err != nil {
		return nil, 0, "", err
	}
	if prodef == nil {
		return nil, 0, "", errors.New("流程【" + name + "】不存在")
	}
	node := &flow.Node{}
	err = util.Str2Struct(prodef.Resource, node)
	return node, prodef.ID, prodef.Name, err
}

// GetResourceByID GetResourceByID
// 根据id查询流程定义
func GetResourceByID(id int) (*flow.Node, int, error) {
	prodef, err := GetProcdefByID(id)
	if err != nil {
		return nil, 0, err
	}
	node := &flow.Node{}
	err = util.Str2Struct(prodef.Resource, node)
	return node, prodef.ID, err
}

// SaveProcdefByToken SaveProcdefByToken
func (p *Procdef) SaveProcdefByToken(token string) (int, error) {
	// 根据 token 获取用户信息
	userinfo, err := GetUserinfoFromRedis(token)
	if err != nil {
		return 0, err
	}
	if len(userinfo.Company) == 0 {
		return 0, errors.New("公司 company 不能为空")
	}
	if len(userinfo.Username) == 0 {
		return 0, errors.New("用户 username 不能为空")
	}
	p.Company = userinfo.Company
	p.Userid = userinfo.Username
	return p.SaveProcdef()
}

// SaveProcdef 保存
func (p *Procdef) SaveProcdef() (id int, err error) {
	// 流程定义有效性检验
	err = IsProdefValid(p.Resource)
	if err != nil {
		return 0, err
	}
	resource, err := util.ToJSONStr(p.Resource)
	if err != nil {
		return 0, err
	}
	// fmt.Println(resource)
	var procdef = model.Procdef{
		Name:     p.Name,
		Userid:   p.Userid,
		Company:  p.Company,
		Resource: resource,
	}
	return SaveProcdef(&procdef)
}

// SaveProcdef 保存
func SaveProcdef(p *model.Procdef) (id int, err error) {
	// 参数是否为空判定
	saveLock.Lock()
	defer saveLock.Unlock()
	old, err := GetProcdefLatestByNameAndCompany(p.Name, p.Company)
	if err != nil {
		return 0, err
	}
	p.DeployTime = util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)
	if old == nil {
		p.Version = 1
		return p.Save()
	}
	tx := model.GetTx()
	// 保存新版本
	p.Version = old.Version + 1
	err = p.SaveTx(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// 转移旧版本
	err = model.MoveProcdefToHistoryByIDTx(old.ID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return p.ID, nil
}

// ExistsProcdefByNameAndCompany if exists
// 查询流程定义是否存在
func ExistsProcdefByNameAndCompany(name, company string) (yes bool, version int, err error) {
	p, err := GetProcdefLatestByNameAndCompany(name, company)
	if p == nil {
		return false, 1, err
	}
	version = p.Version + 1
	return true, version, err
}

// FindAllPageAsJSON find by page and  transform result to string
// 分页查询并将结果转换成 json 字符串
func (p *Procdef) FindAllPageAsJSON() (string, error) {
	datas, count, err := p.FindAll()
	if err != nil {
		return "", err
	}
	return util.ToPageJSON(datas, count, p.PageIndex, p.PageSize)
}

// FindAll FindAll
func (p *Procdef) FindAll() ([]*model.Procdef, int, error) {
	var page = util.Page{}
	page.PageRequest(p.PageIndex, p.PageSize)
	maps := p.getMaps()
	return model.FindProcdefsWithCountAndPaged(page.PageIndex, page.PageSize, maps)
}
func (p *Procdef) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if len(p.Name) > 0 {
		maps["Name"] = p.Name
	}
	if len(p.Company) > 0 {
		maps["company"] = p.Company
	}
	return maps
}

// DelProcdefByID del by id
func DelProcdefByID(id int) error {
	return model.DelProcdefByID(id)
}

// IsProdefValid 流程定义格式是否有效
func IsProdefValid(node *flow.Node) error {

	return flow.IfProcessConifgIsValid(node)
}
