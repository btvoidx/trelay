package trelay

import (
	"bufio"
	"io"
	"sync"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

const globalLuaTableName = "trelay"

type lplugin struct {
	LState   *lua.LState
	mu       sync.Mutex
	bytecode *lua.FunctionProto
}

func (lp *lplugin) Call(fname string, ctx func(*lplugin) int) {
	lp.mu.Lock()
	defer lp.mu.Unlock()

	gt, ok := lp.LState.G.Global.RawGetString(globalLuaTableName).(*lua.LTable)
	if !ok {
		println("ERROR TRELAY NOT A TABLE: ", fname) //todo: logging
		return
	}

	fn, ok := gt.RawGetString(fname).(*lua.LFunction)
	if !ok {
		return
	}

	lp.LState.Push(fn)
	if err := lp.LState.PCall(ctx(lp), lua.MultRet, nil); err != nil {
		println("ERROR IN A FUNCTION: ", fname) //todo: logging
		return
	}
}

func (lp *lplugin) load(t *trelay) error {
	lp.LState = lua.NewState()
	globalTable := lp.LState.NewTable()
	lp.LState.SetGlobal(globalLuaTableName, globalTable)

	lfunc := lp.LState.NewFunctionFromProto(lp.bytecode)
	lp.LState.Push(lfunc)
	return lp.LState.PCall(0, lua.MultRet, nil)
}

func (lp *lplugin) compile(r io.Reader, chunkName string) error {
	chunk, err := parse.Parse(bufio.NewReader(r), chunkName)
	if err != nil {
		return err
	}

	lp.bytecode, err = lua.Compile(chunk, chunkName)
	if err != nil {
		return err
	}

	return nil
}

func luafnSetPacketHandled(L *lua.LState, handled *bool) *lua.LFunction {
	return L.NewFunction(func(l *lua.LState) int {
		*handled = true
		return 0
	})
}
