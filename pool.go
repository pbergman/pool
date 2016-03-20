package pool

import (
	"sync/atomic"
)

// Pool basic struct for the request pool
type Pool struct {
	concurrent uint64            // count of max concurrent
	running    uint64            // count of running
	state      State             // State of crawler
	Success    chan JobInterface // chanel for successful requests
	Error      chan JobInterface // chanel for failed requests
	Client     ClientInterface   // if nul will use Client wrapper
}

func NewPool(c uint64) *Pool {
	return &Pool{concurrent: c, running: 0}
}

// WorkersRunning returns the current running workers
func (p *Pool) WorkersRunning() uint64 {
	return atomic.LoadUint64(&p.running)
}

// MatchOneState will check if one of given (range of) state(s) match the current state, so for example :
// MatchOneState(RUNNING|QUEUE_HAS_SUCCESS) will return true when state is: RUNNING AND OR QUEUE_HAS_FAILURE
func (p *Pool) MatchOneState(s State) bool {
	a := p.GetState()
	if b := s & a; b <= 0 {
		return false
	} else {
		return b == (b & a)
	}
}

// GetState returns the current state of pool loop
func (p *Pool) GetState() State {
	return (State)(atomic.LoadUint32((*uint32)(&p.state)))
}
// HashState checks is a state bit is set
func (p *Pool) HashState(s State) bool {
	return s == (s & p.GetState())
}
// setState sets the state
func (p *Pool) setState(s State) {
	atomic.StoreUint32((*uint32)(&p.state), (uint32)(s))
}
// addState adds a bit to the state
func (p *Pool) addState(s State) {
	atomic.StoreUint32((*uint32)(&p.state), (uint32)(p.GetState())|(uint32)(s))
}
// removeState removes a state bit from the state
func (p *Pool) removeState(s State) {
	atomic.StoreUint32((*uint32)(&p.state), (uint32)(p.GetState())^(uint32)(s))
}
// Close basic implementation of the io.Close and will close the
// channels and set the state so it can be reused again.
func (p *Pool) Close() error {
	p.setState(State(0))
	close(p.Success)
	close(p.Error)
	return nil
}

func (p *Pool) getClient() ClientInterface {
	if p.Client == nil {
		p.Client = NewClient()
	}
	return p.Client
}

func (p *Pool) Start(jobs <-chan JobInterface, bufferSize int) {
	p.Success = make(chan JobInterface, bufferSize)
	p.Error = make(chan JobInterface, bufferSize)
	p.setState(RUNNING)
	// Keep track of queue
	go func() {
		for p.MatchOneState(RUNNING | CHAN_HAS_FAILURE | CHAN_HAS_SUCCESS) {
			if len(p.Error) == 0 && p.HashState(CHAN_HAS_FAILURE|FINISHED) {
				p.removeState(CHAN_HAS_FAILURE)
			}
			if len(p.Success) == 0 && p.HashState(CHAN_HAS_SUCCESS|FINISHED) {
				p.removeState(CHAN_HAS_SUCCESS)
			}
		}
	}()
	for i := uint64(0); i < p.concurrent; i++ {
		go func() {
			atomic.AddUint64(&p.running, 1)
			for job := range jobs {
				if resp, err := p.getClient().Do(job.GetRequest()); err != nil {
					job.SetError(&err)
					p.addState(CHAN_HAS_FAILURE)
					p.Error <- job
				} else {
					job.SetResponse(resp)
					p.addState(CHAN_HAS_SUCCESS)
					p.Success <- job
				}
			}
			if r := atomic.AddUint64(&p.running, ^uint64(0)); r == 0 {
				// Last worker, set state to FINISHED
				p.setState(p.GetState() | FINISHED ^ RUNNING)
			}
		}()
	}
}
