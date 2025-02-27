package worker_manager

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/klog/v2"

	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/otaerror"
	"github.com/edgewize-io/edgewize/edgewize-ota/pkg/types"
)

var workerKeyFmt = "worker@%s"

var defaultTaskManager *TaskManager

// worker represents a worker that executes tasks.
type worker struct {
	WorkerName  string
	TaskChannel chan *Task
	QuitChannel chan bool
	CurrentTask FixedLengthTaskList
	nodeinfo    types.NodeInfo
}

// newWorker creates a new Worker instance. currentTaskListSize is double of channelSize.
func newWorker(workerName string, channelSize int, nodeinfo types.NodeInfo) *worker {
	return &worker{
		WorkerName:  workerName,
		TaskChannel: make(chan *Task, channelSize),
		QuitChannel: make(chan bool),
		CurrentTask: *NewFixedLengthTaskList(2 * channelSize),
		nodeinfo:    nodeinfo,
	}
}

// GetWorkerCurrentTaskList returns the current task list of the worker.
func (w *worker) GetWorkerCurrentTaskList() []*Task {
	var tasks []*Task
	for node := w.CurrentTask.Head; node != nil; node = node.Next {
		tasks = append(tasks, node.Value)
	}
	return tasks
}

// AddTask adds a new task to the worker. 需要增加超时处理
func (w *worker) AddTask(task *Task) {
	w.TaskChannel <- task
	w.CurrentTask.AddTask(task)
}

// Start starts the worker, allowing it to execute tasks.
func (w *worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case <-w.QuitChannel:
				return
			case task := <-w.TaskChannel:
				klog.Infof("Worker %s is executing Task %d\n", w.WorkerName, task.ID)
				task.Status = Running
				task.Time = time.Now().Format("2006-01-02T15:04:05Z")
				w.nodeinfo.UpdateNodeInfoStatus(task.WorkerName, types.NodeStatusWorking)
				if res, err := task.Job(); err != nil {
					klog.Errorf("Worker %s failed to execute Task %d: %v\n", w.WorkerName, task.ID, err)
					task.Status = Failure
					task.Time = time.Now().Format("2006-01-02T15:04:05Z")
					task.Result = res
				} else {
					klog.Infof("Worker %s successfully executed Task %d\n", w.WorkerName, task.ID)
					task.Status = Succeeded
					task.Time = time.Now().Format("2006-01-02T15:04:05Z")
					task.Result = res
				}
				w.nodeinfo.UpdateNodeInfoStatus(task.WorkerName, types.NodeStatusOnline)
			}
		}
	}(wg)
}

// Stop stops the worker, preventing it from executing further tasks.
func (w *worker) Stop() {
	close(w.QuitChannel)
	klog.Infof("Stopped Worker %s\n", w.WorkerName)
}

// TaskManager manages a collection of workers and tasks.
type TaskManager struct {
	Workers   map[string]*worker
	WaitGroup sync.WaitGroup
	Mutex     sync.Mutex
	nodeinfo  types.NodeInfo
}

// newTaskManager creates a new TaskManager instance.
func newTaskManager(nodeinfo types.NodeInfo) *TaskManager {
	return &TaskManager{
		Workers: make(map[string]*worker),
		nodeinfo: nodeinfo,
	}
}

// GetDefaultTaskManager returns the default TaskManager instance.
func GetDefaultTaskManager(nodeinfo types.NodeInfo) *TaskManager {
	if defaultTaskManager == nil {
		defaultTaskManager = newTaskManager(nodeinfo)
	}
	return defaultTaskManager
}

// 生成 workerKey 的函数
func (tm *TaskManager) generateWorkerKey(workerName string) string {
	return fmt.Sprintf(workerKeyFmt, workerName)
}

// AddWorker adds a new worker to the task manager.
func (tm *TaskManager) AddWorker(workerName string, channelSize int) {
	klog.Infof("add worker: %s\n", workerName)
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()

	worker := newWorker(workerName, channelSize, tm.nodeinfo)

	workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
	tm.Workers[workerKey] = worker
	worker.Start(&tm.WaitGroup)
}

// RemoveWorker removes a worker from the task manager.
func (tm *TaskManager) RemoveWorker(workerName string) {
	klog.Infof("remove worker: %s\n", workerName)
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()

	workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
	worker, exists := tm.Workers[workerKey]
	if !exists {
		return
	}

	worker.Stop()
	<-worker.QuitChannel // Wait for confirmation that the worker has stopped.

	// Remove the worker from the map
	delete(tm.Workers, workerKey) // 使用 workerKey 进行删除
}

// AddTask adds a new task to the task manager.
func (tm *TaskManager) AddTask(task *Task) error {
	klog.V(6).Infof("Adding task %s: %d to worker %s\n", task.TaskName, task.ID, task.WorkerName)
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()

	workerKey := tm.generateWorkerKey(task.WorkerName) // 使用新函数生成 workerKey
	worker := tm.Workers[workerKey]
	if worker == nil {
		klog.Errorf("worker not found: %s\n", workerKey)
		return otaerror.ErrWorkerNotFound
	}
	worker.AddTask(task)
	return nil
}

// Stop stops all the workers in the task manager.
func (tm *TaskManager) Stop() {
	for _, worker := range tm.Workers {
		worker.Stop()
	}
	tm.WaitGroup.Wait() // Wait for all workers to finish.
}

func (tm *TaskManager) GetWorkerCurrentTaskList(workerName string) []*Task {
	workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
	worker := tm.Workers[workerKey]
	if worker == nil {
		return nil
	}
	return worker.GetWorkerCurrentTaskList()
}
