package model

import (
	"github.com/jinzhu/gorm"
)

// ProcInstHistory ProcInstHistory
type ProcInstHistory struct {
	ProcInst
}

// SaveProcInstHistory SaveProcInstHistory
func SaveProcInstHistory(p *ProcInst) error {
	return db.Table("proc_inst_history").Create(p).Error
}

// DelProcInstHistoryByID DelProcInstHistoryByID
func DelProcInstHistoryByID(id int) error {
	return db.Where("id=?", id).Delete(&ProcInstHistory{}).Error
}

// SaveProcInstHistoryTx SaveProcInstHistoryTx
func SaveProcInstHistoryTx(p *ProcInst, tx *gorm.DB) error {
	return tx.Table("proc_inst_history").Create(p).Error
}
