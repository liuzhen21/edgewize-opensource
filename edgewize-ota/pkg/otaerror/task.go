package otaerror

import (
	"errors"
)

var (
	// ErrNoSupport is returned when the current status does not allow deletion
	ErrTaskNoSupport           = errors.New("task no support")
	ErrBatchNotSupportNodeName = errors.New("node names cannot be specified for batch operations")
	ErrTaskDefine              = errors.New("task definition is not string")
	ErrTaskNotFound            = errors.New("task not found")

	ErrEdgeServiceResp     = errors.New("edgeservice response result error")
	ErrEdgeServiceRespType = errors.New("edgeservice response result data type error")

	ErrExtConfigEmpty = errors.New("extConfig is empty")
)
