package service

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/go-workflow/go-workflow/workflow-engine/flow"

	"github.com/jinzhu/gorm"

	"github.com/go-workflow/go-workflow/workflow-engine/model"
	"github.com/mumushuiding/util"
)

// TaskReceiver 任务
type TaskReceiver struct {
	TaskID     int    `json:"taskID"`
	UserID     string `json:"userID,omitempty"`
	Pass       string `json:"pass,omitempty"`
	Company    string `json:"company,omitempty"`
	ProcInstID int    `json:"procInstID,omitempty"`
}

var completeLock sync.Mutex

// NewTask 新任务
func NewTask(t *model.Task) (int, error) {
	if len(t.NodeID) == 0 {
		return 0, errors.New("request param nodeID can not be null / 任务当前所在节点nodeId不能为空！")
	}
	t.CreateTime = util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)
	return t.NewTask()
}

// NewTaskTx NewTaskTx
// 开启事务
func NewTaskTx(t *model.Task, tx *gorm.DB) (int, error) {
	if len(t.NodeID) == 0 {
		return 0, errors.New("request param nodeID can not be null / 任务当前所在节点nodeId不能为空！")
	}
	t.CreateTime = util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)
	return t.NewTaskTx(tx)
}

// DeleteTask 删除任务
func DeleteTask(id int) error {
	return model.DeleteTask(id)
}

// GetTaskByID GetTaskById
func GetTaskByID(id int) (task *model.Task, err error) {
	return model.GetTaskByID(id)
}

// GetTaskLastByProInstID GetTaskLastByProInstID
func GetTaskLastByProInstID(procInstID int) (*model.Task, error) {
	return model.GetTaskLastByProInstID(procInstID)
}

