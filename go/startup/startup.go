package startup

import (
	"log"

	"../util"
)

//Task represents a task to be run on startup
type Task struct {
	Time uint8
	Task func() error
}

var tasks []*Task

//NewTask allocates & initializes a Task
func NewTask(time uint8, task func() error) (t *Task) {
	t = new(Task)
	t.Time = time
	t.Task = task
	tasks = append(tasks, t)
	return
}

func init() {
	tasks = []*Task{}
}

//Start starts the program
func Start() {
	log.Println("Initiating startup. . .")
	for i := uint8(0); i < 255; i++ {
		for _, t := range tasks {
			if t.Time == i {
				util.UhOh(t.Task())
			}
		}
	}
	log.Println("Startup complete!")
}
