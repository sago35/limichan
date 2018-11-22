package limichan

import (
	"context"
	"fmt"
	"sync"
)

// Worker ...
type Worker interface {
	Do(context.Context, Job) error
}

// Job ...
type Job interface {
}

// Limichan ...
type Limichan struct {
	worker chan Worker
	ctx    context.Context
	wg     sync.WaitGroup
	errors []error
}

// Option ...
type Option func(*options)

type options struct {
	maxWorker int
}

var defaultLimichanOption = options{
	maxWorker: 1000,
}

// MaxWorker ...
func MaxWorker(m int) Option {
	return func(o *options) {
		o.maxWorker = m
	}
}

// New ...
func New(ctx context.Context, opt ...Option) *Limichan {
	opts := defaultLimichanOption
	for _, o := range opt {
		o(&opts)
	}

	lc := &Limichan{
		worker: make(chan Worker, opts.maxWorker),
		ctx:    ctx,
	}

	return lc
}

// AddWorker ...
func (lc *Limichan) AddWorker(w Worker) error {
	select {
	case lc.worker <- w:
		return nil
	default:
	}

	return fmt.Errorf("AddWorker() failed")
}

// Do ...
func (lc *Limichan) Do(j Job) error {
	select {
	case <-lc.ctx.Done():
		lc.errors = append(lc.errors, lc.ctx.Err())
		return lc.ctx.Err()
	case w := <-lc.worker:
		lc.wg.Add(1)
		go func() {
			defer lc.wg.Done()
			err := w.Do(lc.ctx, j)
			if err != nil {
				lc.errors = append(lc.errors, err)
			}
			lc.worker <- w
		}()
	}
	return nil
}

// Wait ...
func (lc *Limichan) Wait() {
	lc.wg.Wait()
}

// Errors ...
func (lc *Limichan) Errors() []error {
	return lc.errors
}

// ClearErrors ...
func (lc *Limichan) ClearErrors() {
	lc.errors = lc.errors[:0]
}
