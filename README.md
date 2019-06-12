<p align="center">
  <a href="">
    <img width="200" height="200" src="https://github.com/go-workflow/go-workflow/blob/master/images/go.jpg">
  </a>
</p>


<p align="center">
  <a href="https://travis-ci.org/vuetifyjs/vuetify">
    <img src="https://img.shields.io/travis/vuetifyjs/vuetify/dev.svg" alt="travis ci badge">
  </a>
  <a href="https://codecov.io/gh/vuetifyjs/vuetify">
    <img src="https://img.shields.io/codecov/c/github/vuetifyjs/vuetify.svg" alt="Coverage">
  </a>
</p>

<br>

<h1 align="center"><bold>Go-Workflow</bold></h1>

<p>go-workflow 是一个超轻量级的工作流引擎,基本架构同Activiti工作流有些相似，但是它更精简，更轻量，它是一个工作流微服务，具体案例详见：example.md</p>

# 一、特点：

  1.它是一个工作流微服务

  2.将所有的无关流程的数据，包括用户、用户组等信息从服务中解耦出去，go-workflow只纪录流程的流转
  
  3.使用json数组替代bpmn来生成流程定义，简化流程定义的生成

# 二、go-workflow框架
# 1.go-workflow 数据库设计
# 1.1 流程定义表
  表 procdef 用于保存流程的配置，
  主要字段有：

     name: 流程定义的名称，如:"请假流程"
     
     version: 流程定义的版本

     resource: 保存流程定义的具体配置,它是一个json格式的字符串

     company: 保存该流程创建人所在公司

# 1.2 流程实例表

  表 proc_inst 用于保存流程实例，当用户启动一个流程时，就会在这个表存入一个流程实例，

  主要字段有：

    procDefID: 对应表procdef的id，

    title: 标题，如："张三的请假流程"

    department: 用户所在部门

    nodeID: 当前所处于节点的名称

    candidate: 当前审批人或者审批用户组

    taskID: 当前任务id

# 1.3 执行流表
  表 execution 用于保存执行流，当用户启动一个流程时，就会生成一条执行流，之后的流程就会按照执行流的顺序流转，
  
  比如：开始-主管审批-财务审批-人事审批-结束 ，
  
  主要的字段有：

    procInstID： 流程实例id,对应表proc_inst

    procDefID: 流程定义id,对应表procdef

    nodeInfos: 是一个json数组，纪录流程实例会经过的所有节点

# 1.4 关系表
  表 identitylink 用于保存任务task的候选用户组或者候选用户以及用户所参与的流程信息，
  
  主要字段有

    type: 表示关系类型，有："candidate"和"participant"两种

    group: 表示当前审批的用户组

    userID: 表示当前审批的用户

    taskID: 对应任务task表的id

    step: 表示任务对应的执行流位置，比如：有一个执行流：开始-主管审批-财务审批-人事审批-结束，那么
    step=0,则处于【开始】位置，step=1则处于【主管审批】位置

    company: 表示公司

    procInstID： 对应流程实例id

# 1.5 任务表
  表 task 用于保存任务，
  
  主要字段有：

    nodeID: 表示节点，如："主管审批"结点

    step: 表示任务对应的执行流位置

    assignee: 任务的处理人

    memberCount： 表示当前任务需要多少人审批之后才能结束，默认是 1

    unCompleteNum： 表示还有多少人没有审批，默认是1

    agreeNum： 表示通过的人数

    actType： 表示任务类型 "or"表示或签，即一个人通过或者驳回就结束，"and"表示会签，要所有人通过就流
    转到下一步，如果有一个人驳回那么就跳转到上一步

# 1.6 历史数据表
  历史数据表包括 execution_history，identitylink_history，proc_inst_history，task_history这些表字段同正常的表相同，每天0点时，已经结束的流程数据会自动迁移过来
## 2 流程的存储
# 2.1 添加流程资源
  启动 go-workflow 微服务后，可以在浏览器中输入：http://localhost:8080/workflow/procdef/save 进行存储

  具体见 example.md 说明文档

# 3.流程的启动
  通过调用 StartProcessInstanceByID 方法来启动流程实例，
  
  主要涉及：

    获取流程定义
    
    GetResourceByNameAndCompany()->启动流程实例CreateProcInstTx()->生成执行流GenerateExec() -> 生成新任务NewTaskTx() -> 流程流转 MoveStage()

# 4.任务审批
  调用方法 Complete()方法来执行任务的审批，
  涉及方法：

    更新任务 UpdateTaskWhenComplete()-> 流转MoveStageByProcInstID()

  调用方法 WithDrawTask() 方法来执行任务的撤回
   
