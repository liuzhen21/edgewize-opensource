package tasks

import (
	"fmt"
	"os"

	"k8s.io/klog/v2"

	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/ansible"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/otaerror"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/utils"
)

// 定义启动参数
var (
	innerTaskDir = ""
	outerTaskDir = ""
)

var tasks map[string]types.TaskFunc

func InitTasks(innerDir, outerDir string) {
	// 解析启动参数
	innerTaskDir = innerDir
	outerTaskDir = outerDir
	tasks = make(map[string]types.TaskFunc)
	TaskList := GetTaskList()
	for _, task := range TaskList {
		tasks[task] = NewDefaultTask(task)
		klog.Infof("GetTaskList: %s", task)
	}
}

func GetTaskList() []string {
	innerTaskList := getTaskList(innerTaskDir)
	outerTaskList := getTaskList(outerTaskDir)

	// 创建一个映射以去重
	taskSet := make(map[string]struct{})

	// 添加 innerTaskList 中的任务
	for _, task := range innerTaskList {
		taskSet[task] = struct{}{}
	}

	// 添加 outerTaskList 中的任务
	for _, task := range outerTaskList {
		taskSet[task] = struct{}{}
	}

	// 将去重后的任务转换为切片
	uniqueTasks := make([]string, 0, len(taskSet))
	for task := range taskSet {
		uniqueTasks = append(uniqueTasks, task)
	}

	return uniqueTasks
}

func getTaskList(taskDir string) []string {
	if tasks, err := utils.GetImmediateSubdirectories(taskDir); err != nil {
		return nil
	} else {
		return tasks
	}
}

func MakeInventoryFile(config *types.DefaultTaskConfig) string {
	// 创建 Inventory 结构体实例
	inventory := types.Inventory{}
	inventory.All.Hosts = MakeAnsibleHosts(config)
	return inventory.MakeInventoryFile()
}

// defaultTask
func defaultTaskFunc(config interface{}) (*string, error) {
	defaultConfig, ok := config.(*types.DefaultTaskConfig)
	if !ok {
		return nil, fmt.Errorf("config is not *types.DefaultTaskConfig")
	}

	// 优先使用 outer 目录中的 playbook
	playbookFile := fmt.Sprintf("%s/%s/main.yml", outerTaskDir, defaultConfig.TaskName)
	if _, err := os.Stat(playbookFile); os.IsNotExist(err) {
		// 如果 outer 目录不存在，则检查 inner 目录
		playbookFile = fmt.Sprintf("%s/%s/main.yml", innerTaskDir, defaultConfig.TaskName)
	}

	if _, err := os.Stat(playbookFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("playbook file %s not found", playbookFile)
	}

	inventoryFile := MakeInventoryFile(defaultConfig)
	defer os.RemoveAll(inventoryFile)
	klog.Info(inventoryFile)
	return ansible.Run(types.AnsibleConnectionDefault, inventoryFile, playbookFile)
}

type defaultTask struct {
	name string
}

func (d *defaultTask) Run(config interface{}) (*string, error) {
	return defaultTaskFunc(config)
}

func NewDefaultTask(taskName string) types.TaskFunc {
	return &defaultTask{name: taskName}
}

func (d *defaultTask) Name() string {
	return d.name
}

func (d *defaultTask) GetConfig() interface{} {
	return "{}"
}

func (d *defaultTask) CheckConfig(config *types.TaskRequest) (interface{}, error) {
	return DefaultCheckConfig(config, d.name)
}

func RegisterTask(taskFunc types.TaskFunc) {
	klog.Infof("RegisterTask: %s", taskFunc.Name())
	tasks[taskFunc.Name()] = taskFunc
}

func RunTask(taskName string, config interface{}) (*string, error) {
	if taskFunc, ok := tasks[taskName]; ok {
		return taskFunc.Run(config)
	} else {
		return nil, nil
	}
}

func GetTaskFunc(taskName string) types.TaskFunc {
	if taskFunc, ok := tasks[taskName]; ok {
		return taskFunc
	} else {
		return nil
	}
}

func DefaultCheckConfig(defaultConfig *types.TaskRequest, taskName string) (*types.DefaultTaskConfig, error) {
	if defaultConfig.TaskName != taskName {
		return nil, fmt.Errorf("taskName is not %s", taskName)
	}

	if defaultConfig.TargetNode == "" {
		return nil, otaerror.ErrTargetNodeEmpty
	}

	taskConfig := &types.DefaultTaskConfig{}
	taskConfig.TaskName = defaultConfig.TaskName
	taskConfig.TargetNode = defaultConfig.TargetNode
	if defaultConfig.Host != nil {
		taskConfig.UserName = defaultConfig.Host.AnsibleUser
		taskConfig.Password = defaultConfig.Host.AnsiblePassword
		taskConfig.SshPort = defaultConfig.Host.AnsiblePort
		taskConfig.Host = defaultConfig.Host.AnsibleHost
	}
	return taskConfig, nil
}

func MakeAnsibleHosts(config *types.DefaultTaskConfig) map[string]types.Host {
	if config == nil {
		return nil
	}
	hosts := map[string]types.Host{}
	// 添加主机信息
	if config.SshPort != 0 {
		hosts["edgenode"] = types.Host{
			AnsibleHost:     config.Host,
			AnsibleUser:     config.UserName,
			AnsiblePassword: config.Password,
			AnsiblePort:     config.SshPort,
		}
	} else {
		hosts["edgenode"] = types.Host{
			AnsibleHost: config.TargetNode,
		}
	}
	return hosts
}
