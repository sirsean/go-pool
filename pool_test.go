package pool

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewPool(t *testing.T) {
	p := NewPool(1, 2)
	assert.Equal(t, 1, p.queueSize)
	assert.Equal(t, 2, p.numWorkers)
	assert.NotNil(t, p.ch)
}

type Summer struct {
	sync.RWMutex
	Value int
}

func (s *Summer) Plus(n int) {
	s.Lock()
	s.Value += n
	s.Unlock()
}

type PlusUnit struct {
	Sum *Summer
	Num int
}

func (u PlusUnit) Perform() {
	u.Sum.Plus(u.Num)
}

func TestAddBeforeStart(t *testing.T) {
	summer := Summer{}

	p := NewPool(1, 1)
	p.Add(PlusUnit{
		Sum: &summer,
		Num: 1,
	})
	p.Close()

	assert.Equal(t, 0, summer.Value)
}

func TestSharedData(t *testing.T) {
	summer := Summer{}

	p := NewPool(1, 1)
	p.Start()
	p.Add(PlusUnit{
		Sum: &summer,
		Num: 2,
	})
	p.Add(PlusUnit{
		Sum: &summer,
		Num: 3,
	})
	p.Add(PlusUnit{
		Sum: &summer,
		Num: 4,
	})
	p.Close()

	assert.Equal(t, 9, summer.Value)
}

type OutChanUnit struct {
	Chan chan int
	In   int
}

func (u OutChanUnit) Perform() {
	u.Chan <- 2 * u.In
}

func TestOutputChannel(t *testing.T) {
	p := NewPool(5, 5)
	ch := make(chan int, 10)

	p.Start()

	p.Add(OutChanUnit{
		Chan: ch,
		In:   1,
	})
	p.Add(OutChanUnit{
		Chan: ch,
		In:   2,
	})

	p.Close()

	close(ch)

	total := 0
	for n := range ch {
		total += n
	}
	assert.Equal(t, 6, total)
}
