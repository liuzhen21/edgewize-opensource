package types

import (
	"os"

	"k8s.io/klog/v2"
	"gopkg.in/yaml.v2"
)

const (
	AnsibleConnectionSSH  = "ssh"
	AnsibleConnectionMqtt = "mqtt"
)

var AnsibleConnectionDefault = AnsibleConnectionSSH

// Inventory 结构体表示Ansible的主机和变量信息
type Inventory struct {
	All struct {
		Hosts map[string]Host        `yaml:"hosts"`
		Vars  map[string]interface{} `yaml:"vars"`
	} `yaml:"all"`
}

func (i *Inventory) MakeInventoryFile() string {
	// 将 Inventory 结构体编码为 YAML 数据
	yamlData, err := yaml.Marshal(i)
	if err != nil {
		klog.Error("Error encoding YAML:", err)
		return ""
	}
	// 创建一个临时文件
	file, err := os.CreateTemp("", "inventory-*.yaml")
	if err != nil {
		klog.Error("Error creating temporary file:", err)
		return ""
	}
	file.Write(yamlData)
	defer file.Close()
	return file.Name()
}

// Host 结构体表示每个主机的信息
type Host struct {
	AnsibleHost     string `yaml:"ansible_host"`
	AnsiblePort     int    `yaml:"ansible_port"`
	AnsibleUser     string `yaml:"ansible_user"`
	AnsiblePassword string `yaml:"ansible_password,omitempty"`
}

type TaskFunc interface {
	Name() string
	CheckConfig(config *TaskRequest) (interface{}, error)
	Run(config interface{}) (*string, error)
	GetConfig() interface{}
}

// DefaultTaskConfig is the default task config struct
type DefaultTaskConfig struct {
	TaskName    string `json:"taskName"`
	TargetNode  string `json:"targetNode"`
	ClusterName string `json:"clusterName"`
	Host        string `json:"host"`
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	SshPort     int    `json:"sshPort"`
}

// TaskRequest is the http request struct for task
// 对于ssh连接的节点，host不能为空
type TaskRequest struct {
	TaskName   string      `json:"taskName"`
	TargetNode string      `json:"targetNode"`
	Config     interface{} `json:"config"`
	Host       *Host       `json:"host"`
}
