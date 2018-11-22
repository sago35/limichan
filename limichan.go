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
	worker  chan Worker
	ctx     context.Context
	wg      sync.WaitGroup
	err     error
	cancel  func()
	errOnce sync.Once
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
func New(ctx context.Context, opt ...Option) (*Limichan, context.Context) {
	opts := defaultLimichanOption
	for _, o := range opt {
		o(&opts)
	}

	ctx, cancel := context.WithCancel(ctx)

	lc := &Limichan{
		worker: make(chan Worker, opts.maxWorker),
		ctx:    ctx,
		cancel: cancel,
	}

	return lc, ctx
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
		return lc.ctx.Err()
	case w := <-lc.worker:
		lc.wg.Add(1)
		go func() {
			defer lc.wg.Done()
			err := w.Do(lc.ctx, j)
			if err != nil {
				lc.errOnce.Do(func() {
					lc.err = err
					if lc.cancel != nil {
						lc.cancel()
					}
				})
			}
			lc.worker <- w
		}()
	}
	return nil
}

// Wait ...
func (lc *Limichan) Wait() error {
	lc.wg.Wait()
	if lc.cancel != nil {
		lc.cancel()
	}
	return lc.err
}
