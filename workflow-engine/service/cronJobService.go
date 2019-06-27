package service

import (
	"log"

	"github.com/robfig/cron"
)

// CronJobs CronJobs
// 所有定时任务，在启动时会执行
func CronJobs() {
	go func() {
		c := cron.New()
		// 每隔 20 秒执行
		spec := "*/20 * * * * ?"
		//每隔20秒将已经结束的流程数据迁移至历史数据表
		c.AddFunc(spec, func() {
			MoveFinishedProcInstToHistory()
		})
		// c.AddFunc("*/5 * * * * ?", func() {
		// 	log.Println("cron running")
		// })
		// 启动
		c.Start()
		log.Println("----------启动定时任务------------")
		defer c.Stop()
		select {}
	}()
	// c := cron.New()
	// // 每天0点执行
	// spec := "0 0 0 * * ?"
	// //每天0点时将已经结束的流程数据迁移至历史数据表
	// c.AddFunc(spec, func() {
	// 	MoveFinishedProcInstToHistory()
	// })
	// c.AddFunc("*/5 * * * * ?", func() {
	// 	log.Println("cron running")
	// })
	// // 启动
	// c.Start()
	// log.Println("----------启动定时任务------------")
	// defer c.Stop()
	// select {}
}
