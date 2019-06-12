package model

import (
	"github.com/jinzhu/gorm"
)

// Identitylink 用户组同任务的关系
type Identitylink struct {
	Model
	Group      string `json:"group"`
	Type       string `json:"type"`
	UserID     string `json:"userid"`
	TaskID     int    `json:"taskID"`
	Step       int    `json:"step"`
	ProcInstID int    `json:"procInstID"`
	Company    string `json:"company"`
}

// IdentityType 类型
type IdentityType int

const (
	// CANDIDATE 候选
	CANDIDATE IdentityType = iota
	// PARTICIPANT 参与人
	PARTICIPANT
	// MANAGER 上级领导
	MANAGER
)

// IdentityTypes IdentityTypes
var IdentityTypes = [...]string{CANDIDATE: "candidate", PARTICIPANT: "participant", MANAGER: "主管"}

// SaveTx SaveTx
func (i *Identitylink) SaveTx(tx *gorm.DB) error {
	// if len(i.Company) == 0 {
	// 	return errors.New("Identitylink表的company字段不能为空！！")
	// }
	err := tx.Create(i).Error
	return err
}

// DelCandidateByProcInstID DelCandidateByProcInstID
// 删除历史候选人
func DelCandidateByProcInstID(procInstID int, tx *gorm.DB) error {
	return tx.Where("proc_inst_id=? and type=?", procInstID, IdentityTypes[CANDIDATE]).Delete(&Identitylink{}).Error
}

// IfParticipantByTaskID IfParticipantByTaskID
// 针对指定任务判断用户是否已经审批过了
func IfParticipantByTaskID(userID, company string, taskID int) (bool, error) {
	var count int
	err := db.Model(&Identitylink{}).Where("user_id=? and company=? and task_id=?", userID, company, taskID).Count(&count).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
