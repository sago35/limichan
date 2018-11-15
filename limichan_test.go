package limichan

import (
	"log"
	"sync"
	"testing"
	"time"
)

type job struct {
	index   int
	deferFn func()
}

type w1 struct {
	name   string
	waitMs int
}

func (w *w1) Fn(lj LimichanJob) {
	j := lj.(job)
	defer j.deferFn()

	start := time.Now()
	time.Sleep(time.Duration(w.waitMs) * time.Millisecond)
	log.Printf("%3d : %s %3d (%3dms)", j.index, w.name, w.waitMs, time.Since(start).Nanoseconds()/1000000)
}

func TestLimichan(t *testing.T) {
	l := New(5)
	l.AddWorker(&w1{name: "w11", waitMs: 10})
	l.AddWorker(&w1{name: "w12", waitMs: 100})
	l.AddWorker(&w1{name: "w21", waitMs: 30})
	l.AddWorker(&w1{name: "w22", waitMs: 70})
	l.AddWorker(&w1{name: "w23", waitMs: 90})

	var wg sync.WaitGroup
	jobs := 12
	wg.Add(jobs)
	for i := 0; i < jobs; i++ {
		l.Do(job{
			index: i,
			deferFn: func() {
				wg.Done()
			},
		})
	}

	wg.Wait()

	//t.Error("error")
}
