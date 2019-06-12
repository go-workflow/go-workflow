package model

import (
	"sync"

	"github.com/jinzhu/gorm"
)

// ProcInst 流程实例
type ProcInst struct {
	Model
	// 流程定义ID
	ProcDefID int `json:"procDefId"`
	// 流程定义名
	ProcDefName string `json:"procDefName"`
	// title 标题
	Title string `json:"title"`
	// 用户部门
	Department string `json:"department"`
	Company    string `json:"company"`
	// 当前节点
	NodeID string `json:"nodeID"`
	// 审批人
	Candidate string `json:"candidate"`
	// 当前任务
	TaskID      int    `json:"taskID"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	Duration    int64  `json:"duration"`
	StartUserID string `json:"startUserId"`
	IsFinished  bool   `gorm:"default:false" json:"isFinished"`
}

// GroupsNotNull 候选组
func GroupsNotNull(groups []string, company string) func(db *gorm.DB) *gorm.DB {
	if len(groups) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Or("is_finished=0 and candidate in (?) and company=?", groups, company)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// DepartmentsNotNull 分管部门
func DepartmentsNotNull(departments []string, company string) func(db *gorm.DB) *gorm.DB {
	if len(departments) > 0 {
		return func(db *gorm.DB) *gorm.DB {
			return db.Or("is_finished=0 and department in (?) and candidate=? and company=?", departments, IdentityTypes[MANAGER], company)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}

// FindProcInsts FindProcInsts
// 分页查询
func FindProcInsts(userID, company string, groups, departments []string, pageIndex, pageSize int) ([]*ProcInst, int, error) {

	var datas []*ProcInst
	var count int
	//----普通写法--------
	// err := db.Scopes(GroupsNotNull(groups), DepartmentsNotNull(departments)).Or("is_finished=0 and candidate=?", userID).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&datas).Error
	// if err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return nil, 0, nil
	// 	}
	// 	return nil, 0, err
	// }
	// err = db.Scopes(GroupsNotNull(groups), DepartmentsNotNull(departments)).Model(&ProcInst{}).Or("is_finished=0 and candidate=?", userID).Count(&count).Error
	// if err != nil {
	// 	return nil, 0, err
	// }
	// return datas, count, nil
	//--------使用chanel写法
	selectDatas := func(in chan<- error, wg *sync.WaitGroup) {
		go func() {
			err := db.Scopes(GroupsNotNull(groups, company), DepartmentsNotNull(departments, company)).Or("is_finished=0 and candidate=? and company=?", userID, company).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&datas).Error
			in <- err
			wg.Done()
		}()
	}
	selectCount := func(in chan<- error, wg *sync.WaitGroup) {
		go func() {
			err := db.Scopes(GroupsNotNull(groups, company), DepartmentsNotNull(departments, company)).Model(&ProcInst{}).Or("is_finished=0 and candidate=? and company=?", userID, company).Count(&count).Error
			in <- err
			wg.Done()
		}()
	}
	var err1 error
	var wg sync.WaitGroup
	numberOfRoutine := 2
	wg.Add(numberOfRoutine)
	errStream := make(chan error, numberOfRoutine)
	// defer fmt.Println("close channel")
	selectDatas(errStream, &wg)
	selectCount(errStream, &wg)
	wg.Wait()
	defer close(errStream) // 关闭通道

	for i := 0; i < numberOfRoutine; i++ {
		// log.Printf("send: %v", <-errStream)
		if err := <-errStream; err != nil {
			err1 = err
		}
	}
	// fmt.Println("结束")
	return datas, count, err1
}

// Save save
func (p *ProcInst) Save() (int, error) {
	err := db.Create(p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

//SaveTx SaveTx
func (p *ProcInst) SaveTx(tx *gorm.DB) (int, error) {
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return p.ID, nil
}

// DelProcInstByID DelProcInstByID
func DelProcInstByID(id int) error {
	return db.Where("id=?", id).Delete(&ProcInst{}).Error
}

// DelProcInstByIDTx DelProcInstByIDTx
// 事务
func DelProcInstByIDTx(id int, tx *gorm.DB) error {
	return db.Where("id=?", id).Delete(&ProcInst{}).Error
}

// UpdateTx UpdateTx
func (p *ProcInst) UpdateTx(tx *gorm.DB) error {
	return tx.Model(&ProcInst{}).Updates(p).Error
}

// FindFinishedProc FindFinishedProc
func FindFinishedProc() ([]*ProcInst, error) {
	var datas []*ProcInst
	err := db.Where("is_finished=1").Find(&datas).Error
	return datas, err
}
