package dockerinstall

import (
	"encoding/json"
	"os"
	"reflect"

	"k8s.io/klog/v2"

	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/ansible"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/otaerror"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/tasks"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
)

// DockerInstallConfig is the task config struct for docker install task
type DockerInstallConfig struct {
	types.DefaultTaskConfig `json:"-"`
	Version                 string `json:"version"`
}

type DockerInstallTask struct {
	name string
}

func (b *DockerInstallTask) Name() string {
	return b.name
}

func (b *DockerInstallTask) Run(config interface{}) (*string, error) {
	return CommandExec(config)
}

func (b *DockerInstallTask) CheckConfig(config *types.TaskRequest) (interface{}, error) {
	defaultConfig, err := tasks.DefaultCheckConfig(config, b.name)
	if err != nil {
		return nil, err
	}

	taskConfig := &DockerInstallConfig{}

	configBytes, err := json.Marshal(config.Config)
	if err != nil {
		klog.Errorln("Error encoding JSON:", err)
		return nil, err
	}
	err = json.Unmarshal(configBytes, taskConfig)
	if err != nil {
		klog.Errorln("Error decoding JSON:", err)
		return nil, err
	}
	taskConfig.DefaultTaskConfig = *defaultConfig
	return taskConfig, nil
}

func (b *DockerInstallTask) GetConfig() interface{} {
	return &DockerInstallConfig{
		Version: "23.0.6",
	}
}

func NewDockerInstallTask(taskName string) types.TaskFunc {
	return &DockerInstallTask{
		name: taskName,
	}
}

func CommandExec(config interface{}) (*string, error) {
	commandConfig, ok := config.(*DockerInstallConfig)
	if !ok {
		klog.Errorf("config is not *types.DefaultTaskConfig: %v\n", config)
		klog.Errorln(reflect.TypeOf(config))
		return nil, otaerror.ErrCommandTaskConfig
	}

	playbookFile := "roles/docker_install/main.yml"
	inventoryFile := MakeInventoryFile(commandConfig)
	if inventoryFile == "" {
		res := "inventoryFile is empty"
		return &res, otaerror.ErrInventoryFileEmpty
	}
	defer os.RemoveAll(inventoryFile)
	klog.V(8).Infoln("inventory", inventoryFile)
	return ansible.Run(types.AnsibleConnectionDefault, inventoryFile, playbookFile)
}

func MakeInventoryFile(config *DockerInstallConfig) string {
	// 创建 Inventory 结构体实例
	inventory := types.Inventory{}

	// 添加主机信息
	inventory.All.Hosts = tasks.MakeAnsibleHosts(&config.DefaultTaskConfig)
	// 添加变量信息
	inventory.All.Vars = map[string]interface{}{
		"version": config.Version,
	}
	return inventory.MakeInventoryFile()
}

func Init() {
	tasks.RegisterTask(NewDockerInstallTask("docker_install"))
}
