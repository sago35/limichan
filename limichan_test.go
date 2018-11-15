package limichan

import (
	"log"
	"sync"
	"testing"
	"time"
)

type worker struct {
	name   string
	waitMs int
}

func TestLimichan(t *testing.T) {
	l := New(5)

	w1 := worker{name: "w1", waitMs: 10}
	w2 := worker{name: "w2", waitMs: 30}
	l.AddWorker(w1)
	l.AddWorker(w1)
	l.AddWorker(w1)
	l.AddWorker(w2)
	l.AddWorker(w2)

	var wg sync.WaitGroup
	jobs := 12
	wg.Add(jobs)
	for i := 0; i < jobs; i++ {
		l.Do(func(wk Worker) {
			defer wg.Done()

			switch w := wk.(type) {
			case worker:
				start := time.Now()
				time.Sleep(time.Duration(w.waitMs) * time.Millisecond)
				log.Printf("%s (%3dms)", w.name, time.Since(start).Nanoseconds()/1000000)
			}

		})
	}

	wg.Wait()

	//t.Error("error")
}
