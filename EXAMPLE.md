# 1.安装
# 1.1 安装mysql 
  目前只支持mysql数据库，测试之前先安装好数据库
# 1.2 docker 安装最新版 go-workflow 微服务
docker run  -e DbType=mysql -e DbLogMode=false DbName=activiti -e DbHost=localhost -e DbUser=root -e DbPassword=123 -p 8080:8080 registry.cn-hangzhou.aliyuncs.com/mumushuiding/go-workflow
# 1.3 通过 go get 获取
 1.go get github.com/mumushuiding/go-workflow/releases
 2.进入根目录，打开config.json文件,修改数据库连接配置
 3. $ go build
 4. $ go-workflow.exe

# 2.流程存储
# 2.1 存储流程定义。
  通过 Post 访问： http://localhost:8080/workflow/procdef/save

  (参数详解见 ProcessConfig流程定义配置.md)

  Post参数：
  
  {"userid":"11025","name":"请假","company":"A公司","resource":{"name":"发起人","type":"start","nodeId":"sid-startevent","childNode":{"type":"route","prevId":"sid-startevent","nodeId":"8b5c_debb","conditionNodes":[{"name":"条件1","type":"condition","prevId":"8b5c_debb","nodeId":"da89_be76","properties":{"conditions":[[{"type":"dingtalk_actioner_value_condition","paramKey":"DDHolidayField-J2BWEN12__options","paramLabel":"请假类型","paramValue":"","paramValues":["年假"],"oriValue":["年假","事假","病假","调休","产假","婚假","例假","丧假"],"isEmpty":false}]]},"childNode":{"name":"UNKNOWN","type":"approver","prevId":"da89_be76","nodeId":"735c_0854","properties":{"activateType":"ONE_BY_ONE","agreeAll":false,"actionerRules":[{"type":"target_management","level":1,"isEmpty":false,"autoUp":true}],"noneActionerAction":"admin"}}},{"name":"条件2","type":"condition","prevId":"8b5c_debb","nodeId":"a97f_9517","properties":{"conditions":[[{"type":"dingtalk_actioner_value_condition","paramKey":"DDHolidayField-J2BWEN12__options","paramLabel":"请假类型","paramValue":"","paramValues":["调休"],"oriValue":["年假","事假","病假","调休","产假","婚假","例假","丧假"],"isEmpty":false}]]},"childNode":{"name":"UNKNOWN","type":"approver","prevId":"a97f_9517","nodeId":"5891_395b","properties":{"activateType":"ALL","agreeAll":true,"actionerRules":[{"type":"target_label","labelNames":"财务","labels":427529103,"isEmpty":false,"memberCount":2,"actType":"and"}],"noneActionerAction":"auto"}}}],"properties":{},"childNode":{"name":"UNKNOWN","type":"approver","prevId":"8b5c_debb","nodeId":"59ba_8815","properties":{"activateType":"ALL","agreeAll":true,"actionerRules":[{"type":"target_label","labelNames":"人事","labels":427529104,"isEmpty":false,"memberCount":1,"actType":"and"}],"noneActionerAction":"admin"}}}}}

  如果返回：{"data":"1","ok":true} ，1表示流程实例的id,true表示成功了
# 3.启动流程
  通过 POST 访问： http://localhost:8080/workflow/process/start

  POST参数：

  {"procName":"请假","title":"请假-张三","userId":"11025","department":"技术中心","company":"A公司","var":{"DDHolidayField-J2BWEN12__duration":"8","DDHolidayField-J2BWEN12__options":"年假"}}

  返回结果：{"data":"1","ok":true}

# 4.审批

# 4.1 审批
  通过POST访问：http://localhost:8080/workflow/task/complete

  POST参数：{"taskID":2,"pass":"true","userID":"11029","company":"A公司"}

  参数详解： 2代表当前任务id，true表示通过，false表示驳回
# 4.2 撤回

  通过POST访问：http://localhost:8080/workflow/task/withdraw

  POST参数：{"taskID":2,"userID":"11029","procInstID":1,"company":"A公司"}

   参数详解： taskID为当前任务id

# 4.3 任务查询 

  通过POST访问 ：http://localhost:8080/workflow/process/findTask
  
  POST参数：{"userID":"11025","groups":["人事"],"departments":["技术中心"],"company":"A公司"}

  参数详解： groups 表示用户的所有角色，departments表示用户
