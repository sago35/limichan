package limichan

// Worker ...
type Worker interface {
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
func (lc *Limichan) Do(fn func(Worker)) {
	select {
	case w := <-lc.worker:
		go func() {
			fn(w)
			lc.worker <- w
		}()
	}
}
