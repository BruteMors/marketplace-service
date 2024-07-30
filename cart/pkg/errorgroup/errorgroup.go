package errorgroup

import (
	"context"
	"sync"
)

type ErrGroup struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	ctx     context.Context
	cancel  func()
}

func NewErrGroup(ctx context.Context) *ErrGroup {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *ErrGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

func (g *ErrGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

func (g *ErrGroup) Context() context.Context {
	return g.ctx
}
