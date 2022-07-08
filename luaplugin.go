package trelay

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type luaplugin struct {
	path string

	*lua.LState
	sync.Mutex
}
