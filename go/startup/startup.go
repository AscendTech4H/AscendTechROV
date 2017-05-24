//Package startup provides a numerically ordered startup sequence
package startup

import (
	"flag"
	"log"

	"github.com/AscendTech4H/AscendTechROV/go/util"
)

//Task represents a task to be run on startup
type Task struct {
	Time uint8
	Task func() error
}

var tasks []*Task

//NewTask allocates & initializes a Task
//Time:
//	200-249: Setup for web interface
//	250: Start web interface
func NewTask(time uint8, task func() error) (t *Task) {
	t = new(Task)
	t.Time = time
	t.Task = task
	tasks = append(tasks, t)
	return
}

func init() {
	tasks = []*Task{}

	NewTask(2, func() error { //Parse flags
		flag.Parse()
		return nil
	})
}

//Start starts the robot
func Start() {
	log.Println("Initiating startup. . .")
	for i := uint8(0); i < 255; i++ {
		log.Printf("Startup level %d", i)
		for _, t := range tasks {
			if t.Time == i {
				util.UhOh(t.Task())
			}
		}
	}
	log.Println("Startup complete!")
}
