# 配置
整个配置信息参考的是钉钉，钉钉生成的配置信息基本上能用，但是有所精简，只支持 主管审批和角色审批，可以打开钉钉控制平台来生成配置数据

https://github.com/mumushuiding/go-workflow/blob/master/images/processConfig.png

首先配置信息是一个Node对象的嵌套对象

type Node struct {

	Name           string          `json:"name,omitempty"`

	Type           string          `json:"type,omitempty"`

	NodeID         string          `json:"nodeId,omitempty"`

	PrevID         string          `json:"prevId,omitempty"`

	ChildNode      *Node           `json:"childNode,omitempty"`

	ConditionNodes []*Node         `json:"conditionNodes,omitempty"`

	Properties     *NodeProperties `json:"properties,omitempty"`

}

在解析json对象时，先迭代遍历ConditionNodes里面的所有节点，然后再迭代遍历ChildNode

以下是配置信息json对象：

{
  
  "name": "发起人",

  "type": "start", // Node 类型：开始节点

  "nodeId": "开始", // 当前节点名称

  "childNode": {

    "type": "route",  // Node 类型：条件节点

    "prevId": "sid-startevent",

    "nodeId": "8b5c_debb",

    "conditionNodes": [

      {

        "name": "条件1",

        "type": "condition",

        "prevId": "8b5c_debb",

        "nodeId": "da89_be76",

        "properties": {
          
          "conditions": [
            [
              {

                "type": "dingtalk_actioner_value_condition", // 条件类型：代表的是有范围的值

                "paramKey": "DDHolidayField-J2BWEN12__options", // 值的key, EXAMPLE.md文件里面有一个启动流程的案例，你往后台传递的变量key要与之匹配
                
                "paramLabel": "请假类型",
                
                "paramValues": [ // 当前条件值，可以多个，只要包含其中一个，这个条件就满足
                  "年假"
                ],
              }
                <!-- 
                
                "type": "dingtalk_actioner_range_condition", // 代表的是范围类型，比如: 1<a<10
               
                "paramKey": "DDHolidayField-J2BWEN12__duration",
                
                "paramLabel": "时长（天）",
                
                "lowerBound": "10", // lowerBound表示下限 比如：大于等于10
               
                "upperBound": "",   // upperBound表示上限 比如：小于等于20
                
                "unit": "天", -->
            ]
          ]
        },
        
        "childNode": {
         
          "name": "UNKNOWN",
         
          "type": "approver", // Node节点类型 approver 审批人，生成执行流时，只会纪录approver类型的节点
         
          "prevId": "da89_be76",
          
          "nodeId": "735c_0854",
         
          "properties": {
         
            "actionerRules": [
              {
          
                "type": "target_management", //审批人类型 target_management 表示下一级审批人为主管
         
              }
            ],
          }
        }
      },
      {
        
        "name": "条件2",
        
        "type": "condition",
       
        "prevId": "8b5c_debb",
       
        "nodeId": "a97f_9517",
       
        "properties": {
       
          "conditions": [
            [
              {
        
                "type": "dingtalk_actioner_value_condition",
        
                "paramKey": "DDHolidayField-J2BWEN12__options",
      
                "paramLabel": "请假类型",
       
                "paramValue": "",
      
                "paramValues": [
                  "调休"
                ],
              }
            ]
          ]
        },
       
        "childNode": {
        
          "name": "UNKNOWN",
          
          "type": "approver",
         
          "prevId": "a97f_9517",
        
          "nodeId": "5891_395b",
         
          "properties": {
         
            "actionerRules": [
              {
            
                "type": "target_label", // 审批人类型 target_label代表的是角色审批，比如：财务，人事
        
                "labelNames": "财务",
         
                "memberCount": 2, // 表示需要通过的人数，必须有2人审批通过，都会流转到下一环节，只要有一人驳回就流转到上一环节
         
                "actType": "and" // action类型 会签  or表示或签，默认为或签
              }
            ],
            "noneActionerAction": "auto"
          }
        }
      }
    ],
    "properties": {},
   
    "childNode": {
   
      "name": "UNKNOWN",
    
      "type": "approver",
    
      "prevId": "8b5c_debb",
   
      "nodeId": "59ba_8815",
   
      "properties": {
    
        "actionerRules": [
          {
      
            "type": "target_label",
     
            "labelNames": "人事",
    
            "labels": 427529104,
    
            "isEmpty": false,
    
            "memberCount": 1,
     
            "actType": "and"
          }
        ],
      }
    }
  }
}