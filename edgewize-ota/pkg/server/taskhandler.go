package server

import (
	"encoding/json"
	"errors"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"gorm.io/gorm"
	"k8s.io/klog/v2"

	nodeinfodb "github.com/edgewize-io/edgewize/edgewize-ota/pkg/nodeinfodb"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/otaerror"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/tasks"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/utils"
	worker "github.com/edgewize-io/edgewize/edgewize-ota/pkg/worker"
)

func doTask(request *restful.Request, response *restful.Response) {
	// 1. 解析请求参数
	// 2. 校验参数
	// 3. 执行任务
	// 4. 返回结果
	klog.Infof("do task")
	taskRequest := &types.TaskRequest{}
	request.ReadEntity(taskRequest)
	klog.Infof("taskName: %s", taskRequest.TaskName)
	taskFunc := tasks.GetTaskFunc(taskRequest.TaskName)
	if taskFunc == nil {
		klog.Errorf("task %s not found", taskRequest.TaskName)
		response.WriteErrorString(404, "task not found")
		return
	}

	nodeInfo, err := nodeInfoDB.QueryNodeInfo(taskRequest.TargetNode)
	if err != nil {
		klog.Errorf("node %s not found", taskRequest.TargetNode)
		response.WriteErrorString(404, "node not found")
		return
	}
	sshInfo := &SSHInfo{}
	err = json.Unmarshal(nodeInfo.SSHInfo, sshInfo)
	if err != nil {
		klog.Errorf("unmarshal ssh info error: %v", err)
		response.WriteErrorString(400, "unmarshal ssh info error")
		return
	}
	taskRequest.Host = &types.Host{
		AnsibleUser:     sshInfo.UserName,
		AnsiblePassword: sshInfo.Password,
		AnsiblePort:     sshInfo.Port,
		AnsibleHost:     sshInfo.Host,
	}
	klog.Infof("task %s host: %v", taskRequest.TaskName, taskRequest.Host)
	config, err := taskFunc.CheckConfig(taskRequest)
	if err != nil {
		klog.Errorf("task %s config error: %s", taskRequest.TaskName, err)
		response.WriteErrorString(400, "task config error")
		return
	}

	klog.Infof("task %s config: %v", taskRequest.TaskName, config)

	tempTask := &worker.Task{
		ID:         time.Now().UnixNano(),
		WorkerName: taskRequest.TargetNode,
		Job: func() (string, error) {
			res, err := taskFunc.Run(config)
			if err != nil {
				klog.Errorf("task %s run error: %s", taskRequest.TaskName, err)
			}
			resString := ""
			if res != nil {
				resString = *res
			}
			return resString, err
		},
		Config:   config,
		TaskName: taskRequest.TaskName,
		Time:     time.Now().Format("2006-01-02T15:04:05Z"),
	}
	tm.AddTask(tempTask) // Fix: Convert newTask to worker.Task before passing it to AddTask method

	returnSucceeded(response)
}

func getTask(request *restful.Request, response *restful.Response) {
	// 1. 解析请求参数
	// 2. 校验参数
	// 3. 获取任务
	// 4. 返回结果
	klog.Infof("get task")
	tasks := tasks.GetTaskList()
	response.WriteEntity(tasks)
}

func getTaskConfig(request *restful.Request, response *restful.Response) {
	// 1. 解析请求参数
	// 2. 校验参数
	// 3. 获取任务
	// 4. 返回结果
	klog.Infof("get task config")
	taskName := request.PathParameter("taskName")
	taskFunc := tasks.GetTaskFunc(taskName)
	if taskFunc == nil {
		klog.Errorf("task %s not found", taskName)
		response.WriteErrorString(404, "task not found")
		return
	}
	response.WriteEntity(taskFunc.GetConfig())
}

func installAgentHandler(request *restful.Request, response *restful.Response) {
	// 1. 解析请求参数
	// 2. 校验参数
	// 3. 获取任务
	// 4. 返回结果
	klog.Infof("install agent")
	var info EdgeInstallAgent
	err := request.ReadEntity(&info)
	if err == nil {
		klog.Infof("installAgent info : %v \n", info)
		//1 modify edgeagent.yaml
		if len(info.EdgeIpList) > 1 && info.EdgeNodeName != "" {
			klog.Infof("Node names cannot be specified for batch operations")
			returnFaileure(response, otaerror.ErrBatchNotSupportNodeName)
			return
		}

		//2 install
		for _, ip := range info.EdgeIpList {
			var hostName string
			if info.EdgeNodeName != "" {
				hostName = info.EdgeNodeName
			} else {
				hostName = utils.IP2TempEdgeName(ip)
			}

			if _, err := nodeInfoDB.QueryNodeInfo(hostName); !errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}

			tm.AddWorker(hostName, types.DefaultWorkerChanSize)
			sshInfo := SSHInfo{
				Host:     ip,
				UserName: info.UserName,
				Password: info.Password,
				Port:     info.SshPort,
			}
			sshInfoJSON, err := json.Marshal(sshInfo)
			if err != nil {
				klog.Errorf("marshal ssh info error: %v", err)
				returnFaileure(response, err)
				return
			}
			err = nodeInfoDB.InsertNodeInfo(nodeinfodb.NodeInfo{
				NodeName: hostName,
				NodeIP:   ip,
				Status:   types.NodeStatusRegistering,
				SSHInfo:  sshInfoJSON,
			})
			if err != nil {
				klog.Errorf("insert node info error: %v", err)
			}
		}
		returnSucceeded(response)
	}
}
