package context

import (
	"context"
	"sync"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

type Context struct {
	TimerCtx context.Context
	Cancel   context.CancelFunc
	Keys     map[string]interface{}

	lock sync.RWMutex
}

const (
	TimeOut = 1 * time.Second
)

func NewContext() *Context {
	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	return &Context{TimerCtx: ctx, Cancel: cancel, Keys: make(map[string]interface{})}
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.lock.Lock()
	c.Keys[key] = value
	c.lock.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.lock.RLock()
	value, exists = c.Keys[key]
	c.lock.RUnlock()
	return
}

func Transform(c *gin.Context) *Context {
	if c == nil {
		return nil
	}
	if v, ok := c.Get("context"); ok {
		if ctx, isOk := v.(*Context); isOk {
			return ctx
		}
	}
	glog.Error("failed to get context")
	return nil
}
