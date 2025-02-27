package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
	"k8s.io/klog/v2"

	nodeinfodb "github.com/edgewize-io/edgewize/edgewize-ota/pkg/nodeinfodb"
	otaerror "github.com/edgewize-io/edgewize/edgewize-ota/pkg/otaerror"
	types "github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
)

var nodeInfoDB *nodeinfodb.NodeInfoDB

func listNodes(request *restful.Request, response *restful.Response) {
	// 1. 解析请求参数
	// 2. 校验参数
	// 3. 获取任务
	// 4. 返回结果
	nodeName := request.QueryParameter("nodeName")
	status := request.QueryParameter("status")
	sortBy := request.QueryParameter("sortBy")
	limitStr := request.QueryParameter("limit")
	pageStr := request.QueryParameter("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = types.DefaultLimit
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = types.DefaultPage
	}
	nodeList, err := GetEdgeNodeList(nodeName, status, sortBy, limit, page)
	if err != nil {
		returnFaileure(response, err)
		return
	}
	response.AddHeader("Content-Type", "text/json")
	//response.WriteHeader(http.StatusOK)
	response.WriteAsJson(&EdgeOtaResponse{
		Code:   http.StatusOK,
		Status: SUCCEEDED,
		Data:   nodeList,
	})
}

func getNode(request *restful.Request, response *restful.Response) {
	nodeName := request.PathParameter("nodeName")
	nodeInfo, err := nodeInfoDB.QueryNodeInfo(nodeName)
	if err == gorm.ErrRecordNotFound {
		response.WriteHeader(http.StatusNotFound)
		response.WriteAsJson(&EdgeOtaResponse{
			Code:    http.StatusNotFound,
			Status:  FAILURE,
			Message: "节点不存在",
		})
		return
	}
	if err != nil {
		returnFaileure(response, err)
		return
	}
	response.WriteAsJson(&EdgeOtaResponse{
		Code:   http.StatusOK,
		Status: SUCCEEDED,
		Data:   nodeInfo,
	})
}

func delNode(request *restful.Request, response *restful.Response) {
	nodeName := request.PathParameter("nodeName")

	if node, err := nodeInfoDB.QueryNodeInfo(nodeName); err == nil {
		switch node.Status {
		case types.NodeStatusRegisterFailed, types.NodeStatusOffline:
			// 直接删除
			nodeInfoDB.DeleteNodeInfo(nodeName)
		case types.NodeStatusOnline:
			// 删除 agent， 然后删除节点
			nodeInfoDB.UpdateNodeInfoStatus(nodeName, types.NodeStatusDeleting)
			time.AfterFunc(time.Second*60, func() {
				nodeInfoDB.DeleteNodeInfo(nodeName)
			})

		case types.NodeStatusWorking, types.NodeStatusRegistering, types.NodeStatusDeleting:
			// 当前状态不允许删除
			returnFaileure(response, otaerror.ErrTaskNoSupport)
			klog.Errorf("node status is %s, not allowed to delete", node.Status)
			return
		}
	} else {
		klog.Errorf("query node info failed: %v", err)
		returnFaileure(response, err)
		return
	}

	returnSucceeded(response)
}

func GetEdgeNodeList(nodeName, status, sortBy string, limit, page int) (*ListResult, error) {
	// 1. 解析请求参数
	// 2. 校验参数
	// 3. 获取任务
	// 4. 返回结果

	conditions := make(map[string]interface{})
	if nodeName != "" {
		conditions["node_name"] = nodeName
	}
	if status != "" {
		conditions["status"] = status
	}

	if sortBy == "createTime" {
		sortBy = "created_at"
	}

	nodeList, totalCount, err := nodeinfodb.FindNodeInfoList(nodeinfodb.GetDB(), conditions, sortBy, limit, page)

	if err != nil {
		return nil, err

	}

	listResult := &ListResult{TotalItems: totalCount}

	for _, node := range nodeList {
		taskConfig := &EdgeTasksConfig{}
		taskList := tm.GetWorkerCurrentTaskList(node.NodeName)
		klog.Info("taskList: ", taskList)

		for _, task := range taskList {
			taskDefinition, _ := json.Marshal(task.Config)
			taskConfig.Tasks = append(taskConfig.Tasks, &EdgeTask{
				TaskName:       task.TaskName,
				TaskId:         task.ID,
				TaskStatus:     task.Status,
				TaskDefinition: string(taskDefinition),
				TaskTime:       task.Time,
				TaskMessage:    task.Result,
			})
		}
		listResult.Items = append(listResult.Items, &EdgeNodeInfo{
			TasksList:       taskConfig,
			NodeStatus:      node.Status,
			CreateTime:      node.CreatedAt.Format("2006-01-02T15:04:05Z"),
			AgentUpdateTime: node.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			NodeName:        node.NodeName,
		})
	}
	return listResult, nil
}
