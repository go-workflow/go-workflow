package model

import (
	"github.com/jinzhu/gorm"
)

// ExecutionHistory ExecutionHistory
// 执行流历史纪录
type ExecutionHistory struct {
	Execution
}

// CopyExecutionToHistoryByProcInstIDTx CopyExecutionToHistoryByProcInstIDTx
func CopyExecutionToHistoryByProcInstIDTx(procInstID int, tx *gorm.DB) error {
	return tx.Exec("insert into execution_history select * from execution where proc_inst_id=?", procInstID).Error
}
