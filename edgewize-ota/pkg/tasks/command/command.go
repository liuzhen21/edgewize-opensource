package command

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

// CommandTaskConfig is the task config struct for command task
type CommandTaskConfig struct {
	types.DefaultTaskConfig `json:"-"`
	Command                 string `json:"command"`
}

type CommandTask struct {
	name string
}

func (b *CommandTask) Name() string {
	return b.name
}

func (b *CommandTask) Run(config interface{}) (*string, error) {
	return CommandExec(config)
}

func (b *CommandTask) CheckConfig(config *types.TaskRequest) (interface{}, error) {
	defaultConfig, err := tasks.DefaultCheckConfig(config, b.name)
	if err != nil {
		return nil, err
	}

	commandConfig := &CommandTaskConfig{}

	configBytes, err := json.Marshal(config.Config)
	if err != nil {
		klog.Errorln("Error encoding JSON:", err)
		return nil, err
	}
	err = json.Unmarshal(configBytes, commandConfig)
	if err != nil {
		klog.Errorln("Error decoding JSON:", err)
		return nil, err
	}
	commandConfig.DefaultTaskConfig = *defaultConfig
	return commandConfig, nil
}

func (b *CommandTask) GetConfig() interface{} {
	return &CommandTaskConfig{
		Command: "command",
	}
}

func NewCommandTask(taskName string) types.TaskFunc {
	return &CommandTask{
		name: taskName,
	}
}

func CommandExec(config interface{}) (*string, error) {
	commandConfig, ok := config.(*CommandTaskConfig)
	if !ok {
		klog.Errorf("config is not *types.DefaultTaskConfig: %v\n", config)
		klog.Errorln(reflect.TypeOf(config))
		return nil, otaerror.ErrCommandTaskConfig
	}

	playbookFile := "roles/command/main.yml"
	inventoryFile := MakeInventoryFile(commandConfig)
	if inventoryFile == "" {
		res := "inventoryFile is empty"
		return &res, otaerror.ErrInventoryFileEmpty
	}
	defer os.RemoveAll(inventoryFile)
	return ansible.Run(types.AnsibleConnectionDefault, inventoryFile, playbookFile)
}

func MakeInventoryFile(config *CommandTaskConfig) string {
	// 创建 Inventory 结构体实例
	inventory := types.Inventory{}
	inventory.All.Hosts = tasks.MakeAnsibleHosts(&config.DefaultTaskConfig)

	// 添加变量信息
	inventory.All.Vars = map[string]interface{}{
		"command": config.Command,
	}
	return inventory.MakeInventoryFile()
}

func Init() {
	tasks.RegisterTask(NewCommandTask("command"))
}
