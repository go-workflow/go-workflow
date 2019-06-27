package model

import (
	"sync"

	"github.com/jinzhu/gorm"
)

// ProcInstHistory ProcInstHistory
type ProcInstHistory struct {
	ProcInst
}

// FindProcHistory 查询历史纪录
func FindProcHistory(userID, company string, pageIndex, pageSize int) ([]*ProcInstHistory, int, error) {
	var datas []*ProcInstHistory
	var count int
	var err1 error
	var wg sync.WaitGroup
	numberOfRoutine := 2
	errStream := make(chan error, numberOfRoutine)
	selectDatas := func(wg *sync.WaitGroup) {
		go func() {
			err := db.Where("id in (select distinct proc_inst_id from identitylink_history where company=? and user_id=?)", company, userID).
				Offset((pageIndex - 1) * pageSize).Limit(pageSize).
				Order("start_time desc").Find(&datas).Error
			errStream <- err
			wg.Done()
		}()
	}
	selectCount := func(wg *sync.WaitGroup) {
		go func() {
			err := db.Model(&ProcInstHistory{}).
				Where("id in (select distinct proc_inst_id from identitylink_history where company=? and user_id=?)", company, userID).
				Count(&count).Error
			errStream <- err
			wg.Done()
		}()
	}
	wg.Add(numberOfRoutine)
	selectDatas(&wg)
	selectCount(&wg)
	wg.Wait()
	close(errStream)

	for i := 0; i < numberOfRoutine; i++ {
		if err := <-errStream; err != nil {
			err1 = err
		}
	}
	return datas, count, err1
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
