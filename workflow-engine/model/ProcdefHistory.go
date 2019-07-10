package model

// ProcdefHistory 历史流程定义
type ProcdefHistory struct {
	Procdef
}

// Save Save
func (p *ProcdefHistory) Save() (ID int, err error) {
	err = db.Create(p).Error
	if err != nil {
		return 0, err
	}
	return p.ID, nil
}
