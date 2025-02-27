package worker_manager

import (
	"fmt"
	"testing"
	"time"
)

type MockNodeInfo struct {
}

func (m *MockNodeInfo) UpdateNodeInfoStatus(nodeName, status string) error {
	return nil
}

func TestNewTaskManager(t *testing.T) {
	channelSize := 10
	tm := newTaskManager(new(MockNodeInfo))
	workerNameList := []string{"worker1", "worker2", "worker3"}
	for _, workerName := range workerNameList {
		tm.AddWorker(workerName, channelSize)
	}

	defer tm.Stop()

	taskNameList := []string{"task1", "task2", "task3", "task4", "task5"}
	for idx, taskName := range taskNameList {
		for _, workerName := range workerNameList {
			taskNameTemp := taskName
			task := &Task{
				ID:         int64(idx),
				WorkerName: workerName,
				Job: func() (string, error) {
					fmt.Println(taskNameTemp, "is running")
					time.Sleep(1 * time.Second)
					fmt.Println(taskNameTemp, "is done")
					return "", nil
				},
				TaskName: taskNameTemp,
				Config:   nil,
			}
			tm.AddTask(task)
		}
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			for _, workerName := range workerNameList {
				taskList := tm.GetWorkerCurrentTaskList(workerName)
				for _, task := range taskList {
					fmt.Println("worker", workerName, "task", task.TaskName, "status", task.Status)
				}
			}

		}
	}()

	time.Sleep(11 * time.Second)
}

func TestNewTaskManager_AddWorker(t *testing.T) {
	tm := newTaskManager(new(MockNodeInfo))
	workerName := "worker1"
	channelSize := 10

	tm.AddWorker(workerName, channelSize)

	workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
	// Verify that the worker is added to the task manager
	if _, ok := tm.Workers[workerKey]; !ok {
		t.Errorf("Failed to add worker %s to the task manager", workerName)
	}

	// Verify that the worker is started
	worker := tm.Workers[workerKey]
	select {
	case <-worker.QuitChannel:
		t.Errorf("Worker %s is not started", workerName)
	default:
		// Worker is started
	}

	// Verify that the worker's task channel has the correct size
	if cap(worker.TaskChannel) != channelSize {
		t.Errorf("Worker %s has incorrect task channel size. Expected: %d, Actual: %d", workerName, channelSize, cap(worker.TaskChannel))
	}
}

func TestNewTaskManager_RemoveWorker(t *testing.T) {
		tm := newTaskManager(new(MockNodeInfo))
	workerName := "worker1"
	channelSize := 10

	tm.AddWorker(workerName, channelSize)

	// Remove the worker
	tm.RemoveWorker(workerName)

	// Verify that the worker is removed from the task manager
	if _, ok := tm.Workers[workerName]; ok {
		t.Errorf("Failed to remove worker %s from the task manager", workerName)
	}

	// Verify that the worker is stopped
	if worker := tm.Workers[workerName]; worker != nil {
		select {
		case <-worker.QuitChannel:
			// Worker is stopped
		default:
			t.Errorf("Worker %s is not stopped", workerName)
		}
	}
}

func TestNewTaskManager_AddTask(t *testing.T) {
	tm := newTaskManager(new(MockNodeInfo))
	workerName := "worker1"
	channelSize := 10

	tm.AddWorker(workerName, channelSize)

	task := &Task{
		ID:         1,
		WorkerName: workerName,
		Job: func() (string, error) {
			// Test job function
			return "", nil
		},
	}

	tm.AddTask(task)

	// Verify that the task is added to the worker's task channel
	workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
	worker := tm.Workers[workerKey]
	select {
	case <-worker.TaskChannel:
		// Task is added
	default:
		t.Errorf("Failed to add task to worker %s", workerName)
	}
}

func TestNewTaskManager_Start(t *testing.T) {
	channelSize := 10
	tm := newTaskManager(new(MockNodeInfo)	)
	workerNameList := []string{"worker1", "worker2", "worker3"}
	for _, workerName := range workerNameList {
		tm.AddWorker(workerName, channelSize)
	}

	// Verify that all workers are started
	for _, workerName := range workerNameList {
		workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
		worker := tm.Workers[workerKey]
		select {
		case <-worker.QuitChannel:
			t.Errorf("Worker %s is not started", workerName)
		default:
			// Worker is started
		}
	}
}

func TestNewTaskManager_Stop(t *testing.T) {
	channelSize := 10
	tm := newTaskManager(new(MockNodeInfo))
	workerNameList := []string{"worker1", "worker2", "worker3"}
	for _, workerName := range workerNameList {
		tm.AddWorker(workerName, channelSize)
	}

	// 通知所有worker停止
	tm.Stop()

	// Verify that all workers are stopped
	for _, workerName := range workerNameList {
		workerKey := tm.generateWorkerKey(workerName) // 使用新函数生成 workerKey
		worker := tm.Workers[workerKey]
		select {
		case <-worker.QuitChannel:
			// Worker is stopped
		default:
			t.Errorf("Worker %s is not stopped", workerName)
		}
	}
}
