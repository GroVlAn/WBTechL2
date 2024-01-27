package pattern

import "fmt"

/*
Преймущества:
	- Убирает прямую зависимость между объектами, вызывающими операции, и объектами, которые их непосредственно выполняют.
	- Позволяет реализовать простую отмену и повтор операций.
	- Позволяет реализовать отложенный запуск операций.
	- Позволяет собирать сложные команды из простых.
	- Реализует принцип открытости/закрытости.
Недостатки:
	- Усложняет код программы из-за введения множества дополнительных классов.
*/

// TaskService - сервис для управления задачами
type TaskService struct {
	Tasks []string
}

// AddTaskCommand - команда для добавления задачи
type AddTaskCommand struct {
	TaskService *TaskService
	Task        string
}

// RemoveTaskCommand - команда для удаления задачи
type RemoveTaskCommand struct {
	TaskService *TaskService
	IndexTask   int
}

// Invoker - структура, вызывающая команду
type Invoker struct {
	Command Commander
}

func NewAddTaskCommand(taskService *TaskService, Task string) *AddTaskCommand {
	return &AddTaskCommand{
		TaskService: taskService,
		Task:        Task,
	}
}

// Execute - выполняем добавление новой задачи
func (atc *AddTaskCommand) Execute() {
	fmt.Printf("Adding task: %s\n", atc.Task)
	atc.TaskService.Tasks = append(atc.TaskService.Tasks, atc.Task)
}

func NewRemoveTaskCommand(taskService *TaskService, indexTask int) *RemoveTaskCommand {
	return &RemoveTaskCommand{
		TaskService: taskService,
		IndexTask:   indexTask,
	}
}

// Execute - выполняем удаление задачи
func (atc *RemoveTaskCommand) Execute() {
	if len(atc.TaskService.Tasks) < atc.IndexTask || atc.IndexTask < 0 {
		return
	}
	fmt.Printf("Removing task: %s\n", atc.TaskService.Tasks[atc.IndexTask])
	atc.TaskService.Tasks = append(
		atc.TaskService.Tasks[:atc.IndexTask],
		atc.TaskService.Tasks[atc.IndexTask+1:]...)
}

// SetCommand - передаём исподнителю реализацию интерфейса команды
func (i *Invoker) SetCommand(command Commander) {
	i.Command = command
}

// Run - выполняем конкретную задачу
func (i *Invoker) Run() {
	i.Command.Execute()
}

// Commander - интерфейс для определения методов команд
type Commander interface {
	Execute()
}
