package nodeinfodb

import (
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"k8s.io/klog/v2"

	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
)

var (
	ErrDBIsNil = errors.New("db is nil")
)

type NodeInfo struct {
	gorm.Model
	NodeName    string         `json:"nodeName,omitempty" gorm:"uniqueIndex"`
	NodeIP      string         `json:"nodeIP,omitempty"`
	Alias       string         `json:"alias,omitempty"`
	Description string         `json:"description,omitempty"`
	Status      string         `json:"status"`            // 0: online, 1: offline, 2: working
	SSHInfo     datatypes.JSON `json:"sshInfo,omitempty"` // 使用 JSON 字段存储标签
}

type NodeInfoDB struct {
	db *gorm.DB
}

func DefaultNodeInfoDB() *NodeInfoDB {
	return NewNodeInfoDB(db)
}

func NewNodeInfoDB(db *gorm.DB) *NodeInfoDB {
	// 自动迁移确保表结构正确
	//db.AutoMigrate(&NodeInfo{})
	return &NodeInfoDB{db: db}
}

func (n *NodeInfoDB) InsertNodeInfo(nodeInfo NodeInfo) error {
	// 创建记录
	if n.db == nil {
		return ErrDBIsNil
	}
	result := n.db.Create(&nodeInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (n *NodeInfoDB) QueryNodeInfo(nodeName string) (nodeInfo *NodeInfo, err error) {
	if n.db == nil {
		return nil, ErrDBIsNil
	}
	nodeInfo = &NodeInfo{} // 初始化 nodeInfo
	result := n.db.Where("node_name = ?", nodeName).First(nodeInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return nodeInfo, nil
}

func (n *NodeInfoDB) DeleteNodeInfo(nodeName string) error {
	if n.db == nil {
		return ErrDBIsNil
	}

	result := n.db.Unscoped().Where("node_name = ?", nodeName).Delete(&NodeInfo{})
	return result.Error
}

func (n *NodeInfoDB) QueryNodeInfos() (nodeInfo []NodeInfo, err error) {
	result := n.db.Find(&nodeInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return nodeInfo, nil
}

func (n *NodeInfoDB) UpdateNodeInfo(nodeInfo NodeInfo) error {
	result := n.db.Save(&nodeInfo)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (n *NodeInfoDB) InsertOrUpdateNodeInfo(nodeInfo NodeInfo) (bool, error) {

	nodeInfo.Status = types.NodeStatusOnline
	insertFlag := true

	if oldNodeInfo, err := n.QueryNodeInfo(nodeInfo.NodeName); err == nil {
		nodeInfo.ID = oldNodeInfo.ID
		nodeInfo.CreatedAt = oldNodeInfo.CreatedAt

		switch oldNodeInfo.Status {
		case types.NodeStatusWorking, types.NodeStatusDeleting:
			nodeInfo.Status = oldNodeInfo.Status
		}

		insertFlag = false
	}

	result := n.db.Save(&nodeInfo)
	if result.Error != nil {
		return insertFlag, result.Error
	}
	return insertFlag, nil
}

// 状态更新逻辑
func (n *NodeInfoDB) UpdateNodeInfoStatus(nodeName, status string) error {
	if oldNodeInfo, err := n.QueryNodeInfo(nodeName); err == nil {
		if oldNodeInfo.Status == types.NodeStatusDeleting {
			return nil
		}
		oldNodeInfo.Status = status
		result := n.db.Save(&oldNodeInfo)
		return result.Error
	} else {
		return err
	}
}

func FindNodeInfoList(db *gorm.DB, condition map[string]interface{}, sortBy string, limit, page int) ([]NodeInfo, int, error) {
	var nodeList []NodeInfo

	query := db.Where(condition)

	// Apply sorting if sortBy is provided
	if sortBy != "" {
		query = query.Order(sortBy)
	}

	// Apply limiting and pagination
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Find(&nodeList)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	var count int64

	if result = db.Model(&NodeInfo{}).Where(condition).Count(&count); result.Error != nil {
		klog.Error(result.Error)
		return nodeList, int(count), result.Error
	}

	return nodeList, int(count), nil
}
