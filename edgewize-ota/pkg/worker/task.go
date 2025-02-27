package worker_manager

const (
	Pending   = "Pending"
	Running   = "Running"
	Succeeded = "Succeeded"
	Failure   = "Failure"
)

// Task represents a unit of work to be executed by a worker.
type Task struct {
	ID         int64                  // 任务 ID
	WorkerName string                 // 指定执行任务的 worker 的名字, 一般由节点的名称生成
	Job        func() (string, error) // 任务执行的函数
	Status     string                 // 任务状态
	TaskName   string                 // 任务名称
	Config     interface{}            // 任务配置
	Time       string                 // 任务执行时间
	Result     string                 // 任务执行结果
}

// TaskNode is a node in the linked list for storing tasks
type TaskNode struct {
	Value *Task
	Next  *TaskNode
}

// FixedLengthTaskList is a fixed-length linked list for storing tasks
type FixedLengthTaskList struct {
	Head      *TaskNode
	Length    int
	MaxLength int
}

// NewFixedLengthLinkedList creates a new fixed-length linked list with the specified max length
func NewFixedLengthTaskList(maxLength int) *FixedLengthTaskList {
	return &FixedLengthTaskList{
		Head:      nil,
		Length:    0,
		MaxLength: maxLength,
	}
}

// AddNode adds a new node to the linked list，如果长度超过最大长度，则删除最后一个节点
func (fl *FixedLengthTaskList) AddTask(task *Task) {

	newNode := &TaskNode{Value: task, Next: nil}
	if fl.Head == nil {
		fl.Head = newNode
	} else {
		// Add the new node to the front of the list
		newNode.Next = fl.Head
		fl.Head = newNode
	}

	fl.Length++

	// Check if the length exceeds the maximum length
	if fl.Length > fl.MaxLength {
		// Remove the last node to maintain the fixed length
		fl.removeLastTask()
	}
}

// removeLastNode removes the last node from the linked list
func (fl *FixedLengthTaskList) removeLastTask() {
	if fl.Head == nil || fl.Head.Next == nil {
		return
	}

	current := fl.Head
	for current.Next.Next != nil {
		current = current.Next
	}

	// Remove the last node
	current.Next = nil
	fl.Length--
}

// ClearTasks removes all tasks from the list
func (fl *FixedLengthTaskList) ClearTasks() {
	fl.Head = nil
	fl.Length = 0
}
