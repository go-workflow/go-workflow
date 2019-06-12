package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/jinzhu/gorm"

	"github.com/mumushuiding/util"

	"github.com/go-workflow/go-workflow/workflow-engine/flow"
	"github.com/go-workflow/go-workflow/workflow-engine/model"
)

var execLock sync.Mutex

// SaveExecution SaveExecution
func SaveExecution(e *model.Execution) (ID int, err error) {
	execLock.Lock()
	defer execLock.Unlock()
	// check if exists by procInst
	yes, err := model.ExistsExecByProcInst(e.ProcInstID)
	if err != nil {
		return 0, err
	}
	if yes {
		return 0, errors.New("流程实例【" + fmt.Sprintf("%d", e.ProcInstID) + "】已经存在执行流")
	}
	// save
	return e.Save()
}

// SaveExecTx SaveExecTx
func SaveExecTx(e *model.Execution, tx *gorm.DB) (ID int, err error) {
	execLock.Lock()
	defer execLock.Unlock()
	// check if exists by procInst
	yes, err := model.ExistsExecByProcInst(e.ProcInstID)
	if err != nil {
		return 0, err
	}
	if yes {
		return 0, errors.New("流程实例【" + fmt.Sprintf("%d", e.ProcInstID) + "】已经存在执行流")
	}
	// save
	return e.SaveTx(tx)
}

// GetExecByProcInst GetExecByProcInst
// 根据流程实例查询执行流
func GetExecByProcInst(procInst int) (*model.Execution, error) {
	return model.GetExecByProcInst(procInst)
}

// GenerateExec GenerateExec
// 根据流程定义node生成执行流
func GenerateExec(e *model.Execution, node *flow.Node, userID string, variable *map[string]string, tx *gorm.DB) (int, error) {
	list, err := flow.ParseProcessConfig(node, variable)
	if err != nil {
		return 0, err
	}
	list.PushBack(flow.NodeInfo{
		NodeID: "结束",
	})
	list.PushFront(flow.NodeInfo{
		NodeID:  "开始",
		Type:    flow.NodeInfoTypes[flow.STARTER],
		Aprover: userID,
	})
	arr := util.List2Array(list)
	str, err := util.ToJSONStr(arr)
	if err != nil {
		return 0, err
	}
	e.NodeInfos = str
	ID, err := SaveExecTx(e, tx)
	return ID, err
}

// GetExecNodeInfosByProcInstID GetExecNodeInfosByProcInstID
// 获取执行流经过的节点信息
func GetExecNodeInfosByProcInstID(procInstID int) ([]*flow.NodeInfo, error) {
	nodeinfoStr, err := model.GetExecNodeInfosByProcInstID(procInstID)
	if err != nil {
		return nil, err
	}
	var nodeInfos []*flow.NodeInfo
	err = util.Str2Struct(nodeinfoStr, &nodeInfos)
	return nodeInfos, err
}
