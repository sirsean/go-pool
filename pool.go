package pool

import (
	"sync"
)

type WorkUnit interface {
	Perform()
}

type Pool struct {
	queueSize  int
	numWorkers int
	ch         chan WorkUnit
	wg         sync.WaitGroup
}

func NewPool(queueSize, numWorkers int) Pool {
	return Pool{
		queueSize:  queueSize,
		numWorkers: numWorkers,
		ch:         make(chan WorkUnit, queueSize),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for unit := range p.ch {
				unit.Perform()
			}
		}()
	}
}

func (p *Pool) Add(obj WorkUnit) {
	p.ch <- obj
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Close() {
	close(p.ch)
	p.Wait()
}
