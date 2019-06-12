package model

import (
	"github.com/jinzhu/gorm"
)

// IdentitylinkHistory IdentitylinkHistory
type IdentitylinkHistory struct {
	Identitylink
}

// CopyIdentitylinkToHistoryByProcInstID CopyIdentitylinkToHistoryByProcInstID
func CopyIdentitylinkToHistoryByProcInstID(procInstID int, tx *gorm.DB) error {
	return db.Exec("insert into identitylink_history select * from identitylink where proc_inst_id=?", procInstID).Error
}
