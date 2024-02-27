package agent

type status int

const (
	ALIVE   status = iota // значит что можно поручить задачку
	WORKING               // значит что живой, но работает

)

type Agent struct {
	Id     int
	Status status
}

func (a *Agent) Run() {

}
