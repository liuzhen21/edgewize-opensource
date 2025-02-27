# Edgewize OTA

Edgewize OTA 是一个轻量级的节点管理和任务执行系统，支持节点注册和基于 Ansible 的自动化任务执行。

## 功能特性

- 节点管理
  - 节点注册
  - 节点状态查询
  - 节点信息获取
  - 节点删除
- 任务管理
  - 任务创建和执行
  - 任务列表查询
  - 任务配置获取
- RESTful API
  - 完整的 OpenAPI 文档支持
  - 标准化的 HTTP 接口

## API 接口

### swagger 文档
http://x.x.x.x:8082/apidocs/#/

### 节点管理接口

```http
POST /register    # 注册并安装节点代理
GET /nodes        # 获取节点列表
GET /nodes/{nodeName}    # 获取指定节点信息
DELETE /nodes/{nodeName} # 删除指定节点
```

### 任务管理接口

```http
POST /tasks              # 创建并执行任务
GET /tasks              # 获取任务列表
GET /tasks/{taskName}   # 获取任务配置
```

## 任务定义
edgewize-ota 支持的是标准的 ansible 任务, edgewize-ota-server 会自动编写inventory文件，但是不支持配置参数, 如果需要支持配置参数，则需要通过编写插件的形式来支持

### 按任务的复杂程度划分
1. 简单任务，直接定义ansible任务即可，运行时不需要配置可变的参数，不需要添加代码
2. 复杂任务，任务运行时需要额外的参数，编写一定的代码

### 如何自定义一个复杂任务


#### 插件要实现的功能
要定义一个复杂功能要实现如下的函数
```golang
type TaskFunc interface {
	Name() string
	CheckConfig(config *TaskRequest) (interface{}, error)
	Run(config interface{}) (*string, error) 
	GetConfig() interface{}
}
```

- Name()  返回任务的名称
- CheckConfig() 将接口传递的通用配置类型，检查并转为任务支持的配置类型，检查通过会传递给run 函数进行调用
- Run() 执行任务，返回任务的输出，根据配置生成临时的 inventory 文件，然后调用ansible进行执行，返回执行结果，删除临时 inventory 文件
- GetConfig() 获取任务的配置格式及示例

定义任务的配置
``` golang
// CommandTaskConfig is the task config struct for command task
type CommandTaskConfig struct {
	types.DefaultTaskConfig `json:"-"`
	Command                 string `json:"command"`
}

```
DefaultTaskConfig 是必备属性，所有的任务配置都要包含它。

任务注册
```golang
func init() {
	tasks.RegisterTask(NewCommandTask())
}
```