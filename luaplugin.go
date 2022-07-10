package trelay

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type luaplugin struct {
	path string

	lstate *lua.LState
	mu     sync.Mutex
}

func (lp *luaplugin) Call(fname string, ctxfn lua.LGFunction) {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	gt, ok := lp.lstate.G.Global.RawGetString(globalLuaTableName).(*lua.LTable)
	if !ok {
		println("ERROR TRELAY NOT A TABLE: ", fname) //todo logging
		return
	}

	fn, ok := gt.RawGetString(fname).(*lua.LFunction)
	if !ok {
		return
	}

	lp.lstate.Push(fn)
	if err := lp.lstate.PCall(ctxfn(lp.lstate), lua.MultRet, nil); err != nil {
		println("ERROR IN A FUNCTION: ", fname) //todo logging
		return
	}
}

func luafnSetPacketHandled(L *lua.LState, handled *bool) *lua.LFunction {
	return L.NewFunction(func(l *lua.LState) int {
		*handled = true
		return 0
	})
}
