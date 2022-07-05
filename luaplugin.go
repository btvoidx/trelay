package trelay

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type luaplugin struct {
	*lua.LState
	sync.Mutex
}
