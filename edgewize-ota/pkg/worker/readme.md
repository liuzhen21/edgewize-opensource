# WorkerManager 使用文档

`WorkerManager` 是一个用于管理和调度后台任务的系统。它允许你创建工作器（Worker），这些工作器可以异步执行任务，并管理它们的生命周期。以下是如何使用 `WorkerManager` 的详细指南。

## 1. 初始化 TaskManager

`TaskManager` 是管理所有工作器的中心组件。你可以通过调用 `GetDefaultTaskManager` 方法获取默认的 `TaskManager` 实例：

```go
taskManager := worker_manager.GetDefaultTaskManager()
```

这个方法确保全局只有一个 `TaskManager` 实例。

## 2. 创建 Worker

你可以通过 `AddWorker` 方法添加新的工作器。这个方法需要工作器名称和通道大小（用于定义任务通道的容量）：

```go
workerName := "worker1"
channelSize := 10
taskManager.AddWorker(workerName, channelSize)
```

这将创建一个新的工作器，并自动开始监听和执行分配给它的任务。

## 3. 添加任务

要向特定工作器添加任务，你需要创建一个 `Task` 实例并使用 `AddTask` 方法：

```go
task := &worker_manager.Task{
    ID:         1,
    WorkerName: workerName,
    Job: func() (string, error) {
        // 任务逻辑
        return "任务完成", nil
    },
    Status:   worker_manager.Pending,
    TaskName: "示例任务",
}

err := taskManager.AddTask(task)
if err != nil {
    // 处理错误
}
```

这个任务将被发送到指定的工作器，并由该工作器异步执行。

## 4. 查询工作器的任务列表

你可以查询任何工作器当前的任务列表，这通过调用 `GetWorkerCurrentTaskList` 方法实现：

```go
tasks := taskManager.GetWorkerCurrentTaskList(workerName)
for _, task := range tasks {
    fmt.Printf("任务 ID: %d, 状态: %s\n", task.ID, task.Status)
}
```

## 5. 停止工作器

当你需要停止一个工作器时，可以使用 `RemoveWorker` 方法：

```go
taskManager.RemoveWorker(workerName)
```

这将停止工作器并从 `TaskManager` 中移除它。

## 6. 停止所有工作器

如果你需要停止所有工作器，可以调用 `Stop` 方法：

```go
taskManager.Stop()
```

这将逐个停止所有注册的工作器，并等待它们全部停止。

## 总结

`WorkerManager` 提供了一个强大的框架来管理和执行后台任务。通过上述方法，你可以轻松地添加、管理和监控后台任务的执行。这个系统特别适合需要高度并发和任务管理的应用程序。