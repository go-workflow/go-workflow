package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	config "github.com/go-workflow/go-workflow/workflow-config"

	controller "github.com/go-workflow/go-workflow/workflow-controller"
	model "github.com/go-workflow/go-workflow/workflow-engine/model"
	"github.com/go-workflow/go-workflow/workflow-engine/service"
)

// 配置
var conf = *config.Config

func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", conf.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", conf.AccessControlAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AccessControlAllowHeaders)
		h(w, r)
	}
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/workflow/", controller.Index)
	//-------------------------流程定义----------------------
	mux.HandleFunc("/workflow/procdef/save", crossOrigin(controller.SaveProcdef))
	mux.HandleFunc("/workflow/procdef/saveByToken", crossOrigin(controller.SaveProcdefByToken))
	mux.HandleFunc("/workflow/procdef/findAll", crossOrigin(controller.FindAllProcdefPage))
	mux.HandleFunc("/workflow/procdef/delById", crossOrigin(controller.DelProcdefByID))
	// -----------------------流程实例-----------------------
	mux.HandleFunc("/workflow/process/start", crossOrigin(controller.StartProcessInstance))               // 启动流程
	mux.HandleFunc("/workflow/process/startByToken", crossOrigin(controller.StartProcessInstanceByToken)) // 启动流程
	mux.HandleFunc("/workflow/process/findTask", crossOrigin(controller.FindMyProcInstPageAsJSON))        // 查询需要我审批的流程
	mux.HandleFunc("/workflow/process/findTaskByToken", crossOrigin(controller.FindMyProcInstByToken))
	mux.HandleFunc("/workflow/process/startByMyself", crossOrigin(controller.StartByMyself))   // 查询我启动的流程
	mux.HandleFunc("/workflow/process/FindProcNotify", crossOrigin(controller.FindProcNotify)) // 查询抄送我的流程
	// mux.HandleFunc("/workflow/process/moveToHistory", controller.MoveFinishedProcInstToHistory)
	// -----------------------任务--------------------------
	mux.HandleFunc("/workflow/task/complete", crossOrigin(controller.CompleteTask))
	mux.HandleFunc("/workflow/task/completeByToken", crossOrigin(controller.CompleteTaskByToken))
	mux.HandleFunc("/workflow/task/withdraw", crossOrigin(controller.WithDrawTask))
	mux.HandleFunc("/workflow/task/withdrawByToken", crossOrigin(controller.WithDrawTaskByToken))
	// ----------------------- 关系表 -------------------------
	mux.HandleFunc("/workflow/identitylink/findParticipant", crossOrigin(controller.FindParticipantByProcInstID))

	// ******************************** 历史纪录 ***********************************
	// -------------------------- 流程实例 -------------------------------
	mux.HandleFunc("/workflow/procHistory/findTask", crossOrigin(controller.FindProcHistory))
	mux.HandleFunc("/workflow/procHistory/findTaskByToken", crossOrigin(controller.FindProcHistoryByToken))
	mux.HandleFunc("/workflow/procHistory/startByMyself", crossOrigin(controller.StartHistoryByMyself))   // 查询我启动的流程
	mux.HandleFunc("/workflow/procHistory/FindProcNotify", crossOrigin(controller.FindProcHistoryNotify)) // 查询抄送我的流程
	// ----------------------- 关系表 -------------------------
	mux.HandleFunc("/workflow/identitylinkHistory/findParticipant", crossOrigin(controller.FindParticipantHistoryByProcInstID))
	// 启动数据库连接
	model.Setup()
	// 启动redis连接
	model.SetRedis()
	// 启动定时任务
	service.CronJobs()
	// 启动服务
	readTimeout, err := strconv.Atoi(conf.ReadTimeout)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := strconv.Atoi(conf.WriteTimeout)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(readTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(writeTimeout * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("the application start up at port%s", server.Addr)
	if conf.TLSOpen == "true" {
		err = server.ListenAndServeTLS(conf.TLSCrt, conf.TLSKey)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Printf("Server err: %v", err)
	}

}
