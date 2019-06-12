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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/workflow/", controller.Index)
	//-------------------------流程定义----------------------
	mux.HandleFunc("/workflow/procdef/save", controller.SaveProcdef)
	mux.HandleFunc("/workflow/procdef/findAll", controller.FindAllProcdefPage)
	mux.HandleFunc("/workflow/procdef/delById", controller.DelProcdefByID)
	// -----------------------流程实例-----------------------
	mux.HandleFunc("/workflow/process/start", controller.StartProcessInstance)
	mux.HandleFunc("/workflow/process/findTask", controller.FindMyProcInstPageAsJSON)
	// mux.HandleFunc("/workflow/process/moveToHistory", controller.MoveFinishedProcInstToHistory)
	// -----------------------任务--------------------------
	mux.HandleFunc("/workflow/task/complete", controller.CompleteTask)
	mux.HandleFunc("/workflow/task/withdraw", controller.WithDrawTask)
	// 配置
	var config = *config.Config
	// 启动数据库连接
	model.Setup()
	// 启动定时任务
	service.CronJobs()
	// 启动服务
	readTimeout, err := strconv.Atoi(config.ReadTimeout)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := strconv.Atoi(config.WriteTimeout)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(readTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(writeTimeout * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("the application start up at port%s", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
