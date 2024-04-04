package repeaters

import (
	"errors"
	"time"
)

var ErrRetry = errors.New("try again later")

type Strategy interface {
	Start() Retry
}

type Retry interface {
	Completed() bool
}

type UniformRetryStrategy struct {
	RetriesNum int
	Interval   time.Duration
}

func (r UniformRetryStrategy) Start() Retry {
	return &UniformRetry{
		RetriesNum: r.RetriesNum,
		Interval:   r.Interval,
	}
}

type UniformRetry struct {
	RetriesNum int
	Interval   time.Duration
}

func (r *UniformRetry) Completed() bool {
	if r.RetriesNum == 0 {
		return true
	}

	r.RetriesNum -= 1
	time.Sleep(r.Interval)

	return false
}
