package server

const (
	SUCCEEDED = "Succeeded"
	FAILURE   = "Failure"
)

type SSHInfo struct {
	Host     string `json:"host,omitempty"`
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
}

type TaskInfo struct {
	TaskName       string      `json:"taskName,omitempty"`
	TaskDefinition interface{} `json:"taskDefinition,omitempty"`
	TaskId         int64       `json:"taskId,omitempty"`
	TaskResult     string      `json:"taskResult,omitempty"`
	TaskMessage    string      `json:"taskMessage,omitempty"`
	TaskTime       string      `json:"taskTime,omitempty"`
}

type EdgeNodeInfo struct {
	TasksList       *EdgeTasksConfig `json:"taskList,omitempty"`
	NodeStatus      string           `json:"nodeStatus,omitempty"`
	CreateTime      string           `json:"createTime,omitempty"`
	AgentUpdateTime string           `json:"agentUpdateTime,omitempty"`
	//NodeTag         map[string]string `json:"nodeTag,omitempty"`
	NodeName string `json:"nodeName,omitempty"`
}

type EdgeOtaResponse struct {
	Code    uint32      `json:"code,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type EdgeTasksConfig struct {
	Tasks []*EdgeTask `json:"tasks,omitempty"`
}

type EdgeTask struct {
	//which node
	NodeName string `json:"nodeName,omitempty"`
	//which cluster
	ClusterName string `json:"clusterName,omitempty"`
	//what to do
	TaskName string `json:"taskName,omitempty"`
	//how to do
	TaskDefinition interface{} `json:"taskDefinition,omitempty"`
	//TaskDefinition string `json:"taskDefinition",omitempty`
	//match
	TaskId      int64  `json:"taskId,omitempty"`
	TaskStatus  string `json:"taskStatus"`
	TaskMessage string `json:"taskMessage,omitempty"`
	TaskTime    string `json:"taskTime,omitempty"`
}

type EdgeTaskToAgent struct {
	//what to do
	TaskName string `json:"taskName,omitempty"`
	//how to do
	TaskDefinition interface{} `json:"taskDefinition,omitempty"`
	//match
	TaskId int64 `json:"taskId,omitempty"`
}

type EdgeInstallAgent struct {
	//The node name can be defined only for a single node. The system checks whether there is only one IP address in the IP list
	EdgeNodeName string `json:"edgeNodeName,omitempty"`
	//alias
	EdgeNodeAlias string `json:"alias,omitempty"`
	//which nodes need to be installed
	EdgeIpList []string `json:"edgeIpList,omitempty"`
	//what are the unified user name and password
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
	//ssh custom port
	SshPort int `json:"sshPort,omitempty"`
	//description
	Description string `json:"description,omitempty"`
}

type ListResult struct {
	Items      []*EdgeNodeInfo `json:"items"`
	TotalItems int             `json:"totalItems"`
}

// 定义 ResourceList 结构体, 用于返回资源列表, 与 k8s 保持一致, 用于 k8s apiservice 的注册
type ResourceList struct {
	Kind         string        `json:"kind"`
	ApiVersion   string        `json:"apiVersion"`
	GroupVersion string        `json:"groupVersion"`
	Resources    []APIResource `json:"resources"`
}

// 定义 APIResource 结构体
type APIResource struct {
	Name         string   `json:"name"`
	SingularName string   `json:"singularName"`
	Namespaced   bool     `json:"namespaced"`
	Kind         string   `json:"kind"`
	Verbs        []string `json:"verbs"`
}
