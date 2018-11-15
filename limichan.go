package limichan

// Worker ...
type Worker interface {
	Fn(lj LimichanJob)
}

// LimichanJob ...
type LimichanJob interface {
}

// Limichan ...
type Limichan struct {
	worker chan Worker
}

// New ...
func New(maxWorker int) *Limichan {
	lc := &Limichan{
		worker: make(chan Worker, maxWorker),
	}

	return lc
}

// AddWorker ...
func (lc *Limichan) AddWorker(w Worker) {
	lc.worker <- w
}

// Do ...
func (lc *Limichan) Do(lj LimichanJob) {
	select {
	case w := <-lc.worker:
		go func() {
			w.Fn(lj)
			lc.worker <- w
		}()
	}
}
