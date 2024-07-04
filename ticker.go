package timerticker

import (
	"context"
	"time"
)

type TimerTicker struct {
	startTime *time.Time
	duration  time.Duration
	ticker    *time.Ticker
	ch        chan time.Time
}

func NewTimerTicker(duration time.Duration) *TimerTicker {
	if duration <= 0 {
		panic("non-positive interval for NewTicker")
	}
	return &TimerTicker{
		startTime: nil,
		duration:  duration,
		ticker:    nil,
		ch:        make(chan time.Time),
	}
}

func (t *TimerTicker) SetLaunchTime(ctx context.Context, tm time.Time) {
	if t.startTime != nil {
		panic("not allow set start time twice")
	}

	now := time.Now()
	if now.After(tm) {
		panic("set start time before now")
	}

	sleepDuration := tm.Sub(now)
	go func(ctx context.Context) {
		time.Sleep(sleepDuration)
		t.ticker = time.NewTicker(t.duration)
		now := time.Now()
		t.ch <- now
		for {
			select {
			case <-ctx.Done():
				return
			case tm := <-t.ticker.C:
				t.ch <- tm
			}
		}
	}(ctx)
}

func (t *TimerTicker) GetChan() <-chan time.Time {
	return t.ch
}

func (t *TimerTicker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

func (t *TimerTicker) Reset(duration time.Duration) {
	t.duration = duration
	if t.ticker != nil {
		t.ticker.Reset(duration)
	}
	return
}
