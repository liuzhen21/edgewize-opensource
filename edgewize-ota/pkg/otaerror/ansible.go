package otaerror

import "errors"

var (
	ErrAnsibleMqttAgentConfig = errors.New("config is not *types.AnsibleMqttAgentConfig")
	ErrCommandTaskConfig      = errors.New("config is not *types.CommandTaskConfig")
	ErrDefaultTaskConfig      = errors.New("config is not *types.DefaultTaskConfig")
	ErrNodeinfoExtConifg      = errors.New("config is not *types.NodeinfoExtConfig")

	ErrInventoryFileEmpty = errors.New("inventoryFile is empty")
	ErrTargetNodeEmpty    = errors.New("targetNode is empty")
)