// Complete Complete
// 审批
func Complete(taskID int, userID, company string, pass bool) error {
	tx := model.GetTx()
	err := CompleteTaskTx(taskID, userID, company, pass, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// UpdateTaskWhenComplete UpdateTaskWhenComplete
func UpdateTaskWhenComplete(taskID int, userID string, pass bool, tx *gorm.DB) (*model.Task, error) {
	// 获取task
	completeLock.Lock()         // 关锁
	defer completeLock.Unlock() //解锁
	// 查询任务
	task, err := GetTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("任务【" + fmt.Sprintf("%d", task.ID) + "】不存在")
	}
	// 判断是否已经结束
	if task.IsFinished == true {
		if task.NodeID == "结束" {
			return nil, errors.New("流程已经结束")
		}
		return nil, errors.New("任务【" + fmt.Sprintf("%d", taskID) + "】已经被审批过了！！")
	}
	// 设置处理人处理时间
	task.Assignee = userID
	task.ClaimTime = util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)
	// ----------------会签 （默认全部通过才结束），只要存在一个不通过，就结束，然后流转到上一步
	//同意
	if pass {
		task.AgreeNum++
	} else {
		task.IsFinished = true
	}
	// 未审批人数减一
	task.UnCompleteNum--
	// 判断是否结束
	if task.UnCompleteNum == 0 {
		task.IsFinished = true
	}
	err = task.UpdateTx(tx)
	// str, _ := util.ToJSONStr(task)
	// log.Println(str)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// CompleteTaskTx CompleteTaskTx
// 执行任务
func CompleteTaskTx(taskID int, userID, company string, pass bool, tx *gorm.DB) error {

	//更新任务
	task, err := UpdateTaskWhenComplete(taskID, userID, pass, tx)
	if err != nil {
		return err
	}
	// 如果是会签
	if task.ActType == "and" {
		// fmt.Println("------------------是会签，判断用户是否已经审批过了，避免重复审批-------")
		// 判断用户是否已经审批过了（存在会签的情况）
		yes, err := IfParticipantByTaskID(userID, company, taskID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if yes {
			tx.Rollback()
			return errors.New("您已经审批过了，请等待他人审批！）")
		}
	}

	// 查看任务的未审批人数是否为0，不为0就不流转
	if task.UnCompleteNum > 0 && pass == true { // 默认是全部通过
		// 添加参与人
		err := AddParticipantTx(userID, company, task.ID, task.ProcInstID, task.Step, tx)
		if err != nil {
			return err
		}
		return nil
	}
	// 流转到下一流程
	// nodeInfos, err := GetExecNodeInfosByProcInstID(task.ProcInstID)
	// if err != nil {
	// 	return err
	// }
	err = MoveStageByProcInstID(userID, company, task.ID, task.ProcInstID, task.Step, pass, tx)
	if err != nil {
		return err
	}

	return nil
}

// WithDrawTask 撤回任务
func WithDrawTask(taskID, procInstID int, userID, company string) error {
	var err1, err2 error
	var currentTask, lastTask *model.Task
	var timesx time.Time
	var wg sync.WaitGroup
	timesx = time.Now()
	wg.Add(2)
	go func() {
		currentTask, err1 = GetTaskByID(taskID)
		wg.Done()
	}()
	go func() {
		lastTask, err2 = GetTaskLastByProInstID(procInstID)
		wg.Done()
	}()
	wg.Wait()
	if err1 != nil {
		if err1 == gorm.ErrRecordNotFound {
			return errors.New("任务不存在")
		}
		return err1
	}
	if err2 != nil {
		if err2 == gorm.ErrRecordNotFound {
			return errors.New("找不到流程实例id为【" + fmt.Sprintf("%d", procInstID) + "】的任务，无权撤回")
		}
		return err2
	}
	if lastTask.Assignee != userID {
		return errors.New("只能撤回本人审批过的任务！！")
	}
	if currentTask.IsFinished {
		return errors.New("已经审批结束，无法撤回！")
	}
	if currentTask.UnCompleteNum != currentTask.MemberCount {
		return errors.New("已经有人审批过了，无法撤回！")
	}
	sub := currentTask.Step - lastTask.Step
	if math.Abs(float64(sub)) != 1 {
		return errors.New("只能撤回相邻的任务！！")
	}
	var pass = false
	if sub < 0 {
		pass = true
	}
	fmt.Printf("判断是否可以撤回,耗时：%v\n", time.Since(timesx))
	timesx = time.Now()
	tx := model.GetTx()
	// 更新当前的任务
	currentTask.IsFinished = true
	err := currentTask.UpdateTx(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// 撤回
	err = MoveStageByProcInstID(userID, company, currentTask.ID, procInstID, currentTask.Step, pass, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	fmt.Printf("撤回流程耗时：%v\n", time.Since(timesx))
	return nil
}

// MoveStageByProcInstID MoveStageByProcInstID
func MoveStageByProcInstID(userID, company string, taskID, procInstID, step int, pass bool, tx *gorm.DB) (err error) {
	nodeInfos, err := GetExecNodeInfosByProcInstID(procInstID)
	if err != nil {
		return err
	}
	return MoveStage(nodeInfos, userID, company, taskID, procInstID, step, pass, tx)
}

// MoveStage MoveStage
// 流程流转
func MoveStage(nodeInfos []*flow.NodeInfo, userID, company string, taskID, procInstID, step int, pass bool, tx *gorm.DB) (err error) {
	// 添加上一步的参与人
	err = AddParticipantTx(userID, company, taskID, procInstID, step, tx)
	if err != nil {
		return err
	}
	if pass {
		step++
	} else {
		step--
	}
	// 判断下一流程： 如果是审批人是：抄送人
	// fmt.Printf("下一审批人类型：%s\n", nodeInfos[step].AproverType)
	// fmt.Println(nodeInfos[step].AproverType == flow.NodeTypes[flow.NOTIFIER])
	if nodeInfos[step].AproverType == flow.NodeTypes[flow.NOTIFIER] {
		// 生成新的任务
		var task = model.Task{
			NodeID:     flow.NodeTypes[flow.NOTIFIER],
			Step:       step,
			ProcInstID: procInstID,
			IsFinished: true,
		}
		task.IsFinished = true
		_, err := task.NewTaskTx(tx)
		if err != nil {
			return err
		}
		// 添加抄送人
		err = AddNotifierTx(nodeInfos[step].Aprover, company, step, procInstID, tx)
		if err != nil {
			return err
		}
		return MoveStage(nodeInfos, userID, company, taskID, procInstID, step, pass, tx)
	}
	if pass {
		return MoveToNextStage(nodeInfos, userID, company, taskID, procInstID, step, tx)
	}
	return MoveToPrevStage(nodeInfos, userID, company, taskID, procInstID, step, tx)
}

// MoveToNextStage MoveToNextStage
//通过
func MoveToNextStage(nodeInfos []*flow.NodeInfo, userID, company string, currentTaskID, procInstID, step int, tx *gorm.DB) error {
	// fmt.Printf("-----------------流转到流程【%d】-----\n", step)
	// fmt.Printf("step=%d,total=%d\n", step, len(nodeInfos))
	if (step + 1) > len(nodeInfos) {
		return nil
	}
	var currentTime = util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)
	var task = getNewTask(nodeInfos, step, procInstID, currentTime) //新任务
	var procInst = &model.ProcInst{                                 // 流程实例要更新的字段
		NodeID:    nodeInfos[step].NodeID,
		Candidate: nodeInfos[step].Aprover,
	}
	procInst.ID = procInstID
	if (step + 1) != len(nodeInfos) { // 下一步不是【结束】
		// 生成新的任务
		taksID, err := task.NewTaskTx(tx)
		if err != nil {
			return err
		}
		// 添加candidate group
		err = AddCandidateGroupTx(nodeInfos[step].Aprover, company, step, taksID, procInstID, tx)
		if err != nil {
			return err
		}
		// 更新流程实例
		procInst.TaskID = taksID
		err = UpdateProcInst(procInst, tx)
		if err != nil {
			return err
		}
	} else { // 最后一步直接结束
		// 生成新的任务
		task.IsFinished = true
		task.ClaimTime = currentTime
		taksID, err := task.NewTaskTx(tx)
		if err != nil {
			return err
		}
		// 删除候选用户组
		err = DelCandidateByProcInstID(procInstID, tx)
		if err != nil {
			return err
		}
		// 更新流程实例
		procInst.TaskID = taksID
		procInst.EndTime = currentTime
		procInst.IsFinished = true
		err = UpdateProcInst(procInst, tx)
		if err != nil {
			return err
		}
	}
	// 添加上一步的参与人
	// err := AddParticipantTx(userID, company, currentTaskID, procInstID, step-1, tx)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// MoveToPrevStage MoveToPrevStage
// 驳回
func MoveToPrevStage(nodeInfos []*flow.NodeInfo, userID, company string, currentTaskID, procInstID, step int, tx *gorm.DB) error {
	// fmt.Printf("------------------流转到流程【%d】", step)
	if step == -1 {
		return errors.New("流程处于开始位置无法驳回！！")
	}
	// 生成新的任务
	var task = getNewTask(nodeInfos, step, procInstID, util.FormatDate(time.Now(), util.YYYY_MM_DD_HH_MM_SS)) //新任务
	taksID, err := task.NewTaskTx(tx)
	if err != nil {
		return err
	}
	// 添加上一步的参与人
	// err = AddParticipantTx(userID, company, currentTaskID, procInstID, step+1, tx)
	// if err != nil {
	// 	return err
	// }
	// 更新流程实例
	var procInst = &model.ProcInst{ // 流程实例要更新的字段
		NodeID:    nodeInfos[step].NodeID,
		Candidate: nodeInfos[step].Aprover,
	}
	procInst.ID = procInstID
	err = UpdateProcInst(procInst, tx)
	if err != nil {
		return err
	}
	if step == 0 { // 流程回到起始位置，注意起始位置为0,
		err = AddCandidateUserTx(nodeInfos[step].Aprover, company, step, taksID, procInstID, tx)
		if err != nil {
			return err
		}
		return nil
	}
	// 添加candidate group
	err = AddCandidateGroupTx(nodeInfos[step].Aprover, company, step, taksID, procInstID, tx)
	if err != nil {
		return err
	}
	return nil
}
func getNewTask(nodeInfos []*flow.NodeInfo, step, procInstID int, currentTime string) *model.Task {
	var task = &model.Task{ // 新任务
		NodeID:        nodeInfos[step].NodeID,
		Step:          step,
		CreateTime:    currentTime,
		ProcInstID:    procInstID,
		MemberCount:   nodeInfos[step].MemberCount,
		UnCompleteNum: nodeInfos[step].MemberCount,
		ActType:       nodeInfos[step].ActType,
	}
	return task
}
