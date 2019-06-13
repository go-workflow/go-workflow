package service

import (
	"github.com/go-workflow/go-workflow/workflow-engine/model"
	"github.com/jinzhu/gorm"
)

// SaveIdentitylinkTx SaveIdentitylinkTx
func SaveIdentitylinkTx(i *model.Identitylink, tx *gorm.DB) error {
	return i.SaveTx(tx)
}

// AddNotifierTx 添加抄送人候选用户组
func AddNotifierTx(group, company string, step, procInstID int, tx *gorm.DB) error {
	yes, err := ExistsNotifierByProcInstIDAndGroup(procInstID, group)
	if err != nil {
		return err
	}
	if yes {
		return nil
	}
	i := &model.Identitylink{
		Group:      group,
		Type:       model.IdentityTypes[model.NOTIFIER],
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return SaveIdentitylinkTx(i, tx)
}

// AddCandidateGroupTx AddCandidateGroupTx
// 添加候选用户组
func AddCandidateGroupTx(group, company string, step, taskID, procInstID int, tx *gorm.DB) error {
	err := DelCandidateByProcInstID(procInstID, tx)
	if err != nil {
		return err
	}
	i := &model.Identitylink{
		Group:      group,
		Type:       model.IdentityTypes[model.CANDIDATE],
		TaskID:     taskID,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return SaveIdentitylinkTx(i, tx)
}

// AddCandidateUserTx AddCandidateUserTx
// 添加候选用户
func AddCandidateUserTx(userID, company string, step, taskID, procInstID int, tx *gorm.DB) error {
	err := DelCandidateByProcInstID(procInstID, tx)
	if err != nil {
		return err
	}
	i := &model.Identitylink{
		UserID:     userID,
		Type:       model.IdentityTypes[model.CANDIDATE],
		TaskID:     taskID,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return SaveIdentitylinkTx(i, tx)
	// var wg sync.WaitGroup
	// var err1, err2 error
	// numberOfRoutine := 2
	// wg.Add(numberOfRoutine)
	// go func() {
	// 	defer wg.Done()
	// 	err1 = DelCandidateByProcInstID(procInstID, tx)
	// }()
	// go func() {
	// 	defer wg.Done()
	// 	i := &model.Identitylink{
	// 		UserID:     userID,
	// 		Type:       model.IdentityTypes[model.CANDIDATE],
	// 		TaskID:     taskID,
	// 		Step:       step,
	// 		ProcInstID: procInstID,
	// 		Company:    company,
	// 	}
	// 	err2 = SaveIdentitylinkTx(i, tx)
	// }()
	// wg.Wait()
	// fmt.Println("保存identyilink结束")
	// if err1 != nil {
	// 	return err1
	// }
	// return err2
}

//AddParticipantTx AddParticipantTx
// 添加任务参与人
func AddParticipantTx(userID, company string, taskID, procInstID, step int, tx *gorm.DB) error {
	i := &model.Identitylink{
		Type:       model.IdentityTypes[model.PARTICIPANT],
		UserID:     userID,
		ProcInstID: procInstID,
		Step:       step,
		Company:    company,
		TaskID:     taskID,
	}
	return SaveIdentitylinkTx(i, tx)
}

// IfParticipantByTaskID IfParticipantByTaskID
// 针对指定任务判断用户是否已经审批过了
func IfParticipantByTaskID(userID, company string, taskID int) (bool, error) {
	return model.IfParticipantByTaskID(userID, company, taskID)
}

// DelCandidateByProcInstID DelCandidateByProcInstID
// 删除历史候选人
func DelCandidateByProcInstID(procInstID int, tx *gorm.DB) error {
	return model.DelCandidateByProcInstID(procInstID, tx)
}

// ExistsNotifierByProcInstIDAndGroup 抄送人是否已经存在
func ExistsNotifierByProcInstIDAndGroup(procInstID int, group string) (bool, error) {
	return model.ExistsNotifierByProcInstIDAndGroup(procInstID, group)
}
