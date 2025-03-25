# Edgewize OTA

Edgewize OTA 是一个轻量级的节点管理和自动化任务执行系统。它基于 Ansible 构建，提供简单易用的 RESTful API 接口，支持节点的生命周期管理和自动化任务的执行与监控。

## 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [API 文档](#api-文档)
- [任务系统](#任务系统)
- [插件开发](#插件开发)

## 功能特性

### 节点管理
- 自动化节点注册与代理安装
- 节点状态监控
- 详细节点信息查询
- 节点生命周期管理

### 任务系统
- 基于 Ansible 的自动化任务执行
- 灵活的任务配置
- 任务状态跟踪
- 详细的执行日志

### 系统特点
- RESTful API 设计
- OpenAPI/Swagger 文档支持
- 支持任务扩展
- 轻量级部署

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/edgewize-io/edgewize.git

# 进入项目目录
cd edgewize/edgewize-ota

# 运行 或者使用docker run 
go run edgewize-ota/cmd/edgewize-ota-server/edgewize-ota-server.go动服务

```

http://<your-server>:8082/apidocs/#/


### 核心 API 接口

#### 节点管理
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /apis/ota.edgewize.io/v1alpha1/register | 注册并安装节点代理 |
| GET | /apis/ota.edgewize.io/v1alpha1/nodes | 获取节点列表 |
| GET | /apis/ota.edgewize.io/v1alpha1/nodes/{nodeName} | 获取指定节点信息 |
| DELETE | /apis/ota.edgewize.io/v1alpha1/nodes/{nodeName} | 删除指定节点 |

#### 任务管理
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /apis/ota.edgewize.io/v1alpha1/tasks | 创建并执行任务 |
| GET | /apis/ota.edgewize.io/v1alpha1/tasks | 获取任务列表 |
| GET | /apis/ota.edgewize.io/v1alpha1/tasks/{taskName} | 获取任务配置 |

## 任务系统

### 任务类型

1. 简单任务
   - 直接使用 Ansible playbook 定义
   - 无需额外参数配置
   - 适用于标准化的运维操作

2. 复杂任务
   - 支持动态参数配置
   - 需要通过插件方式实现
   - 适用于定制化的业务场景

## 复杂任务开发

### 复杂任务接口
复杂任务需要实现以下接口：

```golang
type TaskFunc interface {
    // 返回任务名称
    Name() string
    
    // 验证任务配置
    CheckConfig(config *TaskRequest) (interface{}, error)
    
    // 执行任务
    Run(config interface{}) (*string, error)
    
    // 获取配置模板
    GetConfig() interface{}
}
```

### 配置定义
```golang
type TaskConfig struct {
    types.DefaultTaskConfig `json:"-"`
    // 自定义配置字段
    CustomField string `json:"customField"`
}
```

### 插件注册
```golang
func init() {
    tasks.RegisterTask(NewCustomTask())
}
```

### 示例插件
参考 `/edgewize-ota/pkg/tasks/command/command.go` 了解完整的插件实现示例。

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进项目。在提交代码前，请确保：

1. 代码风格符合项目规范
2. 添加了必要的测试用例
3. 更新了相关文档

## 许可证


