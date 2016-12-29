package init

type InitTask struct {
	Start func()
	Stop func()
	Dependencies []string
	depsleft int
}

var inits map[string]*InitTask
func init() {
	inits = make(map[string]InitTask)
}

func AddTask(name string) (t *InitTask) {
	t = make(InitTask)
	inits[name]=t
}

func (t *InitTask) SetStart(s func()) {
	t.Start=s
}

func (t *InitTask) SetStop(s func()) {
	t.Stop=s
}

func (t *InitTask) AddDep(d string) {
	t.Dependencies=append(t.Dependencies,d)
	t.depsleft++
}

func (t *InitTask) AddDeps(d ...string) {
	for _,v := range d {
		t.AddDep(v)
	}
}
