package types

const NodeInfoUpdateTopic = "msg/%s/%s/nodeinfo"

const (
	NodeStatusOnline         = "online"
	NodeStatusOffline        = "offline"
	NodeStatusWorking        = "working"
	NodeStatusRegistering    = "registering"
	NodeStatusRegisterFailed = "registerFailed"
	NodeStatusDeleting       = "deleting"
)

const (
	DefaultLimit = 10
	DefaultPage  = 1
)

const DefaultWorkerChanSize = 10