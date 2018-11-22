package limichan

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

type worker struct {
	name   string
	waitMs int
}

type myJob struct {
	index int
}

func (w *worker) Do(ctx context.Context, j Job) error {
	job, ok := j.(myJob)
	if !ok {
		return fmt.Errorf("Do error")
	}

	log.Printf("start %d", job.index)

	time.Sleep(time.Duration(job.index*50) * time.Millisecond)

	select {
	case <-ctx.Done():
		log.Printf("job %2d worker %s %s", job.index, w.name, ctx.Err().Error())
		return ctx.Err()
	default:
	}

	log.Printf("done  %d", job.index)

	return nil
}

func TestLimichan(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	l := New(ctx, MaxWorker(5))

	w1 := &worker{name: "w1", waitMs: 10}
	w2 := &worker{name: "w2", waitMs: 30}
	l.AddWorker(w1)
	l.AddWorker(w1)
	l.AddWorker(w1)
	l.AddWorker(w2)
	l.AddWorker(w2)

	jobs := 12
	for i := 0; i < jobs; i++ {
		i := i
		<-time.After(5 * time.Millisecond)

		err := l.Do(myJob{index: i})

		if err == context.DeadlineExceeded {
			// skip
			log.Printf("job %2d %s", i, err.Error())
		} else if err != nil {
			log.Printf("job %2d %s", i, err.Error())
		}
	}

	l.Wait()

	if len(l.Errors()) > 0 {
		log.Print("--")
		for _, err := range l.Errors() {
			log.Print(err.Error())
		}
	}

	//t.Error("error")
}
