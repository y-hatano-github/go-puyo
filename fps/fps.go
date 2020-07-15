package fps

import (
	"time"
)

type FPS struct {
	startTime int64
	newTime   int64
	fps       int64
	gap       int64
}

func NewFPS(fps int64) *FPS {
	return &FPS{
		startTime: 0,
		newTime:   time.Now().UnixNano() / int64(time.Millisecond),
		fps:       fps,
		gap:       0,
	}
}

func (b *FPS) Update() {
	b.startTime = b.newTime
}

func (b *FPS) Wait() {
	sleepTime := (1000 / b.fps) - (time.Now().UnixNano()/int64(time.Millisecond) - b.startTime) - b.gap
	wait := true
	for wait {
		//time.Sleep(time.Duration(sleepTime))
		wait = (time.Now().UnixNano()/int64(time.Millisecond))-b.startTime <= sleepTime
	}

	b.newTime = time.Now().UnixNano() / int64(time.Millisecond)
	b.gap = b.newTime - b.startTime - sleepTime
}
