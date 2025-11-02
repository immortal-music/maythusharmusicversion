package ubot

import "github.com/immortal-music/maythusharmusicversion/ntgcalls"

func (ctx *Context) Calls() map[int64]*ntgcalls.CallInfo {
	return ctx.binding.Calls()
}
