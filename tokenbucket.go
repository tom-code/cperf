package main

import (
  "time"
)

type tokenBucket struct {
  tokens int64
  speed float64
  currentSpeed float64
  remain float64
  lastTime time.Time
  cancel bool
}
  
func (b *tokenBucket) get() bool {
  for {
    if b.tokens > 0 {
      b.tokens = b.tokens - 1
      return true
    }
    now := time.Now()
    timeDiff := now.Sub(b.lastTime)
    give := timeDiff.Seconds()*b.currentSpeed + b.remain
    if give > 1 {
      whole := int64(give)
      b.remain = give - float64(whole)
      b.tokens = b.tokens + whole
      b.lastTime = now
      continue
    }
    if b.tokens == 0 {
      time.Sleep(5*time.Millisecond)
    }
    if b.cancel {
      return false
    }
  }
}
  
func newBucket(speed int) *tokenBucket {
  bucket := &tokenBucket {
    tokens: 0,
    lastTime: time.Now(),
    speed: float64(speed),
    remain: 0,
    cancel: false,
  }
  go func() {
    for {
      time.Sleep(500*time.Millisecond)
      if bucket.currentSpeed > (bucket.speed - 100) {
        bucket.currentSpeed = bucket.speed
      } else if bucket.currentSpeed < bucket.speed {
        bucket.currentSpeed += 100
      }
    }
  }()
  return bucket
}
