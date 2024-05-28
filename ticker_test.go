package timerticker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimerTicker_SetLaunchTimeTwice(t1 *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t1, "not allow set start time twice", err)
		}
	}()
	ticker := NewTimerTicker(time.Second)
	ticker.SetLaunchTime(context.Background(), time.Now().Add(time.Second))
	ticker.SetLaunchTime(context.Background(), time.Now().Add(time.Second))
}

func TestTimerTicker_SetLaunchTimeBeforeNow(t1 *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t1, "set start time before now", err)
		}
	}()
	ticker := NewTimerTicker(time.Second)
	ticker.SetLaunchTime(context.Background(), time.Now())
}

func TestNewTimerTicker(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, "non-positive interval for NewTicker", err)
		}
	}()
	var duration = -time.Second
	_ = NewTimerTicker(duration)
}

func TestNewTimerTicker2(t *testing.T) {
	now := time.Now()
	stamp := now.UnixMilli()
	t.Log(now)
	closeChan := make(chan struct{})
	go func() {
		time.Sleep(30 * time.Second)
		closeChan <- struct{}{}
	}()
	ticker := NewTimerTicker(time.Second)
	ticker.SetLaunchTime(context.Background(), now.Add(5*time.Second))
	ch := ticker.GetChan()

	exceptStamp := stamp + 5000
	for {
		select {
		case tm := <-ch:
			t.Log(tm)
			// 60ms 误差
			assert.LessOrEqual(t, tm.UnixMilli(), exceptStamp+30)
			assert.GreaterOrEqual(t, tm.UnixMilli(), exceptStamp-30)
			exceptStamp += 1000
		case <-closeChan:
			ticker.Stop()
			return
		}
	}
}

func TestTimerTicker_Reset(t *testing.T) {
	resetChan := make(chan struct{})
	closeChan := make(chan struct{})
	go func() {
		time.Sleep(time.Second * 10)
		resetChan <- struct{}{}
	}()
	go func() {
		time.Sleep(time.Second * 28)
		closeChan <- struct{}{}
	}()
	now := time.Now()
	ticker := NewTimerTicker(2 * time.Second)
	ticker.SetLaunchTime(context.Background(), now.Add(time.Second))
	hitArray := []int64{
		now.UnixMilli() + 1000,
		now.UnixMilli() + 3000,
		now.UnixMilli() + 5000,
		now.UnixMilli() + 7000,
		now.UnixMilli() + 9000,
		now.UnixMilli() + 15000,
		now.UnixMilli() + 20000,
		now.UnixMilli() + 25000,
	}
	i := 0
	ch := ticker.GetChan()
	for {
		select {
		case tm := <-ch:
			t.Log(tm)
			assert.LessOrEqual(t, tm.UnixMilli(), hitArray[i]+30)
			assert.GreaterOrEqual(t, tm.UnixMilli(), hitArray[i]-30)
			i++
		case <-resetChan:
			ticker.Reset(time.Second * 5)
		case <-closeChan:
			return
		}
	}

}
