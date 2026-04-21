package redis

import (
	"github.com/go-redsync/redsync/v4"
)

func NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return Redsync.NewMutex("lock:"+name, options...)
}
