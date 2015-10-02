package hookfs

// Base interface for user-written hooks.
//
// You have to implement HookXXX (e.g. HookOnOpen, HookOnRead, HookOnWrite, ..) interfaces.
type Hook interface{}

// Context objects for interaction between prehooks and posthooks.
type HookContext interface{}

// Called on mount. This also implements Hook.
type HookWithInit interface {
	Init() (err error)
}

// Called on open. This also implements Hook.
type HookOnOpen interface {
	// if hooked is true, the real open() would not be called
	PreOpen(path string, flags uint32) (err error, hooked bool, ctx HookContext)
	PostOpen(realRetCode int32, prehookCtx HookContext) (err error, hooked bool)
}

// Called on read. This also implements Hook.
type HookOnRead interface {
	// if hooked is true, the real read() would not be called
	PreRead(path string, length int64, offset int64) (buf []byte, err error, hooked bool, ctx HookContext)
	PostRead(realRetCode int32, realBuf []byte, prehookCtx HookContext) (buf []byte, err error, hooked bool)
}

// Called on write. This also implements Hook.
// BUG(AkihiroSuda): HookOnWrite is not yet implemented. (Of course, we also need many more things such as Stat.)
type HookOnWrite interface {
	// if hooked is true, the real write() would not be called
	PreWrite(path string, buf []byte, offset int64) (err error, hooked bool, ctx HookContext)
	PostWrite(realRetCode int32, prehookCtx HookContext) (err error, hooked bool)
}
