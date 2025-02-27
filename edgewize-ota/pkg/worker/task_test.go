package worker_manager

import (
	"testing"
)

// TestNewFixedLengthTaskList tests the creation of a new task list
func TestNewFixedLengthTaskList(t *testing.T) {
	maxLength := 5
	taskList := NewFixedLengthTaskList(maxLength)
	if taskList.MaxLength != maxLength {
		t.Errorf("Expected MaxLength %d, got %d", maxLength, taskList.MaxLength)
	}
	if taskList.Length != 0 {
		t.Errorf("Expected initial length to be 0, got %d", taskList.Length)
	}
	if taskList.Head != nil {
		t.Errorf("Expected initial head to be nil, got %v", taskList.Head)
	}
}

// TestAddTask tests adding tasks to the list
func TestAddTask(t *testing.T) {
	taskList := NewFixedLengthTaskList(3)
	task1 := &Task{ID: 1}
	task2 := &Task{ID: 2}
	task3 := &Task{ID: 3}
	task4 := &Task{ID: 4}

	// Add tasks and check list length
	taskList.AddTask(task1)
	if taskList.Length != 1 {
		t.Errorf("Expected length 1, got %d", taskList.Length)
	}

	taskList.AddTask(task2)
	taskList.AddTask(task3)
	taskList.AddTask(task4) // This should remove task1
	if taskList.Length != 3 {
		t.Errorf("Expected length 3, got %d", taskList.Length)
	}

	// Check if task1 was removed
	current := taskList.Head
	found := false
	for current != nil {
		if current.Value.ID == task1.ID {
			found = true
			break
		}
		current = current.Next
	}
	if found {
		t.Errorf("Task with ID %d should have been removed", task1.ID)
	}
}

// TestClearTasks tests clearing the task list
func TestClearTasks(t *testing.T) {
	taskList := NewFixedLengthTaskList(3)
	taskList.AddTask(&Task{ID: 1})
	taskList.AddTask(&Task{ID: 2})
	taskList.ClearTasks()

	if taskList.Length != 0 {
		t.Errorf("Expected length 0 after clearing, got %d", taskList.Length)
	}
	if taskList.Head != nil {
		t.Errorf("Expected head to be nil after clearing, got %v", taskList.Head)
	}
}
