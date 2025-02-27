package types

type NodeInfo interface {
	UpdateNodeInfoStatus(nodeName, status string) error
}
